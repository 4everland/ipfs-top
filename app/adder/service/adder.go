package service

import (
	"encoding/json"
	"fmt"
	"github.com/4everland/ipfs-servers/app/adder/conf"
	"github.com/4everland/ipfs-servers/third_party/coreunix"
	"github.com/4everland/ipfs-servers/third_party/dag"
	httpctx "github.com/go-kratos/kratos/v2/transport/http"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	exchange "github.com/ipfs/go-ipfs-exchange-interface"
	"github.com/ipfs/go-libipfs/files"
	coreiface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/options"
	"mime"
	"mime/multipart"
	"net/http"
	"path"
)

const (
	contentTypeHeader = "Content-Type"
	adderOutChanSize  = 8
)

type AddRequest struct {
	Quiet             bool   `json:"quiet,omitempty"`
	Quieter           bool   `json:"quieter,omitempty"`
	Silent            bool   `json:"silent,omitempty"`
	Progress          bool   `json:"progress,omitempty"`
	Trickle           bool   `json:"trickle,omitempty"`
	OnlyHash          bool   `json:"only-hash,omitempty"`
	WrapWithDirectory bool   `json:"wrap-with-directory,omitempty"`
	Chunker           string `json:"chunker" default:"size-262144"`
	RawLeaves         bool   `json:"raw-leaves,omitempty"`
	CidVersion        int    `json:"cid-version,omitempty"`
	Pin               bool   `json:"pin"`
}

type AddEvent struct {
	Bytes int64  `json:"Bytes,omitempty"`
	Hash  string `json:"Hash,omitempty"`
	Name  string `json:"Name,omitempty"`
	Size  string `json:"Size,omitempty"`
}

func (e AddEvent) Marshal() []byte {
	marshal, _ := json.Marshal(e)
	return marshal
}

type AdderService struct {
	unixfs *coreunix.UnixFsServer
}

func NewAdderService(unixfs *coreunix.UnixFsServer) *AdderService {
	return &AdderService{unixfs}
}

func NewBlockStore(config *conf.Data) blockstore.Blockstore {
	s, err := dag.NewBlockStore(config.BlockstoreUri)
	if err != nil {
		panic(err)
	}
	return s
}

func NewExchange(config *conf.Data) exchange.Interface {
	s, err := dag.NewGrpcRouting(config.ExchangeEndpoint)
	if err != nil {
		panic(err)
	}
	return s
}

func (a *AdderService) Add(ctx httpctx.Context) (err error) {
	w := ctx.Response()
	r := ctx.Request()
	contentType := r.Header.Get(contentTypeHeader)
	mediatype, _, _ := mime.ParseMediaType(contentType)

	var f files.Directory
	if mediatype == "multipart/form-data" {
		var reader *multipart.Reader
		if reader, err = r.MultipartReader(); err != nil {
			return
		}

		if f, err = files.NewFileFromPartReader(reader, mediatype); err != nil {
			return
		}
	}

	if f == nil {
		return fmt.Errorf("file can't be empty")
	}

	var addRequest AddRequest
	if err = ctx.BindQuery(&addRequest); err != nil {
		return
	}

	toadd := f
	if addRequest.WrapWithDirectory {
		toadd = files.NewSliceDirectory([]files.DirEntry{
			files.FileEntry("", f),
		})
	}

	opts := []options.UnixfsAddOption{
		options.Unixfs.Chunker(addRequest.Chunker),

		options.Unixfs.Pin(addRequest.Pin),
		options.Unixfs.HashOnly(addRequest.OnlyHash),

		options.Unixfs.Progress(addRequest.Progress),
		options.Unixfs.Silent(addRequest.Silent),
	}

	if addRequest.CidVersion == 1 {
		opts = append(opts, options.Unixfs.CidVersion(addRequest.CidVersion))
	}

	if addRequest.RawLeaves {
		opts = append(opts, options.Unixfs.RawLeaves(addRequest.RawLeaves))
	}

	if addRequest.Trickle {
		opts = append(opts, options.Unixfs.Layout(options.TrickleLayout))
	}

	opts = append(opts, nil)
	var added int

	quiet := addRequest.Quiet
	quieter := addRequest.Quieter
	quiet = quiet || quieter

	addit := toadd.Entries()
	lastEvent := AddEvent{}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Chunked-Output", "1")
	w.WriteHeader(http.StatusOK)
	for addit.Next() {
		_, dir := addit.Node().(files.Directory)
		errCh := make(chan error, 1)
		events := make(chan interface{}, adderOutChanSize)
		opts[len(opts)-1] = options.Unixfs.Events(events)

		go func() {
			var er error
			defer close(events)
			defer close(errCh)
			_, er = a.unixfs.Add(r.Context(), addit.Node(), opts...)
			if er != nil {
				errCh <- er
				return
			}
		}()

		for event := range events {
			output, ok := event.(*coreiface.AddEvent)
			if !ok {
				return fmt.Errorf("unknown event type")
			}

			h := ""
			if output.Path != nil {
				h = output.Path.Cid().String()
				//h = enc.Encode(output.Path.Cid())
			}

			if !dir && addit.Name() != "" {
				output.Name = addit.Name()
			} else {
				output.Name = path.Join(addit.Name(), output.Name)
			}

			addEvent := AddEvent{
				Name:  output.Name,
				Hash:  h,
				Bytes: output.Bytes,
				Size:  output.Size,
			}
			if len(addEvent.Hash) > 0 {
				lastEvent = addEvent
				if quieter {
					continue
				}
				_, _ = w.Write(addEvent.Marshal())
				_, _ = w.Write([]byte("\n"))
			} else {
				if !addRequest.Progress {
					continue
				}
				_, _ = w.Write(addEvent.Marshal())
				_, _ = w.Write([]byte("\n"))
			}

		}

		if err = <-errCh; err != nil {
			return
		}
		added++
	}

	if addit.Err() != nil {
		return addit.Err()
	}

	if added == 0 {
		return fmt.Errorf("expected a file argument")
	}
	if quieter {
		_, _ = w.Write(lastEvent.Marshal())
		_, _ = w.Write([]byte("\n"))
	}

	return nil
}
