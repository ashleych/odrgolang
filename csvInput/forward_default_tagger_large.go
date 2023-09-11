package csvInput

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func checkForwardDefaultTagEligibility(defaultedRecords []Record, obsDate time.Time) bool {
	for _, record := range defaultedRecords {
		// Check if obsDate is less than any of the lagDate values
		if obsDate.After(record.LagDate) && obsDate.Before((record.ObsDate)) {
			return true
		}

	}
	return false
}

func identify_defaults(inputFilePath string, defaultCustomerRecords map[int][]Record, outputFilePath string) {
	// inputFilePath := "./default_flag_generated.csv"
	// outputFilePath := "output_dataset_1.csv"

	// Open and read the input dataset file
	inputFile, _ := open_file(inputFilePath)

	defer inputFile.Close()
	defaultFile, _ := create_file(outputFilePath)

	defer defaultFile.Close()

	defaultCsvWriter := csv.NewWriter(defaultFile)
	// defer csvWriter.Flush()
	defer defaultCsvWriter.Flush()

	headers := [...]string{"CustID", "FacID", "ObsDate", "DPD", "FwdDefault", "LagDate"}
	defaultCsvWriter.Write(headers[:]) // Write the headers to the CSV file
	// Initialize a map to store records for each customer

	scanner := bufio.NewScanner(inputFile)
	defaultCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")

		if len(parts) != 4 {
			fmt.Println("Invalid data format:", line)
			continue
		}
		dpd, _ := strconv.Atoi(parts[3])
		if dpd > 90 {

			custID, _ := strconv.Atoi(parts[0])
			facID, _ := strconv.Atoi(parts[1])
			obsDateStr := parts[2]

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
				FwdDefault: false,                      // Initialize to false
				LagDate:    obsDate.AddDate(0, -12, 0), //Find date 12 months prior, assuming observaton window is 12 months

			}
			if record.DPD > 90 {
				defaultCount++
			}
			// Store the record in the map for the corresponding customer
			defaultCustomerRecords[custID] = append(defaultCustomerRecords[custID], record)

			writeRecordToCSV(record, defaultCsvWriter)
			// Process the record as needed (e.g., store it, print it, etc.)
			// fmt.Printf("CustID: %d, FacID: %d, ObsDate: %s, DPD: %d, FwdDefault: %v\n", record.CustID, record.FacID, record.ObsDate, record.DPD, record.FwdDefault)
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading input file:", err)
		}
	}

}

func printCustomerRecords(customerRecords map[int][]Record) {
	for custID, records := range customerRecords {
		fmt.Printf("Customer ID: %d\n", custID)
		for _, record := range records {
			fmt.Printf("  FacID: %d, ObsDate: %s, DPD: %d, FwdDefault: %v\n", record.FacID, record.ObsDate.Format("2-Jan-06"), record.DPD, record.FwdDefault)
		}
	}
}

func checkForwardDefaults(customerRecords []Record) {
	for i := range customerRecords {
		// Check if the current record is within the next 12 months
		forwardDate := customerRecords[i].getForwardDate()
		for _, futureRecord := range customerRecords[i+1:] {
			if futureRecord.ObsDate.Before(forwardDate.AddDate(0, 12, 0)) {
				// Check if there is a DPD greater than 90 within 12 months
				if futureRecord.DPD > 90 {
					// Set the FwdDefault field to true for the current record
					customerRecords[i].FwdDefault = true
					break // No need to continue checking
				}
			} else {
				break // No need to check further if future records are outside 12 months
			}
		}
	}
}

func (r Record) getForwardDate() time.Time {
	// obsDate, _ := time.Parse("2-Jan-06", r.ObsDate)
	obsDate := r.ObsDate
	return obsDate.AddDate(0, 12, 0)
}

func writeRecordsToCSV(customerRecords map[int][]Record, csvWriter *csv.Writer) {
	for _, records := range customerRecords {
		for _, record := range records {
			csvWriter.Write([]string{
				strconv.Itoa(record.CustID),
				record.ObsDate.Format("2-Jan-06"),
				strconv.Itoa(record.FacID),
				strconv.Itoa(record.DPD),
				strconv.FormatBool(record.FwdDefault),
			})
		}
	}
}

func writeRecordToCSV(customerRecord Record, csvWriter *csv.Writer) {
	csvWriter.Write([]string{
		strconv.Itoa(customerRecord.CustID),
		strconv.Itoa(customerRecord.FacID),
		customerRecord.ObsDate.Format("2-Jan-06"),
		strconv.Itoa(customerRecord.DPD),
		strconv.FormatBool(customerRecord.FwdDefault),
		customerRecord.LagDate.Format("2-Jan-06"),
	})
}
func writeHeaderToCSV(header, csvWriter *csv.Writer) {

}
