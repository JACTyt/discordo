package logger

import (
	"bytes"
	"log/slog"
	"os"
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

func TestLoad_RealFileSystem_Success(t *testing.T) {
	tmpDir := os.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")

	err := Load(logPath, slog.LevelInfo)
	assert.NoError(t, err, "Load should not return an error")

	slog.Info("hello world")

	data, err := os.ReadFile(logPath)
	assert.NoError(t, err, "Failed to read log file")
	assert.NotEmpty(t, data, "Log file should have content message")
}

func TestLoad_RealFileSystem_Fail_WrongDirectory(t *testing.T) {
	tmpDir := os.TempDir()
	logPath := filepath.Join(tmpDir, "not-a-director:y.txt/test.log")

	err := Load(logPath, slog.LevelInfo)
	assert.Error(t, err, "Load should return an error")
}

func TestLoad_RealFileSystem_Fail_NoFilePermission(t *testing.T) {
	logPath := "C:\\Windows\\system32\\test.log"

	Load(logPath, slog.LevelInfo)

	file, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0000)
	assert.Error(t, err, "Opening log file should fail")
	defer file.Close()
}
