package client

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-list-templ/grpc/config"
	"github.com/valyala/fasthttp"
)

const ContentType = "application/json"

var ErrHTTPStatusNotOK = errors.New("http status not ok")

type Client struct {
	*fasthttp.Client
}

func New(cfg config.Client) *Client {
	return &Client{
		&fasthttp.Client{
			ReadTimeout:                   cfg.ReadTimeout,
			WriteTimeout:                  cfg.WriteTimeout,
			MaxIdleConnDuration:           cfg.MaxIdle,
			NoDefaultUserAgentHeader:      true,
			DisableHeaderNamesNormalizing: true,
			DisablePathNormalizing:        true,
			Dial: (&fasthttp.TCPDialer{
				Concurrency:      fasthttp.DefaultConcurrency,
				DNSCacheDuration: time.Hour,
			}).Dial,
		},
	}
}

func (c *Client) Get(uri string) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	req.SetRequestURI(uri)
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.SetContentType(ContentType)
	err := c.Do(req, res)
	fasthttp.ReleaseRequest(req)

	if err != nil {
		return res, err
	}

	if res.StatusCode() != http.StatusOK {
		return res, fmt.Errorf("%w: %d", ErrHTTPStatusNotOK, res.StatusCode())
	}

	return res, nil
}

func (c *Client) ReleaseGet(res *fasthttp.Response) {
	fasthttp.ReleaseResponse(res)
}
