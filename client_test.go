package endoflife

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return NewClient(WithBaseURL(srv.URL))
}

// writeBody writes b to w, failing the test on error.
func writeBody(t *testing.T, w http.ResponseWriter, b string) {
	t.Helper()
	if _, err := w.Write([]byte(b)); err != nil {
		t.Fatalf("writing response body: %v", err)
	}
}

func TestProduct(t *testing.T) {
	const body = `{
		"schema_version": "1.2.1",
		"generated_at": "2023-03-01T14:05:52+01:00",
		"last_modified": "2023-03-01T14:05:52+01:00",
		"result": {
			"name": "ubuntu",
			"label": "Ubuntu",
			"aliases": ["ubuntu-linux"],
			"category": "os",
			"tags": ["canonical", "os"],
			"versionCommand": "lsb_release --release",
			"identifiers": [{"id": "cpe:/o:canonical:ubuntu_linux", "type": "cpe"}],
			"labels": {"eoas": "Hardware & Maintenance", "discontinued": null, "eol": "Maintenance", "eoes": "Extended Security Maintenance"},
			"links": {"icon": "https://simpleicons.org/icons/ubuntu.svg", "html": "https://endoflife.date/ubuntu", "releasePolicy": null},
			"releases": [{
				"name": "22.04",
				"codename": "Jammy Jellyfish",
				"label": "22.04 'Jammy Jellyfish' (LTS)",
				"releaseDate": "2022-04-21",
				"isLts": true,
				"ltsFrom": "2022-04-21",
				"isEoas": false,
				"eoasFrom": "2024-09-30",
				"isEol": false,
				"eolFrom": "2027-04-01",
				"isEoes": true,
				"eoesFrom": "2032-04-09",
				"isMaintained": true,
				"latest": {"name": "22.04.2", "date": "2022-04-21", "link": "https://wiki.ubuntu.com/"},
				"custom": {"chromeVersion": "M136", "nodeVersion": null}
			}]
		}
	}`

	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/products/ubuntu" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Accept"); got != "application/json" {
			t.Errorf("Accept = %q", got)
		}
		writeBody(t, w, body)
	})

	resp, err := c.Product(context.Background(), "ubuntu")
	if err != nil {
		t.Fatal(err)
	}

	p := resp.Result
	if p.Name != "ubuntu" || p.Label != "Ubuntu" {
		t.Errorf("got name=%q label=%q", p.Name, p.Label)
	}
	if p.VersionCommand == nil || *p.VersionCommand != "lsb_release --release" {
		t.Errorf("versionCommand = %v", p.VersionCommand)
	}
	if p.Labels.Discontinued != nil {
		t.Errorf("discontinued should be nil, got %v", *p.Labels.Discontinued)
	}
	if len(p.Releases) != 1 {
		t.Fatalf("got %d releases", len(p.Releases))
	}

	rel := p.Releases[0]
	if !rel.IsLts || rel.Codename == nil || *rel.Codename != "Jammy Jellyfish" {
		t.Errorf("release lts/codename wrong: %+v", rel)
	}
	want := NewDate(2022, time.April, 21)
	if !rel.ReleaseDate.Equal(want.Time) {
		t.Errorf("releaseDate = %s, want %s", rel.ReleaseDate, want)
	}
	if rel.IsEoas == nil || *rel.IsEoas != false {
		t.Errorf("isEoas = %v", rel.IsEoas)
	}
	if rel.Latest == nil || rel.Latest.Name != "22.04.2" {
		t.Errorf("latest = %+v", rel.Latest)
	}
	if rel.Latest.Date == nil || rel.Latest.Date.String() != "2022-04-21" {
		t.Errorf("latest.date = %v", rel.Latest.Date)
	}
	// custom: nodeVersion is JSON null -> nil entry.
	if v, ok := rel.Custom["nodeVersion"]; !ok || v != nil {
		t.Errorf("custom nodeVersion = %v (ok=%v)", v, ok)
	}
	if v := rel.Custom["chromeVersion"]; v == nil || *v != "M136" {
		t.Errorf("custom chromeVersion = %v", v)
	}
}

func TestProducts(t *testing.T) {
	const body = `{
		"schema_version": "1.2.1",
		"generated_at": "2023-03-01T14:05:52+01:00",
		"total": 1,
		"result": [{"name": "ubuntu", "label": "Ubuntu", "aliases": [], "category": "os", "tags": ["os"], "uri": "https://endoflife.date/api/v1/products/ubuntu/"}]
	}`
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeBody(t, w, body)
	})
	resp, err := c.Products(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if resp.Total != 1 || len(resp.Result) != 1 || resp.Result[0].Name != "ubuntu" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestNotFound(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		writeBody(t, w, "<html>not found</html>")
	})
	_, err := c.Product(context.Background(), "does-not-exist")
	if err == nil {
		t.Fatal("expected error")
	}
	if !IsNotFound(err) {
		t.Errorf("IsNotFound = false for %v", err)
	}
	if IsTooManyRequests(err) {
		t.Error("IsTooManyRequests should be false")
	}
}

func TestTooManyRequests(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "120")
		w.WriteHeader(http.StatusTooManyRequests)
	})
	_, err := c.Products(context.Background())
	if !IsTooManyRequests(err) {
		t.Fatalf("IsTooManyRequests = false for %v", err)
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.RetryAfter != 120*time.Second {
		t.Errorf("RetryAfter = %v", err)
	}
}

func TestRedirectFollowed(t *testing.T) {
	const body = `{"schema_version":"1.2.1","generated_at":"2023-03-01T14:05:52+01:00","total":0,"result":[]}`
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/tags/old" {
			http.Redirect(w, r, "/tags/new", http.StatusMovedPermanently)
			return
		}
		writeBody(t, w, body)
	})
	resp, err := c.ProductsByTag(context.Background(), "old")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Total != 0 {
		t.Errorf("total = %d", resp.Total)
	}
}

func TestPathEscaping(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/products/red%20hat/releases/9" {
			t.Errorf("escaped path = %s", r.URL.EscapedPath())
		}
		writeBody(t, w, `{"schema_version":"1","generated_at":"2023-03-01T14:05:52+01:00","result":{"name":"9","codename":null,"label":"9","releaseDate":"2022-05-17","isLts":false,"ltsFrom":null,"isEol":false,"eolFrom":null,"isMaintained":true,"latest":null}}`)
	})
	_, err := c.ProductRelease(context.Background(), "red hat", "9")
	if err != nil {
		t.Fatal(err)
	}
}
