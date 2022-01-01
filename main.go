package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"time"
)
// https://httpd.apache.org/docs/2.4/programs/ab.html
// run this command while your server is running : ab -n 1000 -c 1000 http://localhost:8080/


// go tool pprof out.dump 

func main() {
	http.HandleFunc("/log", logHandler)
	log.Println("hello milad im running")
	http.ListenAndServe(":8080", nil)
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	ch := make(chan int)
	go func() {
		obj := make(map[string]float64)
		if err := json.NewDecoder(r.Body).Decode(&obj); err != nil {
			ch <- http.StatusBadRequest
			return
		}
		// simulation of a time consuming process like writing logs into db
		time.Sleep(time.Duration(rand.Intn(400)) * time.Millisecond)
		ch <- http.StatusOK
	}()

	select {
	case status := <-ch:
		w.WriteHeader(status)
	case <-time.After(200 * time.Millisecond):
		w.WriteHeader(http.StatusRequestTimeout)
	}
}
