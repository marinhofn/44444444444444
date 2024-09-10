package main

import (
	"fmt"
	"io"
	"os"
)

// fileSum reads a file and computes the sum of chunks of 100 bytes
func fileSum(filePath string) ([]uint64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %v", filePath, err)
	}
	defer file.Close()

	var chunks []uint64
	buffer := make([]byte, 100)

	for {
		bytesRead, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				return nil, fmt.Errorf("error reading file %s: %v", filePath, err)
			}
			break
		}

		sum := sum(buffer[:bytesRead])
		chunks = append(chunks, sum)
	}

	return chunks, nil
}

// sum computes the sum of the byte values in the buffer
func sum(buffer []byte) uint64 {
	var result uint64
	for _, b := range buffer {
		result += uint64(b)
	}
	return result
}

// similarity computes the similarity between two slices of chunk sums
func similarity(base, target []uint64) float64 {
	counter := 0
	targetCopy := append([]uint64(nil), target...)

	for _, value := range base {
		for i, v := range targetCopy {
			if v == value {
				counter++
				targetCopy = append(targetCopy[:i], targetCopy[i+1:]...)
				break
			}
		}
	}

	return float64(counter) / float64(len(base))
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <filepath1> <filepath2> ...")
		os.Exit(1)
	}

	// Map to store fingerprints for each file
	fileFingerprints := make(map[string][]uint64)

	// Calculate fingerprints for each file
	for _, path := range os.Args[1:] {
		fingerprint, err := fileSum(path)
		if err != nil {
			fmt.Printf("Error processing file %s: %v\n", path, err)
			os.Exit(1)
		}
		fileFingerprints[path] = fingerprint
	}

	// Compare each pair of files
	for i := 0; i < len(os.Args)-1; i++ {
		for j := i + 1; j < len(os.Args)-1; j++ {
			file1 := os.Args[i+1]
			file2 := os.Args[j+1]
			fingerprint1 := fileFingerprints[file1]
			fingerprint2 := fileFingerprints[file2]
			similarityScore := similarity(fingerprint1, fingerprint2)
			fmt.Printf("Similarity between %s and %s: %.6f%%\n", file1, file2, similarityScore*100)
		}
	}
}
