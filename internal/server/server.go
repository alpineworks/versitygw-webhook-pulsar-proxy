package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"alpineworks.io/rfc9457"
	"github.com/alpineworks/versitygw-webhook-pulsar-proxy/internal/config"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/gorilla/mux"
	"github.com/versity/versitygw/s3event"
)

type Server struct {
	config   *config.Config
	producer pulsar.Producer
	client   pulsar.Client
}

func New(cfg *config.Config) (*Server, error) {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: cfg.PulsarURL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create pulsar client: %w", err)
	}

	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: cfg.PulsarTopic,
	})
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to create pulsar producer: %w", err)
	}

	return &Server{
		config:   cfg,
		producer: producer,
		client:   client,
	}, nil
}

func (s *Server) Close() {
	if s.producer != nil {
		s.producer.Close()
	}
	if s.client != nil {
		s.client.Close()
	}
}

func (s *Server) Start(ctx context.Context) error {
	r := mux.NewRouter()
	r.HandleFunc("/webhook", s.handleWebhook).Methods("POST")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.ServerPort),
		Handler: r,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	slog.Info("starting server", slog.Int("port", s.config.ServerPort))
	return srv.ListenAndServe()
}

type WebhookResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (s *Server) handleWebhook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("failed to read request body", slog.String("error", err.Error()))
		rfc9457.NewRFC9457(
			rfc9457.WithStatus(http.StatusBadRequest),
			rfc9457.WithDetail("failed to read request body"),
			rfc9457.WithTitle("bad request"),
			rfc9457.WithInstance("/webhook"),
		).ServeHTTP(w, r)
		return
	}

	slog.Debug("received webhook request", slog.String("body", string(body)))

	var eventSchema s3event.EventSchema
	err = json.Unmarshal(body, &eventSchema)
	if err != nil {
		slog.Error("failed to unmarshal event schema", slog.String("error", err.Error()))
		rfc9457.NewRFC9457(
			rfc9457.WithStatus(http.StatusBadRequest),
			rfc9457.WithDetail("failed to unmarshal event schema"),
			rfc9457.WithTitle("bad request"),
			rfc9457.WithInstance("/webhook"),
		).ServeHTTP(w, r)
		return
	}

	eventData, err := json.Marshal(&eventSchema)
	if err != nil {
		slog.Error("failed to marshal event schema", slog.String("error", err.Error()))
		rfc9457.NewRFC9457(
			rfc9457.WithStatus(http.StatusInternalServerError),
			rfc9457.WithDetail("Failed to process the event schema"),
			rfc9457.WithTitle("Internal server error"),
			rfc9457.WithInstance("/webhook"),
		).ServeHTTP(w, r)
		return
	}

	produceCtx, cancel := context.WithTimeout(r.Context(), s.config.PulsarProduceTimeout)
	defer cancel()

	_, err = s.producer.Send(produceCtx, &pulsar.ProducerMessage{
		Payload: eventData,
	})
	if err != nil {
		slog.Error("failed to send message to pulsar", slog.String("error", err.Error()))
		rfc9457.NewRFC9457(
			rfc9457.WithStatus(http.StatusInternalServerError),
			rfc9457.WithDetail("Failed to send message to Pulsar"),
			rfc9457.WithTitle("Message delivery failed"),
			rfc9457.WithInstance("/webhook"),
		).ServeHTTP(w, r)
		return
	}

	slog.Info("event forwarded to pulsar",
		slog.Int("records", len(eventSchema.Records)),
	)

	w.WriteHeader(http.StatusCreated)
	response := WebhookResponse{
		Message: "event/s forwarded to pulsar",
		Code:    http.StatusCreated,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("failed to encode response", slog.String("error", err.Error()))
		rfc9457.NewRFC9457(
			rfc9457.WithStatus(http.StatusInternalServerError),
			rfc9457.WithDetail("Failed to encode response"),
			rfc9457.WithTitle("Internal server error"),
			rfc9457.WithInstance("/webhook"),
		).ServeHTTP(w, r)
		return
	}
	slog.Debug("response sent", slog.String("message", response.Message), slog.Int("code", response.Code))
}
