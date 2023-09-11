package csvInput

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func create_file(inputFilePath string) (*os.File, error) {
	outputFile, err := os.Create(inputFilePath)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return nil, err
	}
	// defer outputFile.Close()
	return outputFile, nil

}

func open_file(inputFilePath string) (*os.File, error) {
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		fmt.Println("Error opening input file:", err)
		return nil, err
	}
	return inputFile, nil

}
func forward_default_tagger_big_data() {
	inputFilePath := "./default_flag_generated.csv"
	defaultFilePath := "./default_output_dataset_1.csv"
	outputFullFilePath := "./output_dataset_full_1.csv"
	defaultCustomerRecords := make(map[int][]Record)

	identify_defaults(inputFilePath, defaultCustomerRecords, defaultFilePath)
	fmt.Printf("%+v", defaultCustomerRecords[1][1])
	assignForwardDefaults(inputFilePath, defaultCustomerRecords, outputFullFilePath, false)
}


func assignForwardDefaults(inputFilePath string, defaultCustomerRecords map[int][]Record, outputFilePath string, writeOutput bool) {

	inputFile, _ := open_file(inputFilePath)
	defer inputFile.Close()

	customerRecords := make(map[int][]Record)
	scanner := bufio.NewScanner(inputFile)

	outputFile, _ := create_file(outputFilePath)
	defer outputFile.Close()
	csvWriter := csv.NewWriter(outputFile)
	defer csvWriter.Flush()

	headers := [...]string{"CustID", "FacID", "ObsDate", "DPD", "FwdDefault", "LagDate"}
	csvWriter.Write(headers[:]) // Write the headers to the CSV file
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
		FwdDefault := false //initialise
		var FwdDefaultDate time.Time
		// Parse the observation date
		obsDate, _ := time.Parse("2-Jan-06", obsDateStr)
		defaultedRecords, exists := defaultCustomerRecords[custID]
		if exists {
			if checkForwardDefaultTagEligibility(defaultedRecords, obsDate) {
				FwdDefault = true
				FwdDefaultDate = obsDate
			}
		} else {
			FwdDefault = false
		}

		record := Record{
			CustID:         custID,
			FacID:          facID,
			ObsDate:        obsDate,
			DPD:            dpd,
			FwdDefault:     FwdDefault,
			LagDate:        obsDate.AddDate(0, -12, 0), //Find date 12 months prior, assuming observaton window is 12 months
			FwdDefaultDate: FwdDefaultDate,             //Find date 12 months prior, assuming observaton window is 12 months
		}
		customerRecords[custID] = append(customerRecords[custID], record)
		if writeOutput {
			writeRecordToCSV(record, csvWriter)
		}
	}
}
