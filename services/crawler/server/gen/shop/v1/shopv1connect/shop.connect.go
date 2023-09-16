// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: shop/v1/shop.proto

package shopv1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	v1 "crawler/server/gen/shop/v1"
	errors "errors"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion0_1_0

const (
	// ShopServiceName is the fully-qualified name of the ShopService service.
	ShopServiceName = "shop.v1.ShopService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// ShopServiceShopListProcedure is the fully-qualified name of the ShopService's ShopList RPC.
	ShopServiceShopListProcedure = "/shop.v1.ShopService/ShopList"
)

// ShopServiceClient is a client for the shop.v1.ShopService service.
type ShopServiceClient interface {
	ShopList(context.Context, *connect.Request[v1.ShopListRequest]) (*connect.Response[v1.ShopListResponse], error)
}

// NewShopServiceClient constructs a client for the shop.v1.ShopService service. By default, it uses
// the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewShopServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) ShopServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &shopServiceClient{
		shopList: connect.NewClient[v1.ShopListRequest, v1.ShopListResponse](
			httpClient,
			baseURL+ShopServiceShopListProcedure,
			opts...,
		),
	}
}

// shopServiceClient implements ShopServiceClient.
type shopServiceClient struct {
	shopList *connect.Client[v1.ShopListRequest, v1.ShopListResponse]
}

// ShopList calls shop.v1.ShopService.ShopList.
func (c *shopServiceClient) ShopList(ctx context.Context, req *connect.Request[v1.ShopListRequest]) (*connect.Response[v1.ShopListResponse], error) {
	return c.shopList.CallUnary(ctx, req)
}

// ShopServiceHandler is an implementation of the shop.v1.ShopService service.
type ShopServiceHandler interface {
	ShopList(context.Context, *connect.Request[v1.ShopListRequest]) (*connect.Response[v1.ShopListResponse], error)
}

// NewShopServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewShopServiceHandler(svc ShopServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	shopServiceShopListHandler := connect.NewUnaryHandler(
		ShopServiceShopListProcedure,
		svc.ShopList,
		opts...,
	)
	return "/shop.v1.ShopService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case ShopServiceShopListProcedure:
			shopServiceShopListHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedShopServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedShopServiceHandler struct{}

func (UnimplementedShopServiceHandler) ShopList(context.Context, *connect.Request[v1.ShopListRequest]) (*connect.Response[v1.ShopListResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("shop.v1.ShopService.ShopList is not implemented"))
}
