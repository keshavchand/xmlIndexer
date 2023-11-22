package main

import (
	"encoding/gob"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
)

type Indexer struct {
	currentElement []Data

	CharaterSet map[string][]Pos
	DataIndex   map[string][]Data
}

func NewIndexer() *Indexer {
	return &Indexer{
		CharaterSet: make(map[string][]Pos),
		DataIndex:   make(map[string][]Data),
	}
}

func DeserializeIndexer(rCharater, rData io.Reader) (*Indexer, error) {
	indexer := NewIndexer()
	err := gob.NewDecoder(rCharater).Decode(&indexer.CharaterSet)
	if err != nil {
		return nil, err
	}

	err = gob.NewDecoder(rData).Decode(&indexer.DataIndex)
	if err != nil {
		return nil, err
	}
	return indexer, nil
}

func (i *Indexer) Serialize(wCharater, wData io.Writer) error {
	err := gob.NewEncoder(wCharater).Encode(i.CharaterSet)
	if err != nil {
		return err
	}
	err = gob.NewEncoder(wData).Encode(i.DataIndex)
	if err != nil {
		return err
	}
	return nil
}

func (i *Indexer) IndexTag(token xml.Token, pos, col int) error {
	switch token := token.(type) {
	case xml.StartElement:
		i.currentElement = append(i.currentElement, Data{
			Token:     token,
			StartLine: pos,
		})

	case xml.EndElement:
		if len(i.currentElement) == 0 {
			return errors.New(fmt.Sprint("unexpected token: ", token.Name.Local, pos, ":", col))
		}
		lastElement := i.currentElement[len(i.currentElement)-1]
		st := lastElement.Token
		if st.Name.Local != token.Name.Local {
			return errors.New("unexpected token")
		}

		lastElement.EndLine = pos
		lastElement.EndCol = col
		i.currentElement = i.currentElement[:len(i.currentElement)-1]
		i.DataIndex[st.Name.Local] = append(i.DataIndex[st.Name.Local], lastElement)
	}
	return nil
}

func (i *Indexer) StoreString(s string, pos, col int) error {
	i.CharaterSet[s] = append(i.CharaterSet[s], Pos{pos, col, len(s)})
	return nil
}

type Pos struct {
	StartLine, StartCol int
	Size                int
}

type Data struct {
	Token               xml.StartElement
	StartLine, StartCol int
	EndLine, EndCol     int
}
