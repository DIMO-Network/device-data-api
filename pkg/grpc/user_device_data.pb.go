// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: pkg/grpc/user_device_data.proto

package grpc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SignalRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IntegrationId string                 `protobuf:"bytes,1,opt,name=integration_id,json=integrationId,proto3" json:"integration_id,omitempty"`
	FromDate      *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=from_date,json=fromDate,proto3" json:"from_date,omitempty"`
	ToDate        *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=to_date,json=toDate,proto3" json:"to_date,omitempty"`
}

func (x *SignalRequest) Reset() {
	*x = SignalRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_grpc_user_device_data_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignalRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignalRequest) ProtoMessage() {}

func (x *SignalRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_grpc_user_device_data_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignalRequest.ProtoReflect.Descriptor instead.
func (*SignalRequest) Descriptor() ([]byte, []int) {
	return file_pkg_grpc_user_device_data_proto_rawDescGZIP(), []int{0}
}

func (x *SignalRequest) GetIntegrationId() string {
	if x != nil {
		return x.IntegrationId
	}
	return ""
}

func (x *SignalRequest) GetFromDate() *timestamppb.Timestamp {
	if x != nil {
		return x.FromDate
	}
	return nil
}

func (x *SignalRequest) GetToDate() *timestamppb.Timestamp {
	if x != nil {
		return x.ToDate
	}
	return nil
}

type SignalResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*SignalItemResponse `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *SignalResponse) Reset() {
	*x = SignalResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_grpc_user_device_data_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignalResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignalResponse) ProtoMessage() {}

func (x *SignalResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_grpc_user_device_data_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignalResponse.ProtoReflect.Descriptor instead.
func (*SignalResponse) Descriptor() ([]byte, []int) {
	return file_pkg_grpc_user_device_data_proto_rawDescGZIP(), []int{1}
}

func (x *SignalResponse) GetItems() []*SignalItemResponse {
	if x != nil {
		return x.Items
	}
	return nil
}

type SignalItemResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Property     string `protobuf:"bytes,1,opt,name=property,proto3" json:"property,omitempty"`
	TotalCount   int32  `protobuf:"varint,2,opt,name=total_count,json=totalCount,proto3" json:"total_count,omitempty"`
	RequestCount int32  `protobuf:"varint,3,opt,name=request_count,json=requestCount,proto3" json:"request_count,omitempty"`
}

func (x *SignalItemResponse) Reset() {
	*x = SignalItemResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_grpc_user_device_data_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignalItemResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignalItemResponse) ProtoMessage() {}

func (x *SignalItemResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_grpc_user_device_data_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignalItemResponse.ProtoReflect.Descriptor instead.
func (*SignalItemResponse) Descriptor() ([]byte, []int) {
	return file_pkg_grpc_user_device_data_proto_rawDescGZIP(), []int{2}
}

func (x *SignalItemResponse) GetProperty() string {
	if x != nil {
		return x.Property
	}
	return ""
}

func (x *SignalItemResponse) GetTotalCount() int32 {
	if x != nil {
		return x.TotalCount
	}
	return 0
}

func (x *SignalItemResponse) GetRequestCount() int32 {
	if x != nil {
		return x.RequestCount
	}
	return 0
}

type UserDeviceDataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserDeviceId       string `protobuf:"bytes,1,opt,name=user_device_id,json=userDeviceId,proto3" json:"user_device_id,omitempty"`
	DeviceDefinitionId string `protobuf:"bytes,2,opt,name=device_definition_id,json=deviceDefinitionId,proto3" json:"device_definition_id,omitempty"`
	DeviceStyleId      string `protobuf:"bytes,3,opt,name=device_style_id,json=deviceStyleId,proto3" json:"device_style_id,omitempty"`
}

func (x *UserDeviceDataRequest) Reset() {
	*x = UserDeviceDataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_grpc_user_device_data_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserDeviceDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserDeviceDataRequest) ProtoMessage() {}

func (x *UserDeviceDataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_grpc_user_device_data_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserDeviceDataRequest.ProtoReflect.Descriptor instead.
func (*UserDeviceDataRequest) Descriptor() ([]byte, []int) {
	return file_pkg_grpc_user_device_data_proto_rawDescGZIP(), []int{3}
}

func (x *UserDeviceDataRequest) GetUserDeviceId() string {
	if x != nil {
		return x.UserDeviceId
	}
	return ""
}

func (x *UserDeviceDataRequest) GetDeviceDefinitionId() string {
	if x != nil {
		return x.DeviceDefinitionId
	}
	return ""
}

func (x *UserDeviceDataRequest) GetDeviceStyleId() string {
	if x != nil {
		return x.DeviceStyleId
	}
	return ""
}

type UserDeviceDataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Charging             bool                   `protobuf:"varint,1,opt,name=charging,proto3" json:"charging,omitempty"`
	FuelPercentRemaining float64                `protobuf:"fixed64,2,opt,name=fuel_percent_remaining,json=fuelPercentRemaining,proto3" json:"fuel_percent_remaining,omitempty"`
	BatteryCapacity      int64                  `protobuf:"varint,3,opt,name=battery_capacity,json=batteryCapacity,proto3" json:"battery_capacity,omitempty"`
	OilLevel             float64                `protobuf:"fixed64,4,opt,name=oil_level,json=oilLevel,proto3" json:"oil_level,omitempty"`
	Odometer             float64                `protobuf:"fixed64,5,opt,name=odometer,proto3" json:"odometer,omitempty"`
	Latitude             float64                `protobuf:"fixed64,6,opt,name=latitude,proto3" json:"latitude,omitempty"`
	Longitude            float64                `protobuf:"fixed64,7,opt,name=longitude,proto3" json:"longitude,omitempty"`
	Range                float64                `protobuf:"fixed64,8,opt,name=range,proto3" json:"range,omitempty"`
	StateOfCharge        float64                `protobuf:"fixed64,9,opt,name=state_of_charge,json=stateOfCharge,proto3" json:"state_of_charge,omitempty"`
	ChargeLimit          float64                `protobuf:"fixed64,10,opt,name=charge_limit,json=chargeLimit,proto3" json:"charge_limit,omitempty"`
	RecordUpdatedAt      *timestamppb.Timestamp `protobuf:"bytes,11,opt,name=record_updated_at,json=recordUpdatedAt,proto3" json:"record_updated_at,omitempty"`
	RecordCreatedAt      *timestamppb.Timestamp `protobuf:"bytes,12,opt,name=record_created_at,json=recordCreatedAt,proto3" json:"record_created_at,omitempty"`
	TirePressure         *TirePressureResponse  `protobuf:"bytes,13,opt,name=tire_pressure,json=tirePressure,proto3" json:"tire_pressure,omitempty"`
	BatteryVoltage       float64                `protobuf:"fixed64,14,opt,name=battery_voltage,json=batteryVoltage,proto3" json:"battery_voltage,omitempty"`
	AmbientTemp          float64                `protobuf:"fixed64,15,opt,name=ambient_temp,json=ambientTemp,proto3" json:"ambient_temp,omitempty"`
}

func (x *UserDeviceDataResponse) Reset() {
	*x = UserDeviceDataResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_grpc_user_device_data_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserDeviceDataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserDeviceDataResponse) ProtoMessage() {}

func (x *UserDeviceDataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_grpc_user_device_data_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserDeviceDataResponse.ProtoReflect.Descriptor instead.
func (*UserDeviceDataResponse) Descriptor() ([]byte, []int) {
	return file_pkg_grpc_user_device_data_proto_rawDescGZIP(), []int{4}
}

func (x *UserDeviceDataResponse) GetCharging() bool {
	if x != nil {
		return x.Charging
	}
	return false
}

func (x *UserDeviceDataResponse) GetFuelPercentRemaining() float64 {
	if x != nil {
		return x.FuelPercentRemaining
	}
	return 0
}

func (x *UserDeviceDataResponse) GetBatteryCapacity() int64 {
	if x != nil {
		return x.BatteryCapacity
	}
	return 0
}

func (x *UserDeviceDataResponse) GetOilLevel() float64 {
	if x != nil {
		return x.OilLevel
	}
	return 0
}

func (x *UserDeviceDataResponse) GetOdometer() float64 {
	if x != nil {
		return x.Odometer
	}
	return 0
}

func (x *UserDeviceDataResponse) GetLatitude() float64 {
	if x != nil {
		return x.Latitude
	}
	return 0
}

func (x *UserDeviceDataResponse) GetLongitude() float64 {
	if x != nil {
		return x.Longitude
	}
	return 0
}

func (x *UserDeviceDataResponse) GetRange() float64 {
	if x != nil {
		return x.Range
	}
	return 0
}

func (x *UserDeviceDataResponse) GetStateOfCharge() float64 {
	if x != nil {
		return x.StateOfCharge
	}
	return 0
}

func (x *UserDeviceDataResponse) GetChargeLimit() float64 {
	if x != nil {
		return x.ChargeLimit
	}
	return 0
}

func (x *UserDeviceDataResponse) GetRecordUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.RecordUpdatedAt
	}
	return nil
}

func (x *UserDeviceDataResponse) GetRecordCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.RecordCreatedAt
	}
	return nil
}

func (x *UserDeviceDataResponse) GetTirePressure() *TirePressureResponse {
	if x != nil {
		return x.TirePressure
	}
	return nil
}

func (x *UserDeviceDataResponse) GetBatteryVoltage() float64 {
	if x != nil {
		return x.BatteryVoltage
	}
	return 0
}

func (x *UserDeviceDataResponse) GetAmbientTemp() float64 {
	if x != nil {
		return x.AmbientTemp
	}
	return 0
}

type TirePressureResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FrontLeft  float64 `protobuf:"fixed64,1,opt,name=front_left,json=frontLeft,proto3" json:"front_left,omitempty"`
	FrontRight float64 `protobuf:"fixed64,2,opt,name=front_right,json=frontRight,proto3" json:"front_right,omitempty"`
	BackLeft   float64 `protobuf:"fixed64,3,opt,name=back_left,json=backLeft,proto3" json:"back_left,omitempty"`
	BackRight  float64 `protobuf:"fixed64,4,opt,name=back_right,json=backRight,proto3" json:"back_right,omitempty"`
	Age        string  `protobuf:"bytes,5,opt,name=age,proto3" json:"age,omitempty"`
	DataAge    string  `protobuf:"bytes,6,opt,name=data_age,json=dataAge,proto3" json:"data_age,omitempty"`
	RequestId  string  `protobuf:"bytes,7,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	UnitSystem string  `protobuf:"bytes,8,opt,name=unit_system,json=unitSystem,proto3" json:"unit_system,omitempty"`
}

func (x *TirePressureResponse) Reset() {
	*x = TirePressureResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_grpc_user_device_data_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TirePressureResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TirePressureResponse) ProtoMessage() {}

func (x *TirePressureResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_grpc_user_device_data_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TirePressureResponse.ProtoReflect.Descriptor instead.
func (*TirePressureResponse) Descriptor() ([]byte, []int) {
	return file_pkg_grpc_user_device_data_proto_rawDescGZIP(), []int{5}
}

func (x *TirePressureResponse) GetFrontLeft() float64 {
	if x != nil {
		return x.FrontLeft
	}
	return 0
}

func (x *TirePressureResponse) GetFrontRight() float64 {
	if x != nil {
		return x.FrontRight
	}
	return 0
}

func (x *TirePressureResponse) GetBackLeft() float64 {
	if x != nil {
		return x.BackLeft
	}
	return 0
}

func (x *TirePressureResponse) GetBackRight() float64 {
	if x != nil {
		return x.BackRight
	}
	return 0
}

func (x *TirePressureResponse) GetAge() string {
	if x != nil {
		return x.Age
	}
	return ""
}

func (x *TirePressureResponse) GetDataAge() string {
	if x != nil {
		return x.DataAge
	}
	return ""
}

func (x *TirePressureResponse) GetRequestId() string {
	if x != nil {
		return x.RequestId
	}
	return ""
}

func (x *TirePressureResponse) GetUnitSystem() string {
	if x != nil {
		return x.UnitSystem
	}
	return ""
}

var File_pkg_grpc_user_device_data_proto protoreflect.FileDescriptor

var file_pkg_grpc_user_device_data_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x04, 0x67, 0x72, 0x70, 0x63, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa4, 0x01, 0x0a, 0x0d, 0x53, 0x69, 0x67,
	0x6e, 0x61, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x25, 0x0a, 0x0e, 0x69, 0x6e,
	0x74, 0x65, 0x67, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0d, 0x69, 0x6e, 0x74, 0x65, 0x67, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49,
	0x64, 0x12, 0x37, 0x0a, 0x09, 0x66, 0x72, 0x6f, 0x6d, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x08, 0x66, 0x72, 0x6f, 0x6d, 0x44, 0x61, 0x74, 0x65, 0x12, 0x33, 0x0a, 0x07, 0x74, 0x6f,
	0x5f, 0x64, 0x61, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x06, 0x74, 0x6f, 0x44, 0x61, 0x74, 0x65, 0x22,
	0x40, 0x0a, 0x0e, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x2e, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x18, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x6c, 0x49, 0x74,
	0x65, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d,
	0x73, 0x22, 0x76, 0x0a, 0x12, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x6c, 0x49, 0x74, 0x65, 0x6d, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x70, 0x65,
	0x72, 0x74, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x70, 0x65,
	0x72, 0x74, 0x79, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x43,
	0x6f, 0x75, 0x6e, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x72, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x97, 0x01, 0x0a, 0x15, 0x55, 0x73,
	0x65, 0x72, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x24, 0x0a, 0x0e, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x64, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x75, 0x73, 0x65,
	0x72, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x12, 0x30, 0x0a, 0x14, 0x64, 0x65, 0x76,
	0x69, 0x63, 0x65, 0x5f, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x44,
	0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x26, 0x0a, 0x0f, 0x64,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x73, 0x74, 0x79, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x53, 0x74, 0x79, 0x6c,
	0x65, 0x49, 0x64, 0x22, 0x86, 0x05, 0x0a, 0x16, 0x55, 0x73, 0x65, 0x72, 0x44, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x63, 0x68, 0x61, 0x72, 0x67, 0x69, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x08, 0x63, 0x68, 0x61, 0x72, 0x67, 0x69, 0x6e, 0x67, 0x12, 0x34, 0x0a, 0x16, 0x66, 0x75,
	0x65, 0x6c, 0x5f, 0x70, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x5f, 0x72, 0x65, 0x6d, 0x61, 0x69,
	0x6e, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x14, 0x66, 0x75, 0x65, 0x6c,
	0x50, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x6d, 0x61, 0x69, 0x6e, 0x69, 0x6e, 0x67,
	0x12, 0x29, 0x0a, 0x10, 0x62, 0x61, 0x74, 0x74, 0x65, 0x72, 0x79, 0x5f, 0x63, 0x61, 0x70, 0x61,
	0x63, 0x69, 0x74, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0f, 0x62, 0x61, 0x74, 0x74,
	0x65, 0x72, 0x79, 0x43, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79, 0x12, 0x1b, 0x0a, 0x09, 0x6f,
	0x69, 0x6c, 0x5f, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08,
	0x6f, 0x69, 0x6c, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x6f, 0x64, 0x6f, 0x6d,
	0x65, 0x74, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x6f, 0x64, 0x6f, 0x6d,
	0x65, 0x74, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x6c, 0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65,
	0x12, 0x1c, 0x0a, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x72,
	0x61, 0x6e, 0x67, 0x65, 0x12, 0x26, 0x0a, 0x0f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x5f, 0x6f, 0x66,
	0x5f, 0x63, 0x68, 0x61, 0x72, 0x67, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0d, 0x73,
	0x74, 0x61, 0x74, 0x65, 0x4f, 0x66, 0x43, 0x68, 0x61, 0x72, 0x67, 0x65, 0x12, 0x21, 0x0a, 0x0c,
	0x63, 0x68, 0x61, 0x72, 0x67, 0x65, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x0a, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x0b, 0x63, 0x68, 0x61, 0x72, 0x67, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x12,
	0x46, 0x0a, 0x11, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x64, 0x5f, 0x61, 0x74, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0f, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x46, 0x0a, 0x11, 0x72, 0x65, 0x63, 0x6f, 0x72,
	0x64, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x0c, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0f,
	0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12,
	0x3f, 0x0a, 0x0d, 0x74, 0x69, 0x72, 0x65, 0x5f, 0x70, 0x72, 0x65, 0x73, 0x73, 0x75, 0x72, 0x65,
	0x18, 0x0d, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x54, 0x69,
	0x72, 0x65, 0x50, 0x72, 0x65, 0x73, 0x73, 0x75, 0x72, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x52, 0x0c, 0x74, 0x69, 0x72, 0x65, 0x50, 0x72, 0x65, 0x73, 0x73, 0x75, 0x72, 0x65,
	0x12, 0x27, 0x0a, 0x0f, 0x62, 0x61, 0x74, 0x74, 0x65, 0x72, 0x79, 0x5f, 0x76, 0x6f, 0x6c, 0x74,
	0x61, 0x67, 0x65, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0e, 0x62, 0x61, 0x74, 0x74, 0x65,
	0x72, 0x79, 0x56, 0x6f, 0x6c, 0x74, 0x61, 0x67, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x61, 0x6d, 0x62,
	0x69, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x65, 0x6d, 0x70, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x0b, 0x61, 0x6d, 0x62, 0x69, 0x65, 0x6e, 0x74, 0x54, 0x65, 0x6d, 0x70, 0x22, 0xff, 0x01, 0x0a,
	0x14, 0x54, 0x69, 0x72, 0x65, 0x50, 0x72, 0x65, 0x73, 0x73, 0x75, 0x72, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x66, 0x72, 0x6f, 0x6e, 0x74, 0x5f, 0x6c,
	0x65, 0x66, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x09, 0x66, 0x72, 0x6f, 0x6e, 0x74,
	0x4c, 0x65, 0x66, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x66, 0x72, 0x6f, 0x6e, 0x74, 0x5f, 0x72, 0x69,
	0x67, 0x68, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0a, 0x66, 0x72, 0x6f, 0x6e, 0x74,
	0x52, 0x69, 0x67, 0x68, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x62, 0x61, 0x63, 0x6b, 0x5f, 0x6c, 0x65,
	0x66, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x62, 0x61, 0x63, 0x6b, 0x4c, 0x65,
	0x66, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x62, 0x61, 0x63, 0x6b, 0x5f, 0x72, 0x69, 0x67, 0x68, 0x74,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x09, 0x62, 0x61, 0x63, 0x6b, 0x52, 0x69, 0x67, 0x68,
	0x74, 0x12, 0x10, 0x0a, 0x03, 0x61, 0x67, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x61, 0x67, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x61, 0x67, 0x65, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x64, 0x61, 0x74, 0x61, 0x41, 0x67, 0x65, 0x12, 0x1d,
	0x0a, 0x0a, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x64, 0x12, 0x1f, 0x0a,
	0x0b, 0x75, 0x6e, 0x69, 0x74, 0x5f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x75, 0x6e, 0x69, 0x74, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x32, 0xa0,
	0x01, 0x0a, 0x15, 0x55, 0x73, 0x65, 0x72, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x44, 0x61, 0x74,
	0x61, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4e, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x55,
	0x73, 0x65, 0x72, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1b, 0x2e,
	0x67, 0x72, 0x70, 0x63, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x44,
	0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x67, 0x72, 0x70,
	0x63, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x37, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x53,
	0x69, 0x67, 0x6e, 0x61, 0x6c, 0x73, 0x12, 0x13, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x53, 0x69,
	0x67, 0x6e, 0x61, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x42, 0x32, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x44, 0x49, 0x4d, 0x4f, 0x2d, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x64, 0x65, 0x76,
	0x69, 0x63, 0x65, 0x2d, 0x64, 0x61, 0x74, 0x61, 0x2d, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x6b, 0x67,
	0x2f, 0x67, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_grpc_user_device_data_proto_rawDescOnce sync.Once
	file_pkg_grpc_user_device_data_proto_rawDescData = file_pkg_grpc_user_device_data_proto_rawDesc
)

func file_pkg_grpc_user_device_data_proto_rawDescGZIP() []byte {
	file_pkg_grpc_user_device_data_proto_rawDescOnce.Do(func() {
		file_pkg_grpc_user_device_data_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_grpc_user_device_data_proto_rawDescData)
	})
	return file_pkg_grpc_user_device_data_proto_rawDescData
}

var file_pkg_grpc_user_device_data_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_pkg_grpc_user_device_data_proto_goTypes = []interface{}{
	(*SignalRequest)(nil),          // 0: grpc.SignalRequest
	(*SignalResponse)(nil),         // 1: grpc.SignalResponse
	(*SignalItemResponse)(nil),     // 2: grpc.SignalItemResponse
	(*UserDeviceDataRequest)(nil),  // 3: grpc.UserDeviceDataRequest
	(*UserDeviceDataResponse)(nil), // 4: grpc.UserDeviceDataResponse
	(*TirePressureResponse)(nil),   // 5: grpc.TirePressureResponse
	(*timestamppb.Timestamp)(nil),  // 6: google.protobuf.Timestamp
}
var file_pkg_grpc_user_device_data_proto_depIdxs = []int32{
	6, // 0: grpc.SignalRequest.from_date:type_name -> google.protobuf.Timestamp
	6, // 1: grpc.SignalRequest.to_date:type_name -> google.protobuf.Timestamp
	2, // 2: grpc.SignalResponse.items:type_name -> grpc.SignalItemResponse
	6, // 3: grpc.UserDeviceDataResponse.record_updated_at:type_name -> google.protobuf.Timestamp
	6, // 4: grpc.UserDeviceDataResponse.record_created_at:type_name -> google.protobuf.Timestamp
	5, // 5: grpc.UserDeviceDataResponse.tire_pressure:type_name -> grpc.TirePressureResponse
	3, // 6: grpc.UserDeviceDataService.GetUserDeviceData:input_type -> grpc.UserDeviceDataRequest
	0, // 7: grpc.UserDeviceDataService.GetSignals:input_type -> grpc.SignalRequest
	4, // 8: grpc.UserDeviceDataService.GetUserDeviceData:output_type -> grpc.UserDeviceDataResponse
	1, // 9: grpc.UserDeviceDataService.GetSignals:output_type -> grpc.SignalResponse
	8, // [8:10] is the sub-list for method output_type
	6, // [6:8] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_pkg_grpc_user_device_data_proto_init() }
func file_pkg_grpc_user_device_data_proto_init() {
	if File_pkg_grpc_user_device_data_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_grpc_user_device_data_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignalRequest); i {
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
		file_pkg_grpc_user_device_data_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignalResponse); i {
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
		file_pkg_grpc_user_device_data_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignalItemResponse); i {
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
		file_pkg_grpc_user_device_data_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserDeviceDataRequest); i {
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
		file_pkg_grpc_user_device_data_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserDeviceDataResponse); i {
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
		file_pkg_grpc_user_device_data_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TirePressureResponse); i {
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
			RawDescriptor: file_pkg_grpc_user_device_data_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_grpc_user_device_data_proto_goTypes,
		DependencyIndexes: file_pkg_grpc_user_device_data_proto_depIdxs,
		MessageInfos:      file_pkg_grpc_user_device_data_proto_msgTypes,
	}.Build()
	File_pkg_grpc_user_device_data_proto = out.File
	file_pkg_grpc_user_device_data_proto_rawDesc = nil
	file_pkg_grpc_user_device_data_proto_goTypes = nil
	file_pkg_grpc_user_device_data_proto_depIdxs = nil
}
