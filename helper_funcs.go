package main

import "fmt"

func log(s any) {
	if debug {
		fmt.Println(s)
	}
}
