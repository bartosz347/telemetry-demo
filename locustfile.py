from locust import HttpUser, task, between, constant, constant_pacing
# from locust.contrib.fasthttp import FastHttpUser
from locust import LoadTestShape

host = "http://localhost:8081"


class SiteUser(HttpUser):
    wait_time = between(1, 2)
    # wait_time = constant(1)
    # wait_time = constant_pacing(1)
    host = host

    @task
    def trafficCamera(self):
        self.client.get("/api/action")
