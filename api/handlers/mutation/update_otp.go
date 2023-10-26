package mutations

import (
	"context"
	"net/http"

	"github.com/Ethiopian-Education/edu-auth-server.git/config"
	"github.com/Ethiopian-Education/edu-auth-server.git/graph"
	"github.com/hasura/go-graphql-client"
	"github.com/sirupsen/logrus"
)

// should be multiple
func UpdateOtp(id string) error {
	var mutation struct {
		DeleteOTP struct {ID string `json:"id" graphql:"id"`} `graphql:"update_user_otp_by_pk(pk_columns:{id:$id}, _set:{used:true}) "`
	}

	client := &http.Client{
		Transport: &graph.Transport{T: http.DefaultTransport},
	}

	graph_client := graphql.NewClient(config.HASURA_GRAPHQL_URL, client)

	err := graph_client.Mutate(context.Background(), &mutation, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		logrus.Error("otp update mutate error", err)
		return  err
	}

	return nil
}