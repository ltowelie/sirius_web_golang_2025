package home_work

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

const group = "/api/v1/home_work"

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) DefineRoutes(mux *http.ServeMux) {
	mux.HandleFunc(group+"/", c.handleRoot)
	mux.HandleFunc(group+"/time", c.handleTime)
	mux.HandleFunc(group+"/echo", c.handleEcho)
	mux.HandleFunc(group+"/greeting", c.handleGreeting)
}

func (c *Controller) handleRoot(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("hello world"))
	fmt.Println()
	if err != nil {
		slog.Error("Failed to write response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (c *Controller) handleTime(w http.ResponseWriter, r *http.Request) {
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

func (c *Controller) handleEcho(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		// Либо можно тут возвращать 404
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Allow", "POST")

		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
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

func (c *Controller) handleGreeting(w http.ResponseWriter, r *http.Request) {
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
