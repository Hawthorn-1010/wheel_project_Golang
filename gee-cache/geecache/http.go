package geecache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const basePath = "/geecache/"

type HTTPPool struct {
	self     string // 记录自己的地址，包括主机名/IP和端口
	basePath string // 节点间通讯地址的前缀
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: basePath,
	}
}

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path:" + r.URL.Path)
	}
	// GET /geecache/scores/Tom
	p.Log("%s %s", r.Method, r.URL.Path)
	// scores/Tom 切片
	// p.Log(r.URL.Path[len(p.basePath):])

	// /<basepath>/<groupname>/<key> required
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request!", http.StatusBadRequest)
		return
	}
	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group:"+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		// 这里会？
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 告知浏览器这是一个字节流
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}
