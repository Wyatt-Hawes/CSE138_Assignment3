FROM golang:1.23.3-alpine

WORKDIR /app

COPY http_server.go helper_funcs.go key_value_ops.go view_ops.go go.mod ./

EXPOSE 8090:8090

# Build an executable since launching the server is faster (We fail the test otherwise)
RUN go build -o asgn3

CMD ["./asgn3"]
#CMD ["go", "run", "http_server.go", "helper_funcs.go", "key_value_ops.go", "view_ops.go"]
