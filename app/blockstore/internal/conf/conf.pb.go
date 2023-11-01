// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: app/blockstore/internal/conf/conf.proto

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

type Data_DBType int32

const (
	Data_LevelDB Data_DBType = 0
	Data_TiKV    Data_DBType = 1
	Data_PG      Data_DBType = 2
)

// Enum value maps for Data_DBType.
var (
	Data_DBType_name = map[int32]string{
		0: "LevelDB",
		1: "TiKV",
		2: "PG",
	}
	Data_DBType_value = map[string]int32{
		"LevelDB": 0,
		"TiKV":    1,
		"PG":      2,
	}
)

func (x Data_DBType) Enum() *Data_DBType {
	p := new(Data_DBType)
	*p = x
	return p
}

func (x Data_DBType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Data_DBType) Descriptor() protoreflect.EnumDescriptor {
	return file_app_blockstore_internal_conf_conf_proto_enumTypes[0].Descriptor()
}

func (Data_DBType) Type() protoreflect.EnumType {
	return &file_app_blockstore_internal_conf_conf_proto_enumTypes[0]
}

func (x Data_DBType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Data_DBType.Descriptor instead.
func (Data_DBType) EnumDescriptor() ([]byte, []int) {
	return file_app_blockstore_internal_conf_conf_proto_rawDescGZIP(), []int{2, 0}
}

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
		mi := &file_app_blockstore_internal_conf_conf_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Bootstrap) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Bootstrap) ProtoMessage() {}

func (x *Bootstrap) ProtoReflect() protoreflect.Message {
	mi := &file_app_blockstore_internal_conf_conf_proto_msgTypes[0]
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
	return file_app_blockstore_internal_conf_conf_proto_rawDescGZIP(), []int{0}
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
	Http *Server_HTTP `protobuf:"bytes,2,opt,name=http,proto3" json:"http,omitempty"`
}

func (x *Server) Reset() {
	*x = Server{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_blockstore_internal_conf_conf_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Server) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Server) ProtoMessage() {}

func (x *Server) ProtoReflect() protoreflect.Message {
	mi := &file_app_blockstore_internal_conf_conf_proto_msgTypes[1]
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
	return file_app_blockstore_internal_conf_conf_proto_rawDescGZIP(), []int{1}
}

func (x *Server) GetGrpc() *Server_GRPC {
	if x != nil {
		return x.Grpc
	}
	return nil
}

func (x *Server) GetHttp() *Server_HTTP {
	if x != nil {
		return x.Http
	}
	return nil
}

type Data struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Storage *Data_S3    `protobuf:"bytes,1,opt,name=storage,proto3" json:"storage,omitempty"`
	Db      *Data_DB    `protobuf:"bytes,2,opt,name=db,proto3" json:"db,omitempty"`
	Cache   *Data_Cache `protobuf:"bytes,3,opt,name=cache,proto3" json:"cache,omitempty"`
	Redis   *Data_Redis `protobuf:"bytes,4,opt,name=redis,proto3" json:"redis,omitempty"`
}

func (x *Data) Reset() {
	*x = Data{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_blockstore_internal_conf_conf_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Data) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Data) ProtoMessage() {}

func (x *Data) ProtoReflect() protoreflect.Message {
	mi := &file_app_blockstore_internal_conf_conf_proto_msgTypes[2]
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
	return file_app_blockstore_internal_conf_conf_proto_rawDescGZIP(), []int{2}
}

func (x *Data) GetStorage() *Data_S3 {
	if x != nil {
		return x.Storage
	}
	return nil
}

func (x *Data) GetDb() *Data_DB {
	if x != nil {
		return x.Db
	}
	return nil
}

func (x *Data) GetCache() *Data_Cache {
	if x != nil {
		return x.Cache
	}
	return nil
}

func (x *Data) GetRedis() *Data_Redis {
	if x != nil {
		return x.Redis
	}
	return nil
}

type Server_GRPC struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Network string               `protobuf:"bytes,1,opt,name=network,proto3" json:"network,omitempty"`
	Addr    string               `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
	Timeout *durationpb.Duration `protobuf:"bytes,3,opt,name=timeout,proto3" json:"timeout,omitempty"`
}

func (x *Server_GRPC) Reset() {
	*x = Server_GRPC{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_blockstore_internal_conf_conf_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Server_GRPC) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Server_GRPC) ProtoMessage() {}

func (x *Server_GRPC) ProtoReflect() protoreflect.Message {
	mi := &file_app_blockstore_internal_conf_conf_proto_msgTypes[3]
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
	return file_app_blockstore_internal_conf_conf_proto_rawDescGZIP(), []int{1, 0}
}

func (x *Server_GRPC) GetNetwork() string {
	if x != nil {
		return x.Network
	}
	return ""
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

type Server_HTTP struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Addr    string               `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	Timeout *durationpb.Duration `protobuf:"bytes,2,opt,name=timeout,proto3" json:"timeout,omitempty"`
}

func (x *Server_HTTP) Reset() {
	*x = Server_HTTP{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_blockstore_internal_conf_conf_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Server_HTTP) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Server_HTTP) ProtoMessage() {}

func (x *Server_HTTP) ProtoReflect() protoreflect.Message {
	mi := &file_app_blockstore_internal_conf_conf_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Server_HTTP.ProtoReflect.Descriptor instead.
func (*Server_HTTP) Descriptor() ([]byte, []int) {
	return file_app_blockstore_internal_conf_conf_proto_rawDescGZIP(), []int{1, 1}
}

func (x *Server_HTTP) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *Server_HTTP) GetTimeout() *durationpb.Duration {
	if x != nil {
		return x.Timeout
	}
	return nil
}

type Data_Redis struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Network      string               `protobuf:"bytes,1,opt,name=network,proto3" json:"network,omitempty"`
	Addr         string               `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
	Password     string               `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
	ReadTimeout  *durationpb.Duration `protobuf:"bytes,4,opt,name=read_timeout,json=readTimeout,proto3" json:"read_timeout,omitempty"`
	WriteTimeout *durationpb.Duration `protobuf:"bytes,5,opt,name=write_timeout,json=writeTimeout,proto3" json:"write_timeout,omitempty"`
}

func (x *Data_Redis) Reset() {
	*x = Data_Redis{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_blockstore_internal_conf_conf_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Data_Redis) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Data_Redis) ProtoMessage() {}

func (x *Data_Redis) ProtoReflect() protoreflect.Message {
	mi := &file_app_blockstore_internal_conf_conf_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Data_Redis.ProtoReflect.Descriptor instead.
func (*Data_Redis) Descriptor() ([]byte, []int) {
	return file_app_blockstore_internal_conf_conf_proto_rawDescGZIP(), []int{2, 0}
}

func (x *Data_Redis) GetNetwork() string {
	if x != nil {
		return x.Network
	}
	return ""
}

func (x *Data_Redis) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *Data_Redis) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *Data_Redis) GetReadTimeout() *durationpb.Duration {
	if x != nil {
		return x.ReadTimeout
	}
	return nil
}

func (x *Data_Redis) GetWriteTimeout() *durationpb.Duration {
	if x != nil {
		return x.WriteTimeout
	}
	return nil
}

type Data_S3 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Endpoint  string `protobuf:"bytes,1,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	Region    string `protobuf:"bytes,2,opt,name=region,proto3" json:"region,omitempty"`
	Bucket    string `protobuf:"bytes,3,opt,name=bucket,proto3" json:"bucket,omitempty"`
	AccessKey string `protobuf:"bytes,4,opt,name=accessKey,proto3" json:"accessKey,omitempty"`
	SecretKey string `protobuf:"bytes,5,opt,name=secretKey,proto3" json:"secretKey,omitempty"`
	Schemes   string `protobuf:"bytes,6,opt,name=schemes,proto3" json:"schemes,omitempty"`
}

func (x *Data_S3) Reset() {
	*x = Data_S3{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_blockstore_internal_conf_conf_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Data_S3) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Data_S3) ProtoMessage() {}

func (x *Data_S3) ProtoReflect() protoreflect.Message {
	mi := &file_app_blockstore_internal_conf_conf_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Data_S3.ProtoReflect.Descriptor instead.
func (*Data_S3) Descriptor() ([]byte, []int) {
	return file_app_blockstore_internal_conf_conf_proto_rawDescGZIP(), []int{2, 1}
}

func (x *Data_S3) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

func (x *Data_S3) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

func (x *Data_S3) GetBucket() string {
	if x != nil {
		return x.Bucket
	}
	return ""
}

func (x *Data_S3) GetAccessKey() string {
	if x != nil {
		return x.AccessKey
	}
	return ""
}

func (x *Data_S3) GetSecretKey() string {
	if x != nil {
		return x.SecretKey
	}
	return ""
}

func (x *Data_S3) GetSchemes() string {
	if x != nil {
		return x.Schemes
	}
	return ""
}

type Data_DB struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type    Data_DBType      `protobuf:"varint,1,opt,name=type,proto3,enum=ipfs.blockstore.api.Data_DBType" json:"type,omitempty"`
	Leveldb *Data_DB_LevelDB `protobuf:"bytes,2,opt,name=leveldb,proto3" json:"leveldb,omitempty"`
	Tikv    *Data_DB_TiKV    `protobuf:"bytes,3,opt,name=tikv,proto3" json:"tikv,omitempty"`
	Pg      *Data_DB_Pg      `protobuf:"bytes,4,opt,name=pg,proto3" json:"pg,omitempty"`
}

func (x *Data_DB) Reset() {
	*x = Data_DB{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_blockstore_internal_conf_conf_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Data_DB) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Data_DB) ProtoMessage() {}

func (x *Data_DB) ProtoReflect() protoreflect.Message {
	mi := &file_app_blockstore_internal_conf_conf_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Data_DB.ProtoReflect.Descriptor instead.
func (*Data_DB) Descriptor() ([]byte, []int) {
	return file_app_blockstore_internal_conf_conf_proto_rawDescGZIP(), []int{2, 2}
}

func (x *Data_DB) GetType() Data_DBType {
	if x != nil {
		return x.Type
	}
	return Data_LevelDB
}

func (x *Data_DB) GetLeveldb() *Data_DB_LevelDB {
	if x != nil {
		return x.Leveldb
	}
	return nil
}

func (x *Data_DB) GetTikv() *Data_DB_TiKV {
	if x != nil {
		return x.Tikv
	}
	return nil
}

func (x *Data_DB) GetPg() *Data_DB_Pg {
	if x != nil {
		return x.Pg
	}
	return nil
}

type Data_Cache struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LruSize   int64  `protobuf:"varint,1,opt,name=lruSize,proto3" json:"lruSize,omitempty"`
	BasePath  string `protobuf:"bytes,2,opt,name=basePath,proto3" json:"basePath,omitempty"`
	IndexPath string `protobuf:"bytes,3,opt,name=indexPath,proto3" json:"indexPath,omitempty"`
}

func (x *Data_Cache) Reset() {
	*x = Data_Cache{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_blockstore_internal_conf_conf_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Data_Cache) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Data_Cache) ProtoMessage() {}

func (x *Data_Cache) ProtoReflect() protoreflect.Message {
	mi := &file_app_blockstore_internal_conf_conf_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Data_Cache.ProtoReflect.Descriptor instead.
func (*Data_Cache) Descriptor() ([]byte, []int) {
	return file_app_blockstore_internal_conf_conf_proto_rawDescGZIP(), []int{2, 3}
}

func (x *Data_Cache) GetLruSize() int64 {
	if x != nil {
		return x.LruSize
	}
	return 0
}

func (x *Data_Cache) GetBasePath() string {
	if x != nil {
		return x.BasePath
	}
	return ""
}

func (x *Data_Cache) GetIndexPath() string {
	if x != nil {
		return x.IndexPath
	}
	return ""
}

type Data_DB_LevelDB struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Path string `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
}

func (x *Data_DB_LevelDB) Reset() {
	*x = Data_DB_LevelDB{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_blockstore_internal_conf_conf_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Data_DB_LevelDB) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Data_DB_LevelDB) ProtoMessage() {}

func (x *Data_DB_LevelDB) ProtoReflect() protoreflect.Message {
	mi := &file_app_blockstore_internal_conf_conf_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Data_DB_LevelDB.ProtoReflect.Descriptor instead.
func (*Data_DB_LevelDB) Descriptor() ([]byte, []int) {
	return file_app_blockstore_internal_conf_conf_proto_rawDescGZIP(), []int{2, 2, 0}
}

func (x *Data_DB_LevelDB) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

type Data_DB_TiKV struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Addrs []string `protobuf:"bytes,1,rep,name=addrs,proto3" json:"addrs,omitempty"`
}

func (x *Data_DB_TiKV) Reset() {
	*x = Data_DB_TiKV{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_blockstore_internal_conf_conf_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Data_DB_TiKV) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Data_DB_TiKV) ProtoMessage() {}

func (x *Data_DB_TiKV) ProtoReflect() protoreflect.Message {
	mi := &file_app_blockstore_internal_conf_conf_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Data_DB_TiKV.ProtoReflect.Descriptor instead.
func (*Data_DB_TiKV) Descriptor() ([]byte, []int) {
	return file_app_blockstore_internal_conf_conf_proto_rawDescGZIP(), []int{2, 2, 1}
}

func (x *Data_DB_TiKV) GetAddrs() []string {
	if x != nil {
		return x.Addrs
	}
	return nil
}

type Data_DB_Pg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SourcesDsn  []string `protobuf:"bytes,1,rep,name=sourcesDsn,proto3" json:"sourcesDsn,omitempty"`
	ReplicasDsn []string `protobuf:"bytes,2,rep,name=replicasDsn,proto3" json:"replicasDsn,omitempty"`
	EnableBloom bool     `protobuf:"varint,3,opt,name=enableBloom,proto3" json:"enableBloom,omitempty"`
}

func (x *Data_DB_Pg) Reset() {
	*x = Data_DB_Pg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_blockstore_internal_conf_conf_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Data_DB_Pg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Data_DB_Pg) ProtoMessage() {}

func (x *Data_DB_Pg) ProtoReflect() protoreflect.Message {
	mi := &file_app_blockstore_internal_conf_conf_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Data_DB_Pg.ProtoReflect.Descriptor instead.
func (*Data_DB_Pg) Descriptor() ([]byte, []int) {
	return file_app_blockstore_internal_conf_conf_proto_rawDescGZIP(), []int{2, 2, 2}
}

func (x *Data_DB_Pg) GetSourcesDsn() []string {
	if x != nil {
		return x.SourcesDsn
	}
	return nil
}

func (x *Data_DB_Pg) GetReplicasDsn() []string {
	if x != nil {
		return x.ReplicasDsn
	}
	return nil
}

func (x *Data_DB_Pg) GetEnableBloom() bool {
	if x != nil {
		return x.EnableBloom
	}
	return false
}

var File_app_blockstore_internal_conf_conf_proto protoreflect.FileDescriptor

var file_app_blockstore_internal_conf_conf_proto_rawDesc = []byte{
	0x0a, 0x27, 0x61, 0x70, 0x70, 0x2f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x2f, 0x63,
	0x6f, 0x6e, 0x66, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13, 0x69, 0x70, 0x66, 0x73, 0x2e,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x1a, 0x1e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x6f,
	0x0a, 0x09, 0x42, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x12, 0x33, 0x0a, 0x06, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x69, 0x70,
	0x66, 0x73, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x06, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x12, 0x2d, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19,
	0x2e, 0x69, 0x70, 0x66, 0x73, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22,
	0xb0, 0x02, 0x0a, 0x06, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x34, 0x0a, 0x04, 0x67, 0x72,
	0x70, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x69, 0x70, 0x66, 0x73, 0x2e,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x53,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x47, 0x52, 0x50, 0x43, 0x52, 0x04, 0x67, 0x72, 0x70, 0x63,
	0x12, 0x34, 0x0a, 0x04, 0x68, 0x74, 0x74, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20,
	0x2e, 0x69, 0x70, 0x66, 0x73, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x48, 0x54, 0x54, 0x50,
	0x52, 0x04, 0x68, 0x74, 0x74, 0x70, 0x1a, 0x69, 0x0a, 0x04, 0x47, 0x52, 0x50, 0x43, 0x12, 0x18,
	0x0a, 0x07, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x64, 0x64, 0x72,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72, 0x12, 0x33, 0x0a, 0x07,
	0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75,
	0x74, 0x1a, 0x4f, 0x0a, 0x04, 0x48, 0x54, 0x54, 0x50, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x64, 0x64,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72, 0x12, 0x33, 0x0a,
	0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f,
	0x75, 0x74, 0x22, 0xe7, 0x08, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x36, 0x0a, 0x07, 0x73,
	0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x69,
	0x70, 0x66, 0x73, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x53, 0x33, 0x52, 0x07, 0x73, 0x74, 0x6f, 0x72,
	0x61, 0x67, 0x65, 0x12, 0x2c, 0x0a, 0x02, 0x64, 0x62, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1c, 0x2e, 0x69, 0x70, 0x66, 0x73, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x74, 0x6f, 0x72,
	0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x44, 0x42, 0x52, 0x02, 0x64,
	0x62, 0x12, 0x35, 0x0a, 0x05, 0x63, 0x61, 0x63, 0x68, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1f, 0x2e, 0x69, 0x70, 0x66, 0x73, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x74, 0x6f,
	0x72, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x43, 0x61, 0x63, 0x68,
	0x65, 0x52, 0x05, 0x63, 0x61, 0x63, 0x68, 0x65, 0x12, 0x35, 0x0a, 0x05, 0x72, 0x65, 0x64, 0x69,
	0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x69, 0x70, 0x66, 0x73, 0x2e, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x44, 0x61,
	0x74, 0x61, 0x2e, 0x52, 0x65, 0x64, 0x69, 0x73, 0x52, 0x05, 0x72, 0x65, 0x64, 0x69, 0x73, 0x1a,
	0xcf, 0x01, 0x0a, 0x05, 0x52, 0x65, 0x64, 0x69, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6e, 0x65, 0x74,
	0x77, 0x6f, 0x72, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6e, 0x65, 0x74, 0x77,
	0x6f, 0x72, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x64, 0x64, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x12, 0x3c, 0x0a, 0x0c, 0x72, 0x65, 0x61, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65,
	0x6f, 0x75, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0b, 0x72, 0x65, 0x61, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75,
	0x74, 0x12, 0x3e, 0x0a, 0x0d, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x6f,
	0x75, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x0c, 0x77, 0x72, 0x69, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75,
	0x74, 0x1a, 0xa6, 0x01, 0x0a, 0x02, 0x53, 0x33, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x70,
	0x6f, 0x69, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x64, 0x70,
	0x6f, 0x69, 0x6e, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06,
	0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x62, 0x75,
	0x63, 0x6b, 0x65, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x4b, 0x65,
	0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x4b,
	0x65, 0x79, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x4b, 0x65, 0x79, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x4b, 0x65, 0x79,
	0x12, 0x18, 0x0a, 0x07, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x65, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x65, 0x73, 0x1a, 0x89, 0x03, 0x0a, 0x02, 0x44,
	0x42, 0x12, 0x34, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x20, 0x2e, 0x69, 0x70, 0x66, 0x73, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x74, 0x6f, 0x72,
	0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x44, 0x42, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x3e, 0x0a, 0x07, 0x6c, 0x65, 0x76, 0x65, 0x6c,
	0x64, 0x62, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x69, 0x70, 0x66, 0x73, 0x2e,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x44,
	0x61, 0x74, 0x61, 0x2e, 0x44, 0x42, 0x2e, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x44, 0x42, 0x52, 0x07,
	0x6c, 0x65, 0x76, 0x65, 0x6c, 0x64, 0x62, 0x12, 0x35, 0x0a, 0x04, 0x74, 0x69, 0x6b, 0x76, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x69, 0x70, 0x66, 0x73, 0x2e, 0x62, 0x6c, 0x6f,
	0x63, 0x6b, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x44, 0x61, 0x74, 0x61,
	0x2e, 0x44, 0x42, 0x2e, 0x54, 0x69, 0x4b, 0x56, 0x52, 0x04, 0x74, 0x69, 0x6b, 0x76, 0x12, 0x2f,
	0x0a, 0x02, 0x70, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x69, 0x70, 0x66,
	0x73, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x44, 0x42, 0x2e, 0x50, 0x67, 0x52, 0x02, 0x70, 0x67, 0x1a,
	0x1d, 0x0a, 0x07, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x44, 0x42, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61,
	0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x1a, 0x1c,
	0x0a, 0x04, 0x54, 0x69, 0x4b, 0x56, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x64, 0x64, 0x72, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x61, 0x64, 0x64, 0x72, 0x73, 0x1a, 0x68, 0x0a, 0x02,
	0x50, 0x67, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x44, 0x73, 0x6e,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x44,
	0x73, 0x6e, 0x12, 0x20, 0x0a, 0x0b, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x73, 0x44, 0x73,
	0x6e, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0b, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61,
	0x73, 0x44, 0x73, 0x6e, 0x12, 0x20, 0x0a, 0x0b, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x42, 0x6c,
	0x6f, 0x6f, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x65, 0x6e, 0x61, 0x62, 0x6c,
	0x65, 0x42, 0x6c, 0x6f, 0x6f, 0x6d, 0x1a, 0x5b, 0x0a, 0x05, 0x43, 0x61, 0x63, 0x68, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x6c, 0x72, 0x75, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x07, 0x6c, 0x72, 0x75, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x62, 0x61, 0x73,
	0x65, 0x50, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x62, 0x61, 0x73,
	0x65, 0x50, 0x61, 0x74, 0x68, 0x12, 0x1c, 0x0a, 0x09, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x50, 0x61,
	0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x50,
	0x61, 0x74, 0x68, 0x22, 0x27, 0x0a, 0x06, 0x44, 0x42, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a,
	0x07, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x44, 0x42, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x54, 0x69,
	0x4b, 0x56, 0x10, 0x01, 0x12, 0x06, 0x0a, 0x02, 0x50, 0x47, 0x10, 0x02, 0x42, 0x2c, 0x5a, 0x2a,
	0x69, 0x70, 0x66, 0x73, 0x2d, 0x74, 0x6f, 0x70, 0x2f, 0x61, 0x70, 0x70, 0x2f, 0x62, 0x6c, 0x6f,
	0x63, 0x6b, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x3b, 0x63, 0x6f, 0x6e, 0x66, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_app_blockstore_internal_conf_conf_proto_rawDescOnce sync.Once
	file_app_blockstore_internal_conf_conf_proto_rawDescData = file_app_blockstore_internal_conf_conf_proto_rawDesc
)

func file_app_blockstore_internal_conf_conf_proto_rawDescGZIP() []byte {
	file_app_blockstore_internal_conf_conf_proto_rawDescOnce.Do(func() {
		file_app_blockstore_internal_conf_conf_proto_rawDescData = protoimpl.X.CompressGZIP(file_app_blockstore_internal_conf_conf_proto_rawDescData)
	})
	return file_app_blockstore_internal_conf_conf_proto_rawDescData
}

var file_app_blockstore_internal_conf_conf_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_app_blockstore_internal_conf_conf_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_app_blockstore_internal_conf_conf_proto_goTypes = []interface{}{
	(Data_DBType)(0),            // 0: ipfs.blockstore.api.Data.DBType
	(*Bootstrap)(nil),           // 1: ipfs.blockstore.api.Bootstrap
	(*Server)(nil),              // 2: ipfs.blockstore.api.Server
	(*Data)(nil),                // 3: ipfs.blockstore.api.Data
	(*Server_GRPC)(nil),         // 4: ipfs.blockstore.api.Server.GRPC
	(*Server_HTTP)(nil),         // 5: ipfs.blockstore.api.Server.HTTP
	(*Data_Redis)(nil),          // 6: ipfs.blockstore.api.Data.Redis
	(*Data_S3)(nil),             // 7: ipfs.blockstore.api.Data.S3
	(*Data_DB)(nil),             // 8: ipfs.blockstore.api.Data.DB
	(*Data_Cache)(nil),          // 9: ipfs.blockstore.api.Data.Cache
	(*Data_DB_LevelDB)(nil),     // 10: ipfs.blockstore.api.Data.DB.LevelDB
	(*Data_DB_TiKV)(nil),        // 11: ipfs.blockstore.api.Data.DB.TiKV
	(*Data_DB_Pg)(nil),          // 12: ipfs.blockstore.api.Data.DB.Pg
	(*durationpb.Duration)(nil), // 13: google.protobuf.Duration
}
var file_app_blockstore_internal_conf_conf_proto_depIdxs = []int32{
	2,  // 0: ipfs.blockstore.api.Bootstrap.server:type_name -> ipfs.blockstore.api.Server
	3,  // 1: ipfs.blockstore.api.Bootstrap.data:type_name -> ipfs.blockstore.api.Data
	4,  // 2: ipfs.blockstore.api.Server.grpc:type_name -> ipfs.blockstore.api.Server.GRPC
	5,  // 3: ipfs.blockstore.api.Server.http:type_name -> ipfs.blockstore.api.Server.HTTP
	7,  // 4: ipfs.blockstore.api.Data.storage:type_name -> ipfs.blockstore.api.Data.S3
	8,  // 5: ipfs.blockstore.api.Data.db:type_name -> ipfs.blockstore.api.Data.DB
	9,  // 6: ipfs.blockstore.api.Data.cache:type_name -> ipfs.blockstore.api.Data.Cache
	6,  // 7: ipfs.blockstore.api.Data.redis:type_name -> ipfs.blockstore.api.Data.Redis
	13, // 8: ipfs.blockstore.api.Server.GRPC.timeout:type_name -> google.protobuf.Duration
	13, // 9: ipfs.blockstore.api.Server.HTTP.timeout:type_name -> google.protobuf.Duration
	13, // 10: ipfs.blockstore.api.Data.Redis.read_timeout:type_name -> google.protobuf.Duration
	13, // 11: ipfs.blockstore.api.Data.Redis.write_timeout:type_name -> google.protobuf.Duration
	0,  // 12: ipfs.blockstore.api.Data.DB.type:type_name -> ipfs.blockstore.api.Data.DBType
	10, // 13: ipfs.blockstore.api.Data.DB.leveldb:type_name -> ipfs.blockstore.api.Data.DB.LevelDB
	11, // 14: ipfs.blockstore.api.Data.DB.tikv:type_name -> ipfs.blockstore.api.Data.DB.TiKV
	12, // 15: ipfs.blockstore.api.Data.DB.pg:type_name -> ipfs.blockstore.api.Data.DB.Pg
	16, // [16:16] is the sub-list for method output_type
	16, // [16:16] is the sub-list for method input_type
	16, // [16:16] is the sub-list for extension type_name
	16, // [16:16] is the sub-list for extension extendee
	0,  // [0:16] is the sub-list for field type_name
}

func init() { file_app_blockstore_internal_conf_conf_proto_init() }
func file_app_blockstore_internal_conf_conf_proto_init() {
	if File_app_blockstore_internal_conf_conf_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_app_blockstore_internal_conf_conf_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_app_blockstore_internal_conf_conf_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_app_blockstore_internal_conf_conf_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
		file_app_blockstore_internal_conf_conf_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
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
		file_app_blockstore_internal_conf_conf_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Server_HTTP); i {
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
		file_app_blockstore_internal_conf_conf_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Data_Redis); i {
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
		file_app_blockstore_internal_conf_conf_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Data_S3); i {
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
		file_app_blockstore_internal_conf_conf_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Data_DB); i {
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
		file_app_blockstore_internal_conf_conf_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Data_Cache); i {
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
		file_app_blockstore_internal_conf_conf_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Data_DB_LevelDB); i {
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
		file_app_blockstore_internal_conf_conf_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Data_DB_TiKV); i {
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
		file_app_blockstore_internal_conf_conf_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Data_DB_Pg); i {
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
			RawDescriptor: file_app_blockstore_internal_conf_conf_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_app_blockstore_internal_conf_conf_proto_goTypes,
		DependencyIndexes: file_app_blockstore_internal_conf_conf_proto_depIdxs,
		EnumInfos:         file_app_blockstore_internal_conf_conf_proto_enumTypes,
		MessageInfos:      file_app_blockstore_internal_conf_conf_proto_msgTypes,
	}.Build()
	File_app_blockstore_internal_conf_conf_proto = out.File
	file_app_blockstore_internal_conf_conf_proto_rawDesc = nil
	file_app_blockstore_internal_conf_conf_proto_goTypes = nil
	file_app_blockstore_internal_conf_conf_proto_depIdxs = nil
}
