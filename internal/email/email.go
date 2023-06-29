// Package email provides a simple email client.
package email

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

var defaultClient = http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          1,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

var propagator = otel.GetTextMapPropagator()

// EmailNotifier is an email notifier service.
type EmailNotifier struct {
	url    string
	client *http.Client
}

// NewEmailNotifier creates a new email notifier service.
func NewEmailNotifier(url string) *EmailNotifier {
	return &EmailNotifier{
		url:    url,
		client: &defaultClient,
	}
}

// Notify sends an email notification.
func (s *EmailNotifier) Notify(ctx context.Context, userID string) error {
	url := fmt.Sprintf("%s/users/%s/notify", s.url, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("creating http request: %w", err)
	}

	// Inject trace context
	carrier := propagation.HeaderCarrier{}
	propagator.Inject(ctx, carrier)

	req.Header = http.Header(carrier)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}
	return nil
}
