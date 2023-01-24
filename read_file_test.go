package go_epo_bdds

import (
	"encoding/xml"
	"fmt"
	"github.com/max-planck-innovation-competition/go-epo-docdb/pkg/docdb"
	"io/ioutil"
	"strings"
	"testing"
)

func TestReadFile(t *testing.T) {

	// process file
	data, err := ioutil.ReadFile("./pkg/docdb/test-data/AP-1206-A_302101161.xml")
	if err != nil {
		t.Error(err)
	}
	// replace bytes
	xmlString := strings.Replace(string(data), "<exch:", "<", -1)
	xmlString = strings.Replace(xmlString, "</exch:", "</", -1)
	// parse data from cml in to exchange object
	var exchangeObject docdb.Exchangedocument
	err = xml.Unmarshal([]byte(xmlString), &exchangeObject)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(exchangeObject)
}
