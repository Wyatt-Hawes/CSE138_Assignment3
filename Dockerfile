FROM golang

WORKDIR /app

COPY http_server.go helper_funcs.go key_value_ops.go view_ops.go go.mod ./

EXPOSE 8090:8090

RUN go mod tidy

CMD ["go", "run", "http_server.go", "helper_funcs.go", "key_value_ops.go", "view_ops.go"]
