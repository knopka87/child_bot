package llm

import (
	"net"
	"net/http"
	"strings"
	"time"
)

// HTTPClient — базовый HTTP клиент для LLM сервера с настроенными таймаутами
type HTTPClient struct {
	Base string
	HC   *http.Client
}

// NewHTTPClient создаёт новый HTTP клиент с оптимизированными настройками для долгих LLM запросов
func NewHTTPClient(base string) *HTTPClient {
	base = strings.TrimRight(base, "/")
	return &HTTPClient{
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
