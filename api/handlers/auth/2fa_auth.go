package auth

import (
	"fmt"
	"net/http"
	"strings"

	mutations "github.com/Ethiopian-Education/edu-auth-server.git/api/handlers/mutation"
	"github.com/Ethiopian-Education/edu-auth-server.git/api/handlers/queries"
	"github.com/Ethiopian-Education/edu-auth-server.git/model"
	"github.com/Ethiopian-Education/edu-auth-server.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Authenticate2FA() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		var err error

		var body struct {
			Input struct {
				Params model.Auth2FA `json:"params"`
			}`json:"input"`
		}

		if err = ctx.BindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "unprocessed_request", Success: false})
			return
		}

		trim_phone := strings.TrimSpace(body.Input.Params.PhoneNumber)

		filters := []string{
			fmt.Sprintf(`phone_number:{_eq: "%s"}`, trim_phone),
		}

		user, err := queries.FindUser(filters)
		// Compare password
		if err != nil {
			logrus.Error("find user error", err)
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "invalid_credentials", Success: false})
			return
		}

		err = CheckUserValidity(user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, model.Response{Message: err.Error(), Success: false})
			return
		}
		// VERIFY OTP
		otp_result, err := utils.VerityOTP(user.ID, body.Input.Params.Code, "authentication")
		if err != nil {
		   logrus.Error("verify otp error", err)
		   ctx.JSON(http.StatusBadRequest, model.Response{Message: "request_can_not_be_processed", Success: false})
		   return
		}

		if !otp_result.IsValid {
		   ctx.JSON(http.StatusBadRequest, model.Response{Message: "invalid_otp_or_expired", Success: false})
		   return
		}

		// Build token
		accessToken, err := BuildToken(user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Response{
				Message: "Failed to sign the token...",
				Success: false,
			})
			return
		}

		if err = mutations.RemoveOTP(otp_result.ID); err != nil {
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "request_can_not_be_processed", Success: false})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"access_token": accessToken,
		})

	}
}