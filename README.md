## About

This demo presents 3 ways of collecting metrics from distributed systems based on microservices:

1. Prometheus
2. OpenTelemetry metrics
3. OpenTelemetry metrics generated from spans (using `spanmetrics` processor)

It also includes OpenTelemetry tracing demonstration (see `Tools` section).

There are 2 variants of this demo:

1. **Basic** – in which requests are sent directly to the microservices.
2. **Extended** – utilizing HAProxy, which works as a loadbalancer for microservices and therefore allows horizontal
   scaling of the microservices.

## Project structure

### Demo microservice

A simple microservice has been created for demonstration purposes. It has two functions: (1) calling other microservices
and
(2) simulating internal processing (using a dummy loop). A list of microservices that should be called by given
microservice can be set in its `SERVICES_TO_CALL` environmental variable (see `docker-compose.yaml`).

#### Endpoints

The demo microservice provides the following HTTP API endpoints:

* `/api/action` – action endpoint that calls services listed in `SERVICES_TO_CALL` variable and runs a dummy loop that
  simulates internal processing. In particular, when `SERVICES_TO_CALL` is empty, only internal processing is executed by given microservice.
  
  This endpoint returns `OK` after receiving a response from all called microservices
  or `ERROR` if at least one microservice does not return success or at least one request times out.
* `/api/health` – health check endpoint that always returns `OK`.

#### Example

Let's assume we have prepared the following configuration:

* 3 microservices: app1 (available publicly at port `8081`), app2 and app3.
* Variable `SERVICES_TO_CALL` for app1 is set to: `app2:8080, app3:8080`.
* Value of `SERVICES_TO_CALL` for app2 and app3 is empty.

After calling the action endpoint for app1 (http://localhost:8081/api/action), app1 will call internally both app2 and
app3. As a result app2 and app3 will execute their 'internal processing' (dummy loop) and respond to app1, app1 will
also execute its 'internal processing', and finally return the response.

It is also possible to change the complexity of dummy loops in each microservice individually by supplying a complexity
parameter for each of them in the following way:
http://localhost:8081/api/action?config=app1:100000,app2:100000,app3:100000

Please note that other scenarios can be easily prepared by creating more microservices and adjusting their `SERVICES_TO_CALL`
variables.

#### Configuration – environmental variables

* `APP_NAME` – name of the app (e.g. `app1`)
* `SERVICES_TO_CALL` – list of services (address and port) that should be called by given service. Be careful, avoid
  loops! (e.g. `app2:8080,app3:8080`)
* `OTEL_AGENT` – address and port of OpenTelemetry agent (e.g. `otel-agent:4317`)

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

Start
```sh
cd deployment
docker-compose -f docker-compose_basic.yaml up -d
```
Rebuild the images
```sh
docker-compose -f docker-compose_basic.yaml up -d --build
```


**Extended variant** (with HAProxy loadbalancer):

Start
```sh
cd deployment
docker-compose -f docker-compose.yaml up -d
```

Rebuild the images
```sh
cd deployment
docker-compose -f docker-compose.yaml up -d --build
```

## Available metrics

The following metrics are available in Prometheus for each microservice:

1. **Native Prometheus metrics**
   `<APP_NAME>_operation_latency_bucket{status="OK|ERROR", type="internal-only|total"}`    
   Labels:  
   * `type=total` – whole request processing time (internal processing + external service calls)
   * `type=internal-only` – internal processing time (only dummy loop)
2. **OpenTelemetry metrics** (metrics collected with OpenTelemetry and exported for Prometheus)
   `otel_<APP_NAME>_operation_latency_bucket{status="OK|ERROR", type="internal-only|total"}`  
   Labels:
   * `type=total` – whole request processing time (internal processing + external service calls)
   * `type=internal-only` – internal processing time (only dummy loop)   
3. **OpenTelemetry metrics generated from spans** (using `spanmetrics` processor, exported for Prometheus). See Jaeger for more information about labels.  
   `otel_latency_bucket{service_name="<APP_NAME>", operation="internal-processing|/api/action|GET, status_code=STATUS_CODE_OK|STATUS_CODE_ERROR"}`

`<APP_NAME>` – name of the application (`APP_NAME` variable), e.g. `app1`

[Click to open an example for `app1`](http://localhost:9090/graph?g0.expr=histogram_quantile(0.95%2C%20sum(rate(app1_operation_latency_bucket%7Bstatus%3D%22OK%22%2C%20type%3D%22internal-only%22%7D%5B1m%5D))%20by%20(le))%20*%201000&g0.tab=0&g0.stacked=0&g0.show_exemplars=0&g0.range_input=15m&g1.expr=histogram_quantile(0.95%2C%20sum(rate(otel_app1_operation_latency_bucket%7Bstatus%3D%22OK%22%2C%20type%3D%22internal-only%22%7D%5B1m%5D))%20by%20(le))%20*%201000&g1.tab=0&g1.stacked=0&g1.show_exemplars=0&g1.range_input=15m&g2.expr=histogram_quantile(0.95%2C%20rate(otel_latency_bucket%7Bservice_name%3D%22app1%22%2C%20operation%3D%22internal-processing%22%7D%5B1m%5D))&g2.tab=0&g2.stacked=0&g2.show_exemplars=0&g2.range_input=15m) 

## Notes

* DNS discovery may cause delays for Prometheus metric updates and HAProxy nodes list.
* Metrics buckets should be adjusted (see `bucketsConfig` in `app/monitoring/common.go` and `latency_histogram_buckets`
  in `otel-collector-config.yaml`).
* In current configuration OpenTelemetry will trace all requests.  
