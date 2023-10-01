package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Ethiopian-Education/edu-auth-server.git/config"
	"github.com/Ethiopian-Education/edu-auth-server.git/graph"
	"github.com/Ethiopian-Education/edu-auth-server.git/model"
	"github.com/Ethiopian-Education/edu-auth-server.git/utils"
	"github.com/Ethiopian-Education/edu-auth-server.git/utils/services"
	"github.com/gin-gonic/gin"
	"github.com/hasura/go-graphql-client"
	"github.com/sirupsen/logrus"
)
type user_users_insert_input struct {
	Email                 *string      `json:"email" graphql:"email,omitempty"`
	PhoneNumber           string      `json:"phone_number" graphql:"phone_number"`
	FirstName             string      `json:"first_name" graphql:"first_name"`
	MiddleName            string      `json:"middle_name" graphql:"middle_name"`
	LastName              string      `json:"last_name" graphql:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
	Password string `json:"password"`
	Gender string `json:"gender"`
	SignupMethod string `json:"signup_method"`
	OTPS struct {
		Data struct {
			Code string `json:"code" graphql:"code"`
			Type string `json:"type" graphql:"type"`
		} `json:"data" graphql:"data"`
	} `json:"otps" graphql:"otps"`
}
func SignUpHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		var body struct {
			Input struct {
				Params model.Signup `json:"params"`
			}`json:"input"`
		}

		if err = ctx.BindJSON(&body); err != nil {
			logrus.Errorf(`Error encountered when decoding signup req body : %v`, err)
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "unprocessed_request", Success: false})
			return
		}

		// logrus.Info("Sighup creds : ", body)

		hasedPassword, err := utils.HashPassword(body.Input.Params.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.Response{Message: "Bad request", Success: false})
		}
		trim_phone := strings.TrimSpace(body.Input.Params.PhoneNumber)

		generated_opt := utils.GenerateOTP(6)
		// logrus.Info("Generated otp : --- ", generated_opt)

		object := user_users_insert_input{
			FirstName: body.Input.Params.FirstName,
			MiddleName: body.Input.Params.MiddleName,
			LastName: body.Input.Params.LastName,
			Email: body.Input.Params.Email,
			Gender: body.Input.Params.Gender,
			Password: hasedPassword,
			PhoneNumber: trim_phone,
			SignupMethod: "regular_signup",
		}
		object.OTPS.Data.Code = generated_opt
		object.OTPS.Data.Type = "phone_verification"

		var mutation struct {
			InserUser model.User `graphql:"insert_user_user(object:$object)"`
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
			ctx.JSON(http.StatusBadRequest, model.Response{Message: "user_signup_mutation_error" , Success: false})
			return
		}

		// logrus.Infoln("User Added -- ", mutation.InserUser)
		if mutation.InserUser.Email != nil {
			logrus.Infoln("send otp via Email")
		}
		
		// Send OTP via phone...
		var m_body = model.TwilioBody{
			To:     trim_phone,
			Message: fmt.Sprintf(`%v - Is your confirmation code and valid for only 20 minutes, please confirm to activate your account.`, generated_opt),
		}
		err = services.TwilioSendSMS(m_body)
		if err != nil {
			logrus.Error("Something went wrong in twilio - ", err.Error())
		}

		ctx.JSON(http.StatusOK, model.Response{Message: fmt.Sprintf("Successful user id - %s", mutation.InserUser.ID), Success: true})
	}
}
