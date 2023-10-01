package utils

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Ethiopian-Education/edu-auth-server.git/config"
	"github.com/Ethiopian-Education/edu-auth-server.git/graph"
	"github.com/hasura/go-graphql-client"
)
type timestamptz string

func VerityOTP(userId string, otp string, otp_type string) (bool,error) {
    var err error
	
	client := &http.Client{
		Transport: &graph.Transport{T: http.DefaultTransport},
	}

	graph_client := graphql.NewClient(config.HASURA_GRAPHQL_URL, client)
	
	filters := []string{
		fmt.Sprintf(`code: {_eq: "%s"}`, otp),
		fmt.Sprintf(`user_id: {_eq: "%s"}`, userId),
		fmt.Sprintf(`type: {_eq:"%s"}`, otp_type),
	}

	query := fmt.Sprintf(`query {
		user_otp(where: {%s}) {
		  id
		  created_at
		  code
		  is_valid
		}
	  }`, strings.Join(filters, ","))

	res := struct {
		UserOTP []struct {
			ID string `json:"id" graphql:"id"`
			Code string `json:"code" graphql:"code"`
			CreatedAt string `json:"created_at" graphql:"created_at"`
			IsValid bool `json:"is_valid" graphql:"is_valid"`
		}`json:"user_otp" graphql:"user_otp"`
	}{}

	err = graph_client.Exec(context.Background(),query, &res, map[string]any{})
	if err != nil {
		return false, err
	}

	if len(res.UserOTP) == 0 {
		return false, nil
	}

	given_otp_data := res.UserOTP[0]

	// now := time.Now().Add(20 * time.Minute).Format(time.RFC3339Nano)

	// formatted_date_1, formatted_date_2, err := FormatDatesToSimilarZone(given_otp_data.CreatedAt, now)
	// if err != nil {
	// 	return false, err
	// }
	// logrus.Info("Formatted ", formatted_date_1, "And the second ", formatted_date_2)

	// logrus.Info("VALUE --- ", formatted_date_2 >= formatted_date_1)

	return given_otp_data.IsValid, nil
}