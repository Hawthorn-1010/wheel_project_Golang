package geecache

import pb "gee-cache/geecache/geecachepb"

// 根据传入的key选择相应的节点PeerGetter
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// 从对应group中查找缓存值，对应于流程中的HTTP客户端
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
