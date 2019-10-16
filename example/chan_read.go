package main

import (
	"bufio"
	"fmt"
	"github.com/xemoe/batteries/pkg/file"
	"github.com/xemoe/batteries/pkg/matcher"
	"io"
	"os"
	"sync"
)

type Metrics struct {
	Filename     file.File
	CountNumeric int
	CountAlpha   int
	CountMixed   int
}

func main() {

	walker := file.FileWalker{
		BaseDir:  "_data",
		FileExt:  ".log",
		MaxDepth: 1,
	}

	var wg sync.WaitGroup

	files := walker.List()
	filesch := make(chan Metrics, len(files))

	Run(files, &wg, filesch)

	wg.Wait()
	close(filesch)

	sum := Metrics{"Summary", 0, 0, 0}
	Reduce(&sum, filesch)

	fmt.Print("%+v\n", sum)
}

// Run worker from given file
func Run(files []file.File, wg *sync.WaitGroup, ch chan Metrics) {
	for _, f := range files {
		wg.Add(1)
		go Worker(f, wg, ch, Counter)
	}
}

// Worker open file and process counter by line
// And pass result to metrics channel(write only)
func Worker(filename file.File, wg *sync.WaitGroup, ch chan<- Metrics, counter func(line []byte, result *Metrics)) {

	result := Metrics{filename, 0, 0, 0}

	f, err := os.Open(string(filename))
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

		counter(line, &result)
	}

	ch <- result
}

// Counter by matcher
// And increase metric counter
func Counter(line []byte, result *Metrics) {
	for _, m := range registeredMatchers(result) {
		m.Count(line)
	}
}

func registeredMatchers(result *Metrics) []matcher.Matcher {
	return []matcher.Matcher{
		matcher.NumericMatcher(func() {
			(*result).CountNumeric += 1
		}),
		matcher.AlphaMatcher(func() {
			(*result).CountAlpha += 1
		}),
		matcher.MixedMatcher(func() {
			(*result).CountMixed += 1
		}),
	}
}

// Reduce summarize from chan metrics
func Reduce(sum *Metrics, ch <-chan Metrics) {
	for c := range ch {
		(*sum).CountNumeric += c.CountNumeric
		(*sum).CountAlpha += c.CountAlpha
		(*sum).CountMixed += c.CountMixed
	}
}
