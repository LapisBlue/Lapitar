package lhttp

import (
	"errors"
	"github.com/LapisBlue/Lapitar/lapitar"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	userAgent = lapitar.DisplayName
	timeout   = 10 * time.Second // Wait max. 10 seconds for the response

	TypeJSON = "application/json"
	TypePNG  = "image/png"
)

var (
	client = &http.Client{Timeout: timeout}
)

func Request(method, url string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url, body)
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", userAgent)
	return
}

func Get(url string) (*http.Request, error) {
	return Request("GET", url, nil)
}

func Do(req *http.Request) (resp *http.Response, err error) {
	return client.Do(req)
}

func NewError(resp *http.Response, err string) error {
	method := resp.Request.Method
	return &url.Error{
		Op:  method[0:1] + strings.ToLower(method[1:]),
		URL: resp.Request.URL.String(),
		Err: errors.New(err),
	}
}

func IsSuccess(resp *http.Response) bool {
	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		return true
	default:
		return false
	}
}

func IsNoContent(resp *http.Response) bool {
	switch resp.StatusCode {
	case http.StatusNotFound, http.StatusNoContent:
		return true
	default:
		return false
	}
}

func ExpectSuccess(resp *http.Response) (err error) {
	if !IsSuccess(resp) {
		err = NewError(resp, "Expected SUCCESS, got "+resp.Status+" instead")
	}

	return
}

func ExpectContent(resp *http.Response, expect string) (err error) {
	if content := resp.Header.Get("Content-Type"); content != expect {
		err = NewError(resp, "Expected "+expect+", got "+content+" instead")
	}

	return
}
