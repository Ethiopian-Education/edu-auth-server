package graph

import (
	"encoding/json"
	"net/http"

	"github.com/Ethiopian-Education/edu-auth-server.git/config"
)

func EncodeMessage(message map[string]interface{})([]byte, error) {
	reply_msg, err := json.Marshal(message)
	if err != nil {
		return nil , err
	}
	return reply_msg, nil
}

type Transport struct {
	T http.RoundTripper
}

func (t *Transport) RoundTrip(req *http.Request)(*http.Response, error) {
	req.Header.Add("x-hasura-admin-secret", config.HASURA_GRAPHQL_ADMIN_SECRET)
	return t.T.RoundTrip(req)
}

// customizable transporter...

type Header map[string]string

type CustomTransport struct {
	CT http.RoundTripper
	Headers Header
}


func (t *CustomTransport) RoundTrip(req *http.Request)(*http.Response, error) {
	for k, v := range t.Headers {
		req.Header.Add(k, v)
	}

	return t.RoundTrip(req)
}