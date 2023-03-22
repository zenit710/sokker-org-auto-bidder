package client

import "fmt"

var _ Client = &httpClient{}

type httpClient struct {
	user string
	pass string
}

func NewHttpClient(user, pass string) *httpClient {
	return &httpClient{user: user, pass: pass}
}

func (s *httpClient) Auth() error {
	fmt.Printf("sokker.org auth. User: %s, Pass: %s", s.user, s.pass)
	return nil
}
