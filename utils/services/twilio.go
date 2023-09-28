package services

import (
	"encoding/json"
	"fmt"

	"github.com/Ethiopian-Education/edu-auth-server.git/config"
	"github.com/Ethiopian-Education/edu-auth-server.git/model"
	twilio "github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

func TwilioSendSMS(body model.TwilioBody) error {

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: config.TWILIO_ACCOUNT_SID,
		Password: config.TWILIO_AUTH_TOKEN,
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(body.To)
	params.SetFrom(config.TWILIO_PHONE_NUMBER)
	params.SetBody(body.Message)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println("Error sending SMS message: " + err.Error())
		return err
	} else {
		response, _ := json.Marshal(*resp)
		fmt.Println("Response: " + string(response))
	}
	return nil
}
