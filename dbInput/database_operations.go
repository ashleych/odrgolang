package dbInput

import (
	"odr/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SubsetDefaultedCustomers(defaultRecords *[]CustomerRecord) {

	db, _ := gorm.Open(sqlite.Open("mydatabase.db"), &gorm.Config{})
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
