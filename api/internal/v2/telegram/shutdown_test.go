package telegram

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

// newTestShutdownManager создаёт менеджер для тестов
func newTestShutdownManager() *ShutdownManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &ShutdownManager{
		ctx:    ctx,
		cancel: cancel,
	}
}

func TestShutdownManager_Basic(t *testing.T) {
	manager := newTestShutdownManager()

	// Проверяем начальное состояние
	if manager.IsShutdown() {
		t.Error("Expected IsShutdown to be false initially")
	}

	// Запускаем тестовую горутину
	var goroutineFinished atomic.Bool
	done := manager.TrackGoroutine()

	go func() {
		defer done()
		select {
		case <-manager.Done():
			goroutineFinished.Store(true)
		case <-time.After(5 * time.Second):
			t.Error("Goroutine did not receive shutdown signal")
		}
	}()

	// Даём горутине время запуститься
	time.Sleep(10 * time.Millisecond)

	// Инициируем shutdown
	err := manager.Shutdown(1 * time.Second)
	if err != nil {
		t.Errorf("Shutdown returned error: %v", err)
	}

	// Проверяем что горутина завершилась
	time.Sleep(50 * time.Millisecond)
	if !goroutineFinished.Load() {
		t.Error("Goroutine did not finish after shutdown")
	}

	// Проверяем состояние после shutdown
	if !manager.IsShutdown() {
		t.Error("Expected IsShutdown to be true after shutdown")
	}
}

func TestShutdownManager_Timeout(t *testing.T) {
	manager := newTestShutdownManager()

	// Запускаем горутину которая не реагирует на shutdown
	done := manager.TrackGoroutine()
	go func() {
		defer done()
		// Игнорируем shutdown сигнал и спим долго
		time.Sleep(10 * time.Second)
	}()

	// Даём горутине время запуститься
	time.Sleep(10 * time.Millisecond)

	// Shutdown с коротким таймаутом должен вернуть ошибку
	err := manager.Shutdown(50 * time.Millisecond)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}

func TestShutdownManager_MultipleGoroutines(t *testing.T) {
	manager := newTestShutdownManager()

	var finishedCount atomic.Int32
	const goroutineCount = 5

	// Запускаем несколько горутин
	for i := 0; i < goroutineCount; i++ {
		done := manager.TrackGoroutine()
		go func() {
			defer done()
			<-manager.Done()
			finishedCount.Add(1)
		}()
	}

	// Даём горутинам время запуститься
	time.Sleep(10 * time.Millisecond)

	// Shutdown
	err := manager.Shutdown(1 * time.Second)
	if err != nil {
		t.Errorf("Shutdown returned error: %v", err)
	}

	// Все горутины должны завершиться
	time.Sleep(50 * time.Millisecond)
	if finishedCount.Load() != goroutineCount {
		t.Errorf("Expected %d goroutines to finish, got %d", goroutineCount, finishedCount.Load())
	}
}

func TestShutdownManager_DoubleShutdown(t *testing.T) {
	manager := newTestShutdownManager()

	// Первый shutdown
	err1 := manager.Shutdown(1 * time.Second)
	if err1 != nil {
		t.Errorf("First shutdown returned error: %v", err1)
	}

	// Второй shutdown должен быть no-op
	err2 := manager.Shutdown(1 * time.Second)
	if err2 != nil {
		t.Errorf("Second shutdown returned error: %v", err2)
	}
}