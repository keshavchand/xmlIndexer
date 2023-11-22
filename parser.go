package main

import "encoding/xml"

type ParsedToken struct {
	Token     xml.Token
	Line, Col int
}

func Next(decoder *xml.Decoder) (ParsedToken, error) {
	line, col := decoder.InputPos()
	token, err := decoder.Token()
	if err != nil {
		return ParsedToken{}, err
	}
	return ParsedToken{
		Token: token,
		Line:  line,
		Col:   col,
	}, nil
}
