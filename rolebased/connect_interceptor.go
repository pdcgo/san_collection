package rolebased

import (
	"context"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pdcgo/schema/services/role_base/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

var RoleHeaderKey = "x-role-identity"

type roleContextKey struct{}

type Claims struct {
	Data []byte `json:"d"`
	jwt.RegisteredClaims
}

type RoleInterceptor connect.Interceptor
type roleInterceptorImpl struct {
	cfg *Config
}

func NewConnectRoleInterceptor(cfg *Config) RoleInterceptor {
	return &roleInterceptorImpl{cfg}
}

// WrapStreamingClient implements [RoleInterceptor].
func (r *roleInterceptorImpl) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

// WrapStreamingHandler implements [RoleInterceptor].
func (r *roleInterceptorImpl) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return next
}

// WrapUnary implements [RoleInterceptor].
func (r *roleInterceptorImpl) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		msg, ok := req.Any().(proto.Message)
		if !ok {
			return nil, errors.New("[role_base] request is not proto message")
		}

		desc := msg.ProtoReflect().Descriptor()
		opts, ok := desc.Options().(*descriptorpb.MessageOptions)
		if !ok || opts == nil {
			// jika options tidak ada sama sekali
			return nil, errors.New("[role_base] request is not have request_policy")
		}

		if !proto.HasExtension(opts, role_base.E_RequestPolicy) {
			return nil, errors.New("[role_base] request is not have extension request_policy")
		}

		ext := proto.GetExtension(opts, role_base.E_RequestPolicy)
		policy, ok := ext.(*role_base.RequestPolicy)

		if !ok || policy == nil {
			return nil, errors.New("[role_base] request_policy schema error or outdated")
		}

		if policy.AllowAll {
			return next(ctx, req)
		}

		identity, err := ExtractRoleFromRequest(req, r.cfg.Secret)

		if err != nil {
			return nil, err
		}

		err = enforcePolicy(policy, identity)
		if err != nil {
			return nil, err
		}

		newCtx := context.WithValue(ctx, roleContextKey{}, identity)
		return next(newCtx, req)
	}
}

func ExtractRoleFromRequest(req connect.AnyRequest, secret string) (*role_base.Identity, error) {
	tokenString := req.Header().Get(RoleHeaderKey)

	var c Claims
	_, err := jwt.ParseWithClaims(tokenString, &c, func(*jwt.Token) (any, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid role_base in %s %s", RoleHeaderKey, err.Error())
	}

	var identity role_base.Identity
	err = proto.Unmarshal(c.Data, &identity)
	if err != nil {
		return nil, fmt.Errorf("invalid role_base in %s: missing identity", RoleHeaderKey)
	}

	if identity.IdentityId == 0 {
		return nil, fmt.Errorf("invalid role_base in %s: missing identity id", RoleHeaderKey)
	}
	return &identity, nil

}

func enforcePolicy(policy *role_base.RequestPolicy, identity *role_base.Identity) error {
	// TODO: implementasi policy
	for _, r := range policy.Roles {
		if r == identity.Role {
			return nil
		}
	}
	return fmt.Errorf("%s role %s not allowed", identity.Username, identity.Role.String())
}
