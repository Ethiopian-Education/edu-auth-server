package mutations

import (
	"context"
	"net/http"

	"github.com/Ethiopian-Education/edu-auth-server.git/config"
	"github.com/Ethiopian-Education/edu-auth-server.git/graph"
	"github.com/hasura/go-graphql-client"
	"github.com/sirupsen/logrus"
)
type user_otp_bool_exp string
func RemoveUnusedOTP(userId string, otp_type string) error {
	var err error

	var mutation struct {
		RemoveOtp struct {
			AffectedRows int `json:"affected_rows" graphql:"affected_rows"`
		}`graphql:"delete_user_otp(where: {_and: {user_id: {_eq: $user_id}, type: {_eq: $type}, is_valid: {_eq: true}}})"`
	}

	client := &http.Client{
		Transport: &graph.Transport{T: http.DefaultTransport},
	}

	graph_client := graphql.NewClient(config.HASURA_GRAPHQL_URL, client)
	// var joined_where = strings.Join(where, ",")
	// var final_joint user_otp_bool_exp = user_otp_bool_exp(joined_where)

	err = graph_client.Mutate(context.Background(), &mutation, map[string]interface{}{
		"user_id": userId,
		"type": otp_type,
	})
	if err != nil {
		logrus.Error("Error on removal mutan_ : ", err.Error())
		return err
	}
	logrus.Infoln("Affected rows : " , mutation.RemoveOtp.AffectedRows)

	return nil
	
}