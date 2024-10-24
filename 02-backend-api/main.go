package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Response struct to hold result
type Response struct {
	Result int `json:"result"`
}

// ErrorResponse struct to hold error messages
type ErrorResponse struct {
	Error string `json:"error"`
}

var logger *slog.Logger

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *responseRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

func main() {
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Getting constants from .env
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		logger.Error("Error loading env variables")
	}

	http.HandleFunc("/add", composeRoute(addHandler))
	http.HandleFunc("/subtract", composeRoute(subtractHandler))
	http.HandleFunc("/multiply", composeRoute(multiplyHandler))
	http.HandleFunc("/divide", composeRoute(divideHandler))
	http.HandleFunc("/sum", composeRoute(sumHandler))

	port = fmt.Sprintf(":%s", port)
	startingMsg := fmt.Sprintf("Staring server on %s", port)
	logger.Info(startingMsg)
	log.Fatal(http.ListenAndServe(port, nil))
}

func composeRoute(handler http.HandlerFunc) http.HandlerFunc {
	middlewares := []Middleware{
		corsMiddleware,
		loggerMiddleware,
	}

	return applyMiddlewares(handler, middlewares...)
}

type Middleware func(http.HandlerFunc) http.HandlerFunc

func applyMiddlewares(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		rec := &responseRecorder{w, http.StatusOK}
		handler(rec, r)
	}
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func loggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		rec := w.(*responseRecorder)
		logger.Info("Request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
			slog.Duration("duration", time.Since(start)),
			slog.Int("status_code", rec.statusCode),
		)
	}
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	var numbers struct {
		Number1 int `json:"number1"`
		Number2 int `json:"number2"`
	}

	if err := json.NewDecoder(r.Body).Decode(&numbers); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	result := numbers.Number1 + numbers.Number2

	writeJSONResponse(w, http.StatusOK, Response{Result: result})
}

func subtractHandler(w http.ResponseWriter, r *http.Request) {
	var numbers struct {
		Number1 int `json:"number1"`
		Number2 int `json:"number2"`
	}

	if err := json.NewDecoder(r.Body).Decode(&numbers); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	result := numbers.Number1 - numbers.Number2

	writeJSONResponse(w, http.StatusOK, Response{Result: result})
}
func multiplyHandler(w http.ResponseWriter, r *http.Request) {
	var numbers struct {
		Number1 int `json:"number1"`
		Number2 int `json:"number2"`
	}

	if err := json.NewDecoder(r.Body).Decode(&numbers); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	result := numbers.Number1 * numbers.Number2

	writeJSONResponse(w, http.StatusOK, Response{Result: result})
}
func divideHandler(w http.ResponseWriter, r *http.Request) {
	var numbers struct {
		Number1 int `json:"number1"`
		Number2 int `json:"number2"`
	}

	if err := json.NewDecoder(r.Body).Decode(&numbers); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if numbers.Number2 == 0 {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	result := numbers.Number1 / numbers.Number2

	writeJSONResponse(w, http.StatusOK, Response{Result: result})
}
func sumHandler(w http.ResponseWriter, r *http.Request) {
	var numbers []int

	if err := json.NewDecoder(r.Body).Decode(&numbers); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var result int
	for _, num := range numbers {
		result += num
	}

	writeJSONResponse(w, http.StatusOK, Response{Result: result})
}

// Utility function to write JSON response
func writeJSONResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// logger.Error("Failed to write response", slog.Error(err))
	}
}
