package auth

import (
	"fmt"
	"net/http"
	"strings"

	mutations "github.com/Ethiopian-Education/edu-auth-server.git/api/handlers/mutation"
	"github.com/Ethiopian-Education/edu-auth-server.git/api/handlers/queries"
	"github.com/Ethiopian-Education/edu-auth-server.git/model"
	otp_types "github.com/Ethiopian-Education/edu-auth-server.git/model/enum"
	"github.com/Ethiopian-Education/edu-auth-server.git/utils"
	"github.com/Ethiopian-Education/edu-auth-server.git/utils/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type forgotPasswordBody struct {
	PhoneNumber string `json:"phone_number"`
}

func ForgotPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error

		var body struct {
			Input struct {
				Params forgotPasswordBody `json:"params"`
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

		generatedOTP := utils.GenerateOTP(6)
		otp_object := model.OTP{
			UserID: user.ID,
			Code:   generatedOTP,
			Type:   otp_types.Forget,
		}

		err = mutations.InsertOTP(otp_object)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "request_can_not_be_processed"})
			return
		}
		// send opt via phone
		var m_body = model.TwilioBody{
			To:      strings.TrimSpace(user.PhoneNumber),
			Message: fmt.Sprintf(`%v - Is your code to reset your password and valid for only 20 minutes.`, generatedOTP),
		}

		err = services.TwilioSendSMS(m_body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "request_can_not_be_processed"})
			return
		}

		ctx.JSON(http.StatusOK, model.Response{
			Message: "Reset password code is sent via your phone.",
			Success: true,
		})

	}
}
