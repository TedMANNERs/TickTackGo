package master

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var slaves = []Slave{}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		slaves = append(slaves, Slave{IP: r.RemoteAddr})
		fmt.Println(r.RemoteAddr + " registered")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Method not allowed"))
	}
}

func gamesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		slavesRef := &slaves
		encodedSlaves, err := json.Marshal(slavesRef)
		if err != nil {
			fmt.Println(err)
			return
		}
		w.Write(encodedSlaves)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Method not allowed"))
	}
}

func Start() {
	fmt.Println("Starting master...")
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/games", gamesHandler)
	http.ListenAndServe("localhost:8080", nil)
}
