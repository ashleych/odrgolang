package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)
func forward_default_tagger() {
	inputFilePath := "./default_flag_generated.csv"
	outputFilePath := "output_dataset_1.csv"

	// Open and read the input dataset file
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		fmt.Println("Error opening input file:", err)
		return
	}
	defer inputFile.Close()

	// Create and open the output dataset file for writing
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	// Initialize a CSV writer for the output file
	csvWriter := csv.NewWriter(outputFile)
	defer csvWriter.Flush()

	// Initialize a map to store records for each customer
	customerRecords := make(map[int][]Record)

	scanner := bufio.NewScanner(inputFile)
	defaultCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")

		if len(parts) != 4 {
			fmt.Println("Invalid data format:", line)
			continue
		}

		custID, _ := strconv.Atoi(parts[0])
		facID, _ := strconv.Atoi(parts[1])
		obsDateStr := parts[2]
		dpd, _ := strconv.Atoi(parts[3])

		// Parse the observation date
		obsDate, err := time.Parse("2-Jan-06", obsDateStr)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			continue
		}

		// Create a record for the current row
		record := Record{
			CustID:     custID,
			FacID:      facID,
			ObsDate:    obsDate,
			DPD:        dpd,
			FwdDefault: false, // Initialize to false
		}
		if record.DPD > 90 {
			defaultCount++
		}
		// Store the record in the map for the corresponding customer
		// customerRecords[custID] = append(customerRecords[custID], record)

		// Process the record as needed (e.g., store it, print it, etc.)
		// fmt.Printf("CustID: %d, FacID: %d, ObsDate: %s, DPD: %d, FwdDefault: %v\n", record.CustID, record.FacID, record.ObsDate, record.DPD, record.FwdDefault)
	}
	fmt.Println("Finished reading")
	fmt.Println("Def count", defaultCount)
	fmt.Println("Customerrecords", customerRecords)
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input file:", err)
	}
	// var wg sync.WaitGroup
	// for _, records := range customerRecords {
	// 	wg.Add(1)
	// 	go func(records []Record) {
	// 		defer wg.Done()
	// 		checkForwardDefaults(records)
	// 	}(records)
	// }

	// // Wait for all workers to finish
	// wg.Wait()
	// // Check for forward defaults for each customer
	// // checkForwardDefaults(customerRecords)
	// // fmt.Println("Output being printed")
	// // printCustomerRecords(customerRecords)
	// // Write the records to the output CSV file
	// writeRecordsToCSV(customerRecords, csvWriter)

}