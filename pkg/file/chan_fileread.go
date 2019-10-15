package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"sync"
)

type Metrics struct {
	Filename     string
	CountNumeric int
	CountAlpha   int
	CountMixed   int
}

func main() {

	files := []string{
		"test1.log",
		"test2.log",
		"test3.log",
		"test4.log",
	}

	Greed(files)
}

func Greed(files []string) {

	ch := make(chan Metrics, len(files))

	// WaitGroup expand from []files size
	var wg sync.WaitGroup

	for _, f := range files {
		wg.Add(1)
		go Worker(f, &wg, ch)
	}

	wg.Wait()
	close(ch)

	// Summarize output from all files
	var reduce = Metrics{"All", 0, 0, 0}
	for c := range ch {
		reduce.CountNumeric += c.CountNumeric
		reduce.CountAlpha += c.CountAlpha
		reduce.CountMixed += c.CountMixed
	}

	fmt.Printf("%+v\n", reduce)
}

func Worker(filename string, wg *sync.WaitGroup, ch chan Metrics) {

	var result = Metrics{filename, 0, 0, 0}

	f, err := os.Open(path.Join("_data", filename))
	if err != nil {
		panic(err)
	}

	defer f.Close()
	defer wg.Done()

	reader := bufio.NewReader(f)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		Pipelines(line, &result)
	}

	ch <- result
}

func Pipelines(line []byte, result *Metrics) {
	// Match numeric
	numeric := regexp.MustCompile(`^\d+$`)
	if numeric.Match(line) {
		(*result).CountNumeric += 1
	}

	// Match alpha
	alpha := regexp.MustCompile(`^[a-zA-Z]+$`)
	if alpha.Match(line) {
		(*result).CountAlpha += 1
	}

	// Match mixed
	mixed := regexp.MustCompile(`^(\[a-zA-Z]+\d|\d+[a-zA-Z])`)
	if mixed.Match(line) {
		(*result).CountMixed += 1
	}
}
