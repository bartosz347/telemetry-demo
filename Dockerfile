FROM golang:1.17rc1-alpine
RUN apk add --no-cache gcc

WORKDIR /app

COPY app/go.mod app/go.sum ./
RUN go mod download

COPY app ./

RUN go build

EXPOSE 8080
CMD [ "./telemetry-demo" ]
