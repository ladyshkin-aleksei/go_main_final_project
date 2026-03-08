FROM ubuntu:latest

RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /app
COPY go_main_final_project /app/go_main_final_project
COPY web /app/web

EXPOSE 7540

ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/todo.db
ENV TODO_PASSWORD=""

CMD ["/app/go_main_final_project"]
