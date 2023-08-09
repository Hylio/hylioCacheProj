package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_hyliocache/"

type HTTPPool struct {
	self     string
	basePath string
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p *HTTPPool) Serve(c *gin.Context) {
	// 限制访问路径
	if !strings.HasPrefix(c.Request.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + c.Request.URL.Path)
	}
	p.Log("%s %s", c.Request.Method, c.Request.URL.Path)
	// /<basepath>/<groupname>/<key>
	parts := strings.SplitN(c.Request.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	groupName, key := parts[0], parts[1]

	group := GetGroup(groupName)

	if group == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "no such group: " + groupName})
		return
	}

	view, err := group.Get(key)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Writer.Write(view.ByteSlice())
}
