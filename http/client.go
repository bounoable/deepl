package http

//go:generate mockgen -source=client.go -destination=./mocks/client.go

import "net/http"

// Client is an interface for the *net/http.Client struct.
type Client interface {
	Do(*http.Request) (*http.Response, error)
}
