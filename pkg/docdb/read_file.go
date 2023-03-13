package docdb

import (
	"encoding/xml"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

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
	err = xml.Unmarshal([]byte(xmlString), &exchangeObject)
	if err != nil {
		log.WithError(err).Error("failed to unmarshal xml")
		return nil, err
	}
	return &exchangeObject, err
}
