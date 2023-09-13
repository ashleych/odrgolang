package dbInput

import (
	"fmt"
	"odr/config"
	"sync"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func OpenDB() *gorm.DB {

	db, _ := gorm.Open(sqlite.Open("mydatabase.db"), &gorm.Config{})
	return db

}
func SubsetDefaultedCustomers(defaultRecords *[]CustomerRecord) {
	db := OpenDB()
	// var defaultRecords []CustomerRecord
	db.Where("DPD > ?", config.DEFAULT_DEFINITION_IN_DPD).Find(&defaultRecords)

}

func ComputeLagDate(defaultRecords *[]CustomerRecord) {
	records := *defaultRecords
	for i := range records {
		// Check if the current record is within the next 12 months
		records[i].LagDate = records[i].ObsDate.AddDate(0, -12, 0)
	}

}
func ConvertToMap(defaultRecords *[]CustomerRecord) map[int][]CustomerRecord {

	customerRecordMap := make(map[int][]CustomerRecord)

	// customerRecords := make(map[int][]Record)
	records := *defaultRecords
	for _, record := range records {
		// Check if the current record is within the next 12 months
		customerRecordMap[record.CustID] = append(customerRecordMap[record.CustID], record)

	}
	return customerRecordMap
}

func ConvertToSlice(customerRecordsMap map[int][]CustomerRecord) []CustomerRecord {
	var customerRecords []CustomerRecord
	for _, custRec := range customerRecordsMap {
		customerRecords = append(customerRecords, custRec...)
	}
	return customerRecords
}
func GetDefaultCustIDs(customerRecordMap map[int][]CustomerRecord) []int {
	// defaultCusts := make([]int)
	var defaultCusts []int
	for key := range customerRecordMap {
		defaultCusts = append(defaultCusts, key)
	}
	return defaultCusts
}

func GetCustomersEverDefaulted(defCusts *[]int) []CustomerRecord {
	db := OpenDB()
	var customers []CustomerRecord
	defs := *defCusts
	batchSize := 1000 // You can adjust the batch size based on your needs

	// Use FindInBatches to retrieve customers in batches
	db.Where("cust_id IN (?)", defs).FindInBatches(&customers, batchSize, func(tx *gorm.DB, batch int) error {
		fmt.Printf("Processing batch %d\n", batch)
		return nil
	})
	// db.Where("cust_id IN (?)", defCusts).Find(&customers)

	return customers
}
func (r CustomerRecord) getForwardDate() time.Time {
	// obsDate, _ := time.Parse("2-Jan-06", r.ObsDate)
	obsDate := r.ObsDate
	return obsDate.AddDate(0, 12, 0)
}

func CheckForwardDefaults(customerRecords []CustomerRecord) {
	for i := range customerRecords {
		// Check if the current record is within the next 12 months
		if customerRecords[i].DPD > config.DEFAULT_DEFINITION_IN_DPD {
			break
		}
		forwardDate := customerRecords[i].getForwardDate()

		for _, futureRecord := range customerRecords[i+1:] {
			if futureRecord.ObsDate.Before(forwardDate.AddDate(0, 12, 0)) {
				// Check if there is a DPD greater than 90 within 12 months
				if futureRecord.DPD > 90 {
					// Set the FwdDefault field to true for the current record
					customerRecords[i].FwdDefault = true
					customerRecords[i].FwdDefaultDate = futureRecord.ObsDate
					break // No need to continue checking
				}
			} else {
				break // No need to check further if future records are outside 12 months
			}
		}
	}
}
func FowardTag(defaultedCustomerRecords map[int][]CustomerRecord) {
	var wg sync.WaitGroup
	for _, records := range defaultedCustomerRecords {
		wg.Add(1)
		go func(records []CustomerRecord) {
			defer wg.Done()
			CheckForwardDefaults(records)
		}(records)
	}
	wg.Wait()

}

// FilterRecordsByFwdDefault filters and keeps only the records with FwdDefault set to true.
func FilterRecordsByFwdDefault(records []CustomerRecord) []CustomerRecord {
	var filteredRecords []CustomerRecord
	for _, record := range records {
		if record.FwdDefault {
			filteredRecords = append(filteredRecords, record)
		}
	}
	return filteredRecords
}


func UpdateCustomerRecords(customerRecords []CustomerRecord) {
	db := OpenDB()

	var wg sync.WaitGroup
	for _, record := range customerRecords {
		wg.Add(1)
		go func(record CustomerRecord) {
			defer wg.Done()

			// Update the database record for ForwardDefTag and ForwardDefaultDate
			result := db.Model(&CustomerRecord{}).Where("id = ?", record.ID).
				Update("fwd_default", record.FwdDefault).
				Update("fwd_default_date", record.FwdDefaultDate)

			if result.Error != nil {
				fmt.Printf("Error updating record with ID %d: %v\n", record.ID, result.Error)
			}
		}(record)
	}

	wg.Wait()
	fmt.Println("Customer records updated successfully.")
}
func BulkInsertCustomerRecords(records []CustomerRecord) {
    // Use the Create method to perform bulk insert

	db := OpenDB()
    result := db.Save(&records)
    if result.Error != nil {
        fmt.Printf("Error performing bulk insert: %v\n", result.Error)
        return
    }

    fmt.Printf("Bulk insert completed successfully. Inserted %d records.\n", len(records))
}