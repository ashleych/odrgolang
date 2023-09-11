// main.go
package main

import (
	"fmt"
	"net/http"
	"odr/dbInput"
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
	LOAD_DATA := true
	if LOAD_DATA {

		dbInput.GenerateDataAndMigrateToDB()
	}

	var defaultedRecords []dbInput.CustomerRecord
	dbInput.SubsetDefaultedCustomers(&defaultedRecords)

	length := len(defaultedRecords)
	fmt.Printf("Length of defaultedRecords: %d\n", length)
	dbInput.ComputeLagDate(&defaultedRecords)
fmt.Printf("%+v",defaultedRecords[1])
}
