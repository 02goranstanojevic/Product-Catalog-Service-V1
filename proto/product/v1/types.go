package productv1

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ = timestamppb.Now

type ProductServiceServer interface {
	CreateProduct(ctx context.Context, req *CreateProductRequest) (*CreateProductReply, error)
	UpdateProduct(ctx context.Context, req *UpdateProductRequest) (*UpdateProductReply, error)
	ActivateProduct(ctx context.Context, req *ActivateProductRequest) (*ActivateProductReply, error)
	DeactivateProduct(ctx context.Context, req *DeactivateProductRequest) (*DeactivateProductReply, error)
	ApplyDiscount(ctx context.Context, req *ApplyDiscountRequest) (*ApplyDiscountReply, error)
	RemoveDiscount(ctx context.Context, req *RemoveDiscountRequest) (*RemoveDiscountReply, error)
	GetProduct(ctx context.Context, req *GetProductRequest) (*GetProductReply, error)
	ListProducts(ctx context.Context, req *ListProductsRequest) (*ListProductsReply, error)
}

type UnimplementedProductServiceServer struct{}

func (UnimplementedProductServiceServer) CreateProduct(context.Context, *CreateProductRequest) (*CreateProductReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
func (UnimplementedProductServiceServer) UpdateProduct(context.Context, *UpdateProductRequest) (*UpdateProductReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
func (UnimplementedProductServiceServer) ActivateProduct(context.Context, *ActivateProductRequest) (*ActivateProductReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
func (UnimplementedProductServiceServer) DeactivateProduct(context.Context, *DeactivateProductRequest) (*DeactivateProductReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
func (UnimplementedProductServiceServer) ApplyDiscount(context.Context, *ApplyDiscountRequest) (*ApplyDiscountReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
func (UnimplementedProductServiceServer) RemoveDiscount(context.Context, *RemoveDiscountRequest) (*RemoveDiscountReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
func (UnimplementedProductServiceServer) GetProduct(context.Context, *GetProductRequest) (*GetProductReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
func (UnimplementedProductServiceServer) ListProducts(context.Context, *ListProductsRequest) (*ListProductsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

func RegisterProductServiceServer(s *grpc.Server, srv ProductServiceServer) {
	s.RegisterService(&_ProductService_serviceDesc, srv)
}

func _ProductService_CreateProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateProductRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServiceServer).CreateProduct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/product.v1.ProductService/CreateProduct",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServiceServer).CreateProduct(ctx, req.(*CreateProductRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductService_UpdateProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateProductRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServiceServer).UpdateProduct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/product.v1.ProductService/UpdateProduct",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServiceServer).UpdateProduct(ctx, req.(*UpdateProductRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductService_ActivateProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ActivateProductRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServiceServer).ActivateProduct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/product.v1.ProductService/ActivateProduct",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServiceServer).ActivateProduct(ctx, req.(*ActivateProductRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductService_DeactivateProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeactivateProductRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServiceServer).DeactivateProduct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/product.v1.ProductService/DeactivateProduct",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServiceServer).DeactivateProduct(ctx, req.(*DeactivateProductRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductService_ApplyDiscount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ApplyDiscountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServiceServer).ApplyDiscount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/product.v1.ProductService/ApplyDiscount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServiceServer).ApplyDiscount(ctx, req.(*ApplyDiscountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductService_RemoveDiscount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveDiscountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServiceServer).RemoveDiscount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/product.v1.ProductService/RemoveDiscount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServiceServer).RemoveDiscount(ctx, req.(*RemoveDiscountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductService_GetProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProductRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServiceServer).GetProduct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/product.v1.ProductService/GetProduct",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServiceServer).GetProduct(ctx, req.(*GetProductRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductService_ListProducts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListProductsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServiceServer).ListProducts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/product.v1.ProductService/ListProducts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServiceServer).ListProducts(ctx, req.(*ListProductsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ProductService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "product.v1.ProductService",
	HandlerType: (*ProductServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "CreateProduct", Handler: _ProductService_CreateProduct_Handler},
		{MethodName: "UpdateProduct", Handler: _ProductService_UpdateProduct_Handler},
		{MethodName: "ActivateProduct", Handler: _ProductService_ActivateProduct_Handler},
		{MethodName: "DeactivateProduct", Handler: _ProductService_DeactivateProduct_Handler},
		{MethodName: "ApplyDiscount", Handler: _ProductService_ApplyDiscount_Handler},
		{MethodName: "RemoveDiscount", Handler: _ProductService_RemoveDiscount_Handler},
		{MethodName: "GetProduct", Handler: _ProductService_GetProduct_Handler},
		{MethodName: "ListProducts", Handler: _ProductService_ListProducts_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/product/v1/product_service.proto",
}

type CreateProductRequest struct {
	Name                 string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description          string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Category             string `protobuf:"bytes,3,opt,name=category,proto3" json:"category,omitempty"`
	BasePriceNumerator   int64  `protobuf:"varint,4,opt,name=base_price_numerator,proto3" json:"base_price_numerator,omitempty"`
	BasePriceDenominator int64  `protobuf:"varint,5,opt,name=base_price_denominator,proto3" json:"base_price_denominator,omitempty"`
}

type CreateProductReply struct {
	ProductId string `protobuf:"bytes,1,opt,name=product_id,proto3" json:"product_id,omitempty"`
}

type UpdateProductRequest struct {
	ProductId   string `protobuf:"bytes,1,opt,name=product_id,proto3" json:"product_id,omitempty"`
	Name        string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Category    string `protobuf:"bytes,4,opt,name=category,proto3" json:"category,omitempty"`
}

type UpdateProductReply struct{}

type ActivateProductRequest struct {
	ProductId string `protobuf:"bytes,1,opt,name=product_id,proto3" json:"product_id,omitempty"`
}

type ActivateProductReply struct{}

type DeactivateProductRequest struct {
	ProductId string `protobuf:"bytes,1,opt,name=product_id,proto3" json:"product_id,omitempty"`
}

type DeactivateProductReply struct{}

type ApplyDiscountRequest struct {
	ProductId  string                 `protobuf:"bytes,1,opt,name=product_id,proto3" json:"product_id,omitempty"`
	Percentage float64                `protobuf:"fixed64,2,opt,name=percentage,proto3" json:"percentage,omitempty"`
	StartDate  *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=start_date,proto3" json:"start_date,omitempty"`
	EndDate    *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=end_date,proto3" json:"end_date,omitempty"`
}

type ApplyDiscountReply struct{}

type RemoveDiscountRequest struct {
	ProductId string `protobuf:"bytes,1,opt,name=product_id,proto3" json:"product_id,omitempty"`
}

type RemoveDiscountReply struct{}

type GetProductRequest struct {
	ProductId string `protobuf:"bytes,1,opt,name=product_id,proto3" json:"product_id,omitempty"`
}

type GetProductReply struct {
	Product *ProductView `protobuf:"bytes,1,opt,name=product,proto3" json:"product,omitempty"`
}

type ListProductsRequest struct {
	Category   string `protobuf:"bytes,1,opt,name=category,proto3" json:"category,omitempty"`
	PageSize   int32  `protobuf:"varint,2,opt,name=page_size,proto3" json:"page_size,omitempty"`
	PageOffset int32  `protobuf:"varint,3,opt,name=page_offset,proto3" json:"page_offset,omitempty"`
}

type ListProductsReply struct {
	Products   []*ProductView `protobuf:"bytes,1,rep,name=products,proto3" json:"products,omitempty"`
	TotalCount int32          `protobuf:"varint,2,opt,name=total_count,proto3" json:"total_count,omitempty"`
}

type ProductView struct {
	Id              string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name            string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description     string   `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Category        string   `protobuf:"bytes,4,opt,name=category,proto3" json:"category,omitempty"`
	BasePrice       string   `protobuf:"bytes,5,opt,name=base_price,proto3" json:"base_price,omitempty"`
	EffectivePrice  string   `protobuf:"bytes,6,opt,name=effective_price,proto3" json:"effective_price,omitempty"`
	DiscountPercent *float64 `protobuf:"fixed64,7,opt,name=discount_percent,proto3,oneof" json:"discount_percent,omitempty"`
	Status          string   `protobuf:"bytes,8,opt,name=status,proto3" json:"status,omitempty"`
	CreatedAt       string   `protobuf:"bytes,9,opt,name=created_at,proto3" json:"created_at,omitempty"`
	UpdatedAt       string   `protobuf:"bytes,10,opt,name=updated_at,proto3" json:"updated_at,omitempty"`
}
