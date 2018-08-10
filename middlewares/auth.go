package middlewares

import (
	"stonesrv/crypto"
	"stonesrv/env"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

//AuthToken TOKEN 验证
func AuthToken() gin.HandlerFunc {
	return func(context *gin.Context) {
		token, err := request.ParseFromRequest(context.Request, request.AuthorizationHeaderExtractor,
			func(token *jwt.Token) (interface{}, error) {
				//硬编码 设置为2
				return crypto.GetSecrct(2), nil
			})
		if err == nil {
			if token.Valid {
				context.Next()
			} else {
				env.GenJSONResponse(context, env.AuthFailed, nil)
				//context.String(http.StatusUnauthorized, "非法的 Token")
			}
		} else {
			env.GenJSONResponse(context, env.AuthFailed, nil)
			//context.String(http.StatusUnauthorized, "该接口未授权")
		}
		context.Abort()
	}
}
