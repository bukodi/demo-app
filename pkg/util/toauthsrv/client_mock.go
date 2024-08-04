package toauthsrv

import (
	"net/http"
	"net/http/cookiejar"
	"strings"
	"testing"
)

type ClientMock struct {
	http.Client
	TestingT *testing.T
}

func NewClientMock(t *testing.T) *ClientMock {
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("%+v", err)
		return nil
	}
	cli := &ClientMock{
		TestingT: t,
		Client: http.Client{
			Jar: jar,
		},
	}
	cli.CheckRedirect = cli.checkRedirect
	return cli
}

func (c *ClientMock) checkRedirect(req *http.Request, via []*http.Request) error {
	urls := []string{}
	for _, r := range via {
		urls = append(urls, r.URL.String())
	}
	urls = append(urls, req.URL.String())
	c.TestingT.Logf("Redirect: %s", strings.Join(urls, " -> "))
	return nil
}
