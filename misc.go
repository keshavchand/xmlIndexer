package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

type BytesReader struct {
	reader io.Reader
	Count  int
}

func (b *BytesReader) Read(p []byte) (n int, err error) {
	r, e := b.reader.Read(p)
	if e != nil {
		b.Count += r
	}

	return r, e
}

func ReadEntireFileLineByLine(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	fileBuffer := bufio.NewScanner(file)
	var data []string
	for fileBuffer.Scan() {
		data = append(data, fileBuffer.Text())
	}

	return data
}

func Verifier(indexer Indexer) {
	for k, v := range indexer.CharaterSet {
		if len(v) == 0 {
			panic(fmt.Sprintf("CharaterSet %s is empty", k))
		}
	}
}
