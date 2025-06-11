package utils

import (
	"fmt"
)

type node struct {
	value string
	key   string
	next  *node
	prev  *node
}

type LRUCache struct {
	capacity int
	size     int
	start    *node
	end      *node
	cache    map[string]*node
}

func (this *LRUCache) addNode(n *node) {
	if this.end != nil {
		this.end.next = n
		n.prev = this.end
		this.end = n
	} else {
		this.start = n
		this.end = n
	}
}

func (this *LRUCache) removeNode(n *node) {
	if n.next != nil {
		n.next.prev = n.prev
	} else {
		this.end = n.prev
	}
	if n.prev != nil {
		n.prev.next = n.next
	} else {
		this.start = n.next
	}
	n.next = nil
	n.prev = nil
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		size:     0,
		start:    nil,
		end:      nil,
		cache:    make(map[string]*node),
	}
}

func (this *LRUCache) Get(key string) (string, error) {
	if n, ok := this.cache[key]; ok {
		this.removeNode(n)
		this.addNode(n)
		return n.value, nil
	}
	return "", fmt.Errorf("Not a cache hit")
}

func (this *LRUCache) Put(key string, value string) {
	n, ok := this.cache[key]
	n = &node{
		value: value,
		key:   key,
	}
	if ok {
		this.removeNode(n)
		delete(this.cache, n.key)
	} else {
		if this.size == this.capacity {
			delete(this.cache, this.start.key)
			this.removeNode(this.start)
			this.size--
		}
		this.size++
	}
	this.cache[key] = n
	this.addNode(n)
}
