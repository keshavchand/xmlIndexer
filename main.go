package main

import (
	"bytes"
	"encoding/xml"
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	method := flag.String("method", "index", "index or search")
	file := flag.String("file", "", "file to index")
	index := flag.String("index", "", "index file")
	flag.Parse()

	if *file == "" {
		log.Fatal("file is required")
	}

	if *index == "" {
		log.Fatal("file is required")
	}

	if *method == "index" {

		startTime := time.Now()
		CreateIndexer(*file, *index)
		endTime := time.Now()

		log.Println("Indexing took", endTime.Sub(startTime))
	} else if *method == "search" {

		indexFileCharacter, err := os.Open("character_" + *index)
		if err != nil {
			log.Fatal(err)
		}
		indexFileData, err := os.Open("data_" + *index)
		if err != nil {
			log.Fatal(err)
		}

		indexer, err := DeserializeIndexer(indexFileCharacter, indexFileData)
		if err != nil {
			log.Fatal(err)
		}

		data := ReadEntireFileLineByLine(*file)
		searcher := NewSearcher(indexer)
		pos, col := searcher.UiSeachString(data)

		log.Printf("[%d:%d]", pos, col)
	}
}

func CreateIndexer(filename, indexerFile string) {
	contents, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	decoder := xml.NewDecoder(bytes.NewReader(contents))
	indexer := NewIndexer()

	for {
		token, err := Next(decoder)
		parserToken := token.Token
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			log.Fatal(err)
		}

		switch st := parserToken.(type) {
		case xml.StartElement, xml.EndElement:
			err := indexer.IndexTag(st, token.Line, token.Col)
			if err != nil {
				log.Fatal(err)
			}
		case xml.CharData:
			err := indexer.StoreString(string(st), token.Line, token.Col)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	log.Println("Indexed", len(indexer.CharaterSet), "characters")

	fileCharacter, err := os.OpenFile("character_"+indexerFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fileData, err := os.OpenFile("data_"+indexerFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	err = indexer.Serialize(fileCharacter, fileData)
	if err != nil {
		log.Fatal(err)
	}
}
