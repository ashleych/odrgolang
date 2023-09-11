package dataCreation

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func GenerateDummyData(customerCount int) {
	// Define the number of customers
	numCustomers := customerCount// Change this to the desired number

	// Create a file for writing
	outputFile, err := os.Create("default_flag_generated.csv")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	// Write the header
	outputFile.WriteString("Cust_ID,Fac_ID,Obs_Date,DPD_new\n")

	// Random number generator
	rand.New(rand.NewSource((time.Now().UnixNano())))
	// rand.Seed(time.Now().UnixNano())

	// Generate dummy data for each customer
	for custID := 1; custID <= numCustomers; custID++ {
		// Define the number of records per customer (you can change this as needed)
		numRecordsPerCustomer := rand.Intn(100) + 1 // Generates a random number between 1 and 100 records per customer

		for i := 0; i < numRecordsPerCustomer; i++ {
			facID := rand.Intn(2) + 1
			obsDate := generateRandomDate()
			dpd := rand.Intn(101) // Generates a random number between 0 and 100 for DPD_new

			row := fmt.Sprintf("%d,%d,%s,%d\n", custID, facID, obsDate, dpd)
			outputFile.WriteString(row)
		}
	}

	fmt.Println("Dummy data generation completed.")
}

func generateRandomDate() string {
	min := time.Date(2014, time.January, 1, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2015, time.February, 28, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min
	sec := rand.Int63n(delta) + min
	randomTime := time.Unix(sec, 0)
	return randomTime.Format("2-Jan-06")
}
