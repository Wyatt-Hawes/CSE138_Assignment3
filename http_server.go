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
type js map[string]interface{}
var kv_pairs = make(js)

// var view []string = strings.Split(os.Getenv("VIEW"),",");
// Above works getting env, for testing, redefine
var view = []string{"localhost:8090", "localhost:8091"}


func main(){
	http.HandleFunc("/kvs/", kvs_handler)
	http.HandleFunc("/update", update_handler)

	fmt.Fprintln(os.Stdout,view);
	fmt.Fprintln(os.Stdout, "Server running!\n---------------")

	// Change from 8090 to 8091 when doing scuffed replication
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
		j_res, status = put_key(key, value_str)
		
		// Replicate if successful
		if (status == http.StatusCreated || status == http.StatusOK){
			log("Replicating PUT")
			go replicate("PUT", key, value_str);
		}

		break

	case "DELETE":
		j_res, status = delete_key(key)
		
		// Replicate
		if (status == http.StatusCreated || status == http.StatusOK){
			log("Replicating DELETE")
			go replicate("DELETE", key, "");
		}
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


func update_handler(w http.ResponseWriter, r *http.Request){
	// Get all we need from the body
	log("Got an update from someone")
	b_data, _ := io.ReadAll(r.Body)
	var body js
	json.Unmarshal(b_data, &body)

	// Body should have Method, Key, Value
	method_d, e1 := body["method"]
	key_d, e2 := body["key"]
	value_d, e3 := body["value"]

	method, e4 := method_d.(string)
	key, e5 := key_d.(string)
	value, e6 := value_d.(string)

	// If any errors, drop request
	if(!e1 || !e2 || !e3 || !e4 || !e5 || !e6){
		return;
	}

	switch method{
		case "PUT":
			log("Received external Update for PUT")
			put_key(key, value)
			break;
		case "DELETE":
			log("Received external Update for DELETE")
			delete_key(key);
			break;
		default:
			break;
	}
	w.WriteHeader(http.StatusOK);
}