package codec

import "io"

/**
消息编解码相关
*/

type Header struct {
	// 服务名+方法名
	ServiceMethod string
	// 序列号，某个请求的id，区分不同的请求 *
	Seq uint64
	// 错误信息
	Error string
}

// 抽象对消息体进行编解码的接口Codec
// 抽象出不同的接口是为了实现不同的Codec实例
type Codec interface {
	io.Closer
	// 读取Header
	ReadHeader(*Header) error
	// 读取消息体
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

// 抽象Codec的构造函数
type NewCodecFunc func(io.ReadWriteCloser) Codec

type Type string

const (
	// Gob协议解析
	GobType Type = "application/gob"
	// Json协议解析
	JsonType Type = "application/json"
)

var NewCodeFuncMap map[Type]NewCodecFunc

func init() {
	NewCodeFuncMap = make(map[Type]NewCodecFunc)
	NewCodeFuncMap[GobType] = NewGobCodec
}
