// key_value_ops.go

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)


func get_key(key string) (js, int) {
	var res js
	var status int

	value, exists := kv_pairs[key]
	version := get_version(key)

	// If key does NOT exist, error
	if !exists {
		res = js{
			"error": "Key does not exist",
		}
		status = http.StatusNotFound
	
	} else {
		// If key DOES exist, return value
		res = js{
			"result": "found",
			"value":  value,
			"causal-metadata": js{
				"key": key,
				"version": version,
			},
		}
		status = http.StatusOK
	}

	return res, status
}


func put_key(key string, value string) (js, int) {
	var res js
	var status int

	_, exists := kv_pairs[key]

	if !exists {
		// If key does NOT exist, add
		kv_pairs[key] = value
		version := get_add_version(key)
		res = js{
			"result": "created",
			"causal-metadata":  js{
				"key": key,
				"version": version,
			},
		}
		status = http.StatusCreated

	} else {
		// If key DOES exists, replace value
		kv_pairs[key] = value
		version := get_add_version(key)
		res = js{
			"result": "replaced",
			"causal-metadata":  js{
				"key": key,
				"version": version,
			},
		}
		status = http.StatusOK
	}

	return res, status
}


func delete_key(key string) (js, int) {
	var res js
	var status int

	_, exists := kv_pairs[key]

	if !exists {
		// If key does NOT exist, error
		res = js{
			"error": "Key does not exist",
		}
		status = http.StatusNotFound
	
	} else {
		// If key DOES exists, delete
		delete(kv_pairs, key)
		version := get_add_version(key)

		res = js{
			"result": "deleted",
			"causal-metadata": js{
				"key": key,
				"version": version,
			},
		}
		status = http.StatusOK
	}
	
	return res, status
}


func replicate(method string, key string, value string, version int){
	for _, v := range VIEW{
		/*
			We should probably have some logic here to NOT communicate with ourselves, but 
			the way it is right now, a request also gets sent to ourself, which has no effect ofc but is annoying
		*/

		// Dont send message to ourself
		if (v == IP){
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
		communicate(req);
	}
}


func communicate(req *http.Request){

	// We need to set a custom redirect time limit or else it will keep re-trying the request
	// Check redirect function from Stack Overflow https://stackoverflow.com/questions/23297520/how-can-i-make-the-go-http-client-not-follow-redirects-automatically
	client := &http.Client{
		Timeout: 500 * time.Millisecond, 
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
        	return http.ErrUseLastResponse // Ensure that no retries ever happen
    	},
	}

	// Probably check for down servers here with the response/error
	resp, err := client.Do(req);
	//client.Do(req);

	// Maybe check for errors here and add/remove servers from the VIEW
	if (err != nil){
		log("Error communicating, server may be down")
		// resp.Body.Close();
		return;
	}
	log(fmt.Sprintf("Replication success : %s", resp.Status))

	resp.Body.Close();
	
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


func set_version(key string, version int){
	log("Setting: "+key+" to "+fmt.Sprint(version))
	kv_version[key] = version;
}