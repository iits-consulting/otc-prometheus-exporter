package otcdoc

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchBytesReturnsBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("hello"))
	}))
	defer srv.Close()

	body, err := fetchBytes(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(body) != "hello" {
		t.Errorf("expected body 'hello', got %q", string(body))
	}
}

func TestFetchBytesErrorOnNon200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	_, err := fetchBytes(srv.URL)
	if err == nil {
		t.Fatal("expected error on 404, got nil")
	}
}

func TestFetchDocumentationSourceEmptyURLReturnsEmpty(t *testing.T) {
	ds := DocumentationSource{Namespace: "ecs", GithubRawUrl: ""}
	page, err := FetchDocumentationSource(ds)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.Namespace != "" || len(page.Metrics) != 0 {
		t.Errorf("expected empty page for empty URL, got %+v", page)
	}
}

func TestFetchDocumentationSourceParsesRST(t *testing.T) {
	rst := []byte(`
Namespace
---------

SYS.NAT

+--------------------+------------------+----------+
| Metric ID          | Metric Name      | Unit     |
+====================+==================+==========+
| nat001_bytes_in    | Inbound Bytes    | Bytes/s  |
+--------------------+------------------+----------+
`)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(rst)
	}))
	defer srv.Close()

	ds := DocumentationSource{Namespace: "nat", GithubRawUrl: srv.URL}
	page, err := FetchDocumentationSource(ds)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.Namespace != "SYS.NAT" {
		t.Errorf("expected namespace SYS.NAT, got %q", page.Namespace)
	}
	if len(page.Metrics) != 1 {
		t.Errorf("expected 1 metric, got %d", len(page.Metrics))
	}
}

func TestFetchDocumentationSourceErrorOnHTTPFailure(t *testing.T) {
	ds := DocumentationSource{Namespace: "nat", GithubRawUrl: "http://127.0.0.1:0/nonexistent"}
	_, err := FetchDocumentationSource(ds)
	if err == nil {
		t.Fatal("expected error on unreachable server, got nil")
	}
}
