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
	reqHeader := c.Request.Header
	path := c.Request.URL.Path

	upstream := routeByHost(c)

	newReq, err := http.NewRequest(method, upstream+path, nil)
	if err != nil {
		logger.Warnf("new request erro: %v", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	copyRequest(newReq, reqHeader)

	resp, err := http.DefaultClient.Do(newReq)
	if err != nil {
		logger.Warnf("send request erro: %v", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer resp.Body.Close()

	authenticateInText := resp.Header.Get(HTTP_HEADER_WWW_AUTHENTICATE)
	authenticate, ok := ParseWwwAuthenticate(authenticateInText)
	if ok {
		// set realm so that client will get token from this service
		url := getTokenUrl(c)
		realm := fmt.Sprintf("%s%s?authenticate=%s", url, "/__token__", base64.URLEncoding.EncodeToString([]byte(authenticateInText)))
		authenticate.Realm = realm
		resp.Header.Set(HTTP_HEADER_WWW_AUTHENTICATE, authenticate.String())
	}

	copyResponse(c, resp)
}

func getTokenUrl(c *gin.Context) string {
	host := c.GetHeader(HTTP_HEADER_X_FORWARDED_HOST)
	proto := c.GetHeader(HTTP_HEADER_X_FORWARDED_PROTO)
	port := c.GetHeader(HTTP_HEADER_X_FORWARDED_PORT)

	if proto == "" || host == "" {
		proto = HTTP_PROTO_HTTP
		if c.Request.TLS != nil {
			proto = HTTP_PROTO_HTTPS
		}
		host = c.Request.Host

		return fmt.Sprintf("%s://%s", proto, host)
	}

	if port == "" ||
		(proto == HTTP_PROTO_HTTP && port == HTTP_PORT_80) ||
		(proto == HTTP_PROTO_HTTPS && port == HTTP_PORT_443) {
		return fmt.Sprintf("%s://%s", proto, host)
	}

	return fmt.Sprintf("%s://%s:%s", proto, host, port)
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
