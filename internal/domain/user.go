package domain

import "net/http"

type UserInfo struct {
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Age    int    `json:"age"`
}

type ResponseUsers struct {
	Users []UserInfo `json:"users"`
}

func (c *ResponseUsers) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}
