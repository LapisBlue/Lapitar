package httputil

import (
	"errors"
	"net/http"
	"net/url"
	"time"
)

const (
	userAgent = "Lapitar"
)

var cancelled = cancelledRequest{}

type cancelledRequest struct{}

func (err cancelledRequest) Error() string {
	return "Request cancelled"
}

type HttpClient struct {
	*http.Client
	RedirectHandler func(req *http.Request, via []*http.Request) bool
}

func Client() (c *HttpClient) {
	c = ForClient(&http.Client{
		Timeout: 30 * time.Second,
	})

	c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) >= 10 {
			return errors.New("stopped after 10 redirects")
		}

		req.Header = via[len(via)-1].Header
		if c.RedirectHandler != nil && !c.RedirectHandler(req, via) {
			return cancelled
		}

		return nil
	}

	return
}

func ForClient(c *http.Client) *HttpClient {
	return &HttpClient{Client: c}
}

func (c *HttpClient) Get(url string) (req *http.Request, err error) {
	req, err = http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", userAgent) // TODO: Version
	return
}

func IsCancelled(err error) (ok bool) {
	if err == nil {
		return
	}

	if _, ok = err.(cancelledRequest); !ok {
		var urlErr *url.Error
		if urlErr, ok = err.(*url.Error); ok {
			_, ok = urlErr.Err.(cancelledRequest)
		}
	}

	return
}
