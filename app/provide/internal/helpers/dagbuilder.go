package helpers

import (
	"context"
	"errors"
	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	dag "github.com/ipfs/go-merkledag"
	ft "github.com/ipfs/go-unixfs"
	pb "github.com/ipfs/go-unixfs/pb"
)

const blockSize = 256 * 1024

var (
	roughLinkBlockSize = 1 << 13    // 8KB
	roughLinkSize      = 34 + 8 + 5 // sha256 multihash + size + no name + protobuf framing

	// DefaultLinksPerBlock governs how the importer decides how many links there
	// will be per block. This calculation is based on expected distributions of:
	//   - the expected distribution of block sizes
	//   - the expected distribution of link sizes
	//   - desired access speed
	//
	// For now, we use:
	//
	//	var roughLinkBlockSize = 1 << 13 // 8KB
	//	var roughLinkSize = 34 + 8 + 5   // sha256 multihash + size + no name
	//	                                 // + protobuf framing
	//	var DefaultLinksPerBlock = (roughLinkBlockSize / roughLinkSize)
	//	                         = ( 8192 / 47 )
	//	                         = (approximately) 174
	DefaultLinksPerBlock = roughLinkBlockSize / roughLinkSize
)

// DagBuilderHelper wraps together a bunch of objects needed to
// efficiently create unixfs dag trees
type DagBuilderHelper struct {
	iter       *DagNodeIter
	dserv      ipld.DAGService
	recvdErr   error
	rawLeaves  bool
	nextData   ipld.Node // the next item to return.
	maxlinks   int
	cidBuilder cid.Builder
}

// DagBuilderParams wraps configuration options to create a DagBuilderHelper
// from a chunker.Splitter.
type DagBuilderParams struct {
	// Maximum number of links per intermediate node
	Maxlinks int

	// RawLeaves signifies that the importer should use raw ipld nodes as leaves
	// instead of using the unixfs TRaw type
	RawLeaves bool

	// CID Builder to use if set
	CidBuilder cid.Builder

	// DAGService to write blocks to (required)
	Dagserv ipld.DAGService
	Iter    *DagNodeIter
}

func (dbp *DagBuilderParams) New() (*DagBuilderHelper, error) {
	db := &DagBuilderHelper{
		iter:       dbp.Iter,
		dserv:      dbp.Dagserv,
		rawLeaves:  dbp.RawLeaves,
		cidBuilder: dbp.CidBuilder,
		maxlinks:   dbp.Maxlinks,
	}

	return db, nil
}

func (db *DagBuilderHelper) prepareNext() {
	if db.nextData != nil || db.recvdErr != nil {
		return
	}

	db.nextData, db.recvdErr = db.iter.Next()
}

func (db *DagBuilderHelper) Done() bool {
	db.prepareNext() // idempotent
	if db.recvdErr != nil {
		return false
	}
	return db.nextData == nil
}

func (db *DagBuilderHelper) Next() (ipld.Node, error) {
	db.prepareNext() // idempotent
	d := db.nextData
	db.nextData = nil // signal we've consumed it
	if db.recvdErr != nil {
		return nil, db.recvdErr
	}
	return d, nil
}

func (db *DagBuilderHelper) GetCidBuilder() cid.Builder {
	return db.cidBuilder
}

func (db *DagBuilderHelper) NewLeafDataNode() (node ipld.Node, dataSize uint64, err error) {
	node, err = db.Next()
	if err != nil {
		return nil, 0, err
	}
	s, _ := node.Size()

	return node, s, err
}

func (db *DagBuilderHelper) Add(node ipld.Node) error {
	if _, ok := node.(*Object); ok {
		return nil
	}
	return db.dserv.Add(context.TODO(), node)
}

func (db *DagBuilderHelper) Maxlinks() int {
	return db.maxlinks
}

type FSNodeOverDag struct {
	dag  *dag.ProtoNode
	file *ft.FSNode
}

func (db *DagBuilderHelper) NewFSNodeOverDag(fsNodeType pb.Data_DataType) *FSNodeOverDag {
	node := new(FSNodeOverDag)
	node.dag = new(dag.ProtoNode)
	node.dag.SetCidBuilder(db.GetCidBuilder())

	node.file = ft.NewFSNode(fsNodeType)

	return node
}

func (db *DagBuilderHelper) NewFSNFromDag(nd *dag.ProtoNode) (*FSNodeOverDag, error) {
	return NewFSNFromDag(nd)
}

func NewFSNFromDag(nd *dag.ProtoNode) (*FSNodeOverDag, error) {
	mb, err := ft.FSNodeFromBytes(nd.Data())
	if err != nil {
		return nil, err
	}

	return &FSNodeOverDag{
		dag:  nd,
		file: mb,
	}, nil
}

func (n *FSNodeOverDag) AddChild(child ipld.Node, fileSize uint64, db *DagBuilderHelper) error {
	err := n.dag.AddNodeLink("", child)
	if err != nil {
		return err
	}

	n.file.AddBlockSize(fileSize)

	return db.Add(child)
}

func (n *FSNodeOverDag) Commit() (ipld.Node, error) {
	fileData, err := n.file.GetBytes()
	if err != nil {
		return nil, err
	}
	n.dag.SetData(fileData)

	return n.dag, nil
}

func (n *FSNodeOverDag) NumChildren() int {
	return n.file.NumChildren()
}

func (n *FSNodeOverDag) FileSize() uint64 {
	return n.file.FileSize()
}

type DagNodeIter struct {
	DAG   ipld.NodeGetter
	roots []ipld.Node

	ch    chan ipld.Node
	errCh chan error
	ctx   context.Context
}

func NewDagNodeIter(ctx context.Context, roots []ipld.Node, ng ipld.NodeGetter) *DagNodeIter {
	iter := &DagNodeIter{
		DAG:   ng,
		roots: roots,
		ch:    make(chan ipld.Node),
		errCh: make(chan error),
		ctx:   ctx,
	}

	go iter.init()

	return iter
}

func (t *DagNodeIter) init() {
	for _, root := range t.roots {
		t.next(root)
	}
	close(t.ch)
	close(t.errCh)
}

func (t *DagNodeIter) next(curr ipld.Node) {
	for _, l := range curr.Links() {
		node, err := t.getNode(l)
		if err != nil {
			select {
			case t.errCh <- err:
				return
			case <-t.ctx.Done():
				return
			}
		}
		if node == nil {
			continue
		}
		if len(node.Links()) > 0 {
			t.next(node)
		} else {
			select {
			case t.ch <- node:
			case <-t.ctx.Done():
				return
			}
		}
	}

	if len(curr.Links()) == 0 {
		select {
		case t.ch <- curr:
		case <-t.ctx.Done():
		}
	}
}

func (t *DagNodeIter) Next() (ipld.Node, error) {
	select {
	case nd := <-t.ch:
		return nd, nil
	case err := <-t.errCh:
		return nil, err
	case <-t.ctx.Done():
		return nil, t.ctx.Err()
	}
}

func (t *DagNodeIter) getNode(link *ipld.Link) (ipld.Node, error) {
	if link.Size <= blockSize {
		return NewObject(link.Name, int64(link.Size), link.Cid), nil
	}
	return link.GetNode(t.ctx, t.DAG)
}

var ErrNotSupport = errors.New("not support")

type Object struct {
	key  string
	cid  cid.Cid
	size int64
}

func NewObject(key string, size int64, c cid.Cid) *Object {
	return &Object{
		key:  key,
		cid:  c,
		size: size,
	}
}

func (o *Object) GetKey() string {
	return o.key
}

func (o *Object) Cid() cid.Cid {
	return o.cid
}

func (o *Object) Size() (uint64, error) {
	return uint64(o.size), nil
}

func (o *Object) String() string {
	return o.Cid().String()
}

func (o *Object) Copy() ipld.Node {
	return &Object{
		key:  o.key,
		cid:  o.cid,
		size: o.size,
	}
}

func (o *Object) Resolve([]string) (interface{}, []string, error) {
	return nil, nil, ErrNotSupport
}

func (o *Object) Tree(string, int) []string {
	return nil
}

func (o *Object) ResolveLink([]string) (*ipld.Link, []string, error) {
	return nil, nil, ErrNotSupport
}

func (o *Object) Links() []*ipld.Link {
	return nil
}

func (o *Object) Loggable() map[string]interface{} {
	return nil
}

func (o *Object) RawData() []byte {
	return nil
}

func (o *Object) Stat() (*ipld.NodeStat, error) {
	return &ipld.NodeStat{}, nil
}
