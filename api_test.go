package gcd

import (
	"net/http"
	"testing"
)

func TestMain(m *testing.M) {
	defaultHTTPClient = &http.Client{}

	m.Run()
}
