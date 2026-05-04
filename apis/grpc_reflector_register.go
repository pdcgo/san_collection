package apis

import (
	"net/http"

	"connectrpc.com/grpcreflect"
)

type RegisterReflectFunc func(grpcReflectNames []string)

func NewRegisterReflect(
	mux *http.ServeMux,
) RegisterReflectFunc {
	return func(grpcReflectNames []string) {
		reflector := grpcreflect.NewStaticReflector(
			grpcReflectNames...,
		// protoc-gen-connect-go generates package-level constants
		// for these fully-qualified protobuf service names, so you'd more likely
		// reference userv1.UserServiceName and groupv1.GroupServiceName.
		)
		mux.Handle(grpcreflect.NewHandlerV1(reflector))
		// Many tools still expect the older version of the server reflection API, so
		// most servers should mount both handlers.
		mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
	}
}
