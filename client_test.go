package tremendous

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateOrganization(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/organizations" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id":"org_123"}`))
	}))
	defer s.Close()

	client := &Client{httpClient: s.Client(), endpoint: s.URL}
	org, err := client.CreateOrganization(context.Background(), &Organization{})
	if err != nil || org == nil {
		t.Errorf("expected valid response, got err: %v", err)
	}
}

func TestListOrganizations(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/organizations" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[]`))
	}))
	defer s.Close()

	client := &Client{httpClient: s.Client(), endpoint: s.URL}
	orgs, err := client.ListOrganizations(context.Background())
	if err != nil || orgs == nil {
		t.Errorf("expected valid response, got err: %v", err)
	}
}

func TestRetrieveOrganization(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/organizations/org_123" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"org_123"}`))
	}))
	defer s.Close()

	client := &Client{httpClient: s.Client(), endpoint: s.URL}
	org, err := client.RetrieveOrganization(context.Background(), "org_123")
	if err != nil || org == nil {
		t.Errorf("expected valid response, got err: %v", err)
	}
}

func TestCreateWebhook(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/webhooks" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id":"webhook_123"}`))
	}))
	defer s.Close()

	client := &Client{httpClient: s.Client(), endpoint: s.URL}
	wh, err := client.CreateWebhook(context.Background(), "https://example.com")
	if err != nil || wh == nil {
		t.Errorf("expected valid response, got err: %v", err)
	}
}

func TestSimulateWebhook(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/webhooks/webhook_123" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer s.Close()

	client := &Client{httpClient: s.Client(), endpoint: s.URL}
	err := client.SimulateWebhook(context.Background(), "webhook_123", "test_event")
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}
