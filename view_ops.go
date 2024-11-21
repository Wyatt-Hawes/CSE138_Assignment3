// view_ops.go

package main

import (
	"net/http"
	"slices"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	return js{"result":"created"}, http.StatusCreated
}


func delete_view(view_to_remove string)(js, int){
	if !slices.Contains(VIEW, view_to_remove){
		return js{"error":"View has no such replica"}, http.StatusNotFound
	}
	VIEW = slices.DeleteFunc(VIEW, func(v string)(bool){return v == view_to_remove})
	return js{"result":"deleted"}, http.StatusOK
}


func notifyInstancesOnStartup() {
	for _, addr := range VIEW {  // ignore index
			if addr == IP {
					// don't notify self
					log("ignoring self")
					continue
			}

			log("Notifying instance: " + addr)

			// Request body
			body := map[string]string{"socket-address": IP}
			bodyJSON, err := json.Marshal(body)
			if err != nil {
					log("Error marshaling body for address: " + addr)
					continue
			}

			// Send the PUT request
			url := fmt.Sprintf("http://%s/view", addr)
			req, err := http.NewRequest("PUT", url, bytes.NewBuffer(bodyJSON))
			if err != nil {
					log("Error creating PUT request to: " + addr)
					continue
			}
			req.Header.Set("Content-Type", "application/json")

			// Execute the request
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
					log("Error sending request to: " + addr)
					continue
			}

			// Read response
			defer resp.Body.Close()  // execute when the function exits
			respBody, _ := io.ReadAll(resp.Body)
			log(fmt.Sprintf("Response from %s: %s", addr, string(respBody)))
	}
}