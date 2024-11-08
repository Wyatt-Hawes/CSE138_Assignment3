package main

import (
	"net/http"
)
type js map[string]interface{}
var kv_pairs = make(js)


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

func post_key(key string, value string) (js, int) {
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