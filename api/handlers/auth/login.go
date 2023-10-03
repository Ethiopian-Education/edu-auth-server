package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	mutations "github.com/Ethiopian-Education/edu-auth-server.git/api/handlers/mutation"
	"github.com/Ethiopian-Education/edu-auth-server.git/api/handlers/queries"
	jwt_jwt "github.com/Ethiopian-Education/edu-auth-server.git/crypto/jwt"
	"github.com/Ethiopian-Education/edu-auth-server.git/model"
	"github.com/Ethiopian-Education/edu-auth-server.git/utils"
	"github.com/Ethiopian-Education/edu-auth-server.git/utils/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoginHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
//  CURRENTLY WE USE PHONE_NUMBER TO LOGIN... IN THE RECENT FUTURE WE'LL UPGRADE THIS TO USE USERNAME, PHONE_NUMBER OR EMAIL AS USERS PREFERENCE
		var err error

		var body struct {
			Input struct {
				Params model.Login `json:"params"`
			}`json:"input"`
		}

		if err = ctx.BindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "unprocessed_request", Success: false})
			return
		}

		trim_phone := strings.TrimSpace(body.Input.Params.Username)

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
		// logrus.Info("User result --- ", user)
		isMatch := utils.CompareHashedPassword(user.Password, body.Input.Params.Password)
		if !isMatch {
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "invalid_credentials", Success: false})
			return
		}
		err = CheckUserValidity(user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, model.Response{Message: err.Error(), Success: false})
			return
		}
	
		// HERE CHECK IF THE USER ENABLE 2FA AUTHENTICATION 
		// ..... 'LL IMPLEMENT IN RECENT FUTURE
		if user.Enable2FA {
			err = Send2faAuthOTP(user)
			if err != nil {
				logrus.Error("2fa sending error -- ", err)
				ctx.JSON(http.StatusBadRequest, model.Response{Message: "2fa_error_encountered", Success: false})
				return
			}
			ctx.JSON(http.StatusOK,gin.H{"message": "Authentication code sent via your phone number."})
			return
		}

		accessToken, err := BuildToken(user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Response{
				Message: "Failed to sign the token...",
				Success: false,
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"access_token": accessToken,
		})
	}
}

func CheckUserValidity(user model.User) error {
		if !user.Active {
			return errors.New("disabled_or_inactive_account")
		}
		if((user.Email != nil && !user.IsEmailConfirmed) || (user.PhoneNumber != "" && !user.IsPhoneConfirmed)) {
			return errors.New("unverified_account")
		}
		return nil
}

// SEND 2FA AUTH OTP - CODE
func Send2faAuthOTP(user model.User) error {
	generatedOTP := utils.GenerateOTP(6)

	// insert otp
	otp_object := model.OTP {
		UserID: user.ID,
		Code: generatedOTP,
		Type: "authentication",
	}
	
	 err := mutations.InsertOTP(otp_object)
	if err != nil {
		logrus.Error("Error on otp mutation", err)
		return err
	}
	
	// send opt via phone
	var m_body = model.TwilioBody{
		To:     strings.TrimSpace(user.PhoneNumber),
		Message: fmt.Sprintf(`%v - Is your 2FA authentication code and valid for only 20 minutes.`, generatedOTP),
	}
	err = services.TwilioSendSMS(m_body)
	if err != nil {
		logrus.Error("Something went wrong in twilio - ", err.Error())
		return err
	}

	return nil
}

// Build accessToken

func BuildToken(user model.User) (string, error) {
	allowed_roles := []string{"user"}
	if len(user.UserRoles) > 0 {
	   for _, val := range user.UserRoles {
		   allowed_roles = append(allowed_roles, val.Role.Name)
	   }
	}

   var claims = &jwt_jwt.JWTClaims{}
   metadata := map[string]interface{}{
	   "roles":           allowed_roles,
	   "x-hasura-user-id": user.ID,
   }
   // set user specific claims data
   var user_email string
   if user.Email != nil {
	   user_email = *user.Email
   }
   claims.Metadata = metadata
   claims.Email = user_email
   // claims.LoginMethod = "regular_login"
   claims.TokenType = "access_token"
   claims.Subject = user.ID
   claims.First_name = user.FirstName
   claims.Last_name = user.LastName
   claims.Middle_name = user.MiddleName
   claims.SignUpMethod = user.SignUpMethod
   claims.Gender = user.Gender

   // HERE SIGN THE TOKEN
   accessToken, err := jwt_jwt.Sign(claims)
   if err != nil {
	return "", err
   }
   return accessToken, nil
}