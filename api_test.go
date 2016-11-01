package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
)

func TestWithFixtures(t *testing.T) {
	var h HTTP
	m := newMock()
	defer m.resetClient()
	m.fileBody("testdata/metadata.json", t)
	if _, err := h.VotingRatesMetadata(); err != nil {
		t.Fatal(err)
	}
	m.fileBody("testdata/query.json", t)
	if _, err := h.VotingRatesQuery(); err != nil {
		t.Fatal(err)
	}
}

func TestWithEmptyBody(t *testing.T) {
	var h HTTP
	m := newMock()
	defer m.resetClient()
	if _, err := h.VotingRatesMetadata(); err == nil {
		t.Fatal("Expected error for empty body")
	}
}

func TestWithUnsuccessfulStatusCode(t *testing.T) {
	var h HTTP
	m := newMock()
	defer m.resetClient()
	m.resp.StatusCode = http.StatusNotImplemented
	_, err := h.VotingRatesMetadata()
	if err == nil {
		t.Fatal("Expected error unsucessful status code")
	}
	want := fmt.Sprintf("HTTP server responded with status: %d", http.StatusNotImplemented)
	if err.Error() != want {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func TestJSONWithBOM(t *testing.T) {
	var foo struct {
		Val string `json:"foo"`
	}
	var b bytes.Buffer
	if _, err := b.Write([]byte{0xef, 0xbb, 0xbf}); err != nil {
		t.Fatalf("error writing BOM: %s", err)
	}
	if _, err := b.WriteString(`{"foo":"bar"}`); err != nil {
		t.Fatalf("error writing JSON: %s", err)
	}
	if err := readJSON(&b, &foo); err != nil {
		log.Fatal(err)
	}
}

type mock struct {
	resp *http.Response
}

func newMock() *mock {
	r := http.Response{
		StatusCode: http.StatusOK,
		Proto:      "HTTP/1.0",
		ProtoMajor: 1,
		ProtoMinor: 0,
		Header:     make(http.Header),
	}
	m := &mock{resp: &r}
	m.body("")
	http.DefaultClient.Transport = m
	return m
}

func (m *mock) resetClient() {
	http.DefaultClient.Transport = nil
}

func (m *mock) body(s string) {
	m.resp.Body = ioutil.NopCloser(strings.NewReader(s))
}

func (m *mock) fileBody(fn string, t *testing.T) {
	if data, err := ioutil.ReadFile(fn); err != nil {
		t.Fatalf("Could not read file %q: %s", fn, err)
	} else {
		m.body(string(data))
	}
}

func (m *mock) RoundTrip(req *http.Request) (*http.Response, error) {
	m.resp.Request = req
	return m.resp, nil
}
