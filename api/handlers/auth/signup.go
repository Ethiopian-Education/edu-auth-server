package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Ethiopian-Education/edu-auth-server.git/config"
	"github.com/Ethiopian-Education/edu-auth-server.git/graph"
	"github.com/Ethiopian-Education/edu-auth-server.git/model"
	otp_types "github.com/Ethiopian-Education/edu-auth-server.git/model/enum"
	auth_method "github.com/Ethiopian-Education/edu-auth-server.git/model/enum/auth_types"
	"github.com/Ethiopian-Education/edu-auth-server.git/utils"
	"github.com/Ethiopian-Education/edu-auth-server.git/utils/services"
	"github.com/gin-gonic/gin"
	"github.com/hasura/go-graphql-client"
	"github.com/sirupsen/logrus"
)

type user_users_insert_input struct {
	Email        *string `json:"email" graphql:"email,omitempty"`
	PhoneNumber  string  `json:"phone_number" graphql:"phone_number"`
	FirstName    string  `json:"first_name" graphql:"first_name"`
	MiddleName   string  `json:"middle_name" graphql:"middle_name"`
	LastName     string  `json:"last_name" graphql:"last_name,omitempty"`
	Username     string  `json:"username,omitempty"`
	Password     string  `json:"password"`
	Gender       string  `json:"gender"`
	SignupMethod string  `json:"signup_method"`
	OTPS         struct {
		Data []struct {
			Code string `json:"code" graphql:"code"`
			Type string `json:"type" graphql:"type"`
		} `json:"data" graphql:"data"`
	} `json:"otps" graphql:"otps"`
}

func SignUpHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error

		var successMessage = model.Response{Message: "confirm_phone_number", Success: true}

		var body struct {
			Input struct {
				Params model.Signup `json:"params"`
			} `json:"input"`
		}

		if err = ctx.BindJSON(&body); err != nil {
			logrus.Errorf(`Error encountered when decoding signup req body : %v`, err)
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "unprocessed_request", Success: false})
			return
		}

		trim_phone := strings.TrimSpace(body.Input.Params.PhoneNumber)
		// check Phone number validity ...
		trim_phone, isValid := utils.ValidatePhone(trim_phone)
		if !isValid {
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "invalid_phone_number", Success: false})
			return
		}
		if body.Input.Params.Email != nil {
			logrus.Info("Params information : -- ", *body.Input.Params.Email)
			if !utils.ValidateEmail(strings.TrimSpace(*body.Input.Params.Email)) {
				ctx.JSON(http.StatusInternalServerError, model.Response{Message: "bad_email_format", Success: false})
			}
		}
		hasedPassword, err := utils.HashPassword(body.Input.Params.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Response{Message: "Bad request", Success: false})
		}

		generated_opt := utils.GenerateOTP(6)
		// logrus.Info("Generated otp : --- ", generated_opt)

		object := user_users_insert_input{
			FirstName:    body.Input.Params.FirstName,
			MiddleName:   body.Input.Params.MiddleName,
			LastName:     body.Input.Params.LastName,
			Email:        body.Input.Params.Email,
			Gender:       body.Input.Params.Gender,
			Password:     hasedPassword,
			PhoneNumber:  trim_phone,
			SignupMethod: auth_method.RegularAuth,
		}
		object.OTPS.Data = append(object.OTPS.Data, struct{Code string "json:\"code\" graphql:\"code\""; Type string "json:\"type\" graphql:\"type\""}{Code: generated_opt, Type: otp_types.PhoneVerification})
		// if body.Input.Params.Email != nil {
		// 	generated_email_otp = utils.GenerateOTP(6)
		// 	object.OTPS.Data = append(object.OTPS.Data, struct{Code string "json:\"code\" graphql:\"code\""; Type string "json:\"type\" graphql:\"type\""}{Code: generated_email_otp, Type: otp_types.EmailVerification})
		// }

		var mutation struct {
			InsertUser model.User `graphql:"insert_user_user(object:$object)"`
		}
		variables := map[string]interface{}{
			"object": object,
		}

		client := &http.Client{
			Transport: &graph.Transport{T: http.DefaultTransport},
		}

		graph_client := graphql.NewClient(config.HASURA_GRAPHQL_URL, client)

		if err = graph_client.Mutate(context.Background(), &mutation, variables); err != nil {
			logrus.Error("Mutation error : ", err.Error())
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "user_signup_mutation_error", Success: false})
			return
		}

		// // logrus.Infoln("User Added -- ", mutation.InsertUser)
		// if body.Input.Params.Email != nil {
		// 	logrus.Infoln("send otp via Email")
		// 	if err = sendEmailConfirmation(mutation.InsertUser, generated_email_otp); err != nil {
		// 		logrus.Error("Error encountered when send email-confirmation",err)
		// 	}
		// 	successMessage.Message = "confirm_phone_number_and_email"
		// }

		// Send OTP via phone...
		var m_body = model.TwilioBody{
			To:      trim_phone,
			Message: fmt.Sprintf(`%v - Is your confirmation code and valid for only 20 minutes, please confirm to activate your account.`, generated_opt),
		}
		err = services.TwilioSendSMS(m_body)
		if err != nil {
			logrus.Error("Something went wrong in twilio - ", err.Error())
		}

		ctx.JSON(http.StatusOK, successMessage)
	}
}

func sendEmailConfirmation(user model.User, code string) error {
	
	email_body_feed_data := map[string]string {
		"name":  fmt.Sprintf(`%s %s`, user.FirstName, user.MiddleName),
		"code": code,
	}

	body, err := services.ParseHtmlTemplate("./templates/email_confirmation.html", email_body_feed_data)
	if err != nil {
		logrus.Error("Error happened while parsing : - ", err)
		panic(1)
	}
	// logrus.Info("returned body : - ", body)

	to := []string{*user.Email}

	if err = services.SendEmailMessage(to, "abemelekmila@gmail.com","Confirm your Email.", body); err != nil{
		logrus.Error("Error happened while parsing : - ", err)
		panic(1)
	}
	return nil
}
