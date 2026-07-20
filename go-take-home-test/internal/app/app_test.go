package app_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"go-take-home-test/internal/app"
)

func TestIngestEndToEndReliable(t *testing.T) {
	examples := []string{
		"person_one.json",
		"person_two.json",
		"person_three.json",
	}

	for _, name := range examples {
		t.Run(name, func(t *testing.T) {
			resp := postIngest(t, ingestOpts{
				exampleFile:      name,
				mockReliable:     "1",
				queueMaxAttempts: 5,
				queueBackoff:     50 * time.Millisecond,
			})
			defer resp.Body.Close()

			assertIngestOK(t, resp)
		})
	}
}

func TestIngestEndToEndUnreliableMocksFail(t *testing.T) {
	// Postcode/email always fail; queue retries then surfaces the error.
	resp := postIngest(t, ingestOpts{
		exampleFile:      "person_one.json",
		mockReliable:     "0",
		queueMaxAttempts: 3,
		queueBackoff:     20 * time.Millisecond,
	})
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status %d when mocks always fail, got %d: %s",
			http.StatusInternalServerError, resp.StatusCode, string(body))
	}
}

func TestIngestEndToEndUnreliableMocksEventuallySucceed(t *testing.T) {
	// Default flaky mocks (~95% success + latency). Queue retries should usually win.
	resp := postIngest(t, ingestOpts{
		exampleFile:      "person_one.json",
		mockReliable:     "", // flaky path
		queueMaxAttempts: 8,
		queueBackoff:     50 * time.Millisecond,
	})
	defer resp.Body.Close()

	assertIngestOK(t, resp)
}

type ingestOpts struct {
	exampleFile      string
	mockReliable     string
	queueMaxAttempts int
	queueBackoff     time.Duration
}

func postIngest(t *testing.T, opts ingestOpts) *http.Response {
	t.Helper()

	if opts.mockReliable == "" {
		t.Setenv("HEALTH_MOCK_RELIABLE", "")
		_ = os.Unsetenv("HEALTH_MOCK_RELIABLE")
	} else {
		t.Setenv("HEALTH_MOCK_RELIABLE", opts.mockReliable)
	}

	var handler http.Handler
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}))
	t.Cleanup(srv.Close)

	dir := t.TempDir()
	e, err := app.NewWithConfig(app.Config{
		DBPath:              filepath.Join(dir, "test.db"),
		MigrationsPath:      migrationsPath(t),
		VersionFilePath:     filepath.Join(dir, "migration_version.txt"),
		BaseURL:             srv.URL,
		QueueMaxAttempts:    opts.queueMaxAttempts,
		QueueInitialBackoff: opts.queueBackoff,
	})
	if err != nil {
		t.Fatalf("failed to create app: %v", err)
	}
	handler = e

	body, err := os.ReadFile(filepath.Join(examplesPath(t), opts.exampleFile))
	if err != nil {
		t.Fatalf("failed to read example %s: %v", opts.exampleFile, err)
	}

	req, err := http.NewRequest(http.MethodPost, srv.URL+"/workers/ingest", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("ingest request failed: %v", err)
	}
	return resp
}

func assertIngestOK(t *testing.T, resp *http.Response) {
	t.Helper()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, resp.StatusCode, string(respBody))
	}

	var ingested map[string]any
	if err := json.Unmarshal(respBody, &ingested); err != nil {
		t.Fatalf("failed to decode ingest response: %v\nbody: %s", err, string(respBody))
	}

	if ingested["id"] == nil {
		t.Fatalf("expected ingested form id in response, got: %s", string(respBody))
	}
	if ingested["session_id"] == "" {
		t.Fatalf("expected session_id in response, got: %s", string(respBody))
	}
}

func migrationsPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(moduleRoot(t), "internal", "migrations")
}

func examplesPath(t *testing.T) string {
	t.Helper()
	path := filepath.Join(moduleRoot(t), "..", "take-home-test", "src", "forms", "examples")
	if _, err := os.Stat(path); err == nil {
		return path
	}
	return filepath.Join(moduleRoot(t), "forms", "examples")
}

func moduleRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to resolve caller path")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}
