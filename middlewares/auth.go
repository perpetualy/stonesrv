package middlewares

import (
    "stonesrv/crypto"
	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"net/http"
)

//AuthToken TOKEN 验证
func AuthToken() gin.HandlerFunc {
    return func(c *gin.Context) {
        token, err := request.ParseFromRequest(c.Request, request.AuthorizationHeaderExtractor,
            func(token *jwt.Token) (interface{}, error) {
                //硬编码 设置为2
                return crypto.GetSecrct(2), nil
            })
        if err == nil {
            if token.Valid {
                c.Next()
            } else {
                c.String(http.StatusUnauthorized, "非法的 Token")
            }
        } else {
            c.String(http.StatusUnauthorized, "该接口未授权")
        }
        c.Abort()
    }
}