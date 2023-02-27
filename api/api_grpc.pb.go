// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: api.proto

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// VMWareDesktopAutoscalerServiceClient is the client API for VMWareDesktopAutoscalerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VMWareDesktopAutoscalerServiceClient interface {
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error)
	Delete(ctx context.Context, in *VirtualMachineRequest, opts ...grpc.CallOption) (*DeleteResponse, error)
	PowerOn(ctx context.Context, in *VirtualMachineRequest, opts ...grpc.CallOption) (*PowerOnResponse, error)
	PowerOff(ctx context.Context, in *VirtualMachineRequest, opts ...grpc.CallOption) (*PowerOffResponse, error)
	ShutdownGuest(ctx context.Context, in *VirtualMachineRequest, opts ...grpc.CallOption) (*ShutdownGuestResponse, error)
	Status(ctx context.Context, in *VirtualMachineRequest, opts ...grpc.CallOption) (*StatusResponse, error)
	WaitForIP(ctx context.Context, in *VirtualMachineRequest, opts ...grpc.CallOption) (*WaitForIPResponse, error)
	WaitForToolsRunning(ctx context.Context, in *VirtualMachineRequest, opts ...grpc.CallOption) (*WaitForToolsRunningResponse, error)
	SetAutoStart(ctx context.Context, in *AutoStartRequest, opts ...grpc.CallOption) (*AutoStartResponse, error)
	VirtualMachineByName(ctx context.Context, in *VirtualMachineRequest, opts ...grpc.CallOption) (*VirtualMachineByNameResponse, error)
	VirtualMachineByUUID(ctx context.Context, in *VirtualMachineRequest, opts ...grpc.CallOption) (*VirtualMachineByNameResponse, error)
	ListVirtualMachines(ctx context.Context, in *VirtualMachinesRequest, opts ...grpc.CallOption) (*VirtualMachinesResponse, error)
}

type vMWareDesktopAutoscalerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewVMWareDesktopAutoscalerServiceClient(cc grpc.ClientConnInterface) VMWareDesktopAutoscalerServiceClient {
	return &vMWareDesktopAutoscalerServiceClient{cc}
}

func (c *vMWareDesktopAutoscalerServiceClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error) {
	out := new(CreateResponse)
	err := c.cc.Invoke(ctx, "/api.VMWareDesktopAutoscalerService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vMWareDesktopAutoscalerServiceClient) Delete(ctx context.Context, in *VirtualMachineRequest, opts ...grpc.CallOption) (*DeleteResponse, error) {
	out := new(DeleteResponse)
	err := c.cc.Invoke(ctx, "/api.VMWareDesktopAutoscalerService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vMWareDesktopAutoscalerServiceClient) PowerOn(ctx context.Context, in *VirtualMachineRequest, opts ...grpc.CallOption) (*PowerOnResponse, error) {
	out := new(PowerOnResponse)
	err := c.cc.Invoke(ctx, "/api.VMWareDesktopAutoscalerService/PowerOn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vMWareDesktopAutoscalerServiceClient) PowerOff(ctx context.Context, in *VirtualMachineRequest, opts ...grpc.CallOption) (*PowerOffResponse, error) {
	out := new(PowerOffResponse)
	err := c.cc.Invoke(ctx, "/api.VMWareDesktopAutoscalerService/PowerOff", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vMWareDesktopAutoscalerServiceClient) ShutdownGuest(ctx context.Context, in *VirtualMachineRequest, opts ...grpc.CallOption) (*ShutdownGuestResponse, error) {
	out := new(ShutdownGuestResponse)
	err := c.cc.Invoke(ctx, "/api.VMWareDesktopAutoscalerService/ShutdownGuest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vMWareDesktopAutoscalerServiceClient) Status(ctx context.Context, in *VirtualMachineRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, "/api.VMWareDesktopAutoscalerService/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vMWareDesktopAutoscalerServiceClient) WaitForIP(ctx context.Context, in *VirtualMachineRequest, opts ...grpc.CallOption) (*WaitForIPResponse, error) {
	out := new(WaitForIPResponse)
	err := c.cc.Invoke(ctx, "/api.VMWareDesktopAutoscalerService/WaitForIP", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vMWareDesktopAutoscalerServiceClient) WaitForToolsRunning(ctx context.Context, in *VirtualMachineRequest, opts ...grpc.CallOption) (*WaitForToolsRunningResponse, error) {
	out := new(WaitForToolsRunningResponse)
	err := c.cc.Invoke(ctx, "/api.VMWareDesktopAutoscalerService/WaitForToolsRunning", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vMWareDesktopAutoscalerServiceClient) SetAutoStart(ctx context.Context, in *AutoStartRequest, opts ...grpc.CallOption) (*AutoStartResponse, error) {
	out := new(AutoStartResponse)
	err := c.cc.Invoke(ctx, "/api.VMWareDesktopAutoscalerService/SetAutoStart", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vMWareDesktopAutoscalerServiceClient) VirtualMachineByName(ctx context.Context, in *VirtualMachineRequest, opts ...grpc.CallOption) (*VirtualMachineByNameResponse, error) {
	out := new(VirtualMachineByNameResponse)
	err := c.cc.Invoke(ctx, "/api.VMWareDesktopAutoscalerService/VirtualMachineByName", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vMWareDesktopAutoscalerServiceClient) VirtualMachineByUUID(ctx context.Context, in *VirtualMachineRequest, opts ...grpc.CallOption) (*VirtualMachineByNameResponse, error) {
	out := new(VirtualMachineByNameResponse)
	err := c.cc.Invoke(ctx, "/api.VMWareDesktopAutoscalerService/VirtualMachineByUUID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vMWareDesktopAutoscalerServiceClient) ListVirtualMachines(ctx context.Context, in *VirtualMachinesRequest, opts ...grpc.CallOption) (*VirtualMachinesResponse, error) {
	out := new(VirtualMachinesResponse)
	err := c.cc.Invoke(ctx, "/api.VMWareDesktopAutoscalerService/ListVirtualMachines", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VMWareDesktopAutoscalerServiceServer is the server API for VMWareDesktopAutoscalerService service.
// All implementations must embed UnimplementedVMWareDesktopAutoscalerServiceServer
// for forward compatibility
type VMWareDesktopAutoscalerServiceServer interface {
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	Delete(context.Context, *VirtualMachineRequest) (*DeleteResponse, error)
	PowerOn(context.Context, *VirtualMachineRequest) (*PowerOnResponse, error)
	PowerOff(context.Context, *VirtualMachineRequest) (*PowerOffResponse, error)
	ShutdownGuest(context.Context, *VirtualMachineRequest) (*ShutdownGuestResponse, error)
	Status(context.Context, *VirtualMachineRequest) (*StatusResponse, error)
	WaitForIP(context.Context, *VirtualMachineRequest) (*WaitForIPResponse, error)
	WaitForToolsRunning(context.Context, *VirtualMachineRequest) (*WaitForToolsRunningResponse, error)
	SetAutoStart(context.Context, *AutoStartRequest) (*AutoStartResponse, error)
	VirtualMachineByName(context.Context, *VirtualMachineRequest) (*VirtualMachineByNameResponse, error)
	VirtualMachineByUUID(context.Context, *VirtualMachineRequest) (*VirtualMachineByNameResponse, error)
	ListVirtualMachines(context.Context, *VirtualMachinesRequest) (*VirtualMachinesResponse, error)
	mustEmbedUnimplementedVMWareDesktopAutoscalerServiceServer()
}

// UnimplementedVMWareDesktopAutoscalerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedVMWareDesktopAutoscalerServiceServer struct {
}

func (UnimplementedVMWareDesktopAutoscalerServiceServer) Create(context.Context, *CreateRequest) (*CreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedVMWareDesktopAutoscalerServiceServer) Delete(context.Context, *VirtualMachineRequest) (*DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedVMWareDesktopAutoscalerServiceServer) PowerOn(context.Context, *VirtualMachineRequest) (*PowerOnResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PowerOn not implemented")
}
func (UnimplementedVMWareDesktopAutoscalerServiceServer) PowerOff(context.Context, *VirtualMachineRequest) (*PowerOffResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PowerOff not implemented")
}
func (UnimplementedVMWareDesktopAutoscalerServiceServer) ShutdownGuest(context.Context, *VirtualMachineRequest) (*ShutdownGuestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShutdownGuest not implemented")
}
func (UnimplementedVMWareDesktopAutoscalerServiceServer) Status(context.Context, *VirtualMachineRequest) (*StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedVMWareDesktopAutoscalerServiceServer) WaitForIP(context.Context, *VirtualMachineRequest) (*WaitForIPResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WaitForIP not implemented")
}
func (UnimplementedVMWareDesktopAutoscalerServiceServer) WaitForToolsRunning(context.Context, *VirtualMachineRequest) (*WaitForToolsRunningResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WaitForToolsRunning not implemented")
}
func (UnimplementedVMWareDesktopAutoscalerServiceServer) SetAutoStart(context.Context, *AutoStartRequest) (*AutoStartResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetAutoStart not implemented")
}
func (UnimplementedVMWareDesktopAutoscalerServiceServer) VirtualMachineByName(context.Context, *VirtualMachineRequest) (*VirtualMachineByNameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VirtualMachineByName not implemented")
}
func (UnimplementedVMWareDesktopAutoscalerServiceServer) VirtualMachineByUUID(context.Context, *VirtualMachineRequest) (*VirtualMachineByNameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VirtualMachineByUUID not implemented")
}
func (UnimplementedVMWareDesktopAutoscalerServiceServer) ListVirtualMachines(context.Context, *VirtualMachinesRequest) (*VirtualMachinesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListVirtualMachines not implemented")
}
func (UnimplementedVMWareDesktopAutoscalerServiceServer) mustEmbedUnimplementedVMWareDesktopAutoscalerServiceServer() {
}

// UnsafeVMWareDesktopAutoscalerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VMWareDesktopAutoscalerServiceServer will
// result in compilation errors.
type UnsafeVMWareDesktopAutoscalerServiceServer interface {
	mustEmbedUnimplementedVMWareDesktopAutoscalerServiceServer()
}

func RegisterVMWareDesktopAutoscalerServiceServer(s grpc.ServiceRegistrar, srv VMWareDesktopAutoscalerServiceServer) {
	s.RegisterService(&VMWareDesktopAutoscalerService_ServiceDesc, srv)
}

func _VMWareDesktopAutoscalerService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VMWareDesktopAutoscalerServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.VMWareDesktopAutoscalerService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VMWareDesktopAutoscalerServiceServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VMWareDesktopAutoscalerService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VirtualMachineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VMWareDesktopAutoscalerServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.VMWareDesktopAutoscalerService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VMWareDesktopAutoscalerServiceServer).Delete(ctx, req.(*VirtualMachineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VMWareDesktopAutoscalerService_PowerOn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VirtualMachineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VMWareDesktopAutoscalerServiceServer).PowerOn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.VMWareDesktopAutoscalerService/PowerOn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VMWareDesktopAutoscalerServiceServer).PowerOn(ctx, req.(*VirtualMachineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VMWareDesktopAutoscalerService_PowerOff_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VirtualMachineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VMWareDesktopAutoscalerServiceServer).PowerOff(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.VMWareDesktopAutoscalerService/PowerOff",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VMWareDesktopAutoscalerServiceServer).PowerOff(ctx, req.(*VirtualMachineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VMWareDesktopAutoscalerService_ShutdownGuest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VirtualMachineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VMWareDesktopAutoscalerServiceServer).ShutdownGuest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.VMWareDesktopAutoscalerService/ShutdownGuest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VMWareDesktopAutoscalerServiceServer).ShutdownGuest(ctx, req.(*VirtualMachineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VMWareDesktopAutoscalerService_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VirtualMachineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VMWareDesktopAutoscalerServiceServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.VMWareDesktopAutoscalerService/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VMWareDesktopAutoscalerServiceServer).Status(ctx, req.(*VirtualMachineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VMWareDesktopAutoscalerService_WaitForIP_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VirtualMachineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VMWareDesktopAutoscalerServiceServer).WaitForIP(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.VMWareDesktopAutoscalerService/WaitForIP",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VMWareDesktopAutoscalerServiceServer).WaitForIP(ctx, req.(*VirtualMachineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VMWareDesktopAutoscalerService_WaitForToolsRunning_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VirtualMachineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VMWareDesktopAutoscalerServiceServer).WaitForToolsRunning(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.VMWareDesktopAutoscalerService/WaitForToolsRunning",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VMWareDesktopAutoscalerServiceServer).WaitForToolsRunning(ctx, req.(*VirtualMachineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VMWareDesktopAutoscalerService_SetAutoStart_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AutoStartRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VMWareDesktopAutoscalerServiceServer).SetAutoStart(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.VMWareDesktopAutoscalerService/SetAutoStart",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VMWareDesktopAutoscalerServiceServer).SetAutoStart(ctx, req.(*AutoStartRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VMWareDesktopAutoscalerService_VirtualMachineByName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VirtualMachineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VMWareDesktopAutoscalerServiceServer).VirtualMachineByName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.VMWareDesktopAutoscalerService/VirtualMachineByName",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VMWareDesktopAutoscalerServiceServer).VirtualMachineByName(ctx, req.(*VirtualMachineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VMWareDesktopAutoscalerService_VirtualMachineByUUID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VirtualMachineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VMWareDesktopAutoscalerServiceServer).VirtualMachineByUUID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.VMWareDesktopAutoscalerService/VirtualMachineByUUID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VMWareDesktopAutoscalerServiceServer).VirtualMachineByUUID(ctx, req.(*VirtualMachineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VMWareDesktopAutoscalerService_ListVirtualMachines_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VirtualMachinesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VMWareDesktopAutoscalerServiceServer).ListVirtualMachines(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.VMWareDesktopAutoscalerService/ListVirtualMachines",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VMWareDesktopAutoscalerServiceServer).ListVirtualMachines(ctx, req.(*VirtualMachinesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// VMWareDesktopAutoscalerService_ServiceDesc is the grpc.ServiceDesc for VMWareDesktopAutoscalerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var VMWareDesktopAutoscalerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.VMWareDesktopAutoscalerService",
	HandlerType: (*VMWareDesktopAutoscalerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _VMWareDesktopAutoscalerService_Create_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _VMWareDesktopAutoscalerService_Delete_Handler,
		},
		{
			MethodName: "PowerOn",
			Handler:    _VMWareDesktopAutoscalerService_PowerOn_Handler,
		},
		{
			MethodName: "PowerOff",
			Handler:    _VMWareDesktopAutoscalerService_PowerOff_Handler,
		},
		{
			MethodName: "ShutdownGuest",
			Handler:    _VMWareDesktopAutoscalerService_ShutdownGuest_Handler,
		},
		{
			MethodName: "Status",
			Handler:    _VMWareDesktopAutoscalerService_Status_Handler,
		},
		{
			MethodName: "WaitForIP",
			Handler:    _VMWareDesktopAutoscalerService_WaitForIP_Handler,
		},
		{
			MethodName: "WaitForToolsRunning",
			Handler:    _VMWareDesktopAutoscalerService_WaitForToolsRunning_Handler,
		},
		{
			MethodName: "SetAutoStart",
			Handler:    _VMWareDesktopAutoscalerService_SetAutoStart_Handler,
		},
		{
			MethodName: "VirtualMachineByName",
			Handler:    _VMWareDesktopAutoscalerService_VirtualMachineByName_Handler,
		},
		{
			MethodName: "VirtualMachineByUUID",
			Handler:    _VMWareDesktopAutoscalerService_VirtualMachineByUUID_Handler,
		},
		{
			MethodName: "ListVirtualMachines",
			Handler:    _VMWareDesktopAutoscalerService_ListVirtualMachines_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}
