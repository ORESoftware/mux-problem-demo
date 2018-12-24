FROM golang:latest
RUN mkdir -p /app
WORKDIR /app
COPY . .

ENV huru_api_host "0.0.0.0"
ENV huru_api_port 3000
ENV GOPATH /app

RUN go install huru

EXPOSE 3000

ENTRYPOINT /app/bin/huru