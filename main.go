package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

func CountWordsInFile(val string, ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(val)
	if err != nil {
		log.Printf("failed to open file %s: %v\n", val, err)
		ch <- 0
		return
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		log.Printf("failed to read file %s: %v\n", val, err)
		ch <- 0
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(string(b)))
	scanner.Split(bufio.ScanWords)

	count := 0
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		log.Printf("scanner error in file %s: %v\n", val, err)
	}

	ch <- count
}

func main() {
	var sv []string
	for i := 1; i <= 44; i++ {
		file := fmt.Sprintf(`C:\TxtFiles\sample%d.txt`, i)
		sv = append(sv, file)
	}

	startTime := time.Now()

	var wg sync.WaitGroup
	// use buffered channel
	ch := make(chan int, len(sv))

	for _, val := range sv {
		wg.Add(1)
		go CountWordsInFile(val, ch, &wg)
	}

	wg.Wait()
	close(ch)

	totalWords := 0
	for count := range ch {
		totalWords += count
	}

	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)

	fmt.Printf("Total word count across all files: %d\n", totalWords)
	fmt.Println("Total time taken to read and count 44 files is", elapsedTime)
}
