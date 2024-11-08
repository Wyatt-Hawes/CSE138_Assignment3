package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)



func get_key(key string) (js, int) {
	if len(key) > 50 {
		return js{"error": "Key is too long"}, http.StatusBadRequest
	}

	value, exists := kv_pairs[key]

	if !exists {
		return js{"error": "Key does not exist"}, http.StatusNotFound
	}

	return js{"result": value}, http.StatusOK
}

func put_key(key string, value string) (js, int) {
	_, exists := kv_pairs[key]

	// Set default response
	res := js{"result": "replaced"}
	status := http.StatusOK

	// Check if key exists & update response
	if !exists {
		res = js{"result": "created"}
		status = http.StatusCreated
	}

	kv_pairs[key] = value
	return res, status
}

func delete_key(key string) (js, int) {
	_, exists := kv_pairs[key]

	if !exists {
		return js{"error": "Key does not exist"}, http.StatusNotFound
	}
	delete(kv_pairs, key)

	return js{"result": "deleted"}, http.StatusOK
}

func replicate(method string, key string, value string){
	for _, v := range view{
		// For each view, make a request to it with the update
		url:= fmt.Sprintf("http://%s/update",v)
		log(fmt.Sprintf("URL:%s",url))
		body := js{"method":method,"key":key,"value":value}
		body_js, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body_js))
		req.Header.Set("Content-Type", "application/json");

		// Asynchronously communicate to server for each value
		go communicate(req);
	}
}

func communicate(req *http.Request){
	client := &http.Client{}

	// Probably check for down servers here with the response/error
	//resp, err := client.Do(req);
	client.Do(req);
}
