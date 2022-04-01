package client

import (
	"fmt"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/agent"
	"github.com/headzoo/surf/browser"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

func SetSurfClient(proxy string) (*browser.Browser, error) {
	client := surf.NewBrowser()
	client.SetUserAgent(agent.Chrome())
	jar, _ := cookiejar.New(nil)

	transport := &http.Transport{}
	// setup proxy
	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			return client, fmt.Errorf("%v %s", err, "Proxy format error")
		}

		// create http client
		transport = &http.Transport{
			Proxy:               http.ProxyURL(proxyURL),
			MaxIdleConns:        1,
			MaxIdleConnsPerHost: 1,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
		}
	}

	client.SetTransport(transport)
	client.SetCookieJar(jar)
	return client, nil
}
