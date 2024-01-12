package epo_docdb

import (
	"bytes"
	"encoding/xml"
	"log/slog"
	"os"
	"strings"
)

// Trimmer is a xml.Decoder that trims the xml.CharData
// https://stackoverflow.com/questions/54096876/trimspaces-for-all-xml-text
type Trimmer struct {
	dec *xml.Decoder
}

// Token returns the next token
// space is trimmed
func (tr Trimmer) Token() (xml.Token, error) {
	t, err := tr.dec.Token()
	if cd, ok := t.(xml.CharData); ok {
		t = xml.CharData(bytes.TrimSpace(cd))
	}
	return t, err
}

// ParseXmlFileToStruct reads a file and returns the ExchangeDocument
func ParseXmlFileToStruct(filepath string) (doc *Exchangedocument, err error) {
	logger := slog.With("filepath", filepath)
	// read file
	data, err := os.ReadFile(filepath)
	if err != nil {
		logger.With("err", err).Error("failed to read file")
		return nil, err
	}
	// replace bytes
	xmlString := strings.Replace(string(data), "<exch:", "<", -1)
	xmlString = strings.Replace(xmlString, "</exch:", "</", -1)
	// parse data from cml in to exchange object
	var exchangeObject Exchangedocument

	// unmarshall xml with a not strict decoder
	d := xml.NewDecoder(strings.NewReader(xmlString))
	d.Strict = false // some parts
	d = xml.NewTokenDecoder(Trimmer{d})
	err = d.Decode(&exchangeObject)
	if err != nil {
		logger.With("err", err).Error("failed to unmarshall xml")
		return nil, err
	}
	return &exchangeObject, err
}
