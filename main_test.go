package main

import (
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRootEndpoint(t *testing.T) {
	go func() {
		if err := startServer(); err != nil {
			t.Fatal(err)
		}
	}()

	timeout := time.After(5 * time.Second)
	for {
		select {
		case <-timeout:
			t.Fatal("Server did not start in time")
		default:
			conn, err := net.Dial("tcp", "localhost:18080")
			if err == nil {
				conn.Close()
				goto ServerReady
			}
			time.Sleep(100 * time.Millisecond)
		}
	}

ServerReady:
	resp, err := http.Get("http://localhost:18080/")
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := ioutil.ReadAll(resp.Body)
	assert.JSONEq(t, `{"message": "Welcome to org-roam-woven!"}`, string(body))
}

