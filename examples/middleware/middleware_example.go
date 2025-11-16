package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"

	"runtime/trace"

	"github.com/google/uuid"
)

type config struct {
	fr     *trace.FlightRecorder
	wg     *sync.WaitGroup
	logger *slog.Logger
}
type contextKey string

const (
	requestIDKey contextKey = "request_id" // type save context keys.
)

// Chain slice of middleware funcs to be applied using slice.Backward(middleware is applied starting from last to first element).
type Chain []func(http.Handler) http.Handler

func (c Chain) thenFunc(h http.HandlerFunc) http.Handler {
	return c.then(h)
}

func (c Chain) then(h http.Handler) http.Handler {
	for _, mw := range slices.Backward(c) {
		h = mw(h)
	}
	return h
}

func GenerateUUID() string {
	return uuid.New().String()
}

// LoggingMiddleware initial middleware starts a log and assigns a request id.
func (cfg *config) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := GenerateUUID()
		fmt.Println("Log mid")

		ctx := context.WithValue(r.Context(), requestIDKey, reqID)
		r = r.WithContext(ctx)

		cfg.logger.Info("request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("request_id", reqID),
		)

		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Auth mid")
		token := r.Header.Get("Authorization")
		if token == "" {
			fmt.Println("unauth access")
			w.WriteHeader(401)
			return
			// TODO: add actual logic duhh
		}
		next.ServeHTTP(w, r)
	})
}

// TraceMiddleware saves a trace file to traces dir for each req that exceeds threshold.
func (cfg *config) TraceMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Trace mid")
		start := time.Now()

		next.ServeHTTP(w, r)

		diff := time.Since(start)
		threshold := 300 * time.Millisecond // trace will be written if time exceeds threshold.

		if diff > threshold {

			reqID, ok := r.Context().Value(requestIDKey).(string)
			if !ok {
				reqID = GenerateUUID() // fallback if not in context.
			}

			if err := writeTrace(cfg.fr, reqID); err != nil {
				cfg.logger.Error("failed to write trace",
					slog.String("request_id", reqID),
					slog.String("error", err.Error()),
				)
				return
			}

			cfg.logger.Warn("trace written",
				slog.String("request_id", reqID),
				slog.String("req_addr", r.RemoteAddr),
				slog.Duration("duration", diff),
			)
		}
	})
}
func CrossOriginProtectMiddleware(next http.Handler) http.Handler {
	cop := http.NewCrossOriginProtection()

	cop.AddTrustedOrigin("https://foo.example.com")

	cop.SetDenyHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("CSRF check failed"))
	}))

	return cop.Handler(next)
}

// writeTrace writes traces to disk.
// - file are in the form reqID.trace .
func writeTrace(fr *trace.FlightRecorder, reqID string) error {
	if !fr.Enabled() {
		return fmt.Errorf("flight recorder not enabled")
	}

	tracesDir := "traces"
	if err := os.MkdirAll(tracesDir, 0755); err != nil {
		return fmt.Errorf("failed to create traces directory: %w", err)
	}

	filename := filepath.Join(tracesDir, fmt.Sprintf("trace-%s.out", reqID))
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create trace file: %w", err)
	}

	defer file.Close()

	_, err = fr.WriteTo(file)
	if err != nil {
		return fmt.Errorf("failed to write trace data: %w", err)
	}

	return nil
}

func main() {
	logg := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})

	frcfg := trace.FlightRecorderConfig{
		MinAge:   time.Second,
		MaxBytes: 3 << 20, // 3MB buffer
	}

	fr := trace.NewFlightRecorder(frcfg)
	if err := fr.Start(); err != nil {
		log.Fatalf("Unable to start trace flight recorder: %v", err.Error())
	}
	defer fr.Stop()

	c := &config{
		fr:     fr,
		wg:     &sync.WaitGroup{},
		logger: slog.New(logg),
	}

	helloHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 2)
		fmt.Fprintln(w, "Hello, World!")
	})

	BaseChain := Chain{
		CrossOriginProtectMiddleware,
		c.TraceMiddleware,
		c.LoggingMiddleware,
	}
	AuthChain := append(BaseChain, AuthMiddleware)
	// hello sims slow req triggering flight trace write to disk.

	http.Handle("/hello", BaseChain.thenFunc(helloHandler))
	http.Handle("/authorized", AuthChain.thenFunc(helloHandler))

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
