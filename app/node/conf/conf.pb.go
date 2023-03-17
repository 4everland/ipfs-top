// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.5
// source: app/node/conf/conf.proto

package conf

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Bootstrap struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Server *Server `protobuf:"bytes,1,opt,name=server,proto3" json:"server,omitempty"`
	Data   *Data   `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Bootstrap) Reset() {
	*x = Bootstrap{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_node_conf_conf_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Bootstrap) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Bootstrap) ProtoMessage() {}

func (x *Bootstrap) ProtoReflect() protoreflect.Message {
	mi := &file_app_node_conf_conf_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Bootstrap.ProtoReflect.Descriptor instead.
func (*Bootstrap) Descriptor() ([]byte, []int) {
	return file_app_node_conf_conf_proto_rawDescGZIP(), []int{0}
}

func (x *Bootstrap) GetServer() *Server {
	if x != nil {
		return x.Server
	}
	return nil
}

func (x *Bootstrap) GetData() *Data {
	if x != nil {
		return x.Data
	}
	return nil
}

type Server struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Grpc *Server_GRPC `protobuf:"bytes,1,opt,name=grpc,proto3" json:"grpc,omitempty"`
	Node *Server_Node `protobuf:"bytes,2,opt,name=node,proto3" json:"node,omitempty"`
}

func (x *Server) Reset() {
	*x = Server{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_node_conf_conf_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Server) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Server) ProtoMessage() {}

func (x *Server) ProtoReflect() protoreflect.Message {
	mi := &file_app_node_conf_conf_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Server.ProtoReflect.Descriptor instead.
func (*Server) Descriptor() ([]byte, []int) {
	return file_app_node_conf_conf_proto_rawDescGZIP(), []int{1}
}

func (x *Server) GetGrpc() *Server_GRPC {
	if x != nil {
		return x.Grpc
	}
	return nil
}

func (x *Server) GetNode() *Server_Node {
	if x != nil {
		return x.Node
	}
	return nil
}

type Data struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BlockstoreUri      string `protobuf:"bytes,1,opt,name=blockstore_uri,json=blockstoreUri,proto3" json:"blockstore_uri,omitempty"`
	ReproviderInterval int64  `protobuf:"varint,2,opt,name=reprovider_interval,json=reproviderInterval,proto3" json:"reprovider_interval,omitempty"`
}

func (x *Data) Reset() {
	*x = Data{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_node_conf_conf_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Data) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Data) ProtoMessage() {}

func (x *Data) ProtoReflect() protoreflect.Message {
	mi := &file_app_node_conf_conf_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Data.ProtoReflect.Descriptor instead.
func (*Data) Descriptor() ([]byte, []int) {
	return file_app_node_conf_conf_proto_rawDescGZIP(), []int{2}
}

func (x *Data) GetBlockstoreUri() string {
	if x != nil {
		return x.BlockstoreUri
	}
	return ""
}

func (x *Data) GetReproviderInterval() int64 {
	if x != nil {
		return x.ReproviderInterval
	}
	return 0
}

type Server_GRPC struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Addr    string               `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	Timeout *durationpb.Duration `protobuf:"bytes,2,opt,name=timeout,proto3" json:"timeout,omitempty"`
}

func (x *Server_GRPC) Reset() {
	*x = Server_GRPC{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_node_conf_conf_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Server_GRPC) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Server_GRPC) ProtoMessage() {}

func (x *Server_GRPC) ProtoReflect() protoreflect.Message {
	mi := &file_app_node_conf_conf_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Server_GRPC.ProtoReflect.Descriptor instead.
func (*Server_GRPC) Descriptor() ([]byte, []int) {
	return file_app_node_conf_conf_proto_rawDescGZIP(), []int{1, 0}
}

func (x *Server_GRPC) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *Server_GRPC) GetTimeout() *durationpb.Duration {
	if x != nil {
		return x.Timeout
	}
	return nil
}

type Server_Node struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MultiAddr   []string `protobuf:"bytes,1,rep,name=multi_addr,json=multiAddr,proto3" json:"multi_addr,omitempty"`
	Peers       []string `protobuf:"bytes,2,rep,name=peers,proto3" json:"peers,omitempty"`
	PrivateKey  string   `protobuf:"bytes,3,opt,name=private_key,json=privateKey,proto3" json:"private_key,omitempty"`
	GracePeriod int64    `protobuf:"varint,4,opt,name=grace_period,json=gracePeriod,proto3" json:"grace_period,omitempty"`
	LowWater    int64    `protobuf:"varint,5,opt,name=low_water,json=lowWater,proto3" json:"low_water,omitempty"`
	HighWater   int64    `protobuf:"varint,6,opt,name=high_water,json=highWater,proto3" json:"high_water,omitempty"`
	LeveldbPath string   `protobuf:"bytes,7,opt,name=leveldb_path,json=leveldbPath,proto3" json:"leveldb_path,omitempty"`
}

func (x *Server_Node) Reset() {
	*x = Server_Node{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_node_conf_conf_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Server_Node) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Server_Node) ProtoMessage() {}

func (x *Server_Node) ProtoReflect() protoreflect.Message {
	mi := &file_app_node_conf_conf_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Server_Node.ProtoReflect.Descriptor instead.
func (*Server_Node) Descriptor() ([]byte, []int) {
	return file_app_node_conf_conf_proto_rawDescGZIP(), []int{1, 1}
}

func (x *Server_Node) GetMultiAddr() []string {
	if x != nil {
		return x.MultiAddr
	}
	return nil
}

func (x *Server_Node) GetPeers() []string {
	if x != nil {
		return x.Peers
	}
	return nil
}

func (x *Server_Node) GetPrivateKey() string {
	if x != nil {
		return x.PrivateKey
	}
	return ""
}

func (x *Server_Node) GetGracePeriod() int64 {
	if x != nil {
		return x.GracePeriod
	}
	return 0
}

func (x *Server_Node) GetLowWater() int64 {
	if x != nil {
		return x.LowWater
	}
	return 0
}

func (x *Server_Node) GetHighWater() int64 {
	if x != nil {
		return x.HighWater
	}
	return 0
}

func (x *Server_Node) GetLeveldbPath() string {
	if x != nil {
		return x.LeveldbPath
	}
	return ""
}

var File_app_node_conf_conf_proto protoreflect.FileDescriptor

var file_app_node_conf_conf_proto_rawDesc = []byte{
	0x0a, 0x18, 0x61, 0x70, 0x70, 0x2f, 0x6e, 0x6f, 0x64, 0x65, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x2f,
	0x63, 0x6f, 0x6e, 0x66, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x69, 0x70, 0x66, 0x73,
	0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x65, 0x0a, 0x09, 0x42, 0x6f,
	0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x12, 0x2e, 0x0a, 0x06, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x69, 0x70, 0x66, 0x73, 0x2e, 0x6e,
	0x6f, 0x64, 0x65, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52,
	0x06, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x28, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x69, 0x70, 0x66, 0x73, 0x2e, 0x6e, 0x6f, 0x64,
	0x65, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x22, 0x9c, 0x03, 0x0a, 0x06, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x2f, 0x0a, 0x04,
	0x67, 0x72, 0x70, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x69, 0x70, 0x66,
	0x73, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x2e, 0x53, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x2e, 0x47, 0x52, 0x50, 0x43, 0x52, 0x04, 0x67, 0x72, 0x70, 0x63, 0x12, 0x2f, 0x0a,
	0x04, 0x6e, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x69, 0x70,
	0x66, 0x73, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x2e, 0x53, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x1a, 0x4f,
	0x0a, 0x04, 0x47, 0x52, 0x50, 0x43, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x64, 0x64, 0x72, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72, 0x12, 0x33, 0x0a, 0x07, 0x74, 0x69,
	0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x1a,
	0xde, 0x01, 0x0a, 0x04, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x6d, 0x75, 0x6c, 0x74,
	0x69, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x6d, 0x75,
	0x6c, 0x74, 0x69, 0x41, 0x64, 0x64, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x65, 0x65, 0x72, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x70, 0x65, 0x65, 0x72, 0x73, 0x12, 0x1f, 0x0a,
	0x0b, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x12, 0x21,
	0x0a, 0x0c, 0x67, 0x72, 0x61, 0x63, 0x65, 0x5f, 0x70, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x67, 0x72, 0x61, 0x63, 0x65, 0x50, 0x65, 0x72, 0x69, 0x6f,
	0x64, 0x12, 0x1b, 0x0a, 0x09, 0x6c, 0x6f, 0x77, 0x5f, 0x77, 0x61, 0x74, 0x65, 0x72, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x6c, 0x6f, 0x77, 0x57, 0x61, 0x74, 0x65, 0x72, 0x12, 0x1d,
	0x0a, 0x0a, 0x68, 0x69, 0x67, 0x68, 0x5f, 0x77, 0x61, 0x74, 0x65, 0x72, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x09, 0x68, 0x69, 0x67, 0x68, 0x57, 0x61, 0x74, 0x65, 0x72, 0x12, 0x21, 0x0a,
	0x0c, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x64, 0x62, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x64, 0x62, 0x50, 0x61, 0x74, 0x68,
	0x22, 0x5e, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x25, 0x0a, 0x0e, 0x62, 0x6c, 0x6f, 0x63,
	0x6b, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x5f, 0x75, 0x72, 0x69, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0d, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x55, 0x72, 0x69, 0x12,
	0x2f, 0x0a, 0x13, 0x72, 0x65, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x12, 0x72, 0x65,
	0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c,
	0x42, 0x21, 0x5a, 0x1f, 0x69, 0x70, 0x66, 0x73, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x73,
	0x2f, 0x61, 0x70, 0x70, 0x2f, 0x6e, 0x6f, 0x64, 0x65, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x3b, 0x63,
	0x6f, 0x6e, 0x66, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_app_node_conf_conf_proto_rawDescOnce sync.Once
	file_app_node_conf_conf_proto_rawDescData = file_app_node_conf_conf_proto_rawDesc
)

func file_app_node_conf_conf_proto_rawDescGZIP() []byte {
	file_app_node_conf_conf_proto_rawDescOnce.Do(func() {
		file_app_node_conf_conf_proto_rawDescData = protoimpl.X.CompressGZIP(file_app_node_conf_conf_proto_rawDescData)
	})
	return file_app_node_conf_conf_proto_rawDescData
}

var file_app_node_conf_conf_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_app_node_conf_conf_proto_goTypes = []interface{}{
	(*Bootstrap)(nil),           // 0: ipfs.node.conf.Bootstrap
	(*Server)(nil),              // 1: ipfs.node.conf.Server
	(*Data)(nil),                // 2: ipfs.node.conf.Data
	(*Server_GRPC)(nil),         // 3: ipfs.node.conf.Server.GRPC
	(*Server_Node)(nil),         // 4: ipfs.node.conf.Server.Node
	(*durationpb.Duration)(nil), // 5: google.protobuf.Duration
}
var file_app_node_conf_conf_proto_depIdxs = []int32{
	1, // 0: ipfs.node.conf.Bootstrap.server:type_name -> ipfs.node.conf.Server
	2, // 1: ipfs.node.conf.Bootstrap.data:type_name -> ipfs.node.conf.Data
	3, // 2: ipfs.node.conf.Server.grpc:type_name -> ipfs.node.conf.Server.GRPC
	4, // 3: ipfs.node.conf.Server.node:type_name -> ipfs.node.conf.Server.Node
	5, // 4: ipfs.node.conf.Server.GRPC.timeout:type_name -> google.protobuf.Duration
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_app_node_conf_conf_proto_init() }
func file_app_node_conf_conf_proto_init() {
	if File_app_node_conf_conf_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_app_node_conf_conf_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Bootstrap); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_app_node_conf_conf_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Server); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_app_node_conf_conf_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Data); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_app_node_conf_conf_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Server_GRPC); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_app_node_conf_conf_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Server_Node); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_app_node_conf_conf_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_app_node_conf_conf_proto_goTypes,
		DependencyIndexes: file_app_node_conf_conf_proto_depIdxs,
		MessageInfos:      file_app_node_conf_conf_proto_msgTypes,
	}.Build()
	File_app_node_conf_conf_proto = out.File
	file_app_node_conf_conf_proto_rawDesc = nil
	file_app_node_conf_conf_proto_goTypes = nil
	file_app_node_conf_conf_proto_depIdxs = nil
}
