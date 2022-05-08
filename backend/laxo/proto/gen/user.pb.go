// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: user.proto

package gen

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Notification struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID              string `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	RedisID         string `protobuf:"bytes,2,opt,name=redisID,proto3" json:"redisID,omitempty"`
	GroupID         string `protobuf:"bytes,3,opt,name=groupID,proto3" json:"groupID,omitempty"`
	Created         int64  `protobuf:"varint,4,opt,name=created,proto3" json:"created,omitempty"`
	Read            int64  `protobuf:"varint,5,opt,name=read,proto3" json:"read,omitempty"`
	CurrentMainStep int64  `protobuf:"varint,6,opt,name=currentMainStep,proto3" json:"currentMainStep,omitempty"`
	CurrentSubStep  int64  `protobuf:"varint,7,opt,name=currentSubStep,proto3" json:"currentSubStep,omitempty"`
	MainMessage     string `protobuf:"bytes,8,opt,name=mainMessage,proto3" json:"mainMessage,omitempty"`
	SubMessage      string `protobuf:"bytes,9,opt,name=subMessage,proto3" json:"subMessage,omitempty"`
}

func (x *Notification) Reset() {
	*x = Notification{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Notification) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Notification) ProtoMessage() {}

func (x *Notification) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Notification.ProtoReflect.Descriptor instead.
func (*Notification) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{0}
}

func (x *Notification) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *Notification) GetRedisID() string {
	if x != nil {
		return x.RedisID
	}
	return ""
}

func (x *Notification) GetGroupID() string {
	if x != nil {
		return x.GroupID
	}
	return ""
}

func (x *Notification) GetCreated() int64 {
	if x != nil {
		return x.Created
	}
	return 0
}

func (x *Notification) GetRead() int64 {
	if x != nil {
		return x.Read
	}
	return 0
}

func (x *Notification) GetCurrentMainStep() int64 {
	if x != nil {
		return x.CurrentMainStep
	}
	return 0
}

func (x *Notification) GetCurrentSubStep() int64 {
	if x != nil {
		return x.CurrentSubStep
	}
	return 0
}

func (x *Notification) GetMainMessage() string {
	if x != nil {
		return x.MainMessage
	}
	return ""
}

func (x *Notification) GetSubMessage() string {
	if x != nil {
		return x.SubMessage
	}
	return ""
}

type NotificationGroup struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID             string `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	UserID         string `protobuf:"bytes,2,opt,name=userID,proto3" json:"userID,omitempty"`
	WorkflowID     string `protobuf:"bytes,3,opt,name=workflowID,proto3" json:"workflowID,omitempty"`
	EntityID       string `protobuf:"bytes,4,opt,name=entityID,proto3" json:"entityID,omitempty"`
	EntityType     string `protobuf:"bytes,5,opt,name=entityType,proto3" json:"entityType,omitempty"`
	TotalMainSteps int64  `protobuf:"varint,6,opt,name=totalMainSteps,proto3" json:"totalMainSteps,omitempty"`
	TotalSubSteps  int64  `protobuf:"varint,7,opt,name=totalSubSteps,proto3" json:"totalSubSteps,omitempty"`
}

func (x *NotificationGroup) Reset() {
	*x = NotificationGroup{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotificationGroup) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotificationGroup) ProtoMessage() {}

func (x *NotificationGroup) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotificationGroup.ProtoReflect.Descriptor instead.
func (*NotificationGroup) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{1}
}

func (x *NotificationGroup) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *NotificationGroup) GetUserID() string {
	if x != nil {
		return x.UserID
	}
	return ""
}

func (x *NotificationGroup) GetWorkflowID() string {
	if x != nil {
		return x.WorkflowID
	}
	return ""
}

func (x *NotificationGroup) GetEntityID() string {
	if x != nil {
		return x.EntityID
	}
	return ""
}

func (x *NotificationGroup) GetEntityType() string {
	if x != nil {
		return x.EntityType
	}
	return ""
}

func (x *NotificationGroup) GetTotalMainSteps() int64 {
	if x != nil {
		return x.TotalMainSteps
	}
	return 0
}

func (x *NotificationGroup) GetTotalSubSteps() int64 {
	if x != nil {
		return x.TotalSubSteps
	}
	return 0
}

type NotificationUpdateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NotificationID      string `protobuf:"bytes,1,opt,name=notificationID,proto3" json:"notificationID,omitempty"`
	NotificationGroupID string `protobuf:"bytes,2,opt,name=notificationGroupID,proto3" json:"notificationGroupID,omitempty"`
}

func (x *NotificationUpdateRequest) Reset() {
	*x = NotificationUpdateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotificationUpdateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotificationUpdateRequest) ProtoMessage() {}

func (x *NotificationUpdateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotificationUpdateRequest.ProtoReflect.Descriptor instead.
func (*NotificationUpdateRequest) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{2}
}

func (x *NotificationUpdateRequest) GetNotificationID() string {
	if x != nil {
		return x.NotificationID
	}
	return ""
}

func (x *NotificationUpdateRequest) GetNotificationGroupID() string {
	if x != nil {
		return x.NotificationGroupID
	}
	return ""
}

type NotificationUpdateReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Notification      *Notification      `protobuf:"bytes,1,opt,name=notification,proto3" json:"notification,omitempty"`
	NotificationGroup *NotificationGroup `protobuf:"bytes,2,opt,name=notificationGroup,proto3" json:"notificationGroup,omitempty"`
}

func (x *NotificationUpdateReply) Reset() {
	*x = NotificationUpdateReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotificationUpdateReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotificationUpdateReply) ProtoMessage() {}

func (x *NotificationUpdateReply) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotificationUpdateReply.ProtoReflect.Descriptor instead.
func (*NotificationUpdateReply) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{3}
}

func (x *NotificationUpdateReply) GetNotification() *Notification {
	if x != nil {
		return x.Notification
	}
	return nil
}

func (x *NotificationUpdateReply) GetNotificationGroup() *NotificationGroup {
	if x != nil {
		return x.NotificationGroup
	}
	return nil
}

var File_user_proto protoreflect.FileDescriptor

var file_user_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x75, 0x73,
	0x65, 0x72, 0x22, 0x94, 0x02, 0x0a, 0x0c, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x49, 0x44, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x65, 0x64, 0x69, 0x73, 0x49, 0x44, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x72, 0x65, 0x64, 0x69, 0x73, 0x49, 0x44, 0x12, 0x18, 0x0a,
	0x07, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x67, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x44, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x65, 0x61, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x04, 0x72, 0x65, 0x61, 0x64, 0x12, 0x28, 0x0a, 0x0f, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74,
	0x4d, 0x61, 0x69, 0x6e, 0x53, 0x74, 0x65, 0x70, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0f,
	0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x4d, 0x61, 0x69, 0x6e, 0x53, 0x74, 0x65, 0x70, 0x12,
	0x26, 0x0a, 0x0e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x53, 0x75, 0x62, 0x53, 0x74, 0x65,
	0x70, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74,
	0x53, 0x75, 0x62, 0x53, 0x74, 0x65, 0x70, 0x12, 0x20, 0x0a, 0x0b, 0x6d, 0x61, 0x69, 0x6e, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6d, 0x61,
	0x69, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x75, 0x62,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73,
	0x75, 0x62, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0xe5, 0x01, 0x0a, 0x11, 0x4e, 0x6f,
	0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12,
	0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44, 0x12,
	0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x12, 0x1e, 0x0a, 0x0a, 0x77, 0x6f, 0x72, 0x6b, 0x66,
	0x6c, 0x6f, 0x77, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x77, 0x6f, 0x72,
	0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x49, 0x44, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x49, 0x44, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x49, 0x44, 0x12, 0x1e, 0x0a, 0x0a, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x54, 0x79, 0x70,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x26, 0x0a, 0x0e, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x4d, 0x61, 0x69, 0x6e,
	0x53, 0x74, 0x65, 0x70, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0e, 0x74, 0x6f, 0x74,
	0x61, 0x6c, 0x4d, 0x61, 0x69, 0x6e, 0x53, 0x74, 0x65, 0x70, 0x73, 0x12, 0x24, 0x0a, 0x0d, 0x74,
	0x6f, 0x74, 0x61, 0x6c, 0x53, 0x75, 0x62, 0x53, 0x74, 0x65, 0x70, 0x73, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0d, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x53, 0x75, 0x62, 0x53, 0x74, 0x65, 0x70,
	0x73, 0x22, 0x75, 0x0a, 0x19, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x26,
	0x0a, 0x0e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x12, 0x30, 0x0a, 0x13, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x44, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x13, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x44, 0x22, 0x98, 0x01, 0x0a, 0x17, 0x4e, 0x6f, 0x74,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x12, 0x36, 0x0a, 0x0c, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x75, 0x73, 0x65,
	0x72, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0c,
	0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x45, 0x0a, 0x11,
	0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75,
	0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x4e,
	0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70,
	0x52, 0x11, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x47, 0x72,
	0x6f, 0x75, 0x70, 0x32, 0x6a, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x5b, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x1f, 0x2e, 0x75, 0x73,
	0x65, 0x72, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x75,
	0x73, 0x65, 0x72, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x30, 0x01, 0x42,
	0x1d, 0x5a, 0x1b, 0x6c, 0x61, 0x78, 0x6f, 0x2e, 0x76, 0x6e, 0x2f, 0x6c, 0x61, 0x78, 0x6f, 0x2f,
	0x6c, 0x61, 0x78, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x65, 0x6e, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_user_proto_rawDescOnce sync.Once
	file_user_proto_rawDescData = file_user_proto_rawDesc
)

func file_user_proto_rawDescGZIP() []byte {
	file_user_proto_rawDescOnce.Do(func() {
		file_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_user_proto_rawDescData)
	})
	return file_user_proto_rawDescData
}

var file_user_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_user_proto_goTypes = []interface{}{
	(*Notification)(nil),              // 0: user.Notification
	(*NotificationGroup)(nil),         // 1: user.NotificationGroup
	(*NotificationUpdateRequest)(nil), // 2: user.NotificationUpdateRequest
	(*NotificationUpdateReply)(nil),   // 3: user.NotificationUpdateReply
}
var file_user_proto_depIdxs = []int32{
	0, // 0: user.NotificationUpdateReply.notification:type_name -> user.Notification
	1, // 1: user.NotificationUpdateReply.notificationGroup:type_name -> user.NotificationGroup
	2, // 2: user.UserService.GetNotificationUpdate:input_type -> user.NotificationUpdateRequest
	3, // 3: user.UserService.GetNotificationUpdate:output_type -> user.NotificationUpdateReply
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_user_proto_init() }
func file_user_proto_init() {
	if File_user_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_user_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Notification); i {
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
		file_user_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotificationGroup); i {
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
		file_user_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotificationUpdateRequest); i {
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
		file_user_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotificationUpdateReply); i {
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
			RawDescriptor: file_user_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_user_proto_goTypes,
		DependencyIndexes: file_user_proto_depIdxs,
		MessageInfos:      file_user_proto_msgTypes,
	}.Build()
	File_user_proto = out.File
	file_user_proto_rawDesc = nil
	file_user_proto_goTypes = nil
	file_user_proto_depIdxs = nil
}
