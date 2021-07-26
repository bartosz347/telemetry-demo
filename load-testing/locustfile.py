from locust import HttpUser, task, between, constant, constant_pacing
from locust import LoadTestShape

host = "http://localhost:8081"

class SiteUser(HttpUser):
    wait_time = between(1, 2)
    # wait_time = constant(1)
    # wait_time = constant_pacing(1)
    host = host

    @task
    def callAction(self):
        self.client.get("/api/action")
        # Custom complexity example:
        # self.client.get("/api/action?config=app1:100000,app2:1000,app3:1000")
