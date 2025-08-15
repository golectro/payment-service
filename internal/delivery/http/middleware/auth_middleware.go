package middleware

import (
	"encoding/json"
	"golectro-payment/internal/constants"
	"golectro-payment/internal/model"
	"golectro-payment/internal/utils"
	"net/http"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

func NewAuth(viper *viper.Viper) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		const bearerPrefix = "Bearer "

		if len(authHeader) <= len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
			res := utils.FailedResponse(ctx, http.StatusUnauthorized, constants.InvalidToken, nil)
			ctx.AbortWithStatusJSON(res.StatusCode, res)
			return
		}

		tokenStr := authHeader[len(bearerPrefix):]
		jwksURL := viper.GetString("KEYCLOAK_URL") + "/realms/golectro/protocol/openid-connect/certs"
		jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{RefreshInterval: time.Hour})
		if err != nil {
			res := utils.FailedResponse(ctx, http.StatusInternalServerError, constants.InternalServerError, nil)
			ctx.AbortWithStatusJSON(res.StatusCode, res)
			return
		}

		token, err := jwt.Parse(tokenStr, jwks.Keyfunc)
		if err != nil || !token.Valid {
			res := utils.FailedResponse(ctx, http.StatusUnauthorized, constants.InvalidToken, nil)
			ctx.AbortWithStatusJSON(res.StatusCode, res)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			res := utils.FailedResponse(ctx, http.StatusUnauthorized, constants.InvalidToken, nil)
			ctx.AbortWithStatusJSON(res.StatusCode, res)
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok {
			res := utils.FailedResponse(ctx, http.StatusUnauthorized, constants.InvalidToken, nil)
			ctx.AbortWithStatusJSON(res.StatusCode, res)
			return
		}
		uid, err := utils.ParseUUID(userID)
		if err != nil {
			res := utils.FailedResponse(ctx, http.StatusUnauthorized, constants.InvalidToken, nil)
			ctx.AbortWithStatusJSON(res.StatusCode, res)
			return
		}

		var roles []string
		if realmAccess, ok := claims["realm_access"].(map[string]any); ok {
			if roleList, ok := realmAccess["roles"].([]any); ok {
				roles = make([]string, 0, len(roleList))
				for _, role := range roleList {
					if roleStr, ok := role.(string); ok {
						roles = append(roles, roleStr)
					}
				}
			}
		}

		rolesJSON, err := json.Marshal(roles)
		if err != nil {
			res := utils.FailedResponse(ctx, http.StatusInternalServerError, constants.InternalServerError, nil)
			ctx.AbortWithStatusJSON(res.StatusCode, res)
			return
		}

		username, _ := claims["preferred_username"].(string)
		email, _ := claims["email"].(string)

		auth := &model.Auth{
			ID:       uid,
			Username: username,
			Email:    email,
			Roles:    rolesJSON,
		}

		ctx.Set("auth", auth)
		ctx.Next()
	}
}

func GetUser(c *gin.Context) *model.Auth {
	if val, exists := c.Get("auth"); exists {
		return val.(*model.Auth)
	}
	return nil
}
