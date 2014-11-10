package httputil

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	userAgent = "Lapitar"
)

var (
	client = &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return errors.New("stopped after 10 redirects")
			}

			req.Header = via[len(via)-1].Header
			return nil
		},
	}

	done = errors.New("done")
)

func Request(method, url string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url, body)
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", userAgent) // TODO: Version
	return
}

func Get(url string) (*http.Request, error) {
	return Request("GET", url, nil)
}

func Do(req *http.Request) (*http.Response, error) {
	return client.Do(req)
}

func GetLocation(req *http.Request, destHost string) (loc *http.Request, err error) {
	c := new(http.Client)
	*c = *client

	c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) >= 10 {
			return errors.New("stopped after 10 redirects")
		}

		req.Header = via[len(via)-1].Header
		if req.URL.Host == destHost {
			loc = req
			return done
		}

		return nil
	}

	resp, err := c.Do(req)
	if loc == nil {
		if err != nil {
			return
		}

		if !IsRedirect(resp) {
			err = NewError(resp, "Expected redirect, got "+resp.Status+" instead")
			return
		}

		err = NewError(resp, "Failed to get location of "+destHost)
		return
	}

	err = nil // Who cares, that's probably our fault
	return
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

func IsRedirect(resp *http.Response) bool {
	switch resp.StatusCode {
	case http.StatusMovedPermanently, http.StatusFound, http.StatusSeeOther, http.StatusTemporaryRedirect:
		return true
	default:
		return false
	}
}
