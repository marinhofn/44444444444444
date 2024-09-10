package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

// read a file from a filepath and return a slice of bytes
func readFile(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		return nil, err
	}
	return data, nil
}

// sum all bytes of a file
func sum(filePath string, wg *sync.WaitGroup, results chan<- struct {
	filePath string
	sum      int
}) {
	defer wg.Done()
	data, err := readFile(filePath)
	if err != nil {
		return
	}

	_sum := 0
	for _, b := range data {
		_sum += int(b)
	}

	results <- struct {
		filePath string
		sum      int
	}{filePath, _sum}
}

// main function
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file1> <file2> ...")
		return
	}

	var wg sync.WaitGroup
	results := make(chan struct {
		filePath string
		sum      int
	}, len(os.Args)-1)

	for _, path := range os.Args[1:] {
		wg.Add(1)
		go sum(path, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results and compute totalSum
	sums := make(map[int][]string)
	var totalSum int64
	for result := range results {
		totalSum += int64(result.sum)
		sums[result.sum] = append(sums[result.sum], result.filePath)
	}

	fmt.Println(totalSum)

	for sum, files := range sums {
		if len(files) > 1 {
			fmt.Printf("Sum %d: %v\n", sum, files)
		}
	}
}
