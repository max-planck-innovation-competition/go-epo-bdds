package epo_docdb

import (
	"bytes"
	"encoding/xml"
	log "github.com/sirupsen/logrus"
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

// ReadFile reads a file and returns the ExchangeDocument
func ReadFile(filepath string) (doc *Exchangedocument, err error) {
	// read file
	data, err := os.ReadFile(filepath)
	if err != nil {
		log.WithError(err).Error("failed to read file")
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
		log.WithError(err).Error("failed to unmarshal xml")
		return nil, err
	}
	return &exchangeObject, err
}
