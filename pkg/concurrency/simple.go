package concurrency

import (
	"bufio"
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
	sync.Mutex
}

func New(name file.File) Metrics {
	var m Metrics
	m.Filename = name
	return m
}

func (m *Metrics) increaseNumericCounter(i int) {
	m.CountNumeric += i
}

func (m *Metrics) increaseAlphaCounter(i int) {
	m.CountAlpha += i
}
func (m *Metrics) increaseMixedCounter(i int) {
	m.CountMixed += i
}

// Args order
// 1. input
// 2. output
// 3. ref
// 4. chan<-
// 5. <-chan
// 6. func
func Example() {
	files, wg, filesch, sum := Setup()

	Run(files, &wg, filesch)

	wg.Wait()
	close(filesch)

	reduce(&sum, filesch)
}

func Setup() ([]file.File, sync.WaitGroup, chan Metrics, Metrics) {

	var wg sync.WaitGroup
	sum := New(file.File("Summary"))

	walker := file.FileWalker{
		BaseDir:  "_data",
		FileExt:  ".log",
		MaxDepth: 1,
	}

	files := walker.List()
	filesch := make(chan Metrics, len(files))

	return files, wg, filesch, sum
}

// Run worker from given file
func Run(files []file.File, wg *sync.WaitGroup, ch chan Metrics) {
	for _, file := range files {
		wg.Add(1)
		go Worker(file, wg, ch, Counter)
	}
}

// Worker open file and process counter by line
// And pass result to metrics channel(write only)
func Worker(filename file.File, wg *sync.WaitGroup, ch chan<- Metrics, counter func(line interface{}, matchers interface{})) {

	var result Metrics
	result.Filename = filename
	matchers := registeredMatchers(&result)

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

		counter(line, matchers)
	}

	ch <- result
}

// Counter by matcher
// And increase metric counter
func Counter(line interface{}, matchers interface{}) {
	for _, m := range matchers.([]matcher.Matcher) {
		m.Count(line.([]byte))
	}
}

func registeredMatchers(result *Metrics) interface{} {
	return []matcher.Matcher{
		matcher.NumericMatcher(func() {
			(*result).Lock()
			defer (*result).Unlock()

			(*result).increaseNumericCounter(1)
		}),
		matcher.AlphaMatcher(func() {
			(*result).Lock()
			defer (*result).Unlock()

			(*result).increaseAlphaCounter(1)
		}),
		matcher.MixedMatcher(func() {
			(*result).Lock()
			defer (*result).Unlock()

			(*result).increaseMixedCounter(1)
		}),
	}
}

// Reduce summarize from chan metrics
func reduce(sum *Metrics, ch <-chan Metrics) {
	for c := range ch {
		(*sum).increaseNumericCounter(c.CountNumeric)
		(*sum).increaseAlphaCounter(c.CountAlpha)
		(*sum).increaseMixedCounter(c.CountMixed)
	}
}
