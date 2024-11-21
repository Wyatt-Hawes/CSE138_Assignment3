// http_server.go

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// This is a shorthand for the MAPS of GO so we dont need to type that long ass type
type js map[string]interface{}

// Debug false disables all prints done with the log() function
const debug bool = true

// Parallel maps for tracking values and versions of keys
var kv_pairs = make(js)
var kv_version = make(js)

// DO NOT EDIT ORIGINAL_VIEW
var ORIGINAL_VIEW []string = strings.Split(os.Getenv("VIEW"),",");
// var ORIGINAL_VIEW = []string{"localhost:8090", "localhost:8091"}

var VIEW []string = strings.Split(os.Getenv("VIEW"),",");
// Above works getting env, for testing, redefine
//var VIEW = []string{"localhost:8090", "localhost:8091"}

var IP = os.Getenv("SOCKET_ADDRESS");


func main(){
	// All method types go to each handler (GET POST PUT DELETE etc.)
	http.HandleFunc("/kvs/", kvs_handler)
	http.HandleFunc("/update", update_handler)
	http.HandleFunc("/view", view_handler)
	http.HandleFunc("/all", all_handler)

	fmt.Fprintln(os.Stdout,"View: ", VIEW)
	fmt.Fprintln(os.Stdout,"IP: ", IP)
	fmt.Fprintln(os.Stdout, "Server running!\n---------------")

	// Change from 8090 to 8091 when doing scuffed replication testing (8090 -> launch 1 server, 8091 -> launch 2nd server)
	//http.ListenAndServe(":8090", nil)

	// Broadcast PUT-View
	fmt.Fprintln(os.Stdout, "notifying instances...\n")

	// Call immediately 
	go notifyInstances()

	// Get all data from other servers
	get_all_data()

	// Call every 2 seconds
	ticker := time.NewTicker(2 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
		select {
			case <- ticker.C:
				go notifyInstances()
			case <- quit:
				ticker.Stop()
				return
			}
		}
	}()

	http.ListenAndServe(IP, nil)
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
	log("--- " + method)

	// Get body
	b_data, _ := io.ReadAll(r.Body)
	var body js
	json.Unmarshal(b_data, &body)
	log("Req Body: " + fmt.Sprint(body))

	// Get value from body (if exists)
	value, val_exists := body["value"]
	value_str, val_success := value.(string)

	// Get meta-data from body
	meta_data, meta_success := body["causal-metadata"].(map[string]interface{})

	// Get meta-data fields
	m_data_key, _ := meta_data["key"]
	m_data_version, _ := meta_data["version"]
	log(fmt.Sprintf("  key: %v", m_data_key))
	log(fmt.Sprintf("  version: %v", m_data_version))

	// Create response variable, j_res is a map of |string -> any|
	var j_res js = js{}
	var status int = http.StatusMethodNotAllowed


	switch method {
	case "GET":
		if(meta_success && len(meta_data) != 0 && key == m_data_key){
			log("GET meta-data is for correct key")

			valid := check_valid_metadata("GET", key, int(m_data_version.(float64)))
			if(!valid){
				j_res  = js{"error": "Causal dependencies not satisfied; try again later"}
				status = http.StatusServiceUnavailable
				break
			}
		}
		// else condition for error if bad meta-data ??
		
		// GET operation
		j_res, status = get_key(key)
		break


	case "PUT":
		if(meta_success && len(meta_data) != 0 && key == m_data_key){
			log("PUT meta-data is for correct key")

			valid := check_valid_metadata("PUT", key, int(m_data_version.(float64)));
			if(!valid){
				j_res = js{"error": "Causal dependencies not satisfied; try again later"}
				status = http.StatusServiceUnavailable
				break
			}
		}
		// else condition for error if bad meta-data ??

		// Check for string value
		if !val_exists || !val_success {
			j_res = js{"error": "PUT request does not specify a value"}
			status = http.StatusBadRequest
			break
		}

		// Check key length
		if len(key) > 50 {
			j_res = js{"error": "Key is too long"}
			status = http.StatusBadRequest
			break;
		}

		// PUT operation
		j_res, status = put_key(key, value_str)
		
		// Replicate if successful
		if (status == http.StatusCreated || status == http.StatusOK){
			log("Replicating PUT")
			go replicate("PUT", key, value_str, get_version(key)); // Go launches a `goroutine` aka async call
		}
		break


	case "DELETE":
		// Did meta_data exist on the body, was it not null, and does the key match the current key?
		if(meta_success && len(meta_data) != 0 && key == m_data_key){
			valid := check_valid_metadata("DELETE", key, int(m_data_version.(float64)));
			if(!valid){
				j_res = js{"error": "Causal dependencies not satisfied; try again later"}
				status = http.StatusServiceUnavailable
				break
			}
		}
		// else condition for error if bad meta-data ??

		// DELETE operation
		j_res, status = delete_key(key)
		
		// Replicate if successful
		if (status == http.StatusOK){
			log("Replicating DELETE")
			go replicate("DELETE", key, "", get_version(key)); // Go launches a `goroutine` aka async call
		}
		break


	default:
		// Break if not a method we have, by default is not_implemented
		break;
	}
	
	// the _ is the error value, we are dropping it
	j_data, _ := json.Marshal(j_res)

	// After operations above complete, send our response
	log("Sending Response")
	w.WriteHeader(status)
	w.Write(j_data)
}


func update_handler(w http.ResponseWriter, r *http.Request){
	// Get all we need from the body
	// log("Got an update from someone")
	b_data, _ := io.ReadAll(r.Body)
	var body js
	json.Unmarshal(b_data, &body)

	// Body should have Method, Key, Value
	method_d, e1 := body["method"]
	key_d, e2 := body["key"]
	value_d, e3 := body["value"]
	version_d, e7 := body["version"]
	req_ip_d, e9 := body["ip"]

	// Convert all to string
	method, e4 := method_d.(string)
	key, e5 := key_d.(string)
	value, e6 := value_d.(string)
	req_ip, e10 := req_ip_d.(string);

	version_f, e8 := version_d.(float64)
	new_version := int(version_f)
	current_version := get_version(key)

	var j_res js = js{}
	j_data, _ := json.Marshal(j_res)

	// If any errors, drop request
	if(!e1 || !e2 || !e3 || !e4  || !e5 || !e6 || !e7 || !e8 || !e9 || !e10){
		log("Error with update")
		w.WriteHeader(http.StatusOK);
		w.Write(j_data)
		return;
	}

	// Message was from ourselves, reject
	if(req_ip == IP){
		w.WriteHeader(http.StatusOK);
		w.Write(j_data)
		return;
	}

	// Check if version is acceptable, tie break with S1 < S2
	if (new_version < current_version){
		// reject request if new_version < current_version
		w.WriteHeader(http.StatusOK);
		w.Write(j_data)
		return;
	}

	// If versions are equal and My ip is less than the request, ignore it, 'lower' IP takes priority
	if(new_version == current_version && IP < req_ip){
		log("Tiebreaker REJECT")
		w.WriteHeader(http.StatusOK);
		w.Write(j_data)
		return;
	}

	// If versions equal and our IP is before the request, block
	// if(new_version == current_version && r.)

	// Just do the same operation on our version
	switch method{
		case "PUT":
			log("Received external Update for PUT")
			put_key(key, value)
			set_version(key, new_version)
			break;

		case "DELETE":
			log("Received external Update for DELETE")
			delete_key(key);
			set_version(key, new_version)
			break;

		default:
			break;
	}
	w.WriteHeader(http.StatusOK);
	w.Write(j_data)
	return;
}


func view_handler(w http.ResponseWriter, r *http.Request){
	// Set return header
	w.Header().Set("Content-Type", "application/json")

	// Get method
	method := r.Method

	// Get body
	b_data, _ := io.ReadAll(r.Body)
	var body js
	json.Unmarshal(b_data, &body)
	log("Body: " + fmt.Sprint(body))

	// Get "socket-address" attribute (if exists)
	socket_address_d, exists := body["socket-address"]
	socket_address, success := socket_address_d.(string)

	var j_res js = js{}
	var status int = http.StatusMethodNotAllowed;

	switch method{
	case "GET":
		// Get all views
		j_res, status = get_all_view()
		break;

	case "PUT":
		// Add new view
		if (!exists || !success){
			j_res = js{"error": "No socket provided in body"}
			status = http.StatusBadRequest
			break
		}
		j_res, status = add_view(socket_address)
		break;

	case "DELETE":
		// Delete view
		if (!exists || !success){
			j_res = js{"error": "No socket provided in body"}
			status = http.StatusBadRequest
			break
		}
		j_res, status = delete_view(socket_address)
		break;

	default:
		break;
	}

	// the _ is the error value, we are dropping it
	j_data, _ := json.Marshal(j_res)

	// After operations above complete, send our response
	log("Sending Response")
	log("----------------")
	w.WriteHeader(status)
	w.Write(j_data)
}

func all_handler(w http.ResponseWriter, r *http.Request){
	all_data := js{"pairs": kv_pairs, "versions": kv_version}
	j_data, _ := json.Marshal(all_data);

	w.WriteHeader(http.StatusOK)
	w.Write(j_data);
}

func get_all_data(){
	for _, addr := range ORIGINAL_VIEW {  // ignore index
		if addr == IP {
				// don't notify self
				log("ignoring self")
				continue
		}

		log("Grabbing data from " + addr)
		get_data(addr);
	}
}

func get_data(addr string){

	url:= fmt.Sprintf("http://%s/all",addr)

	// Create the request
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("Content-Type", "application/json");

	client := &http.Client{
	Timeout: 500 * time.Millisecond, 
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse // Ensure that no retries ever happen
		},
	}

	resp, err := client.Do(req)

	if (err != nil){
		log(fmt.Sprintf("Error grabbing data from %s, server may be down", addr))
		return;
	}

	b_data, _ := io.ReadAll(resp.Body)
	var body js
	json.Unmarshal(b_data, &body)
	log("Response Body: " + fmt.Sprint(body))

	pairs := body["pairs"].(map[string]interface{})
	versions := body["versions"].(map[string]interface{})

	// Loop over all versions, compare it to our version, if its >, use the new value
	for key, version_d := range versions{
		version := int(version_d.(float64))
		our_version := get_version(key)
		if our_version < version{
			kv_pairs[key] = pairs[key]
			set_version(key,version)
		}
	}


}