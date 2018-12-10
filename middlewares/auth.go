package middlewares

import (
	"stonesrv/crypto"
	"stonesrv/env"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

//AuthToken login TOKEN 验证
func AuthToken() gin.HandlerFunc {
	return func(context *gin.Context) {
		token, err := request.ParseFromRequest(context.Request, request.AuthorizationHeaderExtractor,
			func(token *jwt.Token) (interface{}, error) {
				//硬编码 设置为2
				return crypto.GetSecrct(crypto.LoginSalt), nil
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

//PackToken pack TOKEN 验证
func PackToken() gin.HandlerFunc {
	return func(context *gin.Context) {
		token, err := request.ParseFromRequest(context.Request, request.OAuth2Extractor,
			func(token *jwt.Token) (interface{}, error) {
				//硬编码 设置为3
				return crypto.GetSecrct(crypto.PackSalt), nil
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

//SpaceToken SPACE TOKEN 验证
func SpaceToken() gin.HandlerFunc {
	return func(context *gin.Context) {
		token, err := request.ParseFromRequest(context.Request, request.OAuth2Extractor,
			func(token *jwt.Token) (interface{}, error) {
				//硬编码 设置为4
				return crypto.GetSecrct(crypto.SpaceSalt), nil
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

//TableToken TABLE TOKEN 验证
func TableToken() gin.HandlerFunc {
	return func(context *gin.Context) {
		token, err := request.ParseFromRequest(context.Request, request.OAuth2Extractor,
			func(token *jwt.Token) (interface{}, error) {
				//硬编码 设置为5
				return crypto.GetSecrct(crypto.TableSalt), nil
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

//SpaceTableToken SPACE AND TABLE TOKEN 验证
func SpaceTableToken() gin.HandlerFunc {
	return func(context *gin.Context) {
		token, err := request.ParseFromRequest(context.Request, request.OAuth2Extractor,
			func(token *jwt.Token) (interface{}, error) {
				//硬编码 设置为6
				return crypto.GetSecrct(crypto.SpaceAndTableSalt), nil
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
