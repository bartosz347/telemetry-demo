package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Service struct {
	address string
	port    uint
}

func (service Service) call(ctx context.Context, config string) error {
	log.Printf("INFO: calling %s \n", service.address)

	var httpClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	req, _ := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("http://%s:%d/api/internal", service.address, service.port),
		nil,
	)

	q := req.URL.Query()
	q.Add("config", config)
	req.URL.RawQuery = q.Encode()

	response, err := httpClient.Do(req)
	if err != nil {
		log.Printf("WARNING: %s\n", err)
		return err
	}
	body, err := ioutil.ReadAll(response.Body)
	_ = response.Body.Close()

	log.Printf("INFO: Response received from %s: %s\n", service.address, body)
	return nil
}

func (service Service) String() string {
	return fmt.Sprintf("{addr=%s,port=%d}", service.address, service.port)
}
