package http_test

import (
	"bytes"
	"io"
	stdHttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/andybalholm/brotli"
	"github.com/ayn2op/discordo/internal/http"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zstd"
	"github.com/stretchr/testify/assert"
)

// Test plain, non-encoded response
func TestTransport_NoEncoding(t *testing.T) {
	msg := "message_no_encoding"
	ts := httptest.NewServer(stdHttp.HandlerFunc(func(writer stdHttp.ResponseWriter, request *stdHttp.Request) {
		writer.Write([]byte(msg))
	}))
	defer ts.Close()

	tr := http.NewTransport()
	assert.NotNil(t, tr, "Transport should not be nil")

	client := &stdHttp.Client{Transport: tr}
	response, err := client.Get(ts.URL)
	assert.NoError(t, err)

	body, _ := io.ReadAll(response.Body)
	assert.Equal(t, msg, string(body))
	assert.Empty(t, response.Header.Get("Content-Encoding"), "Content-Encoding header should be empty")
}

// Test gzip-encoded response
func TestTransport_GzipEncoding(t *testing.T) {
	msg := []byte("message_gzip_encoding")
	ts := httptest.NewServer(stdHttp.HandlerFunc(func(writer stdHttp.ResponseWriter, request *stdHttp.Request) {
		writer.Header().Set("Content-Encoding", "gzip")
		var buf bytes.Buffer
		gWriter := gzip.NewWriter(&buf)
		gWriter.Write(msg)
		gWriter.Close()
		writer.Write(buf.Bytes())
	}))
	defer ts.Close()

	tr := http.NewTransport()
	assert.NotNil(t, tr, "Transport should not be nil")

	client := &stdHttp.Client{Transport: tr}

	response, err := client.Get(ts.URL)
	assert.NoError(t, err)

	body, _ := io.ReadAll(response.Body)
	assert.Equal(t, msg, body)
	assert.Empty(t, response.Header.Get("Content-Encoding"), "Content-Encoding header should be removed after decoding")
}

// Test zstd-encoded response
func TestTransport_ZstdEncoding(t *testing.T) {
	msg := []byte("message_zstd_encoding")
	ts := httptest.NewServer(stdHttp.HandlerFunc(func(writer stdHttp.ResponseWriter, request *stdHttp.Request) {
		writer.Header().Set("Content-Encoding", "zstd")
		var buf bytes.Buffer
		zw, err := zstd.NewWriter(&buf)
		assert.NoError(t, err)
		zw.Write(msg)
		zw.Close()
		writer.Write(buf.Bytes())
	}))
	defer ts.Close()

	tr := http.NewTransport()
	assert.NotNil(t, tr, "Transport should not be nil")

	client := &stdHttp.Client{Transport: tr}

	response, err := client.Get(ts.URL)
	assert.NoError(t, err)

	body, _ := io.ReadAll(response.Body)
	assert.Equal(t, msg, body)
	assert.Empty(t, response.Header.Get("Content-Encoding"), "Content-Encoding header should be removed after decoding")
}

// Test brotli-encoded response
func TestTransport_BrotliEncoding(t *testing.T) {
	msg := []byte("message_brotli_encoding")
	ts := httptest.NewServer(stdHttp.HandlerFunc(func(writer stdHttp.ResponseWriter, request *stdHttp.Request) {
		writer.Header().Set("Content-Encoding", "br")
		var buf bytes.Buffer
		bw := brotli.NewWriter(&buf)
		bw.Write(msg)
		bw.Close()
		writer.Write(buf.Bytes())
	}))
	defer ts.Close()

	tr := http.NewTransport()
	assert.NotNil(t, tr, "Transport should not be nil")

	client := &stdHttp.Client{Transport: tr}

	response, err := client.Get(ts.URL)
	assert.NoError(t, err)

	body, _ := io.ReadAll(response.Body)
	assert.Equal(t, msg, body)
	assert.Empty(t, response.Header.Get("Content-Encoding"), "Content-Encoding header should be removed after decoding")
}

// Test unfamiliar encoding (should return body as-is)
func TestTransform_UnfamiliarEncoding(t *testing.T) {
	msg := []byte("message_unfamiliar_encoding")
	ts := httptest.NewServer(stdHttp.HandlerFunc(func(writer stdHttp.ResponseWriter, request *stdHttp.Request) {
		writer.Header().Set("Content-Encoding", "bfr")
		writer.Write(msg)
	}))
	defer ts.Close()

	tr := http.NewTransport()
	assert.NotNil(t, tr, "Transport should not be nil")

	client := &stdHttp.Client{Transport: tr}

	response, err := client.Get(ts.URL)
	assert.NoError(t, err)

	body, _ := io.ReadAll(response.Body)
	assert.Equal(t, msg, body)
	assert.Equal(t, "bfr", response.Header.Get("Content-Encoding"))
}

func TestTransport_RawPlainText(t *testing.T) {
	msg_raw := []byte{0x12, 0xAB, 0xCD, 0x34, 0x5F}
	ts := httptest.NewServer(stdHttp.HandlerFunc(func(writer stdHttp.ResponseWriter, response *stdHttp.Request) {
		writer.Write(msg_raw)
	}))
	defer ts.Close()

	tr := http.NewTransport()
	assert.NotNil(t, tr, "Transport should not be nil")

	client := &stdHttp.Client{Transport: tr}
	response, err := client.Get(ts.URL)
	assert.NoError(t, err, "No error should be returned for raw plain text data")

	body, _ := io.ReadAll(response.Body)
	assert.Equal(t, msg_raw, body, "Transport should return raw bytes unchanged")
}

func TestTransport_GzipCorrupted(t *testing.T) {
	msg_corrupted := []byte("plaintext message")

	ts := httptest.NewServer(stdHttp.HandlerFunc(func(writer stdHttp.ResponseWriter, response *stdHttp.Request) {
		writer.Header().Set("Content-Encoding", "gzip")
		writer.Write(msg_corrupted)
	}))
	defer ts.Close()

	tr := http.NewTransport()
	client := &stdHttp.Client{Transport: tr}

	response, err := client.Get(ts.URL)
	assert.NoError(t, err, "HTTP request must succeed")

	_, readErr := io.ReadAll(response.Body)
	assert.Error(t, readErr, "Corrupted gzip should return read error when trying to read body")
}

func TestTransport_ZstdCorrupted(t *testing.T) {
	msg_corrupted := []byte("plaintext message")

	ts := httptest.NewServer(stdHttp.HandlerFunc(func(writer stdHttp.ResponseWriter, response *stdHttp.Request) {
		writer.Header().Set("Content-Encoding", "zstd")
		writer.Write(msg_corrupted)
	}))
	defer ts.Close()

	tr := http.NewTransport()
	client := &stdHttp.Client{Transport: tr}

	response, err := client.Get(ts.URL)
	assert.NoError(t, err, "HTTP request must succeed")

	_, readErr := io.ReadAll(response.Body)
	assert.Error(t, readErr, "Corrupted zstd should return read error when trying to read body")
}

func TestTransport_BrotliCorrupted(t *testing.T) {
	msg_corrupted := []byte("plaintext message")

	ts := httptest.NewServer(stdHttp.HandlerFunc(func(writer stdHttp.ResponseWriter, response *stdHttp.Request) {
		writer.Header().Set("Content-Encoding", "br")
		writer.Write(msg_corrupted)
	}))
	defer ts.Close()

	tr := http.NewTransport()
	client := &stdHttp.Client{Transport: tr}

	response, err := client.Get(ts.URL)
	assert.NoError(t, err, "HTTP request must succeed")

	_, readErr := io.ReadAll(response.Body)
	assert.Error(t, readErr, "Corrupted brotli should return read error when trying to read body")
}
