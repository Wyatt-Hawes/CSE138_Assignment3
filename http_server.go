package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const debug bool = true

// kvs from kv_operations


func main(){
	http.HandleFunc("/kvs/", kvs_handler)
	fmt.Fprintln(os.Stdout, "Server running!\n---------------")
	http.ListenAndServe(":8090", nil)
}

func kvs_handler(w http.ResponseWriter, r *http.Request) {
	// Set return type header
	w.Header().Set("Content-Type", "application/json")

	// Get path
	path := r.URL.Path
	segments := strings.Split(path, "/")
	key := segments[2]

	// Get method
	method := r.Method

	// Get body
	b_data, _ := io.ReadAll(r.Body)
	var body js
	json.Unmarshal(b_data, &body)

	log("Body: " + fmt.Sprint(body))

	// create response variable, j_res is a map of string -> any
	var j_res js = js{}
	var status int = http.StatusMethodNotAllowed

	switch method {
	case "GET":
		j_res, status = get_key(key)
		break
	case "PUT":
		value, exists := body["value"]
		value_str, success := value.(string)
		if !exists || !success {
			j_res, status = js{"error": "PUT request does not specify a value"}, http.StatusBadRequest
			break
		}
		j_res, status = post_key(key, value_str)
		break

	case "DELETE":
		j_res, status = delete_key(key)
		break

	default:
		return
	}
	
	// the _ is the error value
	j_data, _ := json.Marshal(j_res)

	log("Sending Response")
	w.WriteHeader(status)
	w.Write(j_data)
}