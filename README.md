## About

This demo presents 3 ways of collecting metrics from distributed systems based on microservices:

1. Prometheus
2. OpenTelemetry metrics
3. OpenTelemetry metrics generated from spans (using `spanmetrics` processor)

It also includes OpenTelemetry tracing demonstration (see `Tools` section).

There are 2 variants of this demo:

1. **Basic**, in which requests are sent directly to the microservices.
2. **Extended** – utilizing HAProxy, which works as a loadbalancer for microservices and therefore allows vertical
   scaling of the microservices.

## Project structure

### Demo microservice

A simple microservice has been created for demonstration purposes. It has two functions: (1) calling other microservices
and
(2) simulating internal processing (using a dummy loop). A list of microservices that should be called by given
microservice can be set in `SERVICES_TO_CALL` environmental variable (see `docker-compose.yaml`).

#### Endpoints

The demo microservice provides the following HTTP API endpoints:

* `/api/action` – action endpoint that calls services listed in `SERVICES_TO_CALL` variable and runs a dummy loop that
  simulates internal processing. This endpoint returns `OK` after receiving a response from all called microservices
  or `ERROR` if at least one microservice does not respond.
* `/api/health` – health check endpoint that just returns `OK`.

#### Example

Let's assume we have the following configuration:

* 3 microservices: app1 (available publicly at port `8081`), app2 and app3.
* Variable `SERVICES_TO_CALL` for app1 is set to: `app2:8080,app3:8080`.
* Value of `SERVICES_TO_CALL` for app2 and app3 is empty.

After calling the action endpoint for app1 (http://localhost:8081/api/action), app1 will call internally both app2 and
app3. As a result app2 and app3 will execute their 'internal processing' (dummy loop) and respond to app1, app1 will
also execute its 'internal processing', and finally return the response.

It is also possible to change the complexity of dummy loops in each microservice individually by supplying a complexity
parameter for each of them in the following way:
http://localhost:8081/api/action?config=app1:100000,app2:100000,app3:100000

Please note that other scenarios can be easily prepared by creating more microservices and adjusting `SERVICES_TO_CALL`
variables.

#### Configuration – environmental variables

* `APP_NAME` – name of the app
* `SERVICES_TO_CALL` – list of services (address and port) that should be called by given service. Be careful, avoid
  loops!
* `OTEL_AGENT` – address and port of OpenTelemetry agent

### Tools

This demo consists of the following frontend tools:

* **Jaeger** (distributed tracing) – http://localhost:16686
* **Prometheus** (monitoring) – http://localhost:9090
* **Grafana** (data visualization) – http://localhost:3000
* **Locust** (load testing tool) – http://localhost:8089
* **HAProxy stats** (status of the internal loadbalancer) – http://localhost:1936

Additional resources:

* OpenTelemetry metrics in Prometheus format – http://localhost:8889/metrics

## Deployment

**Basic variant**:

```sh
cd deployment
docker-compose -f docker-compose_basic.yaml up -d
```

**Extended variant** (with HAProxy loadbalancer):

```sh
cd deployment
docker-compose -f docker-compose.yaml up -d
```

---

## Notes

* DNS discovery may cause delays for Prometheus metric updates and HAProxy nodes list
* Metrics buckets should be adjusted (see `bucketsConfig` in `app/monitoring/common.go` and `latency_histogram_buckets`
  in `otel-collector-config.yaml`)
* Docker image versions are fixed
