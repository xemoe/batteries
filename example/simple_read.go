package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"sync"
)

func main() {

	files := []string{
		"test1.log",
		"test2.log",
		"test3.log",
	}

	Greed(files)
}

func Greed(files []string) {
	var wg sync.WaitGroup
	for _, f := range files {
		wg.Add(1)
		go OneRead(f, &wg)
	}

	wg.Wait()
}

func OneRead(filename string, wg *sync.WaitGroup) {
	f, err := os.Open(path.Join("_data", filename))
	if err != nil {
		panic(err)
	}

	defer f.Close()

	reader := bufio.NewReader(f)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		fmt.Printf("%s \n", line)
	}

	wg.Done()
}
