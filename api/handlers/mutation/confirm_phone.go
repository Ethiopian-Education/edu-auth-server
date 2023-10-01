package mutations

import (
	"context"
	"net/http"

	"github.com/Ethiopian-Education/edu-auth-server.git/config"
	"github.com/Ethiopian-Education/edu-auth-server.git/graph"
	"github.com/hasura/go-graphql-client"
)

func VerifyPhone(userId string) error {
	var err error
	var mutation struct {
		UpdateUserUser struct {
			ID string `json:"id" graphql:"id"`
		} `graphql:"update_user_user(pk_columns:{id: $id}, _set: {is_phone_confirmed:true})"`
	}

	client := &http.Client{
		Transport: &graph.Transport{T: http.DefaultTransport},
	}

	graph_client := graphql.NewClient(config.HASURA_GRAPHQL_URL, client)

    err = graph_client.Mutate(context.Background(), &mutation, map[string]interface{}{
		"id": userId,
	})
	if err != nil {
		return err
	}
	return nil
}
