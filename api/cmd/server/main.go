package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"child-bot/api/internal/api/router"
	"child-bot/api/internal/config"
	"child-bot/api/internal/llm"
	"child-bot/api/internal/store"

	_ "github.com/lib/pq"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func run() error {
	// Загрузка конфигурации
	cfg := config.Load()

	// Подключение к PostgreSQL
	db, err := sql.Open("postgres", mustEnv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Проверка соединения с БД
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}
	log.Println("✓ Database connection established")

	// Инициализация зависимостей
	st := store.NewStore(db)
	llmClient := llm.NewClient(cfg.LLMServerURL)

	// Создание роутера
	r := router.New(&router.Dependencies{
		Store:      st,
		LLMClient:  llmClient,
		Config:     cfg,
		DefaultLLM: cfg.DefaultLLM,
	})

	// HTTP сервер
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 5 * time.Minute, // Увеличенный timeout для долгих LLM запросов
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("REST API server listening on :%s", cfg.Port)
		serverErrors <- srv.ListenAndServe()
	}()

	// Ожидание сигнала завершения
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Printf("received signal %v, starting graceful shutdown", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			// Принудительное закрытие
			if closeErr := srv.Close(); closeErr != nil {
				log.Printf("force close error: %v", closeErr)
			}
			return fmt.Errorf("graceful shutdown failed: %w", err)
		}

		log.Println("server stopped gracefully")
	}

	return nil
}

func mustEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("missing required environment variable: %s", key)
	}
	return value
}
