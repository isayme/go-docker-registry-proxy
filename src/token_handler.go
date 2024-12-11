package src

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func TokenHandler(c *gin.Context) {
	// 原始 authenticate 是 base64 编码
	encodedAuthenticate := c.Query("authenticate")
	if encodedAuthenticate == "" {
		c.String(http.StatusBadRequest, "query param 'authenticate' not exist")
		return
	}

	authenticate, err := base64.URLEncoding.DecodeString(encodedAuthenticate)
	if err != nil {
		c.String(http.StatusBadRequest, "decode query param 'authenticate' failed")
		return
	}

	authenticateInfo, ok := ParseWwwAuthenticate(string(authenticate))
	if !ok {
		c.String(http.StatusBadRequest, "invalid query param 'authenticate'")
		return
	}

	// 获取 scope
	scope := c.Query("scope")
	if strings.HasPrefix(authenticateInfo.Realm, UPSTREAM_DOCKERHUB) {
		scopeParts := strings.Split(scope, ":")
		if len(scopeParts) == 3 && !strings.Contains(scopeParts[1], "/") {
			scopeParts[1] = "library/" + scopeParts[1]
			scope = strings.Join(scopeParts, ":")
		}
	}
	authenticateInfo.Scope = scope
	getToken(c, authenticateInfo)
}

func getToken(c *gin.Context, authenticate *WwwAuthenticate) {
	url := authenticate.Realm + "?scope=" + authenticate.Scope + "&service=" + authenticate.Service
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	headers := c.Request.Header
	req.Header = make(http.Header)
	for key := range headers {
		req.Header.Set(key, headers.Get(key))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer resp.Body.Close()

	copyResponse(c, resp)
}
