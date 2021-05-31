package http

import (
	"github.com/kirinlabs/HttpRequest"
)

func Get(url string) string {
	res, err := HttpRequest.Get(url)
	if err == nil {
		body, err := res.Body()
		if err == nil {
			return string(body)
		} else {
			return err.Error()
		}
	} else {
		return err.Error()
	}
}

func Post(url string, params string) string {
	res, err := HttpRequest.Post(url, params)
	if err == nil {
		body, err := res.Body()
		if err == nil {
			return string(body)
		} else {
			return err.Error()
		}
	} else {
		return err.Error()
	}
}
