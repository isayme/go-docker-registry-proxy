package src

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isayme/go-logger"
)

func copyRequest(req *http.Request, headers http.Header) {
	reqHeaders := req.Header
	if reqHeaders == nil {
		reqHeaders = make(http.Header)
		req.Header = reqHeaders
	}

	for key := range headers {
		reqHeaders.Set(key, headers.Get(key))
	}
}

func copyResponse(c *gin.Context, resp *http.Response) {
	// c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("content-type"), resp.Body, nil)

	c.Status(resp.StatusCode)
	for key := range resp.Header {
		c.Header(key, resp.Header.Get(key))
	}
	n, err := io.Copy(c.Writer, resp.Body)
	if err != nil {
		logger.Errorf("copy %d bytes, fail: %v", n, err)
	}
}
