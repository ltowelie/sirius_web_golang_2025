package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

type server struct {
	http http.Server
}

func main() {
	ctxSignal, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer cancel()

	loggerSetup()
	a, err := addr()
	if err != nil {
		slog.Error("Failed to parse addr from envs", "error", err)
		os.Exit(1)
	}

	s := newServer(a)

	g, ctxErrGr := errgroup.WithContext(ctxSignal)
	g.Go(func() error {
		slog.Debug("Staring http server")
		err = s.http.ListenAndServe()
		if err != nil {
			return err
		}

		return nil
	})

	g.Go(func() error {
		select {
		case <-ctxErrGr.Done():
			err = ctxErrGr.Err()
			slog.Debug("err group chan done", "error", err)
		case <-ctxSignal.Done():
			err = ctxSignal.Err()
			slog.Debug("signal chan ctxSignal done", "error", err)
		}

		slog.Info("Begin shutdown process gracefully")
		// Даем 5 секунд на завершение работы сервера - заверешить все запросы, а новые запросы сервер
		// после получения сигнала отбрасывает
		// Возникает вопрос - зачем нам в новом контексте использовать context.Backgroung(),
		// если уже есть контексты ctxErrGr и ctxSignal. Новый контекст нужен, так как эти два контекста
		// могут быть завершены (один из них точно завершен - мы же вниз по коду сюда спустились).
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		errSh := s.http.Shutdown(shutdownCtx)
		if errSh != nil {
			slog.Error("Failed to shutdown http server", "error", errSh)
		}

		return err
	})

	if err = g.Wait(); err != nil {
		slog.Error("Exit reason", "error", err)
	}
}

func loggerSetup() {
	lh := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(lh)
}

func addr() (string, error) {
	host := os.Getenv("HOST")
	if host == "" {
		return "", errors.New("env HOST is empty")

	}
	port := os.Getenv("PORT")
	if port == "" {
		return "", errors.New("env PORT is empty")
	}
	a := fmt.Sprintf("%s:%s", host, port)

	return a, nil
}

func newServer(addr string) *server {
	s := &server{http: http.Server{}}
	s.http.Addr = addr
	s.setupRoutes()

	return s
}

func (c *server) setupRoutes() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", c.handleRoot)
	mux.HandleFunc("/time", c.handleTime)
	mux.HandleFunc("/echo", c.handleEcho)
	mux.HandleFunc("/greeting", c.handleGreeting)

	c.http.Handler = mux
}

func (c *server) handleRoot(w http.ResponseWriter, r *http.Request) {
	// Эмуляция долгого запроса для демонстрации работы graceful shutdown
	time.Sleep(3 * time.Second)
	_, err := w.Write([]byte("hello world"))
	if err != nil {
		slog.Error("Failed to write response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (c *server) handleTime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	response := map[string]string{"time": currentTime}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error("Failed to encode response", "err", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (c *server) handleEcho(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)

		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Failed to read body", "err", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}
	defer r.Body.Close()

	_, err = w.Write(body)
	if err != nil {
		slog.Error("Failed to write response", "err", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (c *server) handleGreeting(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	w.Header().Set("Content-Type", "application/json")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(map[string]string{"error": "Укажите параметр name"})
		if err != nil {
			slog.Error("Failed to encode response", "err", err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	}

	greeting := fmt.Sprintf("Привет, %s!", name)
	err := json.NewEncoder(w).Encode(map[string]string{"greeting": greeting})
	if err != nil {
		slog.Error("Failed to write response", "err", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
