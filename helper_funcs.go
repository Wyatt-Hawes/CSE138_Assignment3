package main

import "fmt"

func log(s any) {
	if debug {
		fmt.Println(s)
	}
}

func check_valid_metadata(method string,key string, request_version int)(bool){

	
	server_version := get_version(key);
	log(fmt.Sprintf("Version diff %d | %d",server_version, request_version))

	switch (method){
	case "GET":
		// Check request version. If request is GREATER THAN server version, invalid request(?)
		if(request_version > server_version){
			return false
		}
		break

	case "PUT", "DELETE":
		// Check request version, server_version must be EQUAL or GREATER, if LESS, then reject
		if(request_version < server_version){
			return false
		}
		set_version(key,request_version)
		break
	default:
		return false;
	}

	return true
}
