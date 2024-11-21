// view_ops.go

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"time"
)


func get_all_view() (js, int) {
	return js{"view": VIEW}, http.StatusOK
}


func add_view(new_view_address string)(js, int){
	log("got add_view request (new address: " + new_view_address + ")")

	if slices.Contains(VIEW, new_view_address){
		return js{"result":"already present"}, http.StatusOK
	}
	VIEW = append(VIEW, new_view_address);

	if !slices.Contains(ORIGINAL_VIEW, new_view_address){
		ORIGINAL_VIEW = append(ORIGINAL_VIEW, new_view_address);
	}

	return js{"result":"created"}, http.StatusCreated
}


func delete_view(view_to_remove string)(js, int){
	if !slices.Contains(VIEW, view_to_remove){
		return js{"error":"View has no such replica"}, http.StatusNotFound
	}
	VIEW = slices.DeleteFunc(VIEW, func(v string)(bool){return v == view_to_remove})
	return js{"result":"deleted"}, http.StatusOK
}


func notifyInstances() {
	for _, addr := range ORIGINAL_VIEW {  // ignore index
			if addr == IP {
					// don't notify self
					log("ignoring self")
					continue
			}

			log("Notifying instance: " + addr)
			go notify_server(addr);
	}
}

func notify_server(addr string){
	// Request body
	body := map[string]string{"socket-address": IP}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
			log("Error marshaling body for address: " + addr)
			return
	}

	// Send the PUT request
	url := fmt.Sprintf("http://%s/view", addr)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
			log("Error creating PUT request to: " + addr)
			return
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{
		Timeout: 500 * time.Millisecond, 
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Ensure that no retries ever happen
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		log("Error sending notification request to: " + addr)
		delete_view(addr)
		return
	}

	// Read response
	respBody, _ := io.ReadAll(resp.Body)
	log(fmt.Sprintf("Response from %s: %s", addr, string(respBody)))
	resp.Body.Close()
}