// main.go
package main

import (
	"fmt"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I no love %s!", r.URL.Path[1:])
}
func checkExpire() {
    for {
        // do some job
        fmt.Println(time.Now().UTC())
        time.Sleep(1000 * time.Millisecond)
    }
}

func main() {

	// customerCount := 10
	// GenerateDummyData(customerCount)
	start := time.Now() // Record the start time
	// go checkExpire()
	// http.HandleFunc("/", handler) // http://127.0.0.1:8080/Go
	// http.ListenAndServe(":8080", nil)
	// nonmulti.Old_forward_def_tagger()
	// forward_default_tagger()
	forward_default_tagger_big_data()
	elapsed := time.Since(start) // Calculate the elapsed time
	fmt.Printf("Time taken: %s\n", elapsed)
	// start = time.Now() // Record the start time
	// // GenerateDummyData()
	// // forward_default_tagger()
	// elapsed = time.Since(start) // Calculate the elapsed time
	// fmt.Printf("Time taken: %s\n", elapsed)
}
