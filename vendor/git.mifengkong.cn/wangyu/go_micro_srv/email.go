package go_micro_srv_client

import (
	"golang.org/x/net/context"
	"time"
)

func SendEmailByMicroService(to string, title string, body string, chaosong string, files []*File) (res bool, err error) {
	conn, err := GetSrvHandel("email")
	if conn == nil || err != nil {
		return false, err
	}
	c := NewEmailClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	request := &SendRequest{
		EmailTo:    to,
		EmailTitle: title,
		EmailBody:  body,
		EmailCc:    chaosong,
		EmailFile:  files,
	}
	_, err = c.Send(ctx, request)
	if err != nil {
		return false, err
	}
	return true, nil
}
