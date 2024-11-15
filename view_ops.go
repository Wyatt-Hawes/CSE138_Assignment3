// view_ops.go

package main

import (
	"net/http"
	"slices"
)


func get_all_view() (js, int) {
	return js{"view": VIEW}, http.StatusOK
}


func add_view(new_view_address string)(js, int){
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
