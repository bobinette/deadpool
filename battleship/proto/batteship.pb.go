// Code generated by protoc-gen-go.
// source: batteship.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	batteship.proto

It has these top-level messages:
	Ship
	IdMessage
	EmptyMessage
	ConnectRequest
	Notification
	GameStatus
	PlayRequest
	PlayReply
	PlaceRequest
	PlaceReply
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto1.ProtoPackageIsVersion1

type Tile int32

const (
	Tile_UNKNOWN Tile = 0
	Tile_SEA     Tile = 1
	Tile_SHIP    Tile = 2
	Tile_SUNK    Tile = 3
)

var Tile_name = map[int32]string{
	0: "UNKNOWN",
	1: "SEA",
	2: "SHIP",
	3: "SUNK",
}
var Tile_value = map[string]int32{
	"UNKNOWN": 0,
	"SEA":     1,
	"SHIP":    2,
	"SUNK":    3,
}

func (x Tile) String() string {
	return proto1.EnumName(Tile_name, int32(x))
}
func (Tile) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// Whether the client has won or not
type GameStatus_Status int32

const (
	GameStatus_PLAYING GameStatus_Status = 0
	GameStatus_VICTORY GameStatus_Status = 1
	GameStatus_DEFEAT  GameStatus_Status = 2
)

var GameStatus_Status_name = map[int32]string{
	0: "PLAYING",
	1: "VICTORY",
	2: "DEFEAT",
}
var GameStatus_Status_value = map[string]int32{
	"PLAYING": 0,
	"VICTORY": 1,
	"DEFEAT":  2,
}

func (x GameStatus_Status) String() string {
	return proto1.EnumName(GameStatus_Status_name, int32(x))
}
func (GameStatus_Status) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{5, 0} }

type PlayReply_Status int32

const (
	PlayReply_ACCEPTED         PlayReply_Status = 0
	PlayReply_NOT_YOUR_TURN    PlayReply_Status = 1
	PlayReply_INVALID_POSITION PlayReply_Status = 2
)

var PlayReply_Status_name = map[int32]string{
	0: "ACCEPTED",
	1: "NOT_YOUR_TURN",
	2: "INVALID_POSITION",
}
var PlayReply_Status_value = map[string]int32{
	"ACCEPTED":         0,
	"NOT_YOUR_TURN":    1,
	"INVALID_POSITION": 2,
}

func (x PlayReply_Status) String() string {
	return proto1.EnumName(PlayReply_Status_name, int32(x))
}
func (PlayReply_Status) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{7, 0} }

type Ship struct {
	Pos  int32 `protobuf:"varint,1,opt,name=pos" json:"pos,omitempty"`
	Vert bool  `protobuf:"varint,2,opt,name=vert" json:"vert,omitempty"`
	Size int32 `protobuf:"varint,3,opt,name=size" json:"size,omitempty"`
}

func (m *Ship) Reset()                    { *m = Ship{} }
func (m *Ship) String() string            { return proto1.CompactTextString(m) }
func (*Ship) ProtoMessage()               {}
func (*Ship) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// ---- Common
type IdMessage struct {
	Id int32 `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
}

func (m *IdMessage) Reset()                    { *m = IdMessage{} }
func (m *IdMessage) String() string            { return proto1.CompactTextString(m) }
func (*IdMessage) ProtoMessage()               {}
func (*IdMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type EmptyMessage struct {
}

func (m *EmptyMessage) Reset()                    { *m = EmptyMessage{} }
func (m *EmptyMessage) String() string            { return proto1.CompactTextString(m) }
func (*EmptyMessage) ProtoMessage()               {}
func (*EmptyMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

// ---- Connect
type ConnectRequest struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *ConnectRequest) Reset()                    { *m = ConnectRequest{} }
func (m *ConnectRequest) String() string            { return proto1.CompactTextString(m) }
func (*ConnectRequest) ProtoMessage()               {}
func (*ConnectRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

// ---- Notifications
type Notification struct {
	// Types that are valid to be assigned to Body:
	//	*Notification_ConnectReply
	//	*Notification_GameStatus
	//	*Notification_GameWillStart
	Body isNotification_Body `protobuf_oneof:"body"`
}

func (m *Notification) Reset()                    { *m = Notification{} }
func (m *Notification) String() string            { return proto1.CompactTextString(m) }
func (*Notification) ProtoMessage()               {}
func (*Notification) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type isNotification_Body interface {
	isNotification_Body()
}

type Notification_ConnectReply struct {
	ConnectReply *IdMessage `protobuf:"bytes,1,opt,name=connect_reply,json=connectReply,oneof"`
}
type Notification_GameStatus struct {
	GameStatus *GameStatus `protobuf:"bytes,2,opt,name=game_status,json=gameStatus,oneof"`
}
type Notification_GameWillStart struct {
	GameWillStart *EmptyMessage `protobuf:"bytes,3,opt,name=game_will_start,json=gameWillStart,oneof"`
}

func (*Notification_ConnectReply) isNotification_Body()  {}
func (*Notification_GameStatus) isNotification_Body()    {}
func (*Notification_GameWillStart) isNotification_Body() {}

func (m *Notification) GetBody() isNotification_Body {
	if m != nil {
		return m.Body
	}
	return nil
}

func (m *Notification) GetConnectReply() *IdMessage {
	if x, ok := m.GetBody().(*Notification_ConnectReply); ok {
		return x.ConnectReply
	}
	return nil
}

func (m *Notification) GetGameStatus() *GameStatus {
	if x, ok := m.GetBody().(*Notification_GameStatus); ok {
		return x.GameStatus
	}
	return nil
}

func (m *Notification) GetGameWillStart() *EmptyMessage {
	if x, ok := m.GetBody().(*Notification_GameWillStart); ok {
		return x.GameWillStart
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Notification) XXX_OneofFuncs() (func(msg proto1.Message, b *proto1.Buffer) error, func(msg proto1.Message, tag, wire int, b *proto1.Buffer) (bool, error), func(msg proto1.Message) (n int), []interface{}) {
	return _Notification_OneofMarshaler, _Notification_OneofUnmarshaler, _Notification_OneofSizer, []interface{}{
		(*Notification_ConnectReply)(nil),
		(*Notification_GameStatus)(nil),
		(*Notification_GameWillStart)(nil),
	}
}

func _Notification_OneofMarshaler(msg proto1.Message, b *proto1.Buffer) error {
	m := msg.(*Notification)
	// body
	switch x := m.Body.(type) {
	case *Notification_ConnectReply:
		b.EncodeVarint(1<<3 | proto1.WireBytes)
		if err := b.EncodeMessage(x.ConnectReply); err != nil {
			return err
		}
	case *Notification_GameStatus:
		b.EncodeVarint(2<<3 | proto1.WireBytes)
		if err := b.EncodeMessage(x.GameStatus); err != nil {
			return err
		}
	case *Notification_GameWillStart:
		b.EncodeVarint(3<<3 | proto1.WireBytes)
		if err := b.EncodeMessage(x.GameWillStart); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Notification.Body has unexpected type %T", x)
	}
	return nil
}

func _Notification_OneofUnmarshaler(msg proto1.Message, tag, wire int, b *proto1.Buffer) (bool, error) {
	m := msg.(*Notification)
	switch tag {
	case 1: // body.connect_reply
		if wire != proto1.WireBytes {
			return true, proto1.ErrInternalBadWireType
		}
		msg := new(IdMessage)
		err := b.DecodeMessage(msg)
		m.Body = &Notification_ConnectReply{msg}
		return true, err
	case 2: // body.game_status
		if wire != proto1.WireBytes {
			return true, proto1.ErrInternalBadWireType
		}
		msg := new(GameStatus)
		err := b.DecodeMessage(msg)
		m.Body = &Notification_GameStatus{msg}
		return true, err
	case 3: // body.game_will_start
		if wire != proto1.WireBytes {
			return true, proto1.ErrInternalBadWireType
		}
		msg := new(EmptyMessage)
		err := b.DecodeMessage(msg)
		m.Body = &Notification_GameWillStart{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Notification_OneofSizer(msg proto1.Message) (n int) {
	m := msg.(*Notification)
	// body
	switch x := m.Body.(type) {
	case *Notification_ConnectReply:
		s := proto1.Size(x.ConnectReply)
		n += proto1.SizeVarint(1<<3 | proto1.WireBytes)
		n += proto1.SizeVarint(uint64(s))
		n += s
	case *Notification_GameStatus:
		s := proto1.Size(x.GameStatus)
		n += proto1.SizeVarint(2<<3 | proto1.WireBytes)
		n += proto1.SizeVarint(uint64(s))
		n += s
	case *Notification_GameWillStart:
		s := proto1.Size(x.GameWillStart)
		n += proto1.SizeVarint(3<<3 | proto1.WireBytes)
		n += proto1.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// ---- Game
type GameStatus struct {
	// Whether it is the caller's turn to play or not
	Play   bool              `protobuf:"varint,1,opt,name=play" json:"play,omitempty"`
	Status GameStatus_Status `protobuf:"varint,2,opt,name=status,enum=battleship.GameStatus_Status" json:"status,omitempty"`
}

func (m *GameStatus) Reset()                    { *m = GameStatus{} }
func (m *GameStatus) String() string            { return proto1.CompactTextString(m) }
func (*GameStatus) ProtoMessage()               {}
func (*GameStatus) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

// ---- Play
type PlayRequest struct {
	Id       int32 `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Position int32 `protobuf:"varint,2,opt,name=position" json:"position,omitempty"`
}

func (m *PlayRequest) Reset()                    { *m = PlayRequest{} }
func (m *PlayRequest) String() string            { return proto1.CompactTextString(m) }
func (*PlayRequest) ProtoMessage()               {}
func (*PlayRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

type PlayReply struct {
	Tile   Tile             `protobuf:"varint,1,opt,name=tile,enum=battleship.Tile" json:"tile,omitempty"`
	Status PlayReply_Status `protobuf:"varint,2,opt,name=status,enum=battleship.PlayReply_Status" json:"status,omitempty"`
}

func (m *PlayReply) Reset()                    { *m = PlayReply{} }
func (m *PlayReply) String() string            { return proto1.CompactTextString(m) }
func (*PlayReply) ProtoMessage()               {}
func (*PlayReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

// ---- Place
type PlaceRequest struct {
	Id    int32   `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Ships []*Ship `protobuf:"bytes,2,rep,name=ships" json:"ships,omitempty"`
}

func (m *PlaceRequest) Reset()                    { *m = PlaceRequest{} }
func (m *PlaceRequest) String() string            { return proto1.CompactTextString(m) }
func (*PlaceRequest) ProtoMessage()               {}
func (*PlaceRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *PlaceRequest) GetShips() []*Ship {
	if m != nil {
		return m.Ships
	}
	return nil
}

type PlaceReply struct {
	Valid bool `protobuf:"varint,1,opt,name=valid" json:"valid,omitempty"`
}

func (m *PlaceReply) Reset()                    { *m = PlaceReply{} }
func (m *PlaceReply) String() string            { return proto1.CompactTextString(m) }
func (*PlaceReply) ProtoMessage()               {}
func (*PlaceReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func init() {
	proto1.RegisterType((*Ship)(nil), "battleship.Ship")
	proto1.RegisterType((*IdMessage)(nil), "battleship.IdMessage")
	proto1.RegisterType((*EmptyMessage)(nil), "battleship.EmptyMessage")
	proto1.RegisterType((*ConnectRequest)(nil), "battleship.ConnectRequest")
	proto1.RegisterType((*Notification)(nil), "battleship.Notification")
	proto1.RegisterType((*GameStatus)(nil), "battleship.GameStatus")
	proto1.RegisterType((*PlayRequest)(nil), "battleship.PlayRequest")
	proto1.RegisterType((*PlayReply)(nil), "battleship.PlayReply")
	proto1.RegisterType((*PlaceRequest)(nil), "battleship.PlaceRequest")
	proto1.RegisterType((*PlaceReply)(nil), "battleship.PlaceReply")
	proto1.RegisterEnum("battleship.Tile", Tile_name, Tile_value)
	proto1.RegisterEnum("battleship.GameStatus_Status", GameStatus_Status_name, GameStatus_Status_value)
	proto1.RegisterEnum("battleship.PlayReply_Status", PlayReply_Status_name, PlayReply_Status_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion2

// Client API for Battleship service

type BattleshipClient interface {
	// Connects to the server to get an id to use in other messages
	Connect(ctx context.Context, in *ConnectRequest, opts ...grpc.CallOption) (Battleship_ConnectClient, error)
	// Disconnect from the server
	Disconnect(ctx context.Context, in *IdMessage, opts ...grpc.CallOption) (*EmptyMessage, error)
	// Play
	Play(ctx context.Context, in *PlayRequest, opts ...grpc.CallOption) (*PlayReply, error)
	// Register ship placement
	Place(ctx context.Context, in *PlaceRequest, opts ...grpc.CallOption) (*PlaceReply, error)
}

type battleshipClient struct {
	cc *grpc.ClientConn
}

func NewBattleshipClient(cc *grpc.ClientConn) BattleshipClient {
	return &battleshipClient{cc}
}

func (c *battleshipClient) Connect(ctx context.Context, in *ConnectRequest, opts ...grpc.CallOption) (Battleship_ConnectClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Battleship_serviceDesc.Streams[0], c.cc, "/battleship.Battleship/Connect", opts...)
	if err != nil {
		return nil, err
	}
	x := &battleshipConnectClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Battleship_ConnectClient interface {
	Recv() (*Notification, error)
	grpc.ClientStream
}

type battleshipConnectClient struct {
	grpc.ClientStream
}

func (x *battleshipConnectClient) Recv() (*Notification, error) {
	m := new(Notification)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *battleshipClient) Disconnect(ctx context.Context, in *IdMessage, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := grpc.Invoke(ctx, "/battleship.Battleship/Disconnect", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *battleshipClient) Play(ctx context.Context, in *PlayRequest, opts ...grpc.CallOption) (*PlayReply, error) {
	out := new(PlayReply)
	err := grpc.Invoke(ctx, "/battleship.Battleship/Play", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *battleshipClient) Place(ctx context.Context, in *PlaceRequest, opts ...grpc.CallOption) (*PlaceReply, error) {
	out := new(PlaceReply)
	err := grpc.Invoke(ctx, "/battleship.Battleship/Place", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Battleship service

type BattleshipServer interface {
	// Connects to the server to get an id to use in other messages
	Connect(*ConnectRequest, Battleship_ConnectServer) error
	// Disconnect from the server
	Disconnect(context.Context, *IdMessage) (*EmptyMessage, error)
	// Play
	Play(context.Context, *PlayRequest) (*PlayReply, error)
	// Register ship placement
	Place(context.Context, *PlaceRequest) (*PlaceReply, error)
}

func RegisterBattleshipServer(s *grpc.Server, srv BattleshipServer) {
	s.RegisterService(&_Battleship_serviceDesc, srv)
}

func _Battleship_Connect_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ConnectRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BattleshipServer).Connect(m, &battleshipConnectServer{stream})
}

type Battleship_ConnectServer interface {
	Send(*Notification) error
	grpc.ServerStream
}

type battleshipConnectServer struct {
	grpc.ServerStream
}

func (x *battleshipConnectServer) Send(m *Notification) error {
	return x.ServerStream.SendMsg(m)
}

func _Battleship_Disconnect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IdMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BattleshipServer).Disconnect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/battleship.Battleship/Disconnect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BattleshipServer).Disconnect(ctx, req.(*IdMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Battleship_Play_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlayRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BattleshipServer).Play(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/battleship.Battleship/Play",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BattleshipServer).Play(ctx, req.(*PlayRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Battleship_Place_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlaceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BattleshipServer).Place(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/battleship.Battleship/Place",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BattleshipServer).Place(ctx, req.(*PlaceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Battleship_serviceDesc = grpc.ServiceDesc{
	ServiceName: "battleship.Battleship",
	HandlerType: (*BattleshipServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Disconnect",
			Handler:    _Battleship_Disconnect_Handler,
		},
		{
			MethodName: "Play",
			Handler:    _Battleship_Play_Handler,
		},
		{
			MethodName: "Place",
			Handler:    _Battleship_Place_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Connect",
			Handler:       _Battleship_Connect_Handler,
			ServerStreams: true,
		},
	},
}

var fileDescriptor0 = []byte{
	// 628 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x74, 0x54, 0xdd, 0x6e, 0xd3, 0x4c,
	0x10, 0x8d, 0x63, 0x3b, 0x49, 0x27, 0x69, 0xea, 0x6f, 0xd5, 0xf6, 0x8b, 0x02, 0x48, 0x68, 0x55,
	0x21, 0xc4, 0x45, 0x54, 0x05, 0x90, 0x40, 0x20, 0x95, 0xfc, 0xb8, 0xad, 0x45, 0xb1, 0xa3, 0x8d,
	0xdb, 0xaa, 0xdc, 0x44, 0x4e, 0x62, 0x82, 0x25, 0x37, 0x36, 0xb1, 0x5b, 0x14, 0x5e, 0x00, 0xf1,
	0x26, 0xbc, 0x0f, 0x2f, 0xc4, 0xec, 0xda, 0x75, 0x9d, 0x90, 0x5e, 0x79, 0x76, 0x67, 0xce, 0xcc,
	0x39, 0x33, 0x3b, 0x86, 0x9d, 0xb1, 0x13, 0xc7, 0x6e, 0xf4, 0xd5, 0x0b, 0x5b, 0xe1, 0x22, 0x88,
	0x03, 0x02, 0xfc, 0xc2, 0x17, 0x37, 0xf4, 0x03, 0x28, 0x43, 0xfc, 0x12, 0x0d, 0xe4, 0x30, 0x88,
	0x1a, 0xd2, 0x53, 0xe9, 0xb9, 0xca, 0xb8, 0x49, 0x08, 0x28, 0xb7, 0xee, 0x22, 0x6e, 0x14, 0xf1,
	0xaa, 0xc2, 0x84, 0xcd, 0xef, 0x22, 0xef, 0x87, 0xdb, 0x90, 0x45, 0x98, 0xb0, 0xe9, 0x23, 0xd8,
	0x32, 0xa6, 0x9f, 0xdc, 0x28, 0x72, 0x66, 0x2e, 0xa9, 0x43, 0xd1, 0x9b, 0xa6, 0x59, 0xd0, 0xa2,
	0x75, 0xa8, 0xe9, 0xd7, 0x61, 0xbc, 0x4c, 0xfd, 0xf4, 0x00, 0xea, 0xbd, 0x60, 0x3e, 0x77, 0x27,
	0x31, 0x73, 0xbf, 0xdd, 0xb8, 0x91, 0x48, 0x39, 0x77, 0xae, 0x5d, 0x81, 0xd9, 0x62, 0xc2, 0xa6,
	0x7f, 0x24, 0xa8, 0x99, 0x41, 0xec, 0x7d, 0xf1, 0x26, 0x4e, 0xec, 0x05, 0x73, 0xf2, 0x1e, 0xb6,
	0x27, 0x09, 0x6c, 0xb4, 0x70, 0x43, 0x7f, 0x29, 0xa2, 0xab, 0xed, 0xbd, 0xd6, 0xbd, 0x92, 0x56,
	0x46, 0xe2, 0xb4, 0xc0, 0x6a, 0x93, 0xbb, 0x22, 0x18, 0x4c, 0xde, 0x42, 0x75, 0x86, 0x69, 0x47,
	0x51, 0xec, 0xc4, 0x37, 0x91, 0x10, 0x54, 0x6d, 0xef, 0xe7, 0xb1, 0x27, 0xe8, 0x1e, 0x0a, 0x2f,
	0x82, 0x61, 0x96, 0x9d, 0x48, 0x17, 0x76, 0x04, 0xf4, 0xbb, 0xe7, 0xfb, 0x1c, 0x8f, 0xfd, 0x90,
	0x05, 0xbc, 0x91, 0x87, 0xe7, 0x25, 0x62, 0x82, 0x6d, 0x0e, 0xb9, 0x44, 0xc4, 0x90, 0x03, 0xba,
	0x25, 0x50, 0xc6, 0xc1, 0x74, 0x49, 0x7f, 0x4a, 0x00, 0xf7, 0x85, 0xb8, 0xf0, 0xd0, 0x77, 0x12,
	0x29, 0xd8, 0x5f, 0x6e, 0x93, 0xd7, 0x50, 0xca, 0x91, 0xac, 0xb7, 0x9f, 0x6c, 0x26, 0xd9, 0x4a,
	0x3e, 0x2c, 0x0d, 0xa6, 0x2d, 0x28, 0xa5, 0x49, 0xab, 0x50, 0x1e, 0x9c, 0x75, 0xae, 0x0c, 0xf3,
	0x44, 0x2b, 0xf0, 0xc3, 0x85, 0xd1, 0xb3, 0x2d, 0x76, 0xa5, 0x49, 0x04, 0xa0, 0xd4, 0xd7, 0x8f,
	0xf5, 0x8e, 0xad, 0x15, 0x29, 0x36, 0x64, 0x80, 0xe5, 0xee, 0x46, 0xb0, 0x36, 0x34, 0xd2, 0x84,
	0x0a, 0x3e, 0x00, 0x8f, 0x77, 0x5e, 0xf0, 0x50, 0x59, 0x76, 0xa6, 0xbf, 0x25, 0xd8, 0x4a, 0xb0,
	0xbc, 0xb3, 0x07, 0xa0, 0xc4, 0x9e, 0x9f, 0x0c, 0xaf, 0xde, 0xd6, 0xf2, 0x6c, 0x6d, 0xbc, 0x67,
	0xc2, 0x4b, 0x5e, 0xad, 0xa9, 0x7a, 0x9c, 0x8f, 0xcb, 0x92, 0xad, 0x8b, 0x3a, 0xca, 0x44, 0xd5,
	0xa0, 0xd2, 0xe9, 0xf5, 0xf4, 0x81, 0xad, 0xf7, 0x51, 0xd5, 0x7f, 0xb0, 0x6d, 0x5a, 0xf6, 0xe8,
	0xca, 0x3a, 0x67, 0x23, 0xfb, 0x9c, 0x99, 0xa8, 0x6d, 0x17, 0x34, 0xc3, 0xbc, 0xe8, 0x9c, 0x19,
	0xfd, 0xd1, 0xc0, 0x1a, 0x1a, 0xb6, 0x61, 0x99, 0xa8, 0xf2, 0x18, 0x6a, 0x98, 0x7c, 0xe2, 0x3e,
	0x24, 0xf3, 0x19, 0xa8, 0x9c, 0x01, 0x67, 0x25, 0xe3, 0x44, 0x57, 0xd8, 0xf3, 0x9d, 0x60, 0x89,
	0x9b, 0x52, 0x80, 0x34, 0x0f, 0x97, 0xbc, 0x0b, 0xea, 0xad, 0xe3, 0xa7, 0x89, 0x2a, 0x2c, 0x39,
	0xbc, 0x38, 0x04, 0x85, 0x0b, 0xe6, 0x2d, 0x3f, 0x37, 0x3f, 0x9a, 0xd6, 0xa5, 0x89, 0x4c, 0xcb,
	0x20, 0x0f, 0xf5, 0x0e, 0xf2, 0xab, 0xe0, 0x92, 0x9d, 0x1a, 0x03, 0xad, 0x28, 0x2c, 0x0c, 0xd0,
	0xe4, 0xf6, 0xaf, 0x22, 0x40, 0x37, 0x2b, 0x48, 0x7a, 0x50, 0x4e, 0x17, 0x83, 0x34, 0xf3, 0x44,
	0x56, 0xb7, 0xa5, 0xb9, 0xf2, 0xec, 0xf2, 0x2b, 0x42, 0x0b, 0x87, 0x12, 0x39, 0x02, 0xe8, 0x7b,
	0x51, 0xfa, 0xf6, 0xc9, 0xe6, 0xed, 0x68, 0x3e, 0xf8, 0x72, 0x69, 0x81, 0xbc, 0x01, 0x85, 0xcf,
	0x83, 0xfc, 0xff, 0xef, 0x84, 0x92, 0xfa, 0x7b, 0x1b, 0x47, 0x87, 0xc8, 0x77, 0xa0, 0x8a, 0x26,
	0x91, 0xc6, 0x5a, 0x44, 0xd6, 0xff, 0xe6, 0xfe, 0x06, 0x8f, 0x00, 0x77, 0xcb, 0x9f, 0x55, 0xf1,
	0x67, 0x1a, 0x97, 0xc4, 0xe7, 0xe5, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x52, 0x3e, 0xea, 0xba,
	0xb3, 0x04, 0x00, 0x00,
}
