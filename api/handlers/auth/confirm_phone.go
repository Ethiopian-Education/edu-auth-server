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

func VerifyPhoneNumber() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		var err error
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
		logrus.Info("body -- ", body)

		// FInd user
        trim_phone := strings.TrimSpace(body.Input.Params.PhoneNumber)
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

		 isValid, err := utils.VerityOTP(user.ID, body.Input.Params.Code, "phone_verification")
		 if err != nil {
			logrus.Error("verify otp error", err)
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "request_can_not_be_processed", Success: false})
			return
		 }

		 if !isValid {
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "invalid_otp_expired", Success: false})
			return
		 }

		//  Verify phone ... 
		err = mutations.VerifyPhone(user.ID)
		if err != nil {
			logrus.Error("Verify mutations error", err)
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "request_can_not_be_processed", Success: false})
			return
		}

		ctx.JSON(http.StatusOK, model.Response{
			Message: "Your phone is successfully confirmed",
			Success: true,
		})
	}
}