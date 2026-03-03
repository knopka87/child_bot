package telegram

import (
	"sync"
	"time"
)

// TTLCache — потокобезопасный кэш с автоматическим удалением устаревших записей
type TTLCache struct {
	data    sync.Map
	ttl     time.Duration
	name    string // для логирования
	metrics *cacheMetrics
}

type cacheEntry struct {
	value     interface{}
	expiresAt time.Time
}

type cacheMetrics struct {
	mu      sync.RWMutex
	hits    int64
	misses  int64
	evicted int64
}

// Конфигурация TTL для разных типов данных
const (
	// SessionTTL Сессионные данные (задача, подсказки) — 2 часа
	SessionTTL = 2 * time.Hour

	// PendingTTL Временные данные (pending контексты) — 10 минут
	PendingTTL = 10 * time.Minute

	// UserDataTTL Данные пользователя (класс, настройки) — 24 часа
	UserDataTTL = 24 * time.Hour

	// CleanupInterval Интервал очистки
	CleanupInterval = 5 * time.Minute
)

// NewTTLCache создает новый кэш с заданным TTL
func NewTTLCache(name string, ttl time.Duration) *TTLCache {
	return &TTLCache{
		ttl:     ttl,
		name:    name,
		metrics: &cacheMetrics{},
	}
}

// Store сохраняет значение с автоматическим TTL
func (c *TTLCache) Store(key int64, value interface{}) {
	c.data.Store(key, cacheEntry{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	})
}

// Load загружает значение, возвращает nil если истек TTL
func (c *TTLCache) Load(key int64) (interface{}, bool) {
	v, ok := c.data.Load(key)
	if !ok {
		c.metrics.mu.Lock()
		c.metrics.misses++
		c.metrics.mu.Unlock()
		return nil, false
	}

	entry, ok := v.(cacheEntry)
	if !ok {
		// Unexpected type - delete corrupted entry and return false
		c.data.Delete(key)
		c.metrics.mu.Lock()
		c.metrics.misses++
		c.metrics.mu.Unlock()
		return nil, false
	}

	if time.Now().After(entry.expiresAt) {
		// TTL истек — удаляем и возвращаем false
		c.data.Delete(key)
		c.metrics.mu.Lock()
		c.metrics.evicted++
		c.metrics.misses++
		c.metrics.mu.Unlock()
		return nil, false
	}

	c.metrics.mu.Lock()
	c.metrics.hits++
	c.metrics.mu.Unlock()
	return entry.value, true
}

// Touch обновляет TTL для существующей записи
// Использует атомарную операцию для избежания race condition
// Максимум 3 попытки, чтобы избежать бесконечного цикла при высокой конкуренции
func (c *TTLCache) Touch(key int64) {
	const maxRetries = 3
	for i := 0; i < maxRetries; i++ {
		v, ok := c.data.Load(key)
		if !ok {
			return
		}
		entry, ok := v.(cacheEntry)
		if !ok {
			// Unexpected type - delete and return
			c.data.Delete(key)
			return
		}
		newEntry := cacheEntry{
			value:     entry.value,
			expiresAt: time.Now().Add(c.ttl),
		}
		// CompareAndSwap ensures atomicity - if another goroutine modified
		// the entry between Load and here, this will fail and we retry
		if c.data.CompareAndSwap(key, v, newEntry) {
			return
		}
		// Entry was modified by another goroutine, retry (with limit)
	}
	// После maxRetries попыток просто выходим — TTL не критичен
}

// Delete удаляет запись
func (c *TTLCache) Delete(key int64) {
	c.data.Delete(key)
}

// Cleanup удаляет все устаревшие записи
func (c *TTLCache) Cleanup() int {
	var count int
	now := time.Now()

	c.data.Range(func(key, value interface{}) bool {
		entry, ok := value.(cacheEntry)
		if !ok {
			// Unexpected type - delete corrupted entry
			c.data.Delete(key)
			count++
			return true
		}
		if now.After(entry.expiresAt) {
			c.data.Delete(key)
			count++
		}
		return true
	})

	if count > 0 {
		c.metrics.mu.Lock()
		c.metrics.evicted += int64(count)
		c.metrics.mu.Unlock()
	}

	return count
}

// Size возвращает приблизительный размер кэша
func (c *TTLCache) Size() int {
	var count int
	c.data.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}

// Stats возвращает статистику кэша
func (c *TTLCache) Stats() (hits, misses, evicted int64, size int) {
	c.metrics.mu.RLock()
	hits = c.metrics.hits
	misses = c.metrics.misses
	evicted = c.metrics.evicted
	c.metrics.mu.RUnlock()
	size = c.Size()
	return
}

// CacheManager управляет всеми кэшами и их очисткой
type CacheManager struct {
	mu       sync.RWMutex
	caches   []*TTLCache
	stopChan chan struct{}
	wg       sync.WaitGroup
	stopOnce sync.Once // защита от повторного вызова Stop()
}

var (
	cacheManager     *CacheManager
	cacheManagerOnce sync.Once
)

// GetCacheManager возвращает singleton менеджера кэшей
func GetCacheManager() *CacheManager {
	cacheManagerOnce.Do(func() {
		cacheManager = &CacheManager{
			caches:   make([]*TTLCache, 0),
			stopChan: make(chan struct{}),
		}
	})
	return cacheManager
}

// Register добавляет кэш под управление менеджера
func (m *CacheManager) Register(cache *TTLCache) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.caches = append(m.caches, cache)
}

// Start запускает фоновую очистку
func (m *CacheManager) Start() {
	shutdown := GetShutdownManager()

	// Не запускаем если идёт shutdown
	if shutdown.IsShutdown() {
		return
	}

	// Регистрируем горутину
	goroutineDone := shutdown.TrackGoroutine()

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		defer goroutineDone()

		ticker := time.NewTicker(CleanupInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				m.cleanup()
			case <-m.stopChan:
				return
			case <-shutdown.Done():
				// Graceful shutdown
				return
			}
		}
	}()
}

// Stop останавливает фоновую очистку (безопасен для повторного вызова)
func (m *CacheManager) Stop() {
	m.stopOnce.Do(func() {
		close(m.stopChan)
	})
	m.wg.Wait()
}

func (m *CacheManager) cleanup() {
	m.mu.RLock()
	caches := make([]*TTLCache, len(m.caches))
	copy(caches, m.caches)
	m.mu.RUnlock()

	for _, cache := range caches {
		cache.Cleanup()
	}
	// Также очищаем зависшие батчи фото
	cleanupStaleBatches()
}

// GetAllStats возвращает статистику всех кэшей
func (m *CacheManager) GetAllStats() map[string]map[string]int64 {
	m.mu.RLock()
	caches := make([]*TTLCache, len(m.caches))
	copy(caches, m.caches)
	m.mu.RUnlock()

	stats := make(map[string]map[string]int64)
	for _, cache := range caches {
		hits, misses, evicted, size := cache.Stats()
		stats[cache.name] = map[string]int64{
			"hits":    hits,
			"misses":  misses,
			"evicted": evicted,
			"size":    int64(size),
		}
	}
	return stats
}
