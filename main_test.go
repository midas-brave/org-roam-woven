package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"path/filepath"

	"github.com/stretchr/testify/assert"
)

func setupTestEnvironment(t *testing.T) func() {
	tempDir, err := os.MkdirTemp("", "roam_test")
	assert.NoError(t, err)

	subDir := filepath.Join(tempDir, "subDir")
	assert.NoError(t, os.Mkdir(subDir, 0755))
	testFile := filepath.Join(tempDir, "testFile.txt")
	assert.NoError(t, os.WriteFile(testFile, []byte("hello world!"), 0644))

	os.Setenv("ROAM_DIRECTORY", tempDir)
	cleanup := func() {
		os.RemoveAll(tempDir)
		os.Setenv("ROAM_DIRECTORY", tempDir)
	}
	return cleanup
}

func TestRootEndpoint(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"message": "Welcome to org-roam-woven!"}`, w.Body.String())
}

func TestFiles(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	tests := []struct {
		name         string
		filePath     string
		expectedCode int
		expectedBody any
	}{
		{
			name:         "list path",
			filePath:     "",
			expectedCode: http.StatusOK,
		},
		{
			name:         "view file",
			filePath:     "testFile.txt",
			expectedCode: http.StatusOK,
			expectedBody: []byte("hello world!"),
		},
		{
			name:         "list sub directory",
			filePath:     "subDir",
			expectedCode: http.StatusOK,
		},
		{
			name:         "not exist path",
			filePath:     "404Folder",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "forbidden dir",
			filePath:     "../forbiddenDir",
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "nasty path",
			filePath:     "foo/../../..",
			expectedCode: http.StatusForbidden,
		},
	}

	r := setupRouter()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/api/unstable/files/"+tt.filePath, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedBody != nil {
				assert.Equal(t, tt.expectedBody, w.Body.Bytes())
			}
		})
	}
}

func TestCreateFiles(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	tests := []struct {
		name         string
		filePath     string
		data         []byte
		expectedCode int
	}{
		{
			name:         "create file",
			filePath:     "xxx.txt",
			data:         []byte("Hello!"),
			expectedCode: http.StatusCreated,
		},
		{
			name:         "update file",
			filePath:     "xxx.txt",
			data:         []byte("Hello world!"),
			expectedCode: http.StatusNoContent,
		},
		{
			name:         "create directory failed",
			filePath:     "subDir",
			data:         []byte("Hello!"),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "nasty path",
			filePath:     "foo/../../..",
			data:         []byte("Goodbye world!"),
			expectedCode: http.StatusForbidden,
		},
	}

	r := setupRouter()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			fileWriter, err := writer.CreateFormFile("file", "x")
			assert.NoError(t, err)
			fileReader := bytes.NewReader(tt.data)
			_, err = io.Copy(fileWriter, fileReader)
			assert.NoError(t, err)
			err = writer.Close()
			assert.NoError(t, err)

			req, _ := http.NewRequest(http.MethodPut, "/api/unstable/files/"+tt.filePath, body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestDeleteFiles(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	tests := []struct {
		name          string
		filePath      string
		expectedCode  int
		expectedCheck func(t *testing.T, respBody []byte)
	}{
		{
			name:         "delete path",
			filePath:     "testFile.txt",
			expectedCode: http.StatusOK,
		},
		{
			name:         "delete path after",
			filePath:     "testFile.txt",
			expectedCode: http.StatusNotFound,
		},
	}

	r := setupRouter()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodDelete, "/api/unstable/files/"+tt.filePath, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
