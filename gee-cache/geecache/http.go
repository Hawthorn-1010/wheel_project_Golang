package geecache

import (
	"fmt"
	"gee-cache/geecache/consistenthash"
	pb "gee-cache/geecache/geecachepb"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	basePath        = "/geecache/"
	defaultReplicas = 50
)

type HTTPPool struct {
	self     string // 记录自己的地址，包括主机名/IP和端口
	basePath string // 节点间通讯地址的前缀
	mu       sync.Mutex
	peers    *consistenthash.Map
	// 一个远程节点对应一个httpGetter
	httpGetters map[string]*httpGetter // keyed by e.g. "http://10.0.0.2:8008"
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

	// 7.
	body, err := proto.Marshal(&pb.Response{Value: view.ByteSlice()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 告知浏览器这是一个字节流
	w.Header().Set("Content-Type", "application/octet-stream")
	//w.Write(view.ByteSlice())
	w.Write(body)
}

type httpGetter struct {
	baseUrl string
}

// 实现peers中的Get方法
func (h *httpGetter) Get(in *pb.Request, out *pb.Response) error {
	u := fmt.Sprintf(
		"%v%v%v",
		h.baseUrl,
		url.QueryEscape(in.GetGroup())+"/",
		url.QueryEscape(in.GetKey()),
	)

	// 发送请求
	res, err := http.Get(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}

	if err = proto.Unmarshal(bytes, out); err != nil {
		return fmt.Errorf("decoding response body: %v", err)
	}

	return nil
}

//var _ PeerGetter = (*httpGetter)(nil)

// 更新pool的节点名单
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistenthash.New(nil, defaultReplicas)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		// 只可能端口号不同？
		p.httpGetters[peer] = &httpGetter{baseUrl: peer + p.basePath}
	}
}

// 包装了一致性哈希算法的 Get() 方法，根据具体的 key，选择节点，返回节点对应的 HTTP 客户端
func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		// peer真实节点
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}

//var _ PeerPicker = (*HTTPPool)(nil)
