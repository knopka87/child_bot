package llmclient

import (
	"net"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	Base string
	HC   *http.Client
}

func New(base string) *Client {
	base = strings.TrimRight(base, "/")
	return &Client{
		Base: base,
		HC: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				ResponseHeaderTimeout: 120 * time.Second, // ждём заголовки до 2 минут
			},
			Timeout: 0, // общий таймаут управляем per-request через ctx
		},
	}
}
