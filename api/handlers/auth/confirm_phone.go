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
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func VerifyPhoneNumber() gin.HandlerFunc{
	return func(ctx *gin.Context) {

		var err error

		var generated_email_otp string

		var successMessage = model.Response{Message: "your_phone_is_confirmed_successfully", Success: true}

		var body struct {
			Input struct {
				Params model.ConfirmPhone `json:"params"`
			}`json:"input"`
		}

		if err = ctx.BindJSON(&body); err != nil {
			logrus.Errorf(`Error encountered when decoding signup req body : %v`, err)
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "invalid_request", Success: false})
			return
		}

		// FInd user
        trim_phone := strings.TrimSpace(body.Input.Params.PhoneNumber)
		trim_phone, isValid := utils.ValidatePhone(trim_phone)
		if !isValid {
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "invalid_phone_number", Success: false})
			return
		}
		filters := []string{
			fmt.Sprintf(`phone_number:{_eq: "%s"}`, trim_phone),
		}
		// logrus.Info("Filters -- ", strings.Join(filters, ","))

		user, err := queries.FindUser(filters)
		if err != nil {
			logrus.Error("find user error", err)
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "request_can_not_be_processed", Success: false})
			return
		}
		if user.IsPhoneConfirmed {
			logrus.Error("Already confirmed !!!")
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "already_confirmed", Success: false})
			return
		}

		 otp_result, err := utils.VerityOTP(user.ID, body.Input.Params.Code, otp_types.PhoneVerification)
		 if err != nil {
			logrus.Error("verify otp error", err)
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "request_can_not_be_processed", Success: false})
			return
		 }

		 if !otp_result.IsValid || otp_result.Used {
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "invalid_otp_or_expired", Success: false})
			return
		 }

		//  Verify phone ... 
		err = mutations.VerifyPhone(user.ID)
		if err != nil {
			logrus.Error("Verify mutations error .. ", err)
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "request_can_not_be_processed", Success: false})
			return
		}

		// SEND EMAIL CONFIRMATION
		if user.Email != nil {
			logrus.Info("Inside email sending ... ")
			generated_email_otp = utils.GenerateOTP(6)
			if err = mutations.InsertOTP(model.OTP{Code: generated_email_otp, Type: otp_types.EmailVerification, UserID: user.ID}); err != nil {
				logrus.Error("mutate email confirmation code error : ",err)
				ctx.JSON(http.StatusBadRequest, model.Response{Message: "request_can_not_be_processed_try_again", Success: false})
			}
			if err = sendEmailConfirmation(user, generated_email_otp); err != nil {
				logrus.Error("Error encountered when send email-confirmation",err)
			}
			successMessage.Message = "confirmation sent via your email. confirm your email to log into your account."
		}

		if err = mutations.UpdateOtp(otp_result.ID); err != nil {
			logrus.Error("Update OTP mutations error .. ", err)
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "request_can_not_be_processed", Success: false})
			return
		}

		ctx.JSON(http.StatusOK, successMessage)
	}
}