package srsclient

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

type HttpRequest struct {
	Addr    string
	Method  string
	Encoder httptransport.EncodeRequestFunc
	Decoder httptransport.DecodeResponseFunc
}

func (sr *HttpRequest) MakeRequest(path string) endpoint.Endpoint {
	u, err := url.Parse(sr.Addr)
	if err != nil {
		panic(err)
	}

	if u.Path == "" {
		u.Path = path
	} else {
		u.Path = strings.TrimSuffix(u.Path, "/") + path
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}

	c := httptransport.NewClient(
		sr.Method,
		u,
		sr.Encoder,
		sr.Decoder,
		httptransport.SetClient(client),
		httptransport.ClientBefore(func(ctx context.Context, req *http.Request) context.Context {
			req.Header.Set("Content-type", "application/json")
			return ctx
		}),
	).Endpoint()
	return c
}
