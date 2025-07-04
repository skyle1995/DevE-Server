package cache

import (
	"sync"
	"time"
)

// Item 缓存项结构体
type Item struct {
	Value      interface{}
	Expiration int64
}

// Cache 缓存结构体
type Cache struct {
	items             map[string]Item
	mu                sync.RWMutex
	DefaultExpiration time.Duration
	CleanupInterval   time.Duration
	stopCleanup       chan bool
}

// New 创建一个新的缓存实例
func New(defaultExpiration, cleanupInterval time.Duration) *Cache {
	items := make(map[string]Item)

	cache := &Cache{
		items:             items,
		DefaultExpiration: defaultExpiration,
		CleanupInterval:   cleanupInterval,
		stopCleanup:       make(chan bool),
	}

	// 如果清理间隔大于0，则启动定期清理过期项的goroutine
	if cleanupInterval > 0 {
		go cache.startCleanup()
	}

	return cache
}

// Set 设置缓存项，带过期时间
func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	var expiration int64

	if duration == 0 {
		duration = c.DefaultExpiration
	}

	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.mu.Lock()
	c.items[key] = Item{
		Value:      value,
		Expiration: expiration,
	}
	c.mu.Unlock()
}

// SetDefault 使用默认过期时间设置缓存项
func (c *Cache) SetDefault(key string, value interface{}) {
	c.Set(key, value, c.DefaultExpiration)
}

// Get 获取缓存项
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	item, found := c.items[key]
	if !found {
		c.mu.RUnlock()
		return nil, false
	}

	// 检查是否过期
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			c.mu.RUnlock()
			return nil, false
		}
	}

	c.mu.RUnlock()
	return item.Value, true
}

// Delete 删除缓存项
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}

// Exists 检查缓存项是否存在且未过期
func (c *Cache) Exists(key string) bool {
	_, found := c.Get(key)
	return found
}

// Flush 清空所有缓存项
func (c *Cache) Flush() {
	c.mu.Lock()
	c.items = make(map[string]Item)
	c.mu.Unlock()
}

// Count 获取缓存项数量
func (c *Cache) Count() int {
	c.mu.RLock()
	count := len(c.items)
	c.mu.RUnlock()
	return count
}

// Items 获取所有缓存项的副本
func (c *Cache) Items() map[string]Item {
	c.mu.RLock()
	defer c.mu.RUnlock()

	items := make(map[string]Item, len(c.items))
	for k, v := range c.items {
		// 检查是否过期
		if v.Expiration > 0 && time.Now().UnixNano() > v.Expiration {
			continue
		}
		items[k] = v
	}

	return items
}

// Keys 获取所有缓存键
func (c *Cache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.items))
	for k, v := range c.items {
		// 检查是否过期
		if v.Expiration > 0 && time.Now().UnixNano() > v.Expiration {
			continue
		}
		keys = append(keys, k)
	}

	return keys
}

// startCleanup 启动定期清理过期项的goroutine
func (c *Cache) startCleanup() {
	ticker := time.NewTicker(c.CleanupInterval)
	for {
		select {
		case <-ticker.C:
			c.deleteExpired()
		case <-c.stopCleanup:
			ticker.Stop()
			return
		}
	}
}

// deleteExpired 删除所有过期的缓存项
func (c *Cache) deleteExpired() {
	now := time.Now().UnixNano()
	c.mu.Lock()
	for k, v := range c.items {
		if v.Expiration > 0 && now > v.Expiration {
			delete(c.items, k)
		}
	}
	c.mu.Unlock()
}

// StopCleanup 停止清理goroutine
func (c *Cache) StopCleanup() {
	c.stopCleanup <- true
}

// GetWithTTL 获取缓存项及其剩余生存时间
func (c *Cache) GetWithTTL(key string) (interface{}, time.Duration, bool) {
	c.mu.RLock()
	item, found := c.items[key]
	if !found {
		c.mu.RUnlock()
		return nil, 0, false
	}

	// 检查是否过期
	if item.Expiration > 0 {
		now := time.Now().UnixNano()
		if now > item.Expiration {
			c.mu.RUnlock()
			return nil, 0, false
		}

		// 计算剩余生存时间
		ttl := time.Duration(item.Expiration - now)
		c.mu.RUnlock()
		return item.Value, ttl, true
	}

	c.mu.RUnlock()
	return item.Value, -1, true // -1 表示永不过期
}

// Increment 增加整数缓存项的值
func (c *Cache) Increment(key string, n int64) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.items[key]
	if !found {
		return 0, ErrKeyNotFound
	}

	// 检查是否过期
	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		return 0, ErrKeyNotFound
	}

	value, ok := item.Value.(int64)
	if !ok {
		return 0, ErrNotIntegerValue
	}

	newValue := value + n
	item.Value = newValue
	c.items[key] = item

	return newValue, nil
}

// Decrement 减少整数缓存项的值
func (c *Cache) Decrement(key string, n int64) (int64, error) {
	return c.Increment(key, -n)
}

// GetOrSet 获取缓存项，如果不存在则设置并返回
func (c *Cache) GetOrSet(key string, value interface{}, duration time.Duration) (interface{}, bool) {
	if val, found := c.Get(key); found {
		return val, true
	}

	c.Set(key, value, duration)
	return value, false
}

// GetAndDelete 获取缓存项并删除
func (c *Cache) GetAndDelete(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	// 检查是否过期
	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		delete(c.items, key)
		return nil, false
	}

	delete(c.items, key)
	return item.Value, true
}

// Errors
var (
	ErrKeyNotFound     = NewError("key not found")
	ErrNotIntegerValue = NewError("value is not an integer")
)

// Error 自定义错误类型
type Error struct {
	message string
}

// NewError 创建新的错误
func NewError(message string) *Error {
	return &Error{message: message}
}

// Error 实现error接口
func (e *Error) Error() string {
	return e.message
}
