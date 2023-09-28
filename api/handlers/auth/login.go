package auth

import (
	"net/http"

	jwt_jwt "github.com/Ethiopian-Education/edu-auth-server.git/crypto/jwt"
	"github.com/Ethiopian-Education/edu-auth-server.git/model"
	"github.com/gin-gonic/gin"
)

func LoginHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var err error

		var claims = &jwt_jwt.JWTClaims{}
		metadata := map[string]interface{}{
			"roles":            []string{"user"},
			"x-hasura-user-id": "8cacd89f-9d0d-4035-a1a3-b1a338bef411",
		}
		// set user specific claims data
		claims.Metadata = metadata
		claims.Email = "tester@gmail.com"
		claims.LoginMethod = "regular_login"
		claims.TokenType = "access_token"
		claims.Subject = "4567890987654567-567gau90909-09RAS89"

		accessToken, err := jwt_jwt.Sign(claims)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Response{
				Message: "Failed to sign the token...",
				Success: false,
			})
			return
		}

		ctx.JSON(http.StatusOK, model.Response{
			Message: accessToken,
			Success: true,
		})
	}
}
