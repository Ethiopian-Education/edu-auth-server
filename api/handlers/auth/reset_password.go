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

type resetPasswordBody struct {
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code"`
	Password    string `json:"password"`
}

func ResetPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error

		var body struct {
			Input struct {
				Params resetPasswordBody `json:"params"`
			} `json:"input"`
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
		if err != nil {
			logrus.Error("find user error", err)
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "request_can_not_be_processed", Success: false})
			return
		}

		err = CheckUserValidity(user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, model.Response{Message: err.Error(), Success: false})
			return
		}

		// Verify OTP
		otp_result, err := utils.VerityOTP(user.ID, body.Input.Params.Code, "forget")
		if err != nil {
			logrus.Error("verify otp error", err)
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "request_can_not_be_processed", Success: false})
			return
		}

		if !otp_result.IsValid {
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "invalid_otp_or_expired", Success: false})
			return
		}

		hashedPassword, err := utils.HashPassword(body.Input.Params.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Response{Message: "internal_server_error", Success: false})
			return
		}

		err = mutations.UpdateUserPassword(hashedPassword, user.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Response{Message: "mutation_error", Success: false})
			return
		}

		// remove otp after use
		err = mutations.RemoveOTP(user.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Response{Message: "mutation_error", Success: false})
			return
		}

		ctx.JSON(http.StatusOK, model.Response{
			Message: "password_reset_successfully",
			Success: true,
		})

	}
}
