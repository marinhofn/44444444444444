package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

func readFile(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %v", filePath, err)
	}
	return data, nil
}

func calculateSimilarity(data1, data2 []byte, chunkSize int) float64 {
	chunks1 := len(data1) / chunkSize
	chunks2 := len(data2) / chunkSize
	minChunks := chunks1
	if chunks2 < chunks1 {
		minChunks = chunks2
	}

	similarChunks := 0
	for i := 0; i < minChunks; i++ {
		chunk1 := data1[i*chunkSize : (i+1)*chunkSize]
		chunk2 := data2[i*chunkSize : (i+1)*chunkSize]
		if string(chunk1) == string(chunk2) {
			similarChunks++
		}
	}

	return (float64(similarChunks) / float64(minChunks)) * 100
}

func partialSimilarity(file1, file2 string, chunkSize int, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	data1, err := readFile(file1)
	if err != nil {
		fmt.Printf("Error processing file %s: %v\n", file1, err)
		return
	}
	data2, err := readFile(file2)
	if err != nil {
		fmt.Printf("Error processing file %s: %v\n", file2, err)
		return
	}

	similarity := calculateSimilarity(data1, data2, chunkSize)
	results <- fmt.Sprintf("Similarity between %s and %s: %.6f%%", file1, file2, similarity)
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <file1> <file2> ...")
		return
	}

	var wg sync.WaitGroup
	results := make(chan string, len(os.Args)*(len(os.Args)-1)/2)

	chunkSize := 1024

	for i := 1; i < len(os.Args); i++ {
		for j := i + 1; j < len(os.Args); j++ {
			wg.Add(1)
			go partialSimilarity(os.Args[i], os.Args[j], chunkSize, &wg, results)
		}
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result)
	}
}
