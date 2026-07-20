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
			env := newTestEnv(t, "1", 5, 50*time.Millisecond)
			resp := env.postJSON(t, "/ingest", readExample(t, name))
			defer resp.Body.Close()
			assertIngestOK(t, resp)
		})
	}
}

func TestIngestEndToEndUnreliableMocksFail(t *testing.T) {
	env := newTestEnv(t, "0", 3, 20*time.Millisecond)
	resp := env.postJSON(t, "/ingest", readExample(t, "person_one.json"))
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status %d when mocks always fail, got %d: %s",
			http.StatusInternalServerError, resp.StatusCode, string(body))
	}
}

func TestIngestEndToEndUnreliableMocksEventuallySucceed(t *testing.T) {
	env := newTestEnv(t, "", 8, 50*time.Millisecond)
	resp := env.postJSON(t, "/ingest", readExample(t, "person_one.json"))
	defer resp.Body.Close()
	assertIngestOK(t, resp)
}

func TestIngestDeduplicatesByFingerprint(t *testing.T) {
	env := newTestEnv(t, "1", 5, 50*time.Millisecond)
	payload := readExample(t, "person_one.json")

	first := env.postJSON(t, "/ingest", payload)
	defer first.Body.Close()
	firstBody := assertIngestOK(t, first)

	second := env.postJSON(t, "/ingest", payload)
	defer second.Body.Close()
	secondBody := assertIngestOK(t, second)

	if firstBody["id"] != secondBody["id"] {
		t.Fatalf("expected identical payload to reuse id %v, got %v", firstBody["id"], secondBody["id"])
	}
	if firstBody["fingerprint"] == nil || firstBody["fingerprint"] == "" {
		t.Fatalf("expected fingerprint on ingested form, got: %v", firstBody)
	}
	if firstBody["fingerprint"] != secondBody["fingerprint"] {
		t.Fatalf("expected same fingerprint, got %v and %v", firstBody["fingerprint"], secondBody["fingerprint"])
	}
}

func TestRetryAfterFailedTransform(t *testing.T) {
	env := newTestEnv(t, "0", 2, 20*time.Millisecond)
	payload := readExample(t, "person_two.json")

	failResp := env.postJSON(t, "/ingest", payload)
	failBody, _ := io.ReadAll(failResp.Body)
	failResp.Body.Close()
	if failResp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected initial ingest to fail, got %d: %s", failResp.StatusCode, string(failBody))
	}

	var parsed map[string]any
	_ = json.Unmarshal(payload, &parsed)
	sessionID, _ := parsed["session_id"].(string)

	t.Setenv("HEALTH_MOCK_RELIABLE", "1")
	retryPayload, _ := json.Marshal(map[string]any{"session_id": sessionID})
	retryResp := env.postJSON(t, "/retry", retryPayload)
	defer retryResp.Body.Close()

	retryBody, _ := io.ReadAll(retryResp.Body)
	if retryResp.StatusCode != http.StatusOK {
		t.Fatalf("expected retry status 200, got %d: %s", retryResp.StatusCode, string(retryBody))
	}
}

type testEnv struct {
	server *httptest.Server
	client *http.Client
}

func newTestEnv(t *testing.T, mockReliable string, maxAttempts int, backoff time.Duration) *testEnv {
	t.Helper()

	if mockReliable == "" {
		_ = os.Unsetenv("HEALTH_MOCK_RELIABLE")
	} else {
		t.Setenv("HEALTH_MOCK_RELIABLE", mockReliable)
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
		QueueMaxAttempts:    maxAttempts,
		QueueInitialBackoff: backoff,
	})
	if err != nil {
		t.Fatalf("failed to create app: %v", err)
	}
	handler = e

	return &testEnv{
		server: srv,
		client: &http.Client{Timeout: 60 * time.Second},
	}
}

func (env *testEnv) postJSON(t *testing.T, path string, body []byte) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, env.server.URL+path, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := env.client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	return resp
}

func assertIngestOK(t *testing.T, resp *http.Response) map[string]any {
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
	return ingested
}

func readExample(t *testing.T, name string) []byte {
	t.Helper()
	body, err := os.ReadFile(filepath.Join(examplesPath(t), name))
	if err != nil {
		t.Fatalf("failed to read example %s: %v", name, err)
	}
	return body
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
