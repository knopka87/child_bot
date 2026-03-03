package telegram

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// ShutdownManager управляет graceful shutdown всех горутин
type ShutdownManager struct {
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	isShutdown atomic.Bool
}

var (
	shutdownManager     *ShutdownManager
	shutdownManagerOnce sync.Once
)

// GetShutdownManager возвращает singleton менеджера shutdown
func GetShutdownManager() *ShutdownManager {
	shutdownManagerOnce.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		shutdownManager = &ShutdownManager{
			ctx:    ctx,
			cancel: cancel,
		}
	})
	return shutdownManager
}

// Context возвращает контекст, который отменяется при shutdown
func (m *ShutdownManager) Context() context.Context {
	return m.ctx
}

// Done возвращает канал, который закрывается при shutdown
func (m *ShutdownManager) Done() <-chan struct{} {
	return m.ctx.Done()
}

// IsShutdown проверяет, был ли инициирован shutdown
func (m *ShutdownManager) IsShutdown() bool {
	return m.isShutdown.Load()
}

// TrackGoroutine регистрирует горутину для отслеживания
// Возвращает функцию, которую нужно вызвать при завершении горутины (defer)
func (m *ShutdownManager) TrackGoroutine() func() {
	m.wg.Add(1)
	return func() {
		m.wg.Done()
	}
}

// Shutdown инициирует graceful shutdown
// timeout — максимальное время ожидания завершения всех горутин
func (m *ShutdownManager) Shutdown(timeout time.Duration) error {
	if m.isShutdown.Swap(true) {
		// Уже был shutdown
		return nil
	}

	// Отменяем контекст — все горутины получат сигнал
	m.cancel()

	// Ждём завершения всех горутин с таймаутом
	done := make(chan struct{})
	go func() {
		m.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return context.DeadlineExceeded
	}
}

// ActiveGoroutines возвращает примерное количество активных отслеживаемых горутин
// (для мониторинга, не точное значение)
func (m *ShutdownManager) ActiveGoroutines() int {
	// WaitGroup не предоставляет метод для получения счётчика,
	// но мы можем добавить отдельный счётчик если нужно
	return -1 // -1 означает "неизвестно"
}

// InitShutdownManager инициализирует менеджер shutdown
// Вызывать один раз при старте приложения
func InitShutdownManager() *ShutdownManager {
	return GetShutdownManager()
}

// GracefulShutdown выполняет graceful shutdown всех компонентов
// timeout — максимальное время ожидания
func GracefulShutdown(timeout time.Duration) error {
	// Останавливаем кэш-менеджер
	StopCacheCleanup()

	// Останавливаем все отслеживаемые горутины
	return GetShutdownManager().Shutdown(timeout)
}