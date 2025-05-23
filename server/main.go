package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/ws", handleConnections)
	fmt.Println("--> server runs on 8080...")
	http.ListenAndServe(":8080", nil)
}
