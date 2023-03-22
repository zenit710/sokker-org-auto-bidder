package client

type loginReqBody struct {
	Login string `json:"login"`
	Pass string `json:"password"`
}
