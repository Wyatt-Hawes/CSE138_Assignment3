package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)



func get_key(key string) (js, int) {

	value, exists := kv_pairs[key]
	version := get_version(key)
	if !exists {
		res:= js{"error": "Key does not exist"}
		res["casual-metadata"] = js{"key":key, "version":version}
		if(version == 0){
			res["casual-metadata"] = js{"key":key, "version":nil};
		}
		return res, http.StatusNotFound
	}
	return js{"result": value, "casual-metadata":js{"key":key,"version":version}}, http.StatusOK
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

	version := get_add_version(key)
	res["casual-metadata"] = js{"key":key, "version": version}

	return res, status
}

func delete_key(key string) (js, int) {
	// Check metadata version, version must be EQUAL or GREATER, if LESS, then reject

	_, exists := kv_pairs[key]

	if !exists {
		return js{"error": "Key does not exist"}, http.StatusNotFound
	}
	
	delete(kv_pairs, key)
	res := js{"result": "deleted"}

	version := get_add_version(key)
	res["casual-metadata"] = js{"key":key, "version": version}

	return res, http.StatusOK
}

func replicate(method string, key string, value string, version int){

	for _, v := range VIEW{
		/*
			We should probably have some logic here to NOT communicate with ourselves, but 
			the way it is right now, a request also gets sent to ourself, which has no effect ofc but is annoying
		*/
		if(v == IP){
			continue;
		}

		// For each view, make a request to it with the update
		// Form the full URL
		url:= fmt.Sprintf("http://%s/update",v)
		log(fmt.Sprintf("URL:%s",url))

		// Form body with all needed info (for now)
		body := js{"method":method,"key":key,"value":value, "version": version, "ip":IP}
		body_js, _ := json.Marshal(body)

		// Create the request
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

	// Maybe check for errors here and add/remove servers from the VIEW
}

func get_add_version(key string)(version int){
	vs_d, exists := kv_version[key];
	vs, e := vs_d.(int)
	// If no entry exists, create it with a version 1
	if (!exists || !e){
		log("Version doesnt exist")
		kv_version[key] = 0
		vs = 0
	}

	// Increment version
	vs++

	// Set new version & return
	set_version(key, vs)
	return vs
}

func get_version(key string)(version int){
	v_d, _ := kv_version[key]

	v, _ := v_d.(int)

	return v;
}

//func get_version

func set_version(key string, version int){
	log("Setting: "+key+" to "+fmt.Sprint(version))
	kv_version[key] = version;
}