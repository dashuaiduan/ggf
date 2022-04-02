package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash maps bytes to uint32
type Hash func(data []byte) uint32

// Map constains all hashed keys 一致性哈希算法的主数据结构
type Map struct {
	hash     Hash
	replicas int            // 虚拟节点倍数   每个真实节点对应多少个虚拟节点
	keys     []int          // 哈希环   存入的是int 哈希值
	hashMap  map[int]string // 虚拟节点与真实节点的映射表，键是虚拟节点的哈希值，值是真实节点的名称。
}

// New creates a Map instance
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE // 标准库中 crc32 算法  给定一个字符串 生成一个固定的数字
	}
	return m
}

// Add adds some keys to the hash.  传入一批真实节点的名称  生成一个哈希环和 每个真实节点对应的虚拟节点
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ { // 生成当前真实节点的 虚拟节点
			hash := int(m.hash([]byte(strconv.Itoa(i) + key))) // 生成当前虚拟节点的 哈希值
			m.keys = append(m.keys, hash)                      //添加到哈希环		16 26 6这种格式
			m.hashMap[hash] = key                              //加虚拟节点和真实节点的映射关系。   16 26 6 ===> 6  这种格式map
		}
	}
	sort.Ints(m.keys) // 将虚拟环 排序  本身 会打乱顺序  因为每个虚拟节点 生成的int哈希值 是随机的
}

// 给定一个缓存key  查找 最近的真实节点
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	// Binary search for appropriate replica.
	idx := sort.Search(len(m.keys), func(i int) bool { // 使用二分查找法
		return m.keys[i] >= hash // 大于当前key的 hash值  的 第一个节点 索引	没找到的情况下 返回len
	})

	return m.hashMap[m.keys[idx%len(m.keys)]] //没找到的情况下 返回len  取数据肯定是错的 ，因此把没找到的数据放到第一个节点 也就形成了一个环状
}
