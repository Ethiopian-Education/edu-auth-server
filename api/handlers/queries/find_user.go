package queries

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Ethiopian-Education/edu-auth-server.git/config"
	"github.com/Ethiopian-Education/edu-auth-server.git/graph"
	"github.com/Ethiopian-Education/edu-auth-server.git/model"
	"github.com/hasura/go-graphql-client"
)

func FindUser(filters []string) (model.User, error) {
	var err error
	
	client := &http.Client{
		Transport: &graph.Transport{T: http.DefaultTransport},
	}

	graph_client := graphql.NewClient(config.HASURA_GRAPHQL_URL, client)
	
	query := fmt.Sprintf(`query {
		user_users(where: {%s}){
			id
			first_name
			middle_name
			last_name
			email
			phone_number
			enable_2fa
			is_email_confirmed
			is_phone_confirmed
			active
			password_changed
			password
			gender
			signup_method
			user_roles {
			  role {
				name
			  }
			}
		}
	  }
	`, strings.Join(filters, ","))

	res := struct {
		UserUsers []model.User `json:"user_users" graphql:"user_users"`
	}{}

	err = graph_client.Exec(context.Background(), query, &res, map[string]any{})
	if err != nil {
		return model.User{}, err
	}
	if len(res.UserUsers) == 0 {
		return model.User{}, errors.New("user_not_found")
	}

	return res.UserUsers[0], nil
}