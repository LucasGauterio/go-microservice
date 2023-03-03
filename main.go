package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httpTransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log/level"
)

type helloRequest struct {
	Name string `json:"name"`
}

type helloResponse struct {
	Greeting string `json:"greeting"`
}

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "microservice", "hello")

	helloHandler := httpTransport.NewServer(
		makeHelloEndpoint(),
		decodeHelloRequest,
		encodeHelloResponse,
	)

	http.Handle("/hello", helloHandler)

	level.Info(logger).Log("msg", "starting microservice")

	if err := http.ListenAndServe(":8090", nil); err != nil {
		level.Error(logger).Log("msg", "failed to start miscroservice", "err", err)
		os.Exit(1)
	}
}

func makeHelloEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(helloRequest)
		greeting := fmt.Sprintf("Hello, %s!", req.Name)
		return helloResponse{greeting}, nil
	}
}

func decodeHelloRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request helloRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeHelloResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
