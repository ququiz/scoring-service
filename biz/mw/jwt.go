package jwt

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	gojwt "github.com/golang-jwt/jwt/v4"
	"github.com/hertz-contrib/jwt"
	"go.uber.org/zap"
)

var (
	// JwtMiddleware       *jwt.HertzJWTMiddleware
	IdentityKey         = "sub"
	PublicKeyAuthServer = "-----BEGIN PUBLIC KEY-----\nMIGbMBAGByqGSM49AgEGBSuBBAAjA4GGAAQAodxwFdiFKWTG/ZU7vXPdk8ox+nNU\n1JmxsmI8i8tYrYf6QxmwBz13jS/PZsb8dJbMFY3YTMMih6SKz7e+cQ68IbgA7BnY\n5fYFQET4SNHVX/zaH6J70ERJLsRrarmWSXsNbMbnqXlIkoorYXeAn9vsLbr/RPw9\nDYaoq4JrQ+OGsc4LHMw=\n-----END PUBLIC KEY-----\n"
)

type User struct {
	ID       string `json:"userId"`
	Username string `json:"username"`
}

func GetJwtMiddleware() *jwt.HertzJWTMiddleware {

	publicKeyBlock, _ := pem.Decode([]byte(PublicKeyAuthServer))
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	ECDSAPubKey := publicKey.(*ecdsa.PublicKey)
	JwtMiddleware, err := jwt.New(&jwt.HertzJWTMiddleware{
		Realm:       "dogker digital signature public key auth",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: IdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					IdentityKey: v.ID,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			return &User{
				ID: claims[IdentityKey].(string),
			}
		},
		Authorizator: func(data interface{}, ctx context.Context, c *app.RequestContext) bool {
			if v, ok := data.(*User); ok {
				c.Set("userID", v.ID)
				c.Set("username", v.Username)
				return true
			}

			return false
		},
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			c.JSON(code, map[string]interface{}{
				"code":    code,
				"message": message,
			})
		},

		TokenLookup: "header: Authorization, query: token, cookie: jwt",

		// TokenHeadName is a string in the header. Default value is "Bearer". If you want empty value, use WithoutDefaultTokenHeadName.
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
		KeyFunc: func(t *gojwt.Token) (interface{}, error) {
			// if gojwt.GetSigningMethod(gojwt.SigningMethodECDSA) != t.Method {
			// 	return nil, jwt.ErrInvalidSigningAlgorithm
			// }
			if _, ok := t.Method.(*gojwt.SigningMethodECDSA); !ok {
				return nil, jwt.ErrInvalidSigningAlgorithm
			}

			return ECDSAPubKey, nil
		},
	})
	if err != nil {
		zap.L().Fatal("JWT Error:"+err.Error(), zap.Error(err))
	}
	return JwtMiddleware
}
