package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int            // 虚拟节点个数
	keys     []int          // 哈希环keys
	hashMap  map[int]string // 虚拟节点和真实节点的映射
}

func New(hash Hash, replicas int) *Map {
	m := &Map{
		hash:     hash,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// 添加机器(节点)
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			// 计算uint32
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key // hash: 确定虚拟节点; key: 确定真实节点
		}
	}
	sort.Ints(m.keys) //
}

// 选择节点
func (m *Map) Get(key string) string {
	if len(key) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))

	index := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	return m.hashMap[m.keys[index%len(m.keys)]]
}
