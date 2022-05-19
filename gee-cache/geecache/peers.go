package geecache

// 根据传入的key选择相应的节点PeerGetter
type PeerPicker interface {
	PickPeer(key string) (peer PeerPicker, ok bool)
}

// 从对应group中查找缓存值，对应于流程中的HTTP客户端
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
