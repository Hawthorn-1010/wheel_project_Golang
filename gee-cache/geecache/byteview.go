package geecache

/*
使用sync.Mutex封装LRU的几个方法，使之支持并发读写
*/

type ByteView struct {
	b []byte
}

// 实现Len方法
func (b ByteView) Len() int {
	return len(b.b)
}

func (b ByteView) ByteSlice() []byte {
	return CloneBytes(b.b)
}

func (b ByteView) String() string {
	return string(b.b)
}

func CloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
