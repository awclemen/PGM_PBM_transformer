////////////////////////////////////////////////////////////////////////////////
// Assignment: Project: Learn a New (to You!) Programming Language - part 2
// Filename: segmentation.go
// Author: Andy Clements, Cora Clements
// 
// Course: CSc 372
// Instructor: L. McCann
// TAs: Muaz Ali, Daniel Reynaldo
// Due Date: April 20, 2026
// 
// Description: Segmentation.go takes a PGM file, determines a good segmentation 
//              value from the PGM pixel values, then generates a PBM file based
//              on the segmentation value and saves the PBM file to disk.
//              
// 
// Language: Go
// Ex. Packages: fmt, os, bufio, strconv, strings, math.rand.v2, math
// 
// Deficiencies: written by novices in the Go programming language.
////////////////////////////////////////////////////////////////////////////////
package main

import (
       "fmt"
	"os"
	"bufio"
	"strconv"
	"strings"
	"math/rand/v2"
	"math"
)


////////////////////////////////////////////////////////////////////////////////
// Method Name:createPBMFile
// Method Purpose: go through original pixels, create PBM, save to file
//                 PBM has P1 instead of P2 on first line
//                 PBM has no max value
//                 PBM has .pbm file extension
//                 
// Method Pre-conditions: 
// Method Post-conditions: 
// Method Return Value: 
//                      
// Method Parameters: 
////////////////////////////////////////////////////////////////////////////////
func createPBMFile(filename string, pixelList []int, rows int, columns int, segValue float64) {

}

////////////////////////////////////////////////////////////////////////////////
// Method Name: bestSegmentation
// Method Purpose: Determine the best segmentation value for generating a PBM
//                 from the PGM values.  The Algorithm is the following:
//                 1. randomly pick 10 pixels, generat a mean from values
//                 2. perform the following in a loop
//                 3. partition pixels into two groups, those less than the mean
//                    and the others
//	           4. compute the mean of each partition
//	           5. compute the means of the two means = new_mean
//                 6. end loop when the difference between new_mean and
//                    previous_mean are less than 0.001 or looped 100 times,
//                    otherwise previous_mean = new mean)
// Method Pre-conditions: none
// Method Post-conditions: none
// Method Return Value: float64 segmentation value.
// Method Parameters: int array of pixel values.
////////////////////////////////////////////////////////////////////////////////
func bestSegmentation(pixelList []int) (segValue float64) {
	segValue = 0.0
	previousMean := 0.0
	loopCount, sum := 0, 0
	topSum, topCount, bottomSum, bottomCount := 0, 0, 0, 0
	
	for i := 0; i < 10; i++ {
		randomIndex := rand.IntN(len(pixelList))
		sum = sum + pixelList[randomIndex]
	}
	previousMean = float64(sum)/10

	for {
		for i := 0; i < len(pixelList); i++ {
			if float64(pixelList[i]) < previousMean {
				bottomSum += pixelList[i]
				bottomCount++
			} else {
				topSum += pixelList[i]
				topCount++
			}
		}
		segValue = ((float64(topSum)/float64(topCount)) + (float64(bottomSum)/float64(bottomCount)))/2.0
		//fmt.Fprintf(os.Stdout, "segValue: %f, preMean: %f, loop count: %d\n", segValue, previousMean, loopCount)
		if math.Abs(segValue - previousMean) < 0.001 || loopCount >= 100 {
			return segValue
		} else {
			previousMean = segValue
			loopCount++
			topCount, topSum, bottomCount, bottomSum = 0,0,0,0
		}
	}
	
	return segValue
}

////////////////////////////////////////////////////////////////////////////////
// Method Name: parseIntegers
// Method Purpose: pull integers from string(line) and put them in an int array.
// Method Pre-conditions: none
// Method Post-conditions: none
// Method Return Value: array of pixel values from string, error if present.
// Method Parameters: String of pixel values
////////////////////////////////////////////////////////////////////////////////
func parseIntegers(line string) ([]int, error) {
    fields := strings.Fields(line) 
    nums := make([]int, 0, len(fields)) 

    for _, field := range fields {
        n, err := strconv.Atoi(field) 
        if err != nil {
            return nil, err
        }
        nums = append(nums, n)
    }
    return nums, nil
}

////////////////////////////////////////////////////////////////////////////////
// Method Name: readPGMFile
// Method Purpose: Read pixel data, max brightness value, number of picture rows
//                 and columns from file and return values.
// Method Pre-conditions: working file handle to PGM file.
// Method Post-conditions: none
// Method Return Value: array of pixel values, max brightness, number of rows,
//                      and number of columns
// Method Parameters: File handle to PGM file.
////////////////////////////////////////////////////////////////////////////////
func readPGMFile(f *os.File) ([]int, int, int, int) {
	var pixelList []int
	var maxBright int
	var rows int
	var columns int

	input := bufio.NewScanner(f)

	//TODO: handle potential erros from input.Err()
	// check to see if P2 is firstline
	if input.Scan() {
		fileID := input.Text()
		if fileID != "P2" {
			fmt.Fprintf(os.Stderr, "Error: file not identified as PGM file.\n")
			os.Exit(1)
		}
	}

	// skip comments, find row & column values
	for input.Scan() {
		line := input.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		_, err := fmt.Sscan(line, &rows, &columns)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: row and column designations are not numbers. %v\n", err)
		}
		break;
	}

	// read in max brightness designation
	if input.Scan() {
		var err error
		maxBright, err = strconv.Atoi(input.Text()) 
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: max brightness designation is not a number. %v\n", err)
		}
	}

	// read through the rest of the file and get pixel data.
	for input.Scan() {
		line := input.Text()
		pixelValues, err := parseIntegers(line)
		if err == nil {
			pixelList = append(pixelList, pixelValues...);
		}
	}
	return pixelList, maxBright, rows, columns
}

////////////////////////////////////////////////////////////////////////////////
// Method Name: transformToPBM
// Method Purpose: Open file and pull pixel information from file, use pixel
//                 information to determine the best segmentation value.  Use
//                 that value to generate a PBM file.
// Method Pre-conditions: PGM must be present in working directory
// Method Post-conditions: PBM file written to disk
// Method Return Value: none
//                      
// Method Parameters: filename - PGM file filename
////////////////////////////////////////////////////////////////////////////////
func transformToPBM(filename string) {
	f, err := os.Open(filename)
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: segmentation: %v\n", err)
		os.Exit(1)
	}
	pixelList, maxBright, rows, columns := readPGMFile(f)
	
	//fmt.Fprintf(os.Stdout, "pixelList: %v\n", pixelList)
	//fmt.Fprintf(os.Stdout, "maxBright: %d\n", maxBright)
	//fmt.Fprintf(os.Stdout, "rows: %d\n", rows)
	//fmt.Fprintf(os.Stdout, "columns: %d\n", columns)
	segValue := bestSegmentation(pixelList)
	//fmt.Fprintf(os.Stdout, "Segmentation Value: %f\n", segValue)	
		
	//TODO: parsefile from extension
	//TODO: create PBM file
	createPBMFile(filename, pixelList, rows, columns, segValue) 
	f.Close()
}

////////////////////////////////////////////////////////////////////////////////
// Method Name: main
// Method Purpose: starting program, getting PGM file name from program
//                 arguments
// Method Pre-conditions: PGM file present in directory
// Method Post-conditions: Write ability in the directory
// Method Return Value: none
//                      
// Method Parameters: none
////////////////////////////////////////////////////////////////////////////////
func main() {
	filename := ""
	// first argument is program name, second is the PGM Filename.
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Error: Usage: segmentation <PGMFilename>\n")
		os.Exit(1)
	}
	filename = os.Args[1]
	transformToPBM(filename)
	fmt.Println("File transfer from PGM to PBM complete\n")
}





