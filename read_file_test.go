package go_epo_bdds

import (
	"encoding/xml"
	"github.com/max-planck-innovation-competition/go-epo-docdb/pkg/docdb"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strconv"
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

	ass := assert.New(t)
	ass.NoError(err)

	// Exchange-document tag
	ass.Equal("AP", exchangeObject.CountryAttr)
	ass.Equal(20041016, exchangeObject.DateaddeddocdbAttr)
	ass.Equal(20220630, exchangeObject.DateofpreviousexchangeAttr)
	ass.Equal(20221027, exchangeObject.DateoflastexchangeAttr)
	ass.Equal(20030918, exchangeObject.DatepublAttr)
	ass.Equal("1206", exchangeObject.DocnumberAttr)
	ass.Equal("22179393", exchangeObject.FamilyidAttr)
	ass.Equal("381754736", exchangeObject.DocidAttr)
	ass.Equal("YES", exchangeObject.IsrepresentativeAttr)
	ass.Equal("A", exchangeObject.KindAttr)
	ass.Equal("EP", exchangeObject.OriginatingofficeAttr)
	ass.Equal("A", exchangeObject.StatusAttr)

	// bibliographic-data
	// publication-reference
	ass.Equal("publication-reference", exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].XMLName.Local)
	ass.Equal("docdb", exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].DataformatAttr)
	ass.Equal("document-id", exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].Documentid.XMLName.Local)
	ass.Equal("en", exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].Documentid.LangAttr)
	ass.Equal("AP", exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].Documentid.Country)
	ass.Equal("doc-number", exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].Documentid.Docnumber.XMLName.Local)
	ass.Equal("1206", exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].Documentid.Docnumber.Value)
	ass.Equal("kind", exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].Documentid.Kind.XMLName.Local)
	ass.Equal("A", exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].Documentid.Kind.Value)
	ass.Equal(20030918, exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].Documentid.Date)

	ass.Equal("publication-reference", exchangeObject.ExchBibliographicdata.ExchPublicationreference[1].XMLName.Local)
	ass.Equal("epodoc", exchangeObject.ExchBibliographicdata.ExchPublicationreference[1].DataformatAttr)
	ass.Equal("doc-number", exchangeObject.ExchBibliographicdata.ExchPublicationreference[1].Documentid.Docnumber.XMLName.Local)
	ass.Equal("AP1206", exchangeObject.ExchBibliographicdata.ExchPublicationreference[1].Documentid.Docnumber.Value)

	ass.Equal("publication-reference", exchangeObject.ExchBibliographicdata.ExchPublicationreference[2].XMLName.Local)
	ass.Equal("original", exchangeObject.ExchBibliographicdata.ExchPublicationreference[2].DataformatAttr)
	ass.Equal("doc-number", exchangeObject.ExchBibliographicdata.ExchPublicationreference[2].Documentid.Docnumber.XMLName.Local)
	ass.Equal("AP 1206", exchangeObject.ExchBibliographicdata.ExchPublicationreference[2].Documentid.Docnumber.Value)

	// classification-ipc
	ass.Equal("classification-ipc", exchangeObject.ExchBibliographicdata.ExchClassificationipc.XMLName.Local)
	ass.Equal("edition", exchangeObject.ExchBibliographicdata.ExchClassificationipc.Edition.XMLName.Local)
	ass.Equal("7", exchangeObject.ExchBibliographicdata.ExchClassificationipc.Edition.Value)
	ass.Equal("main-classification", exchangeObject.ExchBibliographicdata.ExchClassificationipc.Mainclassification[0].XMLName.Local)
	ass.Equal(" 7A 61K 9/16 A", exchangeObject.ExchBibliographicdata.ExchClassificationipc.Mainclassification[0].Value)

	for i := 0; i <= 17; i++ {
		ass.Equal(strconv.Itoa(i+1), exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[i].SequenceAttr)
		ass.Equal("text", exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[i].Text.XMLName.Local)
	}

	ass.Equal("classifications-ipcr", exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.XMLName.Local)
	ass.Equal("A61K 9/16 20060101A I20051008RMEP ", exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[0].Text.Value)
	ass.Equal("A61K 9/32 20060101ALI20030127BMRU ", exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[1].Text.Value)
	ass.Equal("A61K 9/48 20060101ALI20030127BMRU ", exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[2].Text.Value)
	ass.Equal("A61K 9/50 20060101A I20051008RMEP ", exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[3].Text.Value)
	ass.Equal("A61K 9/54 20060101A I20051110RMEP ", exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[4].Text.Value)
	ass.Equal("A61K 9/62 20060101A I20051110RMEP ", exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[5].Text.Value)
	ass.Equal("A61K 9/64 20060101A I20060521RMUS ", exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[6].Text.Value)
	ass.Equal("A61K 31/22 20060101ALI20051220RMJP ", exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[7].Text.Value)
	ass.Equal("A61K 31/522 20060101A I20051110RMEP ", exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[8].Text.Value)
	ass.Equal("A61K 31/704 20060101ALI20051220RMJP ", exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[9].Text.Value)
	ass.Equal("A61K 31/7048 20060101ALI20051220RMJP ", exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[10].Text.Value)
	ass.Equal("A61K 31/708 20060101ALI20051220RMJP ", exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[11].Text.Value)
	ass.Equal("A61K 47/02 20060101ALI20051220RMJP ", exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[12].Text.Value)
	ass.Equal("A61K 47/14 20060101ALI20051220RMJP ", exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[13].Text.Value)
	ass.Equal("A61K 47/32 20060101ALI20051220RMJP ", exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[14].Text.Value)
	ass.Equal("A61K 47/36 20060101ALI20051220RMJP ", exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[15].Text.Value)
	ass.Equal("A61K 47/38 20060101ALI20051220RMJP ", exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[16].Text.Value)
	ass.Equal("A61P 31/18 20060101A I20051110RMEP ", exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[17].Text.Value)

	// patent classification
	for i := 0; i <= 7; i++ {
		ass.Equal(strconv.Itoa(i+1), exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[i].SequenceAttr)
		ass.Equal("patent-classification", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[i].XMLName.Local)
		ass.Equal("classification-scheme", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[i].Classificationscheme.XMLName.Local)
		ass.Equal("CPCI", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[i].Classificationscheme.SchemeAttr)
		ass.Equal("B", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[i].Classificationstatus)
		ass.Equal("H", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[i].Classificationdatasource)
		ass.Equal("action-date", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[i].Actiondate.XMLName.Local)
	}

	ass.Equal("A61K 9/1652 ", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[0].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[0].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[0].Symbolposition)
	ass.Equal("EP", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[0].Classificationscheme.OfficeAttr)
	ass.Equal(20130101, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[0].Actiondate.Date)

	ass.Equal("A61K 9/485 ", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[1].Classificationsymbol)
	ass.Equal("A", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[1].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[1].Symbolposition)
	ass.Equal("EP", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[1].Classificationscheme.OfficeAttr)
	ass.Equal(20130101, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[1].Actiondate.Date)

	ass.Equal("A61K 9/501 ", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[2].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[2].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[2].Symbolposition)
	ass.Equal("EP", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[2].Classificationscheme.OfficeAttr)
	ass.Equal(20130101, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[2].Actiondate.Date)

	ass.Equal("A61K 9/5015 ", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[3].Classificationsymbol)
	ass.Equal("A", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[3].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[3].Symbolposition)
	ass.Equal("EP", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[3].Classificationscheme.OfficeAttr)
	ass.Equal(20130101, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[3].Actiondate.Date)

	ass.Equal("A61K 9/5026 ", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[4].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[4].Classificationvalue)
	ass.Equal("F", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[4].Symbolposition)
	ass.Equal("EP", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[4].Classificationscheme.OfficeAttr)
	ass.Equal(20130101, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[4].Actiondate.Date)

	ass.Equal("A61K 9/5073 ", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[5].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[5].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[5].Symbolposition)
	ass.Equal("EP", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[5].Classificationscheme.OfficeAttr)
	ass.Equal(20130101, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[5].Actiondate.Date)

	ass.Equal("A61P 31/18 ", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[6].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[6].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[6].Symbolposition)
	ass.Equal("EP", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[6].Classificationscheme.OfficeAttr)
	ass.Equal(20200327, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[6].Actiondate.Date)

	ass.Equal("A61K 9/16 ", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[7].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[7].Classificationvalue)
	ass.Equal("F", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[7].Symbolposition)
	ass.Equal("EP", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[7].Classificationscheme.OfficeAttr)
	ass.Equal(20160901, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[7].Actiondate.Date)

	// application reference
	ass.Equal("application-reference", exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].XMLName.Local)
	ass.Equal("NO", exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].IsrepresentativeAttr)
	ass.Equal("470406", exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].DocidAttr)
	ass.Equal("docdb", exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].DataformatAttr)
	ass.Equal("document-id", exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].Documentid.XMLName.Local)
	ass.Equal("AP", exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].Documentid.Country)
	ass.Equal("doc-number", exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].Documentid.Docnumber.XMLName.Local)
	ass.Equal("2000001988", exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].Documentid.Docnumber.Value)
	ass.Equal("kind", exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].Documentid.Kind.XMLName.Local)
	ass.Equal("A", exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].Documentid.Kind.Value)
	ass.Equal(19980804, exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].Documentid.Date)

	ass.Equal("application-reference", exchangeObject.ExchBibliographicdata.ExchApplicationreference[1].XMLName.Local)
	ass.Equal("epodoc", exchangeObject.ExchBibliographicdata.ExchApplicationreference[1].DataformatAttr)
	ass.Equal("document-id", exchangeObject.ExchBibliographicdata.ExchApplicationreference[1].Documentid.XMLName.Local)
	ass.Equal("doc-number", exchangeObject.ExchBibliographicdata.ExchApplicationreference[1].Documentid.Docnumber.XMLName.Local)
	ass.Equal("AP19200001988", exchangeObject.ExchBibliographicdata.ExchApplicationreference[1].Documentid.Docnumber.Value)

	ass.Equal("application-reference", exchangeObject.ExchBibliographicdata.ExchApplicationreference[2].XMLName.Local)
	ass.Equal("original", exchangeObject.ExchBibliographicdata.ExchApplicationreference[2].DataformatAttr)
	ass.Equal("document-id", exchangeObject.ExchBibliographicdata.ExchApplicationreference[2].Documentid.XMLName.Local)
	ass.Equal("doc-number", exchangeObject.ExchBibliographicdata.ExchApplicationreference[2].Documentid.Docnumber.XMLName.Local)
	ass.Equal("AP/P/2000/001988", exchangeObject.ExchBibliographicdata.ExchApplicationreference[2].Documentid.Docnumber.Value)

	// language of publication
	ass.Equal("language-of-publication", exchangeObject.ExchBibliographicdata.ExchLanguageofpublication.XMLName.Local)
	ass.Equal("en", exchangeObject.ExchBibliographicdata.ExchLanguageofpublication.Value)

	ass.Equal("family-member", exchangeObject.ExchPatentfamily.ExchFamilymember[0].XMLName.Local)

	ass.Equal("en", exchangeObject.ExchAbstract[0].LangAttr)
	ass.Equal("docdba", exchangeObject.ExchAbstract[0].DataformatAttr)
	ass.Equal("national office", exchangeObject.ExchAbstract[0].AbstractsourceAttr)
	ass.Equal(682, len(exchangeObject.ExchAbstract[0].ExchP[0].Value))
	ass.Equal("p", exchangeObject.ExchAbstract[0].ExchP[0].XMLName.Local)
}
