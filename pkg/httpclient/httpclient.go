package httpclient

import (
	"time"

	"github.com/go-list-templ/grpc/config"
	"github.com/valyala/fasthttp"
)

const (
	UserAgent   = "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:15.0) Gecko/20100101 Firefox/15.0.1"
	ContentType = "application/json"
)

type Client struct {
	*fasthttp.Client
}

func New(cfg config.Client) *Client {
	return &Client{
		&fasthttp.Client{
			ReadTimeout:                   cfg.ReadTimeout,
			WriteTimeout:                  cfg.WriteTimeout,
			MaxIdleConnDuration:           cfg.MaxIdle,
			NoDefaultUserAgentHeader:      true, // Don't send: User-Agent: fasthttp
			DisableHeaderNamesNormalizing: true, // If you set the case on your headers correctly you can enable this
			DisablePathNormalizing:        true,
			// increase DNS cache time to an hour instead of default minute
			Dial: (&fasthttp.TCPDialer{
				Concurrency:      4096,
				DNSCacheDuration: time.Hour,
			}).Dial,
		},
	}
}

func (c *Client) Get(uri string) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	req.SetRequestURI(uri)
	req.Header.SetUserAgent(UserAgent)
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.SetContentType(ContentType)
	err := c.Do(req, resp)
	fasthttp.ReleaseRequest(req)

	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (c *Client) ReleaseGet(res *fasthttp.Response) {
	fasthttp.ReleaseResponse(res)
}
