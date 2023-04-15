// Package storedata implements a set of utility functions, which are capable of importing/exporting data
// from/to JSON or CSV file
package storedata

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// exportDataToJSON stores generated during benchmarking data to JSON file
func exportDataToJSON(path, filename string, data map[int]map[int]map[int]float64, prefix, indent string) error {

	// marshal data to JSON
	out, err := json.MarshalIndent(data, prefix, indent)
	if err != nil {
		log.Panicf("Something went wrong during marshalling of data into JSON.. %v\n", err)
	}

	// export JSON data to file
	err = os.WriteFile(path+filename+".json", out, 0644)
	if err != nil {
		log.Panicf("Something went wrong when the data were written to the file... %v\n", err)
	}

	return nil
}

// importDataFromJSON imports data from JSON file
func importDataFromJSON(path, filename string) (map[int]map[int]map[int]float64, error) {
	sourceFile, err := os.Open(path + filename + ".json")
	if err != nil {
		return nil, err
	}

	var out map[int]map[int]map[int]float64
	if err := json.NewDecoder(sourceFile).Decode(&out); err != nil {
		return nil, err
	}

	err = sourceFile.Close()
	if err != nil {
		return nil, err
	}

	return out, nil
}

// exportDataToCSV stores generated during benchmarking data to CVS file
func exportDataToCSV(path, filename string, data map[int]map[int]map[int]float64, names ...string) error {

	// creating a new file to store CSV data
	outputFile, err := os.Create(path + filename + ".csv")
	if err != nil {
		return err
	}

	// write the header of the CSV file
	writer := csv.NewWriter(outputFile)
	// setting a delimiter
	writer.Comma = ';'

	header := make([]string, 0)
	header = append(header, names...)
	// write headers to the file
	if err = writer.Write(header); err != nil {
		return err
	}

	// write the rows by iterating over the map
	for d, map1 := range data {
		for a, map2 := range map1 {
			for i, val := range map2 {
				csvRow := make([]string, 0)
				csvRow = append(csvRow, strconv.Itoa(d), strconv.Itoa(a), strconv.Itoa(i),
					strconv.FormatFloat(val, 'f', -1, 64))
				if err = writer.Write(csvRow); err != nil {
					return err
				}
			}
		}
	}

	writer.Flush()

	err = outputFile.Close()
	if err != nil {
		return err
	}

	return nil
}

// importDataFromCSV imports data from CVS file
func importDataFromCSV(path, filename string) (map[int]map[int]map[int]float64, error) {

	// open a file
	file, err := os.Open(path + filename + ".csv")
	if err != nil {
		log.Fatal(err)
	}

	// read csv values
	csvReader := csv.NewReader(file)
	// setting a delimiter
	csvReader.Comma = ';'

	data := make(map[int]map[int]map[int]float64, 0)
	// reading a file
	for {
		rec, err := csvReader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// check if it is a first line (shouldn't be parsed)
		if strings.Contains(rec[0], "Depth") || strings.Contains(rec[1], "Application") ||
			strings.Contains(rec[2], "Instances") || strings.Contains(rec[3], "Time") {
			// go to the next line
			continue
		}
		// convert data
		depth, err := strconv.Atoi(rec[0])
		if err != nil {
			return nil, err
		}
		appNumber, err := strconv.Atoi(rec[1])
		if err != nil {
			return nil, err
		}
		instances, err := strconv.Atoi(rec[2])
		if err != nil {
			return nil, err
		}
		time, err := strconv.ParseFloat(rec[3], 64)
		if err != nil {
			return nil, err
		}

		// check if the map entry exists
		if _, ok := data[depth]; !ok {
			data[depth] = make(map[int]map[int]float64, 0)
		}

		// check if the map entry exists
		if _, ok := data[depth][appNumber]; !ok {
			data[depth][appNumber] = make(map[int]float64, 0)
		}

		// store line in a map
		data[depth][appNumber][instances] = time
	}

	// closing file
	err = file.Close()
	if err != nil {
		return nil, err
	}

	return data, nil
}

// SaveData saves data to a file (both, .csv and .json)
func SaveData(benchmarkedData map[int]map[int]map[int]float64, name string) error {

	err := exportDataToJSON("data/", name, benchmarkedData, "", " ")
	if err != nil {
		log.Panicf("Something went wrong during storing of the data in JSON file... %v\n", err)
		return err
	}

	err = exportDataToCSV("data/", name, benchmarkedData, "Fractal MAIS Depth [-]",
		"Application Number in Fractal MAIS [-]", "Maximum Number of Instances Deployed by Application [-]",
		"Time [us]")
	if err != nil {
		log.Panicf("Something went wrong during storing of the data in CSV file... %v\n", err)
		return err
	}

	return nil
}

// ImportData function imports data from a file
func ImportData(path, fileName string) (map[int]map[int]map[int]float64, error) {

	data := make(map[int]map[int]map[int]float64, 0)
	var err error
	if isJSON(fileName) {
		// cutting out extension
		fileName = strings.ReplaceAll(fileName, ".json", "")
		data, err = importDataFromJSON(path, fileName)
		if err != nil {
			return nil, err
		}
	}

	if isCSV(fileName) {
		// cutting out extension
		fileName = strings.ReplaceAll(fileName, ".csv", "")
		data, err = importDataFromCSV(path, fileName)
		if err != nil {
			return nil, err
		}
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("couldn't read file %v - either a bad descriptor, or"+
			" unmatched extension (accepted only .csv and .json)\n", fileName)
	}

	return data, nil
}

// isJSON checks if file has .json extension
func isJSON(name string) bool {
	return strings.Contains(name, ".json")
}

// isCSV checks if file has .csv extension
func isCSV(name string) bool {
	return strings.Contains(name, ".csv")
}
