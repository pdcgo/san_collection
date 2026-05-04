package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/pdcgo/san_collection/apis"
	"github.com/pdcgo/san_collection/rolebased"
	"github.com/pdcgo/schema/services/example_iface/v1"
	"github.com/pdcgo/schema/services/example_iface/v1/example_ifaceconnect"
	"github.com/pdcgo/shared/custom_connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type HelloService struct {
	cfg *rolebased.Config
}

// CreateToken implements [example_ifaceconnect.ExampleHelloServiceHandler].
func (h *HelloService) CreateToken(ctx context.Context, req *connect.Request[example_iface.CreateTokenRequest]) (*connect.Response[example_iface.CreateTokenResponse], error) {
	tokenString, err := rolebased.GenerateToken(req.Msg.Identity, h.cfg.Secret, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&example_iface.CreateTokenResponse{
		Token: tokenString,
	}), nil
}

// HelloAdmin implements [example_ifaceconnect.ExampleHelloServiceHandler].
func (h *HelloService) HelloAdmin(ctx context.Context, req *connect.Request[example_iface.HelloAdminRequest]) (*connect.Response[example_iface.HelloAdminResponse], error) {
	return connect.NewResponse(&example_iface.HelloAdminResponse{
		Message: "Hello " + req.Msg.Name,
	}), nil
}

// HelloWorld implements [example_ifaceconnect.ExampleHelloServiceHandler].
func (h *HelloService) HelloWorld(ctx context.Context, req *connect.Request[example_iface.HelloWorldRequest]) (*connect.Response[example_iface.HelloWorldResponse], error) {
	return connect.NewResponse(&example_iface.HelloWorldResponse{
		Message: req.Msg.Name,
	}), nil
}

func main() {
	server := http.NewServeMux()
	registerReflect := apis.NewRegisterReflect(server)

	roleCfg := &rolebased.Config{
		Secret: "hello_secret",
	}

	iterceptor := connect.WithInterceptors(
		rolebased.NewConnectRoleInterceptor(roleCfg),
	)

	path, handler := example_ifaceconnect.NewExampleHelloServiceHandler(&HelloService{
		cfg: roleCfg,
	}, iterceptor)
	server.Handle(path, handler)

	registerReflect([]string{
		example_ifaceconnect.ExampleHelloServiceName,
	})
	slog.Info("running example",
		"host", "localhost:8080",
	)
	http.ListenAndServe(
		"localhost:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(
			custom_connect.WithCORS(server),
			&http2.Server{}),
	)

}
