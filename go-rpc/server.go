package go_rpc

import (
	"encoding/json"
	"fmt"
	"go-rpc/codec"
	"io"
	"log"
	"net"
	"reflect"
	"sync"
)

const MagicNumber = 0xaaaaaa

// 消息的编解码方式
type Option struct {
	MagicNumber int
	CodeType    codec.Type
}

var DefaultOption = &Option{
	MagicNumber: MagicNumber,
	CodeType:    codec.GobType,
}

// 服务端首先使用 JSON 解码 Option，然后通过 Option 的 CodeType 解码剩余的内容

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

var DefaultServer = NewServer()

// 接受连接，服务request
func (server *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc server: Accept error: ", err)
			return
		}
		go server.ServeConn(conn)
	}
}

func Accept(lis net.Listener) {
	DefaultServer.Accept(lis)
}

// serve a single connection
func (server *Server) ServeConn(conn io.ReadWriteCloser) {
	// Json反序列化得到Option实例
	defer func() {
		_ = conn.Close()
	}()
	var opt Option
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Println("rec server: decode option error:", err)
		return
	}
	// 检查MagicNumber和CodeType是否正确
	if opt.MagicNumber != MagicNumber {
		log.Println("rpc server: MagicNumber error: %x", opt.MagicNumber)
		return
	}
	//根据opt中传来的CodecType来获取到构造方法
	f := codec.NewCodeFuncMap[opt.CodeType]
	if f == nil {
		log.Println("rpc server: CodecType error: %s", opt.CodeType)
		return
	}
	// 调用serveCodec
	server.serveCodec(f(conn))
}

var invalidRequest = struct {
}{}

func (server *Server) serveCodec(cc codec.Codec) {
	// 读取请求
	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	for {
		req, err := server.readRequest(cc)
		if err != nil {
			if req == nil {
				break
			}
			req.h.Error = err.Error()
			server.sendResponse(cc, req.h, invalidRequest, sending)
			continue
		}
		wg.Add(1)
		go server.handleRequest(cc, req, sending, wg)
	}
	// 处理请求
	wg.Wait()
	_ = cc.Close()

	// 回复请求
}

// 一个call的信息
type request struct {
	h            *codec.Header // header of request
	argv, replyV reflect.Value // argv and reply of request
}

func (server *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server: read header error:", err)
		}
		return nil, err
	}
	return &h, nil
}

func (server *Server) readRequest(cc codec.Codec) (*request, error) {
	h, err := server.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}
	req := &request{h: h}

	req.argv = reflect.New(reflect.TypeOf(""))
	if err = cc.ReadBody(req.argv.Interface()); err != nil {
		log.Println("rpc server: read argv err:", err)
	}
	return req, nil
}

func (server *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(h, body); err != nil {
		log.Println("rpc server: write response error: ", err)
	}
}

func (server *Server) handleRequest(c codec.Codec, req *request, sendLock *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println(req.h, req.argv.Elem())
	req.replyV = reflect.ValueOf(fmt.Sprintf("rpc response: %d", req.h.Seq))
	server.sendResponse(c, req.h, req.replyV.Interface(), sendLock)

}
