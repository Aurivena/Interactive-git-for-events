package initialization

import (
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func newVpnClient() *http.Client {
	proxyURL, _ := url.Parse("http://127.0.0.1:8888")

	vpnHosts := map[string]struct{}{
		"generativelanguage.googleapis.com": {},
		// если надо — добавляй хосты сюда
	}

	tr := &http.Transport{
		Proxy: func(r *http.Request) (*url.URL, error) {
			if _, ok := vpnHosts[strings.ToLower(r.URL.Hostname())]; ok {
				return proxyURL, nil
			}
			return nil, nil
		},
		ForceAttemptHTTP2:   true,
		MaxIdleConns:        200,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     60 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}
	return &http.Client{Transport: tr, Timeout: 30 * time.Second}
}
