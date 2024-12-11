package src

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isayme/go-logger"
)

func V2Handler(c *gin.Context) {
	method := c.Request.Method
	host := c.Request.Host
	reqHeader := c.Request.Header
	path := c.Request.URL.Path

	upstream := routeByHost(c)

	newReq, err := http.NewRequest(method, upstream+path, nil)
	if err != nil {
		logger.Warnf("new request erro: %v", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	newReq.Header = make(http.Header)
	for key := range reqHeader {
		newReq.Header.Set(key, reqHeader.Get(key))
	}

	resp, err := http.DefaultClient.Do(newReq)
	if err != nil {
		logger.Warnf("send request erro: %v", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer resp.Body.Close()

	authenticateInText := resp.Header.Get("WWW-Authenticate")
	authenticate, ok := ParseWwwAuthenticate(authenticateInText)
	if ok {
		// set realm so that client will get token from this service
		proto := c.GetHeader("X-Forwarded-Proto")
		if proto == "" {
			proto = "http"
		}
		realm := fmt.Sprintf("%s://%s%s?authenticate=%s", proto, host, "/__token__", base64.URLEncoding.EncodeToString([]byte(authenticateInText)))
		authenticate.Realm = realm
		resp.Header.Set("WWW-Authenticate", authenticate.String())
	}

	copyResponse(c, resp)
}

func routeByHost(c *gin.Context) string {
	host := c.Request.Host

	conf := GetConfig()
	for _, route := range conf.Routes {
		if route.Host == host {
			return route.Upstream
		}
	}

	// default upstream
	return UPSTREAM_DOCKERHUB
}
