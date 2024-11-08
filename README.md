# How to Run

`go run http_server.go helper_funcs.go key_value_ops.go`
`set VIEW="localhost:8090,localhost:8091" `

`docker build -t app .`
`docker run --rm -p 8090:8090 -e=VIEW=localhost:8090,localhost:8091 app`
