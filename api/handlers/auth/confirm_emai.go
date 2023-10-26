package auth

import (
	"fmt"
	"net/http"

	mutations "github.com/Ethiopian-Education/edu-auth-server.git/api/handlers/mutation"
	"github.com/Ethiopian-Education/edu-auth-server.git/api/handlers/queries"
	"github.com/Ethiopian-Education/edu-auth-server.git/model"
	otp_types "github.com/Ethiopian-Education/edu-auth-server.git/model/enum"
	"github.com/Ethiopian-Education/edu-auth-server.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func VerifyEmail() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		var err error
		var body struct {
			Input struct {
				Params model.ConfirmEmail `json:"params"`
			}`json:"input"`
		}

		if err = ctx.BindJSON(&body); err != nil {
			logrus.Errorf(`Error encountered when decoding signup req body : %v`, err)
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "invalid_request", Success: false})
			return
		}

		// FInd user
        
		isValid := utils.ValidateEmail(body.Input.Params.Email)
		if !isValid {
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "invalid_email_format", Success: false})
			return
		}
		filters := []string{
			fmt.Sprintf(`email:{_eq: "%s"}`, body.Input.Params.Email),
		}
		// logrus.Info("Filters -- ", strings.Join(filters, ","))

		user, err := queries.FindUser(filters)
		if err != nil {
			logrus.Error("find user error", err)
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "request_can_not_be_processed", Success: false})
			return
		}
		if user.IsEmailConfirmed {
			logrus.Error("Already confirmed !!!")
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "already_confirmed", Success: false})
			return
		}

		 otp_result, err := utils.VerityOTP(user.ID, body.Input.Params.Code, otp_types.EmailVerification)
		 if err != nil {
			logrus.Error("verify otp error", err)
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "request_can_not_be_processed", Success: false})
			return
		 }

		 if !otp_result.IsValid || otp_result.Used {
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "invalid_otp_or_expired", Success: false})
			return
		 }

		//  Verify Email ... 
		err = mutations.VerifyEmail(user.ID)
		if err != nil {
			logrus.Error("Verify mutations error", err)
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "request_can_not_be_processed", Success: false})
			return
		}

		if err = mutations.UpdateOtp(otp_result.ID); err != nil {
			logrus.Error("Update OTP mutations error .. ", err)
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "request_can_not_be_processed", Success: false})
			return
		}

		ctx.JSON(http.StatusOK, model.Response{
			Message: "Your email is successfully confirmed!",
			Success: true,
		})
	}
}