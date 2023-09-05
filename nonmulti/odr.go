package nonmulti

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Record struct {
    CustID        int
    FacID         int
    ObsDate       time.Time
    DPD           int
    FwdDefault    bool
}

func Old_forward_def_tagger() {
    inputFilePath := "./default_flag_generated.csv"
    outputFilePath := "output_dataset.csv"

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

        // Store the record in the map for the corresponding customer
        customerRecords[custID] = append(customerRecords[custID], record)

        // Process the record as needed (e.g., store it, print it, etc.)
        // fmt.Printf("CustID: %d, FacID: %d, ObsDate: %s, DPD: %d, FwdDefault: %v\n", record.CustID, record.FacID, record.ObsDate, record.DPD, record.FwdDefault)
    }

    if err := scanner.Err(); err != nil {
        fmt.Println("Error reading input file:", err)
    }

    // Check for forward defaults for each customer
    checkForwardDefaults(customerRecords)

    // Write the records to the output CSV file
    writeRecordsToCSV(customerRecords, csvWriter)
}

func checkForwardDefaults(customerRecords map[int][]Record) {
    for _, records := range customerRecords {
        for i := range records {
            // Check if the current record is within the next 12 months
            forwardDate := records[i].getForwardDate()
            for _, futureRecord := range records[i+1:] {
                if futureRecord.ObsDate.Before(forwardDate.AddDate(0, 12, 0)) {
                    // Check if there is a DPD greater than 90 within 12 months
                    if futureRecord.DPD > 90 {
                        // Set the fwd_default_flag to true for the current record
                        records[i].FwdDefault = true
                        break // No need to continue checking
                    }
                } else {
                    break // No need to check further if future records are outside 12 months
                }
            }
        }
    }
}

func (r Record) getForwardDate() time.Time {
    // obsDate, _ := time.Parse("2-Jan-06", r.ObsDate)
    obsDate:= r.ObsDate
    return obsDate.AddDate(0, 12, 0)
}

func writeRecordsToCSV(customerRecords map[int][]Record, csvWriter *csv.Writer) {
    for _, records := range customerRecords {
        for _, record := range records {
            csvWriter.Write([]string{
                strconv.Itoa(record.CustID),
                strconv.Itoa(record.FacID),
                record.ObsDate.Format("2-Jan-06"),
                strconv.Itoa(record.DPD),
                strconv.FormatBool(record.FwdDefault),
            })
        }
    }
}
