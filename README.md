### App environmental variables

* APP_NAME
* OTEL_AGENT
* SERVICES_TO_CALL

### Services

* http://localhost:16686 – Jaeger
* http://localhost:9090 – Prometheus
* http://localhost:3000 – Grafana
* http://localhost:1936 – HAProxy stats
* http://localhost:8081 – App1

http://localhost:8889/metrics

http://localhost:8081/api/action?config=app1:10,app2:10,app3:10

### Notes

* DNS discovery may cause delays for Prometheus metric updates and HAProxy nodes list
* Adjust buckets
* Docker image versions are fixed!


Disallow calling itself?

empty services to call