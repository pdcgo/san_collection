package rolebased

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pdcgo/schema/services/role_base/v1"
	"google.golang.org/protobuf/proto"
)

func GenerateToken(identity *role_base.Identity, secret string, expiry time.Duration) (string, error) {

	rawIdentity, err := proto.Marshal(identity)
	if err != nil {
		return "", err
	}

	claims := &Claims{
		Data: rawIdentity,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
