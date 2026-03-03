package telegram

import (
	"testing"
	"time"
)

func TestTTLCache_BasicOperations(t *testing.T) {
	cache := NewTTLCache("test", 1*time.Hour)

	// Test Store and Load
	cache.Store(123, "value1")
	val, ok := cache.Load(123)
	if !ok {
		t.Error("Expected to find value")
	}
	if val.(string) != "value1" {
		t.Errorf("Expected 'value1', got '%v'", val)
	}

	// Test non-existent key
	_, ok = cache.Load(999)
	if ok {
		t.Error("Expected not to find non-existent key")
	}

	// Test Delete
	cache.Delete(123)
	_, ok = cache.Load(123)
	if ok {
		t.Error("Expected not to find deleted key")
	}
}

func TestTTLCache_Expiration(t *testing.T) {
	// Короткий TTL для теста
	cache := NewTTLCache("test_expiry", 50*time.Millisecond)

	cache.Store(1, "will_expire")

	// Сразу должно быть доступно
	_, ok := cache.Load(1)
	if !ok {
		t.Error("Expected to find value immediately")
	}

	// Ждём истечения TTL
	time.Sleep(100 * time.Millisecond)

	// Теперь должно быть удалено
	_, ok = cache.Load(1)
	if ok {
		t.Error("Expected value to be expired")
	}
}

func TestTTLCache_Touch(t *testing.T) {
	cache := NewTTLCache("test_touch", 100*time.Millisecond)

	cache.Store(1, "value")

	// Ждём половину TTL
	time.Sleep(60 * time.Millisecond)

	// Touch обновляет TTL
	cache.Touch(1)

	// Ждём ещё половину исходного TTL
	time.Sleep(60 * time.Millisecond)

	// Значение должно быть доступно благодаря Touch
	_, ok := cache.Load(1)
	if !ok {
		t.Error("Expected value to still exist after Touch")
	}
}

func TestTTLCache_Cleanup(t *testing.T) {
	cache := NewTTLCache("test_cleanup", 50*time.Millisecond)

	// Добавляем несколько записей
	cache.Store(1, "a")
	cache.Store(2, "b")
	cache.Store(3, "c")

	// Проверяем размер
	if cache.Size() != 3 {
		t.Errorf("Expected size 3, got %d", cache.Size())
	}

	// Ждём истечения
	time.Sleep(100 * time.Millisecond)

	// Запускаем очистку
	evicted := cache.Cleanup()
	if evicted != 3 {
		t.Errorf("Expected 3 evicted, got %d", evicted)
	}

	if cache.Size() != 0 {
		t.Errorf("Expected size 0 after cleanup, got %d", cache.Size())
	}
}

func TestTTLCache_Stats(t *testing.T) {
	cache := NewTTLCache("test_stats", 1*time.Hour)

	cache.Store(1, "value")

	// Hit
	cache.Load(1)
	cache.Load(1)

	// Miss
	cache.Load(999)

	hits, misses, _, size := cache.Stats()

	if hits != 2 {
		t.Errorf("Expected 2 hits, got %d", hits)
	}
	if misses != 1 {
		t.Errorf("Expected 1 miss, got %d", misses)
	}
	if size != 1 {
		t.Errorf("Expected size 1, got %d", size)
	}
}

func TestCacheManager(t *testing.T) {
	// Создаём новый менеджер для теста (не используем singleton)
	manager := &CacheManager{
		caches:   make([]*TTLCache, 0),
		stopChan: make(chan struct{}),
	}

	cache1 := NewTTLCache("cache1", 1*time.Hour)
	cache2 := NewTTLCache("cache2", 1*time.Hour)

	manager.Register(cache1)
	manager.Register(cache2)

	cache1.Store(1, "a")
	cache2.Store(2, "b")

	stats := manager.GetAllStats()

	if len(stats) != 2 {
		t.Errorf("Expected 2 caches in stats, got %d", len(stats))
	}

	if stats["cache1"]["size"] != 1 {
		t.Errorf("Expected cache1 size 1, got %d", stats["cache1"]["size"])
	}
}
