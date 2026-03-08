FROM golang:1.21 AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o go_main_final_project main.go

FROM ubuntu:latest

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /app

COPY --from=builder /app/go_main_final_project .
COPY web ./web

EXPOSE 7540
ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/todo.db
ENV TODO_PASSWORD=""

CMD ["./go_main_final_project"]