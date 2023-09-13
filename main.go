// main.go
package main

import (
	"fmt"
	"odr/dbInput"
	"time"
)

func TestRoutine(){
	LOAD_DATA := true
	if LOAD_DATA {

		dbInput.GenerateDataAndMigrateToDB(10000)
	}

	var defaultedRecords []dbInput.CustomerRecord
	dbInput.SubsetDefaultedCustomers(&defaultedRecords)

	length := len(defaultedRecords)
	fmt.Printf("Length of defaultedRecords: %d\n", length)
	dbInput.ComputeLagDate(&defaultedRecords)
	defaultedRecordsMap := dbInput.ConvertToMap(&defaultedRecords)
	defCusts := dbInput.GetDefaultCustIDs(defaultedRecordsMap)
	customersEverDefaulted := dbInput.GetCustomersEverDefaulted(&defCusts)
	customersEverDefaultedMap := dbInput.ConvertToMap(&customersEverDefaulted)
	dbInput.FowardTag(customersEverDefaultedMap)
	fmt.Printf("%+v", customersEverDefaulted[0])
	customerRecordsTobeUpdated := dbInput.ConvertToSlice(customersEverDefaultedMap)
	customersWithFwdDefaultTrue := dbInput.FilterRecordsByFwdDefault(customerRecordsTobeUpdated)

	// dbInput.UpdateCustomerRecords(customersWithFwdDefaultTrue)
	dbInput.BulkInsertCustomerRecords(customersWithFwdDefaultTrue)
	fmt.Printf("customer example with fwd default true %+v \n", customersWithFwdDefaultTrue[0])

}

func main() {
	go TestRoutine()
	time.Sleep(5* time.Second)
	fmt.Println("************************************************NEW TEST*************************************")
	time.Sleep(5* time.Second)

}
