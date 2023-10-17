package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gitaepark/pha/controller/response"
	"github.com/gitaepark/pha/util"
	"github.com/gitaepark/pha/util/jwt"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "user"
)

func AuthMiddleware(config util.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(AuthorizationHeaderKey)

		// authorization header 값 존재 검증
		if len(authorizationHeader) == 0 {
			response.NewErrResponse(ctx, errEmptyAuthorizationHeader)
			return
		}

		fields := strings.Fields(authorizationHeader)
		// authorization header bearer 존재 검증
		if len(fields) < 2 {
			response.NewErrResponse(ctx, errInvalidAuthorizationHeader)
			return
		}

		authorizationType := strings.ToLower(fields[0])
		// authorization header bearer 타입 검증
		if authorizationType != AuthorizationTypeBearer {
			response.NewErrResponse(ctx, errInvalidAuthorizationBearer)
			return
		}

		accessToken := fields[1]
		// 토큰 검증
		payload, err := jwt.VerifyToken(accessToken, config.JWTSecret)
		if err != nil {
			response.NewErrResponse(ctx, errToken(err))
			return
		}

		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}
