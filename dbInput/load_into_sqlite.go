package dbInput

import (
	"encoding/csv"
	"fmt"
	"odr/dataCreation"
	"os"
	"strconv"
	"sync"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type CustomerRecord struct {
	ID         uint `gorm:"primaryKey"`
	CustID     int
	FacID      int
	ObsDate    time.Time
	DPD        int
	FwdDefault bool
	LagDate time.Time
}

func parseRecord(recordFields []string) CustomerRecord {

	record := CustomerRecord{}

	// Parse CustID from recordFields[0]
	custID, err := strconv.Atoi(recordFields[0])
	if err != nil {
		fmt.Println("Error parsing CustID:", err)
		// Skip this record and move to the next one
	}
	record.CustID = custID

	// Parse FacID from recordFields[1]
	facID, err := strconv.Atoi(recordFields[1])
	if err != nil {
		fmt.Println("Error parsing FacID:", err)
		// Skip this record and move to the next one
	}
	record.FacID = facID

	// Parse ObsDate from recordFields[2]
	obsDateStr := recordFields[2]
	obsDate, err := time.Parse("2-Jan-06", obsDateStr)
	if err != nil {
		fmt.Println("Error parsing ObsDate:", err)
		// Skip this record and move to the next one
	}
	record.ObsDate = obsDate

	// Parse DPD from recordFields[3]
	dpd, err := strconv.Atoi(recordFields[3])
	if err != nil {
		fmt.Println("Error parsing DPD:", err)
	}
	record.DPD = dpd

	// Initialize FwdDefault to false
	record.FwdDefault = false

	return record

}
func MigrateToDb(csvFilePath string) {
	// Initialize GORM and open a connection to the SQLite database
	db, err := gorm.Open(sqlite.Open("mydatabase.db"), &gorm.Config{SkipDefaultTransaction: true,})
	if err != nil {
		fmt.Println("Failed to connect to the database:", err)
		return
	}

	// AutoMigrate will create the "records" table if it doesn't exist
	// Check if the table exists
	if db.Migrator().HasTable(&CustomerRecord{}) {
		// Delete the table
		err = db.Migrator().DropTable(&CustomerRecord{})
		if err != nil {
			fmt.Println("Error dropping the table:", err)
			return
		}
		fmt.Println("Table deleted successfully.")
	} else {
		fmt.Println("Table doesn't exist.")
	}

	db.AutoMigrate(&CustomerRecord{})
	// Open the CSV file
	file, err := os.Open(csvFilePath)
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return
	}
	defer file.Close()
	recordCount := 0
	batchSize := 5000 // Adjust this as needed
	recordsToInsert := make([]CustomerRecord, 0, batchSize)
	targetIncrement := 10000
	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read and insert records into the database
	for {
		recordFields, err := reader.Read()
		if err != nil {
			break // End of file
		}
		recordCount++

		record := parseRecord(recordFields)
		recordsToInsert = append(recordsToInsert, record)

		if len(recordsToInsert) == batchSize {
			// Batch insert
			result := db.Create(&recordsToInsert)
			if result.Error != nil {
				fmt.Println("Error inserting records:", result.Error)
			}
			// fmt.Printf("Processed %d records\n", recordCount)
			recordsToInsert = recordsToInsert[:0] // Clear the slice
		}
		if recordCount%targetIncrement == 0 {
			fmt.Printf("Target Processed %d records\n", recordCount)
		}
	}
	// Insert any remaining records
	if len(recordsToInsert) > 0 {
		result := db.Create(&recordsToInsert)
		if result.Error != nil {
			fmt.Println("Error inserting records:", result.Error)
		}
	}

	fmt.Println("Data loaded into the database successfully.")
}

// func migrateToDb(csvFilePath string) {
// 	// Initialize GORM and open a connection to the SQLite database
// 	db, err := gorm.Open(sqlite.Open("mydatabase.db"), &gorm.Config{})
// 	if err != nil {
// 		fmt.Println("Failed to connect to the database:", err)
// 		return
// 	}

// 	// AutoMigrate will create the "records" table if it doesn't exist
// 	// Check if the table exists
// 	if db.Migrator().HasTable(&CustomerRecord{}) {
// 		// Delete the table
// 		err = db.Migrator().DropTable(&CustomerRecord{})
// 		if err != nil {
// 			fmt.Println("Error dropping the table:", err)
// 			return
// 		}
// 		fmt.Println("Table deleted successfully.")
// 	} else {
// 		fmt.Println("Table doesn't exist.")
// 	}

// 	db.AutoMigrate(&CustomerRecord{})
// 	// Open the CSV file
// 	file, err := os.Open(csvFilePath)
// 	if err != nil {
// 		fmt.Println("Error opening CSV file:", err)
// 		return
// 	}
// 	defer file.Close()
// 	recordCount := 0
// 	targetIncrement := 1000000
// 	// Create a CSV reader
// 	reader := csv.NewReader(file)
// 	// Read and insert records into the database
// 	for {
// 		recordFields, err := reader.Read()
// 		if err != nil {
// 			break // End of file
// 		}
// 		recordCount++
// 		if recordCount%targetIncrement == 0 {
// 			fmt.Printf("Processed %d records\n", recordCount)
// 		}
// 		record := parseRecord(recordFields)
// 		// Insert the record into the database
// 		result := db.Create(&record)
// 		if result.Error != nil {
// 			fmt.Println("Error inserting record:", result.Error)
// 		}
// 	}

// 	fmt.Println("Data loaded into the database successfully.")
// }

func processRecord(fields []string, wg *sync.WaitGroup, db *gorm.DB) {
	defer wg.Done() // Decrement the WaitGroup when done

	record := parseRecord(fields)

	result := db.Create(&record)
	if result.Error != nil {
		fmt.Println("Error inserting record:", result.Error)
	}
}

func migrateToDb_concurrent(csvFilePath string) {
	// Initialize GORM and open a connection to the SQLite database

	dbName := "mydatabase.db"
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		fmt.Println("Failed to connect to the database:", err)
		return
	}

	// AutoMigrate will create the "records" table if it doesn't exist
	// Check if the table exists
	if db.Migrator().HasTable(&CustomerRecord{}) {
		// Delete the table
		err = db.Migrator().DropTable(&CustomerRecord{})
		if err != nil {
			fmt.Println("Error dropping the table:", err)
			return
		}
		fmt.Println("Table deleted successfully.")
	} else {
		fmt.Println("Table doesn't exist.")
	}

	db.AutoMigrate(&CustomerRecord{})
	// Open the CSV file
	file, err := os.Open(csvFilePath)
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)
	var wg sync.WaitGroup

	// Read and insert records into the database
	for {
		recordFields, err := reader.Read()
		if err != nil {
			break // End of file
		}
		wg.Add(1) // Increment the WaitGroup for each record

		go processRecord(recordFields, &wg, db)
	}

	wg.Wait() // Wait for all goroutines to finish
	fmt.Println("Data loaded into the database successfully.")
}



func GenerateDataAndMigrateToDB() {

	customerCount := 1000 
	dataCreation.GenerateDummyData(customerCount)
	// go checkExpire()
	// http.HandleFunc("/", handler) // http://127.0.0.1:8080/Go
	// http.ListenAndServe(":8080", nil)
	// nonmulti.Old_forward_def_tagger()
	// forward_default_tagger()
	// migrateToDb("default_flag_generated_large.csv")
	// forward_default_tagger_big_data()
	start := time.Now() // Record the start time
	// migrateToDb_concurrent("default_flag_generated.csv")
	MigrateToDb("default_flag_generated.csv")
	// migrateToDb("default_flag_generated_large.csv")
	elapsed := time.Since(start) // Calculate the elapsed time
	fmt.Printf("Time taken: %s\n", elapsed)
	// start = time.Now() // Record the start time
	// elapsed = time.Since(start) // Calculate the elapsed time
	// fmt.Printf("Time taken: %s\n", elapsed)
	// start = time.Now() // Record the start time
	// // GenerateDummyData()
	// // forward_default_tagger()
	// elapsed = time.Since(start) // Calculate the elapsed time
	// fmt.Printf("Time taken: %s\n", elapsed)
}