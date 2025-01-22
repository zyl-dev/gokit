package grpcgw

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

func CustomMatcher(key string) (string, bool) {
	switch strings.ToLower(key) {
	case "xjx-modulegid":
		return "modulegid", true
	case "xjx-areagid":
		return "areagid", true
	case "operator":
		return "useremail", true
	case "operatornameascii":
		return "usercnname", true
	case "operatoruid":
		return "userid", true
	case "language":
		return "lang", true
	case "x-request-id":
		return "request-id", true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

func CustomErrorHandler(ctx context.Context, serveMux *runtime.ServeMux, m runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}
}

type CustomMarshaler struct {
	runtime.JSONPb
}

func (c *CustomMarshaler) Marshal(v interface{}) ([]byte, error) {
	b, err := c.JSONPb.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	err = c.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return json.Marshal(map[string]any{
		"code": 200,
		"msg":  "OK",
		"data": m,
	})
}

func withMiddlewares(handler http.Handler, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.Handler {
	if len(middlewares) == 0 {
		return handler
	}
	var hf http.HandlerFunc
	hf = func(writer http.ResponseWriter, request *http.Request) {
		handler.ServeHTTP(writer, request)
	}
	for _, m := range middlewares {
		hf = m(hf)
	}
	return hf
}

func Handler(ctx context.Context, grpcEndpoint string, middlewares []func(http.HandlerFunc) http.HandlerFunc, registerFuncs ...func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)) (http.Handler, error) {
	mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(CustomMatcher), runtime.WithErrorHandler(CustomErrorHandler), runtime.WithMarshalerOption(runtime.MIMEWildcard, &CustomMarshaler{
		JSONPb: runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		},
	}),
	)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	for _, reg := range registerFuncs {
		err := reg(ctx, mux, grpcEndpoint, opts)
		if err != nil {
			return nil, err
		}
	}
	return withMiddlewares(mux, middlewares...), nil
}
