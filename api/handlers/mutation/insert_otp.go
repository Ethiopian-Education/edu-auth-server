package mutations

import (
	"context"
	"net/http"

	"github.com/Ethiopian-Education/edu-auth-server.git/config"
	"github.com/Ethiopian-Education/edu-auth-server.git/graph"
	"github.com/Ethiopian-Education/edu-auth-server.git/model"
	"github.com/hasura/go-graphql-client"
	"github.com/sirupsen/logrus"
)
type user_otp_insert_input model.OTP

func InsertOTP(otp model.OTP) ( error) {

	otp_creds := user_otp_insert_input{
		UserID: otp.UserID,
		Code: otp.Code,
		Type: otp.Type,
	}

	logrus.Info("otp creds ", otp_creds)

	var mutation struct {
		InsertUserOtpOne struct {ID string `json:"id" graphql:"id"`} `graphql:"insert_user_otp_one(object: $object)"`
	}

	client := &http.Client{
		Transport: &graph.Transport{T: http.DefaultTransport},
	}

	graph_client := graphql.NewClient(config.HASURA_GRAPHQL_URL, client)

	if err := graph_client.Mutate(context.Background(), &mutation, map[string]interface{}{
		"object": otp_creds,
	}); err != nil {
		logrus.Error("otp muta error", err)
		return  err
	}

	return nil
}