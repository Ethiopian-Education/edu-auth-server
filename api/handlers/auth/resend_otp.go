package auth

import (
	"fmt"
	"net/http"
	"strings"

	mutations "github.com/Ethiopian-Education/edu-auth-server.git/api/handlers/mutation"
	"github.com/Ethiopian-Education/edu-auth-server.git/api/handlers/queries"
	"github.com/Ethiopian-Education/edu-auth-server.git/model"
	"github.com/Ethiopian-Education/edu-auth-server.git/utils"
	"github.com/Ethiopian-Education/edu-auth-server.git/utils/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type resendOTPBody struct {
	Email string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Type string `json:"otp_type"`
}

func ResendOTP() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		var err error

		var body struct {
			Input struct {
				Params resendOTPBody `json:"params"`
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


		// Search otp with a given user and otp type .... and descending order from otp storage ... 
		// where_2 := []string{
		// 	fmt.Sprintf(`user_id: {_eq: "%s"}`, user.ID),
		// 	fmt.Sprintf(`type: {_eq: "%s"}`, user.ID),
		// 	fmt.Sprintf(`is_valid: {_eq: true}`),
		// }
		// where := fmt.Sprintf(`
		// {_and: {user_id: {_eq: "%s"}, type: {_eq: "%s"}, is_valid: {_eq: true}}
		// `, user.ID, body.Input.Params.Type)
		if err = mutations.RemoveUnusedOTP(user.ID, body.Input.Params.Type); err != nil {
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "error_encountered_when_removing_unused_otp's", Success: false})
			return
		}
		//  Here we send the opt code to a given phone_number
		generatedOTP := utils.GenerateOTP(6)
		// send opt via phone
		var m_body = model.TwilioBody{
			To:      strings.TrimSpace(user.PhoneNumber),
			Message: fmt.Sprintf(`%v - Is your code for %s service and valid for only 20 minutes.`, generatedOTP, body.Input.Params.Type),
		}

		err = services.TwilioSendSMS(m_body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "request_can_not_be_processed"})
			return
		}
			ctx.JSON(http.StatusOK, model.Response{Message: "otp_sent_successfully", Success: false})
		}
}