package logger

import (
	"bytes"
	"log/slog"
	"path/filepath"
	"testing"

	"github.com/lmittmann/tint"
	"github.com/stretchr/testify/assert"
)

// Mock os level write with fakeFileWriter

type fakeFileWriter struct {
	content *bytes.Buffer
}

func (f *fakeFileWriter) Write(p []byte) (int, error) {
	return f.content.Write(p)
}

// Mock Load to write to fake writer instead of disk
func mockLoad(level slog.Level) *fakeFileWriter {
	fileWriter := &fakeFileWriter{content: &bytes.Buffer{}}
	opts := &tint.Options{Level: level}
	handler := tint.NewHandler(fileWriter, opts)
	slog.SetDefault(slog.New(handler))
	return fileWriter
}

func TestLoad_CreatesFileAndLogs(t *testing.T) {
	// Mocked Load
	mockWriter := mockLoad(slog.LevelInfo)

	msg := "Message"
	slog.Info(msg)

	// Read msg
	content := mockWriter.content.String()
	assert.NotEmpty(t, content, "Log should not be empty")
	assert.Contains(t, content, msg, "Log should contain the message")
}

func TestDefaultPath(t *testing.T) {
	path := DefaultPath()
	assert.Equal(t, filepath.Base(path), fileName, "Default File name must be %s", fileName)
}
