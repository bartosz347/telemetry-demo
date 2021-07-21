### Demo application (microservice)

http://localhost:8081  
http://localhost:8081/api/action?config=app1:100000,app2:100000,app3:100000

#### Endpoints

* `/api/action` – action endpoint that calls services listed in `SERVICES_TO_CALL` variable and runs a dummy loop that
  simulates internal processing and returns `OK` or `ERROR` if at least one service does not respond.
* `/api/health` – health check endpoint that just returns `OK`.

#### Configuration – environmental variables

* `APP_NAME` – name of the app
* `SERVICES_TO_CALL` – list of services (address and port) that should be called by given service. Be careful, avoid
  loops!
* `OTEL_AGENT` – address and port of OpenTelemetry agent

---

### Tools

This demo consists of the following tools:

* Jaeger – http://localhost:16686
* Prometheus – http://localhost:9090
* Grafana – http://localhost:3000
* HAProxy stats – http://localhost:1936
* Locust – http://localhost:8089

Additional resources:

* OpenTelemetry metrics in Prometheus format – http://localhost:8889/metrics

---

### Deployment

Basic version (without HAProxy):

```sh
cd deployment
docker-compose -f docker-compose_basic.yaml up -d
```

Extended version (with HAProxy):

```sh
cd deployment
docker-compose -f docker-compose.yaml up -d
```

---

### Notes

* DNS discovery may cause delays for Prometheus metric updates and HAProxy nodes list
* Metrics buckets should be adjusted (see `bucketsConfig` in `app/monitoring/common.go` and `latency_histogram_buckets`
  in `otel-collector-config.yaml`)
* Docker image versions are fixed!

docker-compose -f docker-compose_basic.yaml up -d

TODO pomiar netto i brutto - z siecią?

dopisać komentarze, zmienić nazwy handlerów

list endpointów

uspójnienie yaml

fix, żeby action umożliwiało wywołanie np. app2 która wywoła jakieś usługi

drools?