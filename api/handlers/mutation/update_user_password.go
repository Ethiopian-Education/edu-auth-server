package mutations

import (
	"context"
	"net/http"

	"github.com/Ethiopian-Education/edu-auth-server.git/config"
	"github.com/Ethiopian-Education/edu-auth-server.git/graph"
	"github.com/hasura/go-graphql-client"
	"github.com/sirupsen/logrus"
)

type user_users_set_input struct {
	PasswordChanged bool   `json:"password_changed" graphql:"password_changed"`
	Password        string `json:"password" graphql:"password"`
}

func UpdateUserPassword(password string, userId string) error {

	update_obj := user_users_set_input{
		Password:        password,
		PasswordChanged: true,
	}

	var mutation struct {
		UpdateUser struct {
			ID string `json:"id" graphql:"id"`
		} `graphql:"update_user_user(pk_columns:{id:$id}, _set: $set)"`
	}

	client := &http.Client{
		Transport: &graph.Transport{T: http.DefaultTransport},
	}

	graph_client := graphql.NewClient(config.HASURA_GRAPHQL_URL, client)

	err := graph_client.Mutate(context.Background(), &mutation, map[string]interface{}{
		"id":  userId,
		"set": update_obj,
	})
	if err != nil {
		logrus.Error("password mutate error : ", err)
		return err
	}

	return nil
}
