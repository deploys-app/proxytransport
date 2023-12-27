package proxytransport

import (
	"encoding/base64"
	"net/http"
	"sync"
)

const Host = "proxy.deploys.app"

type Auth struct {
	User     string
	Password string
}

func (a Auth) basicAuth() string {
	if a.User == "" {
		return ""
	}
	auth := a.User + ":" + a.Password
	s := base64.StdEncoding.EncodeToString([]byte(auth))
	return "Basic " + s
}

type Transport struct {
	http.RoundTripper
	Auth Auth

	init      sync.Once
	basicAuth string
}

func (t *Transport) tr() http.RoundTripper {
	if t.RoundTripper == nil {
		return http.DefaultTransport
	}
	return t.RoundTripper
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.init.Do(func() {
		t.basicAuth = t.Auth.basicAuth()
	})
	req.Header.Set("X-Proxy-Method", req.Method)
	req.Header.Set("X-Proxy-URL", req.URL.String())
	if t.basicAuth != "" {
		req.Header.Set("X-Proxy-Authorization", t.basicAuth)
	}
	req.Method = http.MethodPost
	req.Host = Host
	req.URL.Scheme = "https"
	req.URL.Host = Host
	req.URL.Path = "/"
	req.URL.RawQuery = ""
	return t.tr().RoundTrip(req)
}
