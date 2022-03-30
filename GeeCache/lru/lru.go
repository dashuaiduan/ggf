package lru

import "container/list"

// Cache is a LRU cache. It is not safe for concurrent access.
type Cache struct {
	maxBytes int64                    //当前队列 限制的最大使用内存
	nbytes   int64                    // 当前缓存队列已使用内存
	ll       *list.List               // 标准库的双向链表
	cache    map[string]*list.Element // 保存key和*list.Element 的映射关系 方便与查找的时候降低时间复杂度，o(n)--->o(1)
	// optional and executed when an entry is purged.
	OnEvicted func(key string, value Value) //   删除 缓存时的回调  支持自定义操作
}

//  存入到list.Element的数据结构Element.value 获取
type entry struct {
	key   string
	value Value //这里才是真正存数据的地方
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}

// New is the Constructor of Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted, // 删除缓存时的回调
	}
}

// Add adds a value to the cache.
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok { // 如果在map中 表明当前key已经有元素在缓存中 但未必是和原来一样的value
		c.ll.MoveToFront(ele)                                  // 将当前使用到的元素 移动到列表的最前面
		kv := ele.Value.(*entry)                               // 从Element中获取旧的 *entry 对象
		c.nbytes += int64(value.Len()) - int64(kv.value.Len()) // 新数据的value长度 - 旧数据vaule的长度 += 用新的数据长度替换旧数据的长度
		kv.value = value                                       //将 旧的entry.value 替换成 新的value
	} else { // map中不存在key 则表明是新增数据  直接添加到最前面
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele                               // 新数据 添加到map中
		c.nbytes += int64(len(key)) + int64(value.Len()) // 占用内存 = key的长度 = value的长度
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest() // 如果使用的内存 大于最大限度内存 则删除最后一个数据 一直删除 直到占用内存 < 最大限制
	}
}

// Get look ups a key's value
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele) //查询到数据之后  将数据移动到最前 表示活跃数据
		kv := ele.Value.(*entry)
		return kv.value, true // 返回 *entry.value 的真实数据
	}
	return
}

// RemoveOldest removes the oldest item
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}
