// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: user_auth/user_auth.proto

package userauthpb

import (
	userpb "github.com/mamataliev-dev/social-platform/api/gen/user"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RegisterRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserName      string                 `protobuf:"bytes,1,opt,name=user_name,json=userName,proto3" json:"user_name,omitempty"`
	Email         string                 `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	Password      string                 `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"` // plain-text; service hashes internally
	Nickname      string                 `protobuf:"bytes,4,opt,name=nickname,proto3" json:"nickname,omitempty"`
	Bio           string                 `protobuf:"bytes,5,opt,name=bio,proto3" json:"bio,omitempty"`
	AvatarUrl     string                 `protobuf:"bytes,6,opt,name=avatar_url,json=avatarUrl,proto3" json:"avatar_url,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegisterRequest) Reset() {
	*x = RegisterRequest{}
	mi := &file_user_auth_user_auth_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterRequest) ProtoMessage() {}

func (x *RegisterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_auth_user_auth_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterRequest.ProtoReflect.Descriptor instead.
func (*RegisterRequest) Descriptor() ([]byte, []int) {
	return file_user_auth_user_auth_proto_rawDescGZIP(), []int{0}
}

func (x *RegisterRequest) GetUserName() string {
	if x != nil {
		return x.UserName
	}
	return ""
}

func (x *RegisterRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *RegisterRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *RegisterRequest) GetNickname() string {
	if x != nil {
		return x.Nickname
	}
	return ""
}

func (x *RegisterRequest) GetBio() string {
	if x != nil {
		return x.Bio
	}
	return ""
}

func (x *RegisterRequest) GetAvatarUrl() string {
	if x != nil {
		return x.AvatarUrl
	}
	return ""
}

type RegisterResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	User          *userpb.UserProfile    `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegisterResponse) Reset() {
	*x = RegisterResponse{}
	mi := &file_user_auth_user_auth_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterResponse) ProtoMessage() {}

func (x *RegisterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_auth_user_auth_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterResponse.ProtoReflect.Descriptor instead.
func (*RegisterResponse) Descriptor() ([]byte, []int) {
	return file_user_auth_user_auth_proto_rawDescGZIP(), []int{1}
}

func (x *RegisterResponse) GetUser() *userpb.UserProfile {
	if x != nil {
		return x.User
	}
	return nil
}

type LoginRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Email         string                 `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Password      string                 `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"` // plain-text; service verifies internally
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoginRequest) Reset() {
	*x = LoginRequest{}
	mi := &file_user_auth_user_auth_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoginRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginRequest) ProtoMessage() {}

func (x *LoginRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_auth_user_auth_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginRequest.ProtoReflect.Descriptor instead.
func (*LoginRequest) Descriptor() ([]byte, []int) {
	return file_user_auth_user_auth_proto_rawDescGZIP(), []int{2}
}

func (x *LoginRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *LoginRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type LoginResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	User          *userpb.UserProfile    `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoginResponse) Reset() {
	*x = LoginResponse{}
	mi := &file_user_auth_user_auth_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoginResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginResponse) ProtoMessage() {}

func (x *LoginResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_auth_user_auth_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginResponse.ProtoReflect.Descriptor instead.
func (*LoginResponse) Descriptor() ([]byte, []int) {
	return file_user_auth_user_auth_proto_rawDescGZIP(), []int{3}
}

func (x *LoginResponse) GetUser() *userpb.UserProfile {
	if x != nil {
		return x.User
	}
	return nil
}

var File_user_auth_user_auth_proto protoreflect.FileDescriptor

const file_user_auth_user_auth_proto_rawDesc = "" +
	"\n" +
	"\x19user_auth/user_auth.proto\x12\tuser_auth\x1a\x0fuser/user.proto\x1a\x1cgoogle/api/annotations.proto\"\xad\x01\n" +
	"\x0fRegisterRequest\x12\x1b\n" +
	"\tuser_name\x18\x01 \x01(\tR\buserName\x12\x14\n" +
	"\x05email\x18\x02 \x01(\tR\x05email\x12\x1a\n" +
	"\bpassword\x18\x03 \x01(\tR\bpassword\x12\x1a\n" +
	"\bnickname\x18\x04 \x01(\tR\bnickname\x12\x10\n" +
	"\x03bio\x18\x05 \x01(\tR\x03bio\x12\x1d\n" +
	"\n" +
	"avatar_url\x18\x06 \x01(\tR\tavatarUrl\"9\n" +
	"\x10RegisterResponse\x12%\n" +
	"\x04user\x18\x01 \x01(\v2\x11.user.UserProfileR\x04user\"@\n" +
	"\fLoginRequest\x12\x14\n" +
	"\x05email\x18\x01 \x01(\tR\x05email\x12\x1a\n" +
	"\bpassword\x18\x02 \x01(\tR\bpassword\"6\n" +
	"\rLoginResponse\x12%\n" +
	"\x04user\x18\x01 \x01(\v2\x11.user.UserProfileR\x04user2\xb7\x01\n" +
	"\vAuthService\x12Y\n" +
	"\bRegister\x12\x1a.user_auth.RegisterRequest\x1a\x1b.user_auth.RegisterResponse\"\x14\x82\xd3\xe4\x93\x02\x0e:\x01*\"\t/register\x12M\n" +
	"\x05Login\x12\x17.user_auth.LoginRequest\x1a\x18.user_auth.LoginResponse\"\x11\x82\xd3\xe4\x93\x02\v:\x01*\"\x06/loginBIZGgithub.com/mamataliev-dev/social-platform/api/gen/userauthpb;userauthpbb\x06proto3"

var (
	file_user_auth_user_auth_proto_rawDescOnce sync.Once
	file_user_auth_user_auth_proto_rawDescData []byte
)

func file_user_auth_user_auth_proto_rawDescGZIP() []byte {
	file_user_auth_user_auth_proto_rawDescOnce.Do(func() {
		file_user_auth_user_auth_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_user_auth_user_auth_proto_rawDesc), len(file_user_auth_user_auth_proto_rawDesc)))
	})
	return file_user_auth_user_auth_proto_rawDescData
}

var file_user_auth_user_auth_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_user_auth_user_auth_proto_goTypes = []any{
	(*RegisterRequest)(nil),    // 0: user_auth.RegisterRequest
	(*RegisterResponse)(nil),   // 1: user_auth.RegisterResponse
	(*LoginRequest)(nil),       // 2: user_auth.LoginRequest
	(*LoginResponse)(nil),      // 3: user_auth.LoginResponse
	(*userpb.UserProfile)(nil), // 4: user.UserProfile
}
var file_user_auth_user_auth_proto_depIdxs = []int32{
	4, // 0: user_auth.RegisterResponse.user:type_name -> user.UserProfile
	4, // 1: user_auth.LoginResponse.user:type_name -> user.UserProfile
	0, // 2: user_auth.AuthService.Register:input_type -> user_auth.RegisterRequest
	2, // 3: user_auth.AuthService.Login:input_type -> user_auth.LoginRequest
	1, // 4: user_auth.AuthService.Register:output_type -> user_auth.RegisterResponse
	3, // 5: user_auth.AuthService.Login:output_type -> user_auth.LoginResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_user_auth_user_auth_proto_init() }
func file_user_auth_user_auth_proto_init() {
	if File_user_auth_user_auth_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_user_auth_user_auth_proto_rawDesc), len(file_user_auth_user_auth_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_user_auth_user_auth_proto_goTypes,
		DependencyIndexes: file_user_auth_user_auth_proto_depIdxs,
		MessageInfos:      file_user_auth_user_auth_proto_msgTypes,
	}.Build()
	File_user_auth_user_auth_proto = out.File
	file_user_auth_user_auth_proto_goTypes = nil
	file_user_auth_user_auth_proto_depIdxs = nil
}
