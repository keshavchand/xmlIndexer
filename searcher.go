package main

import (
	"fmt"
	"log"

	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
)

type Searcher struct {
	Index *Indexer
}

func NewSearcher(index *Indexer) Searcher {
	return Searcher{Index: index}
}

func (s *Searcher) UiSeachString(content []string) (pos, col int) {
	var strings []string
	for str := range s.Index.CharaterSet {
		strings = append(strings, str)
	}
	idx, err := fuzzyfinder.Find(strings, func(i int) string {
		return strings[i]
	})

	if err != nil {
		log.Fatal(err)
	}

	dataPoints := s.Index.CharaterSet[strings[idx]]

	idx, err = fuzzyfinder.Find(dataPoints, func(i int) string {
		line, col := dataPoints[i].StartLine-1, dataPoints[i].StartCol
		return fmt.Sprintf("[%d:%d] %s", line, col, content[line])
	})

	data := dataPoints[idx]
	return data.StartLine, data.Size
}
