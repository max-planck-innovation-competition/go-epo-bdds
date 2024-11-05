package epo_docdb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestReadFile(t *testing.T) {

	// process file
	exchangeObject, err := ParseXmlFileToStruct("./test-data/AP-1206-A_302101161.xml")
	if err != nil {
		t.Error(err)
	}

	// jsonObject, err := json.MarshalIndent(exchangeObject, "", "  ")
	// fmt.Println(string(jsonObject))

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
	ass.Equal("document-id", exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].Documentid.XMLName.Local)
	ass.Equal("en", exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].Documentid.LangAttr)
	ass.Equal("AP", exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].Documentid.Country)
	ass.Equal("A", *exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].Documentid.Kind)
	ass.Equal(20030918, exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].Documentid.Date)

	ass.Equal("docdb", exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].DataformatAttr)
	ass.Equal("1206", *exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].Documentid.Docnumber)

	ass.Equal("epodoc", exchangeObject.ExchBibliographicdata.ExchPublicationreference[1].DataformatAttr)
	ass.Equal("AP1206", *exchangeObject.ExchBibliographicdata.ExchPublicationreference[1].Documentid.Docnumber)

	ass.Equal("original", exchangeObject.ExchBibliographicdata.ExchPublicationreference[2].DataformatAttr)
	ass.Equal("AP 1206", *exchangeObject.ExchBibliographicdata.ExchPublicationreference[2].Documentid.Docnumber)

	// classification-ipc
	ass.Equal("classification-ipc", exchangeObject.ExchBibliographicdata.ExchClassificationipc.XMLName.Local)
	ass.Equal("edition", exchangeObject.ExchBibliographicdata.ExchClassificationipc.Edition.XMLName.Local)
	ass.Equal("7", exchangeObject.ExchBibliographicdata.ExchClassificationipc.Edition.Value)
	ass.Equal("main-classification", exchangeObject.ExchBibliographicdata.ExchClassificationipc.Mainclassification[0].XMLName.Local)
	ass.Equal("7A 61K   9/16   A", exchangeObject.ExchBibliographicdata.ExchClassificationipc.Mainclassification[0].Value)

	for i := 0; i <= 17; i++ {
		ass.Equal(strconv.Itoa(i+1), exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[i].SequenceAttr)
	}

	ass.Equal("classifications-ipcr", exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.XMLName.Local)
	ass.Equal("A61K   9/16        20060101A I20051008RMEP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[0].Text)
	ass.Equal("A61K   9/32        20060101ALI20030127BMRU", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[1].Text)
	ass.Equal("A61K   9/48        20060101ALI20030127BMRU", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[2].Text)
	ass.Equal("A61K   9/50        20060101A I20051008RMEP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[3].Text)
	ass.Equal("A61K   9/54        20060101A I20051110RMEP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[4].Text)
	ass.Equal("A61K   9/62        20060101A I20051110RMEP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[5].Text)
	ass.Equal("A61K   9/64        20060101A I20060521RMUS", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[6].Text)
	ass.Equal("A61K  31/22        20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[7].Text)
	ass.Equal("A61K  31/522       20060101A I20051110RMEP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[8].Text)
	ass.Equal("A61K  31/704       20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[9].Text)
	ass.Equal("A61K  31/7048      20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[10].Text)
	ass.Equal("A61K  31/708       20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[11].Text)
	ass.Equal("A61K  47/02        20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[12].Text)
	ass.Equal("A61K  47/14        20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[13].Text)
	ass.Equal("A61K  47/32        20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[14].Text)
	ass.Equal("A61K  47/36        20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[15].Text)
	ass.Equal("A61K  47/38        20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[16].Text)
	ass.Equal("A61P  31/18        20060101A I20051110RMEP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[17].Text)

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

	ass.Equal("A61K   9/1652", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[0].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[0].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[0].Symbolposition)
	ass.Equal("EP", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[0].Classificationscheme.OfficeAttr)
	ass.Equal(20130101, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[0].Actiondate.Date)

	ass.Equal("A61K   9/485", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[1].Classificationsymbol)
	ass.Equal("A", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[1].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[1].Symbolposition)
	ass.Equal("EP", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[1].Classificationscheme.OfficeAttr)
	ass.Equal(20130101, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[1].Actiondate.Date)

	ass.Equal("A61K   9/501", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[2].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[2].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[2].Symbolposition)
	ass.Equal("EP", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[2].Classificationscheme.OfficeAttr)
	ass.Equal(20130101, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[2].Actiondate.Date)

	ass.Equal("A61K   9/5015", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[3].Classificationsymbol)
	ass.Equal("A", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[3].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[3].Symbolposition)
	ass.Equal("EP", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[3].Classificationscheme.OfficeAttr)
	ass.Equal(20130101, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[3].Actiondate.Date)

	ass.Equal("A61K   9/5026", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[4].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[4].Classificationvalue)
	ass.Equal("F", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[4].Symbolposition)
	ass.Equal("EP", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[4].Classificationscheme.OfficeAttr)
	ass.Equal(20130101, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[4].Actiondate.Date)

	ass.Equal("A61K   9/5073", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[5].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[5].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[5].Symbolposition)
	ass.Equal("EP", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[5].Classificationscheme.OfficeAttr)
	ass.Equal(20130101, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[5].Actiondate.Date)

	ass.Equal("A61P  31/18", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[6].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[6].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[6].Symbolposition)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[6].Symbolposition)
	ass.Equal("EP", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[6].Classificationscheme.OfficeAttr)
	ass.Equal(20200327, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[6].Actiondate.Date)

	ass.Equal("A61K   9/16", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[7].Classificationsymbol)
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
	ass.Equal("2000001988", *exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].Documentid.Docnumber)
	ass.Equal("A", *exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].Documentid.Kind)
	ass.Equal(19980804, exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].Documentid.Date)

	ass.Equal("application-reference", exchangeObject.ExchBibliographicdata.ExchApplicationreference[1].XMLName.Local)
	ass.Equal("epodoc", exchangeObject.ExchBibliographicdata.ExchApplicationreference[1].DataformatAttr)
	ass.Equal("document-id", exchangeObject.ExchBibliographicdata.ExchApplicationreference[1].Documentid.XMLName.Local)
	ass.Equal("AP19200001988", *exchangeObject.ExchBibliographicdata.ExchApplicationreference[1].Documentid.Docnumber)

	ass.Equal("application-reference", exchangeObject.ExchBibliographicdata.ExchApplicationreference[2].XMLName.Local)
	ass.Equal("original", exchangeObject.ExchBibliographicdata.ExchApplicationreference[2].DataformatAttr)
	ass.Equal("document-id", exchangeObject.ExchBibliographicdata.ExchApplicationreference[2].Documentid.XMLName.Local)
	ass.Equal("AP/P/2000/001988", *exchangeObject.ExchBibliographicdata.ExchApplicationreference[2].Documentid.Docnumber)

	// language of publication
	ass.Equal("language-of-publication", exchangeObject.ExchBibliographicdata.ExchLanguageofpublication.XMLName.Local)
	ass.Equal("en", exchangeObject.ExchBibliographicdata.ExchLanguageofpublication.Value)

	// priority claims
	ass.Equal("priority-claims", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.XMLName.Local)
	ass.Equal("priority-claim", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[0].XMLName.Local)
	ass.Equal("1", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[0].SequenceAttr)
	ass.Equal("docdb", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[0].DataformatAttr)
	ass.Equal("document-id", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[0].Documentid.XMLName.Local)
	ass.Equal("8359798", *exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[0].Documentid.Docnumber)
	ass.Equal("A", *exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[0].Documentid.Kind)
	ass.Equal("US", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[0].Documentid.Country)
	ass.Equal(19980522, exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[0].Documentid.Date)
	ass.Equal("Y", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[0].ExchPriorityactiveindicator)

	ass.Equal("priority-claim", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[1].XMLName.Local)
	ass.Equal("1", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[1].SequenceAttr)
	ass.Equal("epodoc", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[1].DataformatAttr)
	ass.Equal("document-id", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[1].Documentid.XMLName.Local)
	ass.Equal("US19980083597", *exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[1].Documentid.Docnumber)

	ass.Equal("priority-claim", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[2].XMLName.Local)
	ass.Equal("2", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[2].SequenceAttr)
	ass.Equal("docdb", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[2].DataformatAttr)
	ass.Equal("document-id", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[2].Documentid.XMLName.Local)
	ass.Equal("9816128", *exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[2].Documentid.Docnumber)
	ass.Equal("W", *exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[2].Documentid.Kind)
	ass.Equal("US", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[2].Documentid.Country)
	ass.Equal(19980804, exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[2].Documentid.Date)
	ass.Equal("W", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[2].ExchPrioritylinkagetype)
	ass.Equal("N", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[2].ExchPriorityactiveindicator)

	ass.Equal("priority-claim", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[3].XMLName.Local)
	ass.Equal("2", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[3].SequenceAttr)
	ass.Equal("epodoc", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[3].DataformatAttr)
	ass.Equal("document-id", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[3].Documentid.XMLName.Local)
	ass.Equal("WO1998US16128", *exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[3].Documentid.Docnumber)

	ass.Equal("priority-claim", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[4].XMLName.Local)
	ass.Equal("1", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[4].SequenceAttr)
	ass.Equal("original", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[4].DataformatAttr)
	ass.Equal("document-id", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[4].Documentid.XMLName.Local)
	ass.Equal("PCT/US98/16128", *exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[4].Documentid.Docnumber)

	ass.Equal("priority-claim", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[5].XMLName.Local)
	ass.Equal("2", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[5].SequenceAttr)
	ass.Equal("original", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[5].DataformatAttr)
	ass.Equal("document-id", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[5].Documentid.XMLName.Local)
	ass.Equal("09/083,597", *exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[5].Documentid.Docnumber)

	// parties
	ass.Equal("parties", exchangeObject.ExchBibliographicdata.ExchParties.XMLName.Local)
	ass.Equal("applicants", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.XMLName.Local)
	ass.Equal("applicant", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[0].XMLName.Local)
	ass.Equal("1", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[0].SequenceAttr)
	ass.Equal("docdb", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[0].DataformatAttr)
	ass.Equal("A", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[0].StatusAttr)
	ass.Equal("applicant-name", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[0].ExchApplicantname[0].XMLName.Local)
	ass.Equal("name", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[0].ExchApplicantname[0].Name.XMLName.Local)
	ass.Equal("BRISTOL MYERS SQUIBB CO", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[0].ExchApplicantname[0].Name.Value)
	ass.Equal("residence", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[0].Residence.XMLName.Local)
	ass.Equal("US", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[0].Residence.Country)

	ass.Equal("applicant", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[1].XMLName.Local)
	ass.Equal("1", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[1].SequenceAttr)
	ass.Equal("docdba", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[1].DataformatAttr)
	ass.Equal("applicant-name", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[1].ExchApplicantname[0].XMLName.Local)
	ass.Equal("name", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[1].ExchApplicantname[0].Name.XMLName.Local)
	ass.Equal("BRISTOL-MYERS SQUIBB COMPANY", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[1].ExchApplicantname[0].Name.Value)

	ass.Equal("inventors", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.XMLName.Local)
	ass.Equal("inventor", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[0].XMLName.Local)
	ass.Equal("inventor-name", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[0].ExchInventorname[0].XMLName.Local)
	ass.Equal("1", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[0].SequenceAttr)
	ass.Equal("docdb", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[0].DataformatAttr)

	ass.Equal("invention-title", exchangeObject.ExchBibliographicdata.ExchInventiontitle[0].XMLName.Local)
	ass.Equal("en", exchangeObject.ExchBibliographicdata.ExchInventiontitle[0].LangAttr)
	ass.Equal("docdba", exchangeObject.ExchBibliographicdata.ExchInventiontitle[0].DataformatAttr)
	ass.Equal("dates-of-public-availability", exchangeObject.ExchBibliographicdata.ExchDatesofpublicavailability.XMLName.Local)
	ass.Equal("printed-with-grant", exchangeObject.ExchBibliographicdata.ExchDatesofpublicavailability.ExchPrintedwithgrant.XMLName.Local)
	ass.Equal("document-id", exchangeObject.ExchBibliographicdata.ExchDatesofpublicavailability.ExchPrintedwithgrant.Documentid.XMLName.Local)
	ass.Equal(20030918, exchangeObject.ExchBibliographicdata.ExchDatesofpublicavailability.ExchPrintedwithgrant.Documentid.Date)

	ass.Equal("references-cited", exchangeObject.ExchBibliographicdata.ExchReferencescited.XMLName.Local)
	ass.Equal("citation", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[0].XMLName.Local)
	ass.Equal("SEA", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[0].CitedphaseAttr)
	ass.Equal("1", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[0].SequenceAttr)
	ass.Equal("patcit", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[0].Patcit.XMLName.Local)
	ass.Equal("1", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[0].Patcit.NumAttr)
	ass.Equal("US5556839A", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[0].Patcit.DnumAttr)
	ass.Equal("publication number", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[0].Patcit.DnumtypeAttr)
	ass.Equal("document-id", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[0].Patcit.Documentid.XMLName.Local)
	ass.Equal(298340634, exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[0].Patcit.Documentid.DocidAttr)
	ass.Equal("5556839", *exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[0].Patcit.Documentid.Docnumber)
	ass.Equal("US", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[0].Patcit.Documentid.Country)
	ass.Equal("5556839", *exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[0].Patcit.Documentid.Docnumber)
	// ass.Equal("kind", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[0].Patcit.Documentid.Kind.XMLName.Local)
	ass.Equal("A", *exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[0].Patcit.Documentid.Kind)
	ass.Equal("name", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[0].Patcit.Documentid.Name.XMLName.Local)
	ass.Equal("GREENE JAMES M [US], et al", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[0].Patcit.Documentid.Name.Value)
	ass.Equal(19960917, exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[0].Patcit.Documentid.Date)

	ass.Equal("citation", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[1].XMLName.Local)
	ass.Equal("SEA", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[1].CitedphaseAttr)
	ass.Equal("2", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[1].SequenceAttr)
	ass.Equal("patcit", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[1].Patcit.XMLName.Local)
	ass.Equal("2", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[1].Patcit.NumAttr)
	ass.Equal("US5510114A", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[1].Patcit.DnumAttr)
	ass.Equal("publication number", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[1].Patcit.DnumtypeAttr)
	ass.Equal("document-id", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[1].Patcit.Documentid.XMLName.Local)
	ass.Equal(301191403, exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[1].Patcit.Documentid.DocidAttr)
	ass.Equal("US", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[1].Patcit.Documentid.Country)
	ass.Equal("5510114", *exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[1].Patcit.Documentid.Docnumber)

	// ass.Equal("kind", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[1].Patcit.Documentid.Kind.XMLName.Local)
	ass.Equal("A", *exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[1].Patcit.Documentid.Kind)
	ass.Equal("name", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[1].Patcit.Documentid.Name.XMLName.Local)
	ass.Equal("BORELLA FABIO [IT], et al", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[1].Patcit.Documentid.Name.Value)
	ass.Equal(19960423, exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[1].Patcit.Documentid.Date)

	ass.Equal("citation", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[2].XMLName.Local)
	ass.Equal("SEA", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[2].CitedphaseAttr)
	ass.Equal("3", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[2].SequenceAttr)
	ass.Equal("patcit", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[2].Patcit.XMLName.Local)
	ass.Equal("3", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[2].Patcit.NumAttr)
	ass.Equal("US5326570A", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[2].Patcit.DnumAttr)
	ass.Equal("publication number", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[2].Patcit.DnumtypeAttr)
	ass.Equal("document-id", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[2].Patcit.Documentid.XMLName.Local)
	ass.Equal(302101161, exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[2].Patcit.Documentid.DocidAttr)
	ass.Equal("US", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[2].Patcit.Documentid.Country)
	ass.Equal("5326570", *exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[2].Patcit.Documentid.Docnumber)
	// ass.Equal("kind", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[2].Patcit.Documentid.Kind.XMLName.Local)
	ass.Equal("A", *exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[2].Patcit.Documentid.Kind)
	ass.Equal("name", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[2].Patcit.Documentid.Name.XMLName.Local)
	ass.Equal("RUDNIC EDWARD M [US], et al", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[2].Patcit.Documentid.Name.Value)
	ass.Equal(19940705, exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[2].Patcit.Documentid.Date)

	// abstract tag
	ass.Equal("abstract", exchangeObject.ExchAbstract[0].XMLName.Local)
	ass.Equal("docdba", exchangeObject.ExchAbstract[0].DataformatAttr)
	ass.Equal("national office", exchangeObject.ExchAbstract[0].AbstractsourceAttr)
	ass.Equal("p", exchangeObject.ExchAbstract[0].ExchP[0].XMLName.Local)
	ass.Equal(598, len(exchangeObject.ExchAbstract[0].ExchP[0].Value))
	fmt.Println(exchangeObject.ExchAbstract[0].ExchP[0].Value)

	// family member tag 1 (46)
	ass.Equal("family-member", exchangeObject.ExchPatentfamily.ExchFamilymember[0].XMLName.Local)
	ass.Equal("application-reference", exchangeObject.ExchPatentfamily.ExchFamilymember[0].ExchApplicationreference[0].XMLName.Local)
	ass.Equal("docdb", exchangeObject.ExchPatentfamily.ExchFamilymember[0].ExchApplicationreference[0].DataformatAttr)
	ass.Equal("NO", exchangeObject.ExchPatentfamily.ExchFamilymember[0].ExchApplicationreference[0].IsrepresentativeAttr)
	ass.Equal("2000001988", *exchangeObject.ExchPatentfamily.ExchFamilymember[0].ExchApplicationreference[0].Documentid.Docnumber)
	ass.Equal("publication-reference", exchangeObject.ExchPatentfamily.ExchFamilymember[0].ExchPublicationreference[0].XMLName.Local)
	ass.Equal("docdb", exchangeObject.ExchPatentfamily.ExchFamilymember[0].ExchPublicationreference[0].DataformatAttr)
	ass.Equal("document-id", exchangeObject.ExchPatentfamily.ExchFamilymember[0].ExchPublicationreference[0].Documentid.XMLName.Local)
	ass.Equal("AP", exchangeObject.ExchPatentfamily.ExchFamilymember[0].ExchPublicationreference[0].Documentid.Country)
	ass.Equal("1206", *exchangeObject.ExchPatentfamily.ExchFamilymember[0].ExchPublicationreference[0].Documentid.Docnumber)
	// ass.Equal("kind", exchangeObject.ExchPatentfamily.ExchFamilymember[0].ExchPublicationreference[0].Documentid.Kind.XMLName.Local)
	ass.Equal("A", *exchangeObject.ExchPatentfamily.ExchFamilymember[0].ExchPublicationreference[0].Documentid.Kind)
	ass.Equal("publication-reference", exchangeObject.ExchPatentfamily.ExchFamilymember[0].ExchPublicationreference[1].XMLName.Local)
	ass.Equal("epodoc", exchangeObject.ExchPatentfamily.ExchFamilymember[0].ExchPublicationreference[1].DataformatAttr)
	ass.Equal("document-id", exchangeObject.ExchPatentfamily.ExchFamilymember[0].ExchPublicationreference[1].Documentid.XMLName.Local)
	ass.Equal("AP1206", *exchangeObject.ExchPatentfamily.ExchFamilymember[0].ExchPublicationreference[1].Documentid.Docnumber)
	ass.Equal("publication-reference", exchangeObject.ExchPatentfamily.ExchFamilymember[0].ExchPublicationreference[2].XMLName.Local)
	ass.Equal("docdb", exchangeObject.ExchPatentfamily.ExchFamilymember[0].ExchPublicationreference[2].DataformatAttr)
	ass.Equal("document-id", exchangeObject.ExchPatentfamily.ExchFamilymember[0].ExchPublicationreference[2].Documentid.XMLName.Local)
	ass.Equal("AP", exchangeObject.ExchPatentfamily.ExchFamilymember[0].ExchPublicationreference[2].Documentid.Country)
	ass.Equal("2000001988", *exchangeObject.ExchPatentfamily.ExchFamilymember[0].ExchPublicationreference[2].Documentid.Docnumber)
	// ass.Equal("kind", exchangeObject.ExchPatentfamily.ExchFamilymember[0].ExchPublicationreference[2].Documentid.Kind.XMLName.Local)
	ass.Equal("A0", *exchangeObject.ExchPatentfamily.ExchFamilymember[0].ExchPublicationreference[2].Documentid.Kind)

	for i := 0; i <= 45; i++ {
		ass.Equal("family-member", exchangeObject.ExchPatentfamily.ExchFamilymember[i].XMLName.Local)
		ass.Equal("application-reference", exchangeObject.ExchPatentfamily.ExchFamilymember[i].ExchApplicationreference[0].XMLName.Local)
		ass.Equal("docdb", exchangeObject.ExchPatentfamily.ExchFamilymember[i].ExchApplicationreference[0].DataformatAttr)
		ass.Equal("document-id", exchangeObject.ExchPatentfamily.ExchFamilymember[i].ExchApplicationreference[0].Documentid.XMLName.Local)
		// ass.Equal("doc-number", exchangeObject.ExchPatentfamily.ExchFamilymember[i].ExchApplicationreference[0].Documentid.Docnumber.XMLName.Local)
		ass.Equal("publication-reference", exchangeObject.ExchPatentfamily.ExchFamilymember[i].ExchPublicationreference[0].XMLName.Local)
		ass.Equal("docdb", exchangeObject.ExchPatentfamily.ExchFamilymember[i].ExchPublicationreference[0].DataformatAttr)
		ass.Equal("document-id", exchangeObject.ExchPatentfamily.ExchFamilymember[i].ExchPublicationreference[0].Documentid.XMLName.Local)
		// ass.Equal("doc-number", exchangeObject.ExchPatentfamily.ExchFamilymember[i].ExchPublicationreference[0].Documentid.Docnumber.XMLName.Local)
		// ass.Equal("kind", exchangeObject.ExchPatentfamily.ExchFamilymember[i].ExchPublicationreference[0].Documentid.Kind.XMLName.Local)
		ass.Equal("publication-reference", exchangeObject.ExchPatentfamily.ExchFamilymember[i].ExchPublicationreference[1].XMLName.Local)
		ass.Equal("epodoc", exchangeObject.ExchPatentfamily.ExchFamilymember[i].ExchPublicationreference[1].DataformatAttr)
		ass.Equal("document-id", exchangeObject.ExchPatentfamily.ExchFamilymember[i].ExchPublicationreference[1].Documentid.XMLName.Local)
		// ass.Equal("doc-number", exchangeObject.ExchPatentfamily.ExchFamilymember[i].ExchPublicationreference[1].Documentid.Docnumber.XMLName.Local)
	}
	// 2
	ass.Equal("NO", exchangeObject.ExchPatentfamily.ExchFamilymember[2].ExchApplicationreference[0].IsrepresentativeAttr)
	ass.Equal("4342002", *exchangeObject.ExchPatentfamily.ExchFamilymember[2].ExchApplicationreference[0].Documentid.Docnumber)
	ass.Equal("application-reference", exchangeObject.ExchPatentfamily.ExchFamilymember[2].ExchApplicationreference[1].XMLName.Local)
	ass.Equal("epodoc", exchangeObject.ExchPatentfamily.ExchFamilymember[2].ExchApplicationreference[1].DataformatAttr)
	ass.Equal("document-id", exchangeObject.ExchPatentfamily.ExchFamilymember[2].ExchApplicationreference[1].Documentid.XMLName.Local)
	//ass.Equal("doc-number", exchangeObject.ExchPatentfamily.ExchFamilymember[2].ExchApplicationreference[1].Documentid.Docnumber.XMLName.Local)
	ass.Equal("AT20020000434U", *exchangeObject.ExchPatentfamily.ExchFamilymember[2].ExchApplicationreference[1].Documentid.Docnumber)
	ass.Equal("AT", exchangeObject.ExchPatentfamily.ExchFamilymember[2].ExchPublicationreference[0].Documentid.Country)
	ass.Equal("6311", *exchangeObject.ExchPatentfamily.ExchFamilymember[2].ExchPublicationreference[0].Documentid.Docnumber)
	ass.Equal("U1", *exchangeObject.ExchPatentfamily.ExchFamilymember[2].ExchPublicationreference[0].Documentid.Kind)
	ass.Equal("AT6311U", *exchangeObject.ExchPatentfamily.ExchFamilymember[2].ExchPublicationreference[1].Documentid.Docnumber)
	ass.Equal("application-reference", exchangeObject.ExchPatentfamily.ExchFamilymember[2].ExchApplicationreference[1].XMLName.Local)
	ass.Equal("epodoc", exchangeObject.ExchPatentfamily.ExchFamilymember[2].ExchApplicationreference[1].DataformatAttr)
	ass.Equal("document-id", exchangeObject.ExchPatentfamily.ExchFamilymember[2].ExchApplicationreference[1].Documentid.XMLName.Local)
	//ass.Equal("doc-number", exchangeObject.ExchPatentfamily.ExchFamilymember[2].ExchApplicationreference[1].Documentid.Docnumber.XMLName.Local)
	ass.Equal("AT20020000434U", *exchangeObject.ExchPatentfamily.ExchFamilymember[2].ExchApplicationreference[1].Documentid.Docnumber)

	// 3
	ass.Equal("NO", exchangeObject.ExchPatentfamily.ExchFamilymember[3].ExchApplicationreference[0].IsrepresentativeAttr)
	ass.Equal("98938302", *exchangeObject.ExchPatentfamily.ExchFamilymember[3].ExchApplicationreference[0].Documentid.Docnumber)
	ass.Equal("application-reference", exchangeObject.ExchPatentfamily.ExchFamilymember[3].ExchApplicationreference[1].XMLName.Local)
	ass.Equal("epodoc", exchangeObject.ExchPatentfamily.ExchFamilymember[3].ExchApplicationreference[1].DataformatAttr)
	ass.Equal("document-id", exchangeObject.ExchPatentfamily.ExchFamilymember[3].ExchApplicationreference[1].Documentid.XMLName.Local)
	ass.Equal("AT19980938302T", *exchangeObject.ExchPatentfamily.ExchFamilymember[3].ExchApplicationreference[1].Documentid.Docnumber)
	ass.Equal("AT", exchangeObject.ExchPatentfamily.ExchFamilymember[3].ExchPublicationreference[0].Documentid.Country)
	ass.Equal("311859", *exchangeObject.ExchPatentfamily.ExchFamilymember[3].ExchPublicationreference[0].Documentid.Docnumber)
	ass.Equal("T", *exchangeObject.ExchPatentfamily.ExchFamilymember[3].ExchPublicationreference[0].Documentid.Kind)
	ass.Equal("AT311859T", *exchangeObject.ExchPatentfamily.ExchFamilymember[3].ExchPublicationreference[1].Documentid.Docnumber)
	ass.Equal("application-reference", exchangeObject.ExchPatentfamily.ExchFamilymember[3].ExchApplicationreference[1].XMLName.Local)
	ass.Equal("epodoc", exchangeObject.ExchPatentfamily.ExchFamilymember[3].ExchApplicationreference[1].DataformatAttr)
	ass.Equal("document-id", exchangeObject.ExchPatentfamily.ExchFamilymember[3].ExchApplicationreference[1].Documentid.XMLName.Local)
	ass.Equal("AT19980938302T", *exchangeObject.ExchPatentfamily.ExchFamilymember[3].ExchApplicationreference[1].Documentid.Docnumber)

	// 4
	ass.Equal("NO", exchangeObject.ExchPatentfamily.ExchFamilymember[4].ExchApplicationreference[0].IsrepresentativeAttr)
	ass.Equal("8685498", *exchangeObject.ExchPatentfamily.ExchFamilymember[4].ExchApplicationreference[0].Documentid.Docnumber)
	ass.Equal("application-reference", exchangeObject.ExchPatentfamily.ExchFamilymember[4].ExchApplicationreference[1].XMLName.Local)
	ass.Equal("epodoc", exchangeObject.ExchPatentfamily.ExchFamilymember[4].ExchApplicationreference[1].DataformatAttr)
	ass.Equal("document-id", exchangeObject.ExchPatentfamily.ExchFamilymember[4].ExchApplicationreference[1].Documentid.XMLName.Local)
	ass.Equal("AU19980086854", *exchangeObject.ExchPatentfamily.ExchFamilymember[4].ExchApplicationreference[1].Documentid.Docnumber)
	ass.Equal("AU", exchangeObject.ExchPatentfamily.ExchFamilymember[4].ExchPublicationreference[0].Documentid.Country)
	ass.Equal("750911", *exchangeObject.ExchPatentfamily.ExchFamilymember[4].ExchPublicationreference[0].Documentid.Docnumber)
	ass.Equal("B2", *exchangeObject.ExchPatentfamily.ExchFamilymember[4].ExchPublicationreference[0].Documentid.Kind)
	ass.Equal("1", exchangeObject.ExchPatentfamily.ExchFamilymember[4].ExchPublicationreference[0].SequenceAttr)
	ass.Equal("AU750911B", *exchangeObject.ExchPatentfamily.ExchFamilymember[4].ExchPublicationreference[1].Documentid.Docnumber)
	ass.Equal("1", exchangeObject.ExchPatentfamily.ExchFamilymember[4].ExchPublicationreference[1].SequenceAttr)
	ass.Equal("publication-reference", exchangeObject.ExchPatentfamily.ExchFamilymember[4].ExchPublicationreference[2].XMLName.Local)
	ass.Equal("docdb", exchangeObject.ExchPatentfamily.ExchFamilymember[4].ExchPublicationreference[2].DataformatAttr)
	ass.Equal("2", exchangeObject.ExchPatentfamily.ExchFamilymember[4].ExchPublicationreference[2].SequenceAttr)
	ass.Equal("document-id", exchangeObject.ExchPatentfamily.ExchFamilymember[4].ExchPublicationreference[2].Documentid.XMLName.Local)
	// ass.Equal("kind", exchangeObject.ExchPatentfamily.ExchFamilymember[4].ExchPublicationreference[2].Documentid.Kind.XMLName.Local)
	ass.Equal("AU", exchangeObject.ExchPatentfamily.ExchFamilymember[4].ExchPublicationreference[2].Documentid.Country)
	ass.Equal("8685498", *exchangeObject.ExchPatentfamily.ExchFamilymember[4].ExchPublicationreference[2].Documentid.Docnumber)
	ass.Equal("A", *exchangeObject.ExchPatentfamily.ExchFamilymember[4].ExchPublicationreference[2].Documentid.Kind)
	ass.Equal("publication-reference", exchangeObject.ExchPatentfamily.ExchFamilymember[4].ExchPublicationreference[3].XMLName.Local)
	ass.Equal("epodoc", exchangeObject.ExchPatentfamily.ExchFamilymember[4].ExchPublicationreference[3].DataformatAttr)
	ass.Equal("2", exchangeObject.ExchPatentfamily.ExchFamilymember[4].ExchPublicationreference[3].SequenceAttr)
	ass.Equal("document-id", exchangeObject.ExchPatentfamily.ExchFamilymember[4].ExchPublicationreference[3].Documentid.XMLName.Local)
	ass.Equal("AU8685498", *exchangeObject.ExchPatentfamily.ExchFamilymember[4].ExchPublicationreference[3].Documentid.Docnumber)

}

func TestReadFileWO(t *testing.T) {

	// process file
	exchangeObject, err := ParseXmlFileToStruct("./test-data/WO-2022259205-A1_544370561.xml")
	if err != nil {
		t.Error(err)
	}

	// jsonObject, err := json.MarshalIndent(exchangeObject, "", "  ")
	// fmt.Println(string(jsonObject))

	ass := assert.New(t)
	ass.NoError(err)

	// Exchange-document tag
	ass.Equal("WO", exchangeObject.CountryAttr)
	ass.Equal(20221216, exchangeObject.DateaddeddocdbAttr)
	ass.Equal(20230219, exchangeObject.DateoflastexchangeAttr)
	ass.Equal(20221215, exchangeObject.DatepublAttr)
	ass.Equal("2022259205", exchangeObject.DocnumberAttr)
	ass.Equal("82558028", exchangeObject.FamilyidAttr)
	ass.Equal("584641170", exchangeObject.DocidAttr)
	ass.Equal("YES", exchangeObject.IsrepresentativeAttr)
	ass.Equal("A1", exchangeObject.KindAttr)
	ass.Equal("EP", exchangeObject.OriginatingofficeAttr)
	ass.Equal("", exchangeObject.StatusAttr)

	// bibliographic-data
	// publication-reference
	ass.Equal("publication-reference", exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].XMLName.Local)
	ass.Equal("document-id", exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].Documentid.XMLName.Local)
	ass.Equal("en", exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].Documentid.LangAttr)
	ass.Equal("WO", exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].Documentid.Country)
	ass.Equal("A1", *exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].Documentid.Kind)
	ass.Equal(20221215, exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].Documentid.Date)

	ass.Equal("docdb", exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].DataformatAttr)
	ass.Equal("2022259205", *exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].Documentid.Docnumber)
	ass.Equal("epodoc", exchangeObject.ExchBibliographicdata.ExchPublicationreference[1].DataformatAttr)
	ass.Equal("WO2022259205", *exchangeObject.ExchBibliographicdata.ExchPublicationreference[1].Documentid.Docnumber)
	ass.Equal("original", exchangeObject.ExchBibliographicdata.ExchPublicationreference[2].DataformatAttr)
	ass.Equal("2022/259205", *exchangeObject.ExchBibliographicdata.ExchPublicationreference[2].Documentid.Docnumber)

	ass.Equal("x", exchangeObject.ExchBibliographicdata.ExchExtendedkindcode.Value)

	// classification-ipc
	for i := 0; i <= 10; i++ {
		ass.Equal(strconv.Itoa(i+1), exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[i].SequenceAttr)
	}

	ass.Equal("classifications-ipcr", exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.XMLName.Local)
	ass.Equal("A61K   8/04        20060101ALI20221215BHEP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[0].Text)
	ass.Equal("A61K   9/00        20060101AFI20221215BHEP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[1].Text)
	ass.Equal("A61K   9/107       20060101ALI20221215BHEP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[2].Text)
	ass.Equal("A61K  31/366       20060101ALI20221215BHEP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[3].Text)
	ass.Equal("A61K  31/4166      20060101ALI20221215BHEP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[4].Text)
	ass.Equal("A61K  31/57        20060101ALI20221215BHEP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[5].Text)
	ass.Equal("A61K  31/575       20060101ALI20221215BHEP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[6].Text)
	ass.Equal("A61K  31/675       20060101ALI20221215BHEP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[7].Text)
	ass.Equal("A61K  47/10        20170101ALI20221215BHEP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[8].Text)
	ass.Equal("A61K  47/14        20170101ALI20221215BHEP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[9].Text)
	ass.Equal("A61K  47/42        20170101ALI20221215BHEP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[10].Text)

	// patent classification
	for i := 0; i <= 15; i++ {
		ass.Equal(strconv.Itoa(i+1), exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[i].SequenceAttr)
		ass.Equal("patent-classification", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[i].XMLName.Local)
		ass.Equal("classification-scheme", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[i].Classificationscheme.XMLName.Local)
		ass.Equal("CPCI", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[i].Classificationscheme.SchemeAttr)
		ass.Equal("B", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[i].Classificationstatus)
		ass.Equal("H", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[i].Classificationdatasource)
		ass.Equal("EP", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[i].Classificationscheme.OfficeAttr)
		ass.Equal("action-date", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[i].Actiondate.XMLName.Local)
	}

	ass.Equal("A61K   9/1075", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[0].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[0].Classificationvalue)
	ass.Equal("F", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[0].Symbolposition)
	ass.Equal(20220809, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[0].Actiondate.Date)
	ass.Equal("A61K   8/064", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[1].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[1].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[1].Symbolposition)
	ass.Equal(20221020, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[1].Actiondate.Date)
	ass.Equal("A61K   8/375", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[2].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[2].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[2].Symbolposition)
	ass.Equal(20221020, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[2].Actiondate.Date)
	ass.Equal("A61K   8/922", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[3].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[3].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[3].Symbolposition)
	ass.Equal(20221020, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[3].Actiondate.Date)
	ass.Equal("A61K   9/0043", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[4].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[4].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[4].Symbolposition)
	ass.Equal(20220809, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[4].Actiondate.Date)
	ass.Equal("A61K  31/366", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[5].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[5].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[5].Symbolposition)
	ass.Equal(20220928, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[5].Actiondate.Date)
	ass.Equal("A61K  31/4166", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[6].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[6].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[6].Symbolposition)
	ass.Equal(20220928, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[6].Actiondate.Date)
	ass.Equal("A61K  31/57", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[7].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[7].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[7].Symbolposition)
	ass.Equal(20220928, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[7].Actiondate.Date)
	ass.Equal("A61K  31/575", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[8].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[8].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[8].Symbolposition)
	ass.Equal(20220928, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[8].Actiondate.Date)
	ass.Equal("A61K  31/675", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[9].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[9].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[9].Symbolposition)
	ass.Equal(20220928, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[9].Actiondate.Date)
	ass.Equal("A61K  47/10", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[10].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[10].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[10].Symbolposition)
	ass.Equal(20220809, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[10].Actiondate.Date)
	ass.Equal("A61K  47/14", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[11].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[11].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[11].Symbolposition)
	ass.Equal(20220809, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[11].Actiondate.Date)
	ass.Equal("A61K  47/42", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[12].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[12].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[12].Symbolposition)
	ass.Equal(20220809, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[12].Actiondate.Date)
	ass.Equal("A61K2800/10", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[13].Classificationsymbol)
	ass.Equal("A", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[13].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[13].Symbolposition)
	ass.Equal(20221020, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[13].Actiondate.Date)
	ass.Equal("A61K2800/21", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[14].Classificationsymbol)
	ass.Equal("A", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[14].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[14].Symbolposition)
	ass.Equal(20221020, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[14].Actiondate.Date)
	ass.Equal("A61Q  19/00", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[15].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[15].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[15].Symbolposition)
	ass.Equal(20221020, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[15].Actiondate.Date)

	// application reference
	ass.Equal("application-reference", exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].XMLName.Local)
	ass.Equal("YES", exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].IsrepresentativeAttr)
	ass.Equal("574950146", exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].DocidAttr)
	ass.Equal("document-id", exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].Documentid.XMLName.Local)
	ass.Equal("IB", exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].Documentid.Country)
	ass.Equal("W", *exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].Documentid.Kind)
	ass.Equal(20220609, exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].Documentid.Date)

	ass.Equal("docdb", exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].DataformatAttr)
	ass.Equal("2022055385", *exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].Documentid.Docnumber)
	ass.Equal("epodoc", exchangeObject.ExchBibliographicdata.ExchApplicationreference[1].DataformatAttr)
	ass.Equal("WO2022IB55385", *exchangeObject.ExchBibliographicdata.ExchApplicationreference[1].Documentid.Docnumber)
	ass.Equal("original", exchangeObject.ExchBibliographicdata.ExchApplicationreference[2].DataformatAttr)
	ass.Equal("IB2022/055385", *exchangeObject.ExchBibliographicdata.ExchApplicationreference[2].Documentid.Docnumber)

	// language of filing (instead of language of publication)
	ass.Equal("language-of-filing", exchangeObject.ExchBibliographicdata.ExchLanguageoffiling.XMLName.Local)
	ass.Equal("en", exchangeObject.ExchBibliographicdata.ExchLanguageoffiling.Value)

	// priority claims
	ass.Equal("priority-claims", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.XMLName.Local)
	ass.Equal("priority-claim", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[0].XMLName.Local)
	ass.Equal("1", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[0].SequenceAttr)
	ass.Equal("docdb", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[0].DataformatAttr)
	ass.Equal("document-id", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[0].Documentid.XMLName.Local)
	ass.Equal("11728321", *exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[0].Documentid.Docnumber)

	ass.Equal("A", *exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[0].Documentid.Kind)
	ass.Equal("PT", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[0].Documentid.Country)
	ass.Equal(20210611, exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[0].Documentid.Date)
	ass.Equal("Y", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[0].ExchPriorityactiveindicator)
	ass.Equal("1", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[1].SequenceAttr)
	ass.Equal("epodoc", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[1].DataformatAttr)
	ass.Equal("document-id", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[1].Documentid.XMLName.Local)
	ass.Equal("PT20210117283", *exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[1].Documentid.Docnumber)
	ass.Equal("1", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[2].SequenceAttr)
	ass.Equal("original", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[2].DataformatAttr)
	ass.Equal("document-id", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[2].Documentid.XMLName.Local)
	ass.Equal("117283", *exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[2].Documentid.Docnumber)

	// parties
	ass.Equal("parties", exchangeObject.ExchBibliographicdata.ExchParties.XMLName.Local)
	// applicants
	ass.Equal("applicants", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.XMLName.Local)
	ass.Equal("applicant", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[0].XMLName.Local)
	ass.Equal("1", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[0].SequenceAttr)
	ass.Equal("docdb", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[0].DataformatAttr)
	ass.Equal("applicant-name", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[0].ExchApplicantname[0].XMLName.Local)
	ass.Equal("name", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[0].ExchApplicantname[0].Name.XMLName.Local)
	ass.Equal("UNIV DA BEIRA INTERIOR", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[0].ExchApplicantname[0].Name.Value)
	ass.Equal("residence", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[0].Residence.XMLName.Local)
	ass.Equal("PT", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[0].Residence.Country)

	ass.Equal("applicant", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[1].XMLName.Local)
	ass.Equal("1", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[1].SequenceAttr)
	ass.Equal("docdba", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[1].DataformatAttr)
	ass.Equal("applicant-name", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[1].ExchApplicantname[0].XMLName.Local)
	ass.Equal("name", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[1].ExchApplicantname[0].Name.XMLName.Local)
	ass.Equal("UNIVERSIDADE DA BEIRA INTERIOR", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[1].ExchApplicantname[0].Name.Value)

	// inventors
	ass.Equal("inventors", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.XMLName.Local)

	for i := 0; i <= 7; i++ {
		ass.Equal(strconv.Itoa(i+1), exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[i].SequenceAttr)
		ass.Equal("docdb", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[i].DataformatAttr)
		ass.Equal("PT", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[i].Residence.Country)
	}
	ass.Equal("OLIVEIRA DOS SANTOS ADRIANA", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[0].ExchInventorname[0].Name.Value)
	ass.Equal("CARVALHO FERNANDES MARIANA", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[1].ExchInventorname[0].Name.Value)
	ass.Equal("CABRAL PIRES PATRÍCIA SOFIA", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[2].ExchInventorname[0].Name.Value)
	ass.Equal("LOURENÇO ALVES GILBERTO", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[3].ExchInventorname[0].Name.Value)
	ass.Equal("MARICOTO FAZENDEIRO ANA CAROLINA", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[4].ExchInventorname[0].Name.Value)
	ass.Equal("MATOS SILVA PEREIRA NINA FRANCISCA", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[5].ExchInventorname[0].Name.Value)
	ass.Equal("DA SILVA FERREIRA GOMES MARIA DE FÁTIMA", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[6].ExchInventorname[0].Name.Value)
	ass.Equal("ESTEVES RODRIGUES LINA ISABEL", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[7].ExchInventorname[0].Name.Value)

	for i := 8; i <= 15; i++ {
		ass.Equal(strconv.Itoa(i-7), exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[i].SequenceAttr)
		ass.Equal("docdba", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[i].DataformatAttr)
	}
	ass.Equal("OLIVEIRA DOS SANTOS, Adriana", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[8].ExchInventorname[0].Name.Value)
	ass.Equal("CARVALHO FERNANDES, Mariana", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[9].ExchInventorname[0].Name.Value)
	ass.Equal("CABRAL PIRES, Patrícia Sofia", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[10].ExchInventorname[0].Name.Value)
	ass.Equal("LOURENÇO ALVES, Gilberto", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[11].ExchInventorname[0].Name.Value)
	ass.Equal("MARICOTO FAZENDEIRO, Ana Carolina", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[12].ExchInventorname[0].Name.Value)
	ass.Equal("MATOS SILVA PEREIRA NINA, Francisca", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[13].ExchInventorname[0].Name.Value)
	ass.Equal("DA SILVA FERREIRA GOMES, Maria De Fátima", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[14].ExchInventorname[0].Name.Value)
	ass.Equal("ESTEVES RODRIGUES, Lina Isabel", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[15].ExchInventorname[0].Name.Value)

	// invention title
	ass.Equal("invention-title", exchangeObject.ExchBibliographicdata.ExchInventiontitle[0].XMLName.Local)
	ass.Equal("en", exchangeObject.ExchBibliographicdata.ExchInventiontitle[0].LangAttr)
	ass.Equal("docdba", exchangeObject.ExchBibliographicdata.ExchInventiontitle[0].DataformatAttr)
	ass.Equal("SELF-EMULSIFYING COMPOSITION, PRODUCTION METHODS AND USES THEREOF", exchangeObject.ExchBibliographicdata.ExchInventiontitle[0].Value)
	ass.Equal("fr", exchangeObject.ExchBibliographicdata.ExchInventiontitle[1].LangAttr)
	ass.Equal("docdba", exchangeObject.ExchBibliographicdata.ExchInventiontitle[1].DataformatAttr)
	ass.Equal("COMPOSITION AUTO-ÉMULSIFIANTE, SES PROCÉDÉS DE PRODUCTION ET SES UTILISATIONS", exchangeObject.ExchBibliographicdata.ExchInventiontitle[1].Value)

	ass.Equal("dates-of-public-availability", exchangeObject.ExchBibliographicdata.ExchDatesofpublicavailability.XMLName.Local)
	// examined/ unexamined (not) printed with(out) grant
	ass.Equal("examined-printed-without-grant", exchangeObject.ExchBibliographicdata.ExchDatesofpublicavailability.ExchExaminedprintedwithoutgrant.XMLName.Local)
	ass.Equal(20221215, exchangeObject.ExchBibliographicdata.ExchDatesofpublicavailability.ExchExaminedprintedwithoutgrant.Documentid.Date)

	ass.Equal("references-cited", exchangeObject.ExchBibliographicdata.ExchReferencescited.XMLName.Local)
	ass.Equal("citation", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[0].XMLName.Local)

	for i := 0; i <= 11; i++ {
		ass.Equal("ISR", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[i].CitedphaseAttr)
		ass.Equal(strconv.Itoa(i+1), exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[i].SequenceAttr)
	}
	for i := 12; i <= 22; i++ {
		ass.Equal("APP", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[i].CitedphaseAttr)
		ass.Equal(strconv.Itoa(i-11), exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[i].SequenceAttr)
	}

	ass.Equal("s", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[0].Nplcit.NpltypeAttr)
	ass.Equal("055966431", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[0].Nplcit.ExtractedxpAttr)
	ass.Equal(242, len(*exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[0].Nplcit.Text))
	ass.Equal("BEZERRA-SOUZA ADRIANA ET AL", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[0].Nplcit.Article.Author[0].Name.Value)
	ass.Equal("Repurposing Butenafine as An Oral Nanomedicine for Visceral Leishmaniasis", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[0].Nplcit.Article.Atl.Value)
	ass.Equal("PHARMACEUTICS", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[0].Nplcit.Article.Serial.Sertitle.Value)
	ass.Equal("CH", *exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[0].Nplcit.Article.Serial.Imprint.Address.Text)
	ass.Equal("20190701", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[0].Nplcit.Article.Serial.Pubdate.Sdate[0].Value)
	ass.Equal("11", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[0].Nplcit.Article.Serial.Vid.Value)
	ass.Equal("7", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[0].Nplcit.Article.Serial.Ino.Value)
	ass.Equal("10.3390/pharmaceutics11070353", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[0].Nplcit.Article.Serial.Doi.Value)

	ass.Equal("citation", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[6].XMLName.Local)
	ass.Equal("patcit", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[6].Patcit.XMLName.Local)
	ass.Equal("1", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[6].Patcit.NumAttr)
	ass.Equal("CN102908333A", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[6].Patcit.DnumAttr)
	ass.Equal("publication number", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[6].Patcit.DnumtypeAttr)
	ass.Equal("document-id", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[6].Patcit.Documentid.XMLName.Local)
	ass.Equal(381263225, exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[6].Patcit.Documentid.DocidAttr)
	ass.Equal("102908333", *exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[6].Patcit.Documentid.Docnumber)
	ass.Equal("CN", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[6].Patcit.Documentid.Country)
	ass.Equal("102908333", *exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[6].Patcit.Documentid.Docnumber)
	ass.Equal("A", *exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[6].Patcit.Documentid.Kind)
	ass.Equal("name", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[6].Patcit.Documentid.Name.XMLName.Local)
	ass.Equal("UNIV CHINA PHARMA", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[6].Patcit.Documentid.Name.Value)
	ass.Equal(20130206, exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[6].Patcit.Documentid.Date)

	ass.Equal("s", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[18].Nplcit.NpltypeAttr)
	ass.Equal("027210215", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[18].Nplcit.ExtractedxpAttr)
	ass.Equal(123, len(*exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[18].Nplcit.Text))
	ass.Equal("NAZIR ASCHROEN KBOOM R", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[18].Nplcit.Article.Author[0].Name.Value)
	ass.Equal("Premix emulsification: A review", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[18].Nplcit.Article.Atl.Value)
	ass.Equal("J Member Sci", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[18].Nplcit.Article.Serial.Sertitle.Value)
	ass.Equal("20100000", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[18].Nplcit.Article.Serial.Pubdate.Sdate[0].Value)
	ass.Equal("362", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[18].Nplcit.Article.Serial.Vid.Value)
	ass.Equal("1-2", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[18].Nplcit.Article.Serial.Ino.Value)
	ass.Equal("1", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[18].Nplcit.Article.Location.Pp.Ppf[0].Value)
	ass.Equal("11", exchangeObject.ExchBibliographicdata.ExchReferencescited.ExchCitation[18].Nplcit.Article.Location.Pp.Ppl[0].Value)

	// abstract tag
	ass.Equal("abstract", exchangeObject.ExchAbstract[0].XMLName.Local)
	ass.Equal("en", exchangeObject.ExchAbstract[0].LangAttr)
	ass.Equal("docdba", exchangeObject.ExchAbstract[0].DataformatAttr)
	ass.Equal("national office", exchangeObject.ExchAbstract[0].AbstractsourceAttr)
	ass.Equal("p", exchangeObject.ExchAbstract[0].ExchP[0].XMLName.Local)
	ass.Equal(532, len(exchangeObject.ExchAbstract[0].ExchP[0].Value))

	ass.Equal("abstract", exchangeObject.ExchAbstract[1].XMLName.Local)
	ass.Equal("fr", exchangeObject.ExchAbstract[1].LangAttr)
	ass.Equal("docdba", exchangeObject.ExchAbstract[1].DataformatAttr)
	ass.Equal("national office", exchangeObject.ExchAbstract[1].AbstractsourceAttr)
	ass.Equal("p", exchangeObject.ExchAbstract[1].ExchP[0].XMLName.Local)
	ass.Equal(581, len(exchangeObject.ExchAbstract[1].ExchP[0].Value))

	fmt.Println(exchangeObject.ExchAbstract[0].ExchP[0].Value)
}

func TestReadFileYU(t *testing.T) {

	// process file
	exchangeObject, err := ParseXmlFileToStruct("./test-data/YU-6701-A_381353354.xml")
	if err != nil {
		t.Error(err)
	}

	// jsonObject, err := json.MarshalIndent(exchangeObject, "", "  ")
	// fmt.Println(string(jsonObject))

	ass := assert.New(t)
	ass.NoError(err)

	// Exchange-document
	ass.Equal("YU", exchangeObject.CountryAttr)
	ass.Equal(20140328, exchangeObject.DateaddeddocdbAttr)
	ass.Equal(20230219, exchangeObject.DateoflastexchangeAttr)
	ass.Equal(20031231, exchangeObject.DatepublAttr)
	ass.Equal("6701", exchangeObject.DocnumberAttr)
	ass.Equal("22248526", exchangeObject.FamilyidAttr)
	ass.Equal("404650208", exchangeObject.DocidAttr)
	ass.Equal("YES", exchangeObject.IsrepresentativeAttr)
	ass.Equal("A", exchangeObject.KindAttr)
	ass.Equal("EP", exchangeObject.OriginatingofficeAttr)
	ass.Equal("", exchangeObject.StatusAttr)

	// bibliographic-data
	// publication-reference
	ass.Equal("publication-reference", exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].XMLName.Local)
	ass.Equal("document-id", exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].Documentid.XMLName.Local)
	ass.Equal("sh", exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].Documentid.LangAttr)
	ass.Equal("YU", exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].Documentid.Country)
	ass.Equal("A", *exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].Documentid.Kind)
	ass.Equal(20031231, exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].Documentid.Date)

	ass.Equal("docdb", exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].DataformatAttr)
	ass.Equal("6701", *exchangeObject.ExchBibliographicdata.ExchPublicationreference[0].Documentid.Docnumber)
	ass.Equal("epodoc", exchangeObject.ExchBibliographicdata.ExchPublicationreference[1].DataformatAttr)
	ass.Equal("YU6701", *exchangeObject.ExchBibliographicdata.ExchPublicationreference[1].Documentid.Docnumber)
	ass.Equal("original", exchangeObject.ExchBibliographicdata.ExchPublicationreference[2].DataformatAttr)
	ass.Equal("6701", *exchangeObject.ExchBibliographicdata.ExchPublicationreference[2].Documentid.Docnumber)

	//////////

	// classification-ipc
	for i := 0; i <= 38; i++ {
		ass.Equal(strconv.Itoa(i+1), exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[i].SequenceAttr)
	}

	ass.Equal("classifications-ipcr", exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.XMLName.Local)
	ass.Equal("A61K  31/18        20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[0].Text)
	ass.Equal("A61K  31/198       20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[1].Text)
	ass.Equal("A61K  31/34        20060101A I20061125RMEP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[2].Text)
	ass.Equal("A61K  31/343       20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[3].Text)
	ass.Equal("A61K  31/38        20060101A I20061125RMEP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[4].Text)
	ass.Equal("A61K  31/381       20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[5].Text)
	ass.Equal("A61K  31/403       20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[6].Text)
	ass.Equal("A61K  31/4035      20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[7].Text)
	ass.Equal("A61K  31/404       20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[8].Text)
	ass.Equal("A61K  31/4178      20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[9].Text)
	ass.Equal("A61P   1/02        20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[10].Text)
	ass.Equal("A61P   9/04        20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[11].Text)
	ass.Equal("A61P   9/10        20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[12].Text)
	ass.Equal("A61P  13/12        20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[13].Text)
	ass.Equal("A61P  17/02        20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[14].Text)
	ass.Equal("A61P  19/02        20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[15].Text)
	ass.Equal("A61P  19/10        20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[16].Text)
	ass.Equal("A61P  21/04        20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[17].Text)
	ass.Equal("A61P  25/04        20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[18].Text)
	ass.Equal("A61P  25/14        20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[19].Text)
	ass.Equal("A61P  25/16        20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[20].Text)
	ass.Equal("A61P  25/28        20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[21].Text)
	ass.Equal("A61P  27/02        20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[22].Text)
	ass.Equal("A61P  29/00        20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[23].Text)
	ass.Equal("A61P  31/18        20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[24].Text)
	ass.Equal("A61P  35/00        20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[25].Text)
	ass.Equal("A61P  37/00        20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[26].Text)
	ass.Equal("A61P  43/00        20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[27].Text)
	ass.Equal("C07C 311/19        20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[28].Text)
	ass.Equal("C07C 311/21        20060101A I20051008RMEP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[29].Text)
	ass.Equal("C07D 209/88        20060101A I20051008RMEP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[30].Text)
	ass.Equal("C07D 307/91        20060101A I20051008RMEP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[31].Text)
	ass.Equal("C07D 307/93        20060101A I20051008RMEP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[32].Text)
	ass.Equal("C07D 333/76        20060101A I20051008RMEP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[33].Text)
	ass.Equal("C07D 405/04        20060101A I20051008RMEP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[34].Text)
	ass.Equal("C07D 405/12        20060101ALI20051220RMJP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[35].Text)
	ass.Equal("H04M  15/00        20060101A I20051008RMEP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[36].Text)
	ass.Equal("H04M  15/28        20060101A I20051008RMEP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[37].Text)
	ass.Equal("H04M  15/30        20060101A I20051008RMEP", *exchangeObject.ExchBibliographicdata.ExchClassificationsipcr.Classificationipcr[38].Text)

	// patent classification
	for i := 0; i <= 15; i++ {
		ass.Equal(strconv.Itoa(i+1), exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[i].SequenceAttr)
		ass.Equal("patent-classification", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[i].XMLName.Local)
		ass.Equal("classification-scheme", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[i].Classificationscheme.XMLName.Local)
		ass.Equal("CPCI", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[i].Classificationscheme.SchemeAttr)
		ass.Equal("B", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[i].Classificationstatus)
		ass.Equal("H", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[i].Classificationdatasource)
		ass.Equal("EP", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[i].Classificationscheme.OfficeAttr)
		ass.Equal("action-date", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[i].Actiondate.XMLName.Local)
	}

	ass.Equal("C07D 307/91", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[0].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[0].Classificationvalue)
	ass.Equal("F", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[0].Symbolposition)
	ass.Equal(20130101, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[0].Actiondate.Date)
	ass.Equal("A61P   1/02", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[1].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[1].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[1].Symbolposition)
	ass.Equal(20200318, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[1].Actiondate.Date)
	ass.Equal("A61P   9/04", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[2].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[2].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[2].Symbolposition)
	ass.Equal(20200331, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[2].Actiondate.Date)
	ass.Equal("A61P   9/10", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[3].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[3].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[3].Symbolposition)
	ass.Equal(20200331, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[3].Actiondate.Date)
	ass.Equal("A61P  13/12", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[4].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[4].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[4].Symbolposition)
	ass.Equal(20200319, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[4].Actiondate.Date)
	ass.Equal("A61P  17/02", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[5].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[5].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[5].Symbolposition)
	ass.Equal(20200319, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[5].Actiondate.Date)
	ass.Equal("A61P  19/02", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[6].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[6].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[6].Symbolposition)
	ass.Equal(20200320, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[6].Actiondate.Date)
	ass.Equal("A61P  19/10", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[7].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[7].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[7].Symbolposition)
	ass.Equal(20200320, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[7].Actiondate.Date)
	ass.Equal("A61P  21/04", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[8].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[8].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[8].Symbolposition)
	ass.Equal(20200320, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[8].Actiondate.Date)
	ass.Equal("A61P  25/04", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[9].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[9].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[9].Symbolposition)
	ass.Equal(20200323, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[9].Actiondate.Date)
	ass.Equal("A61P  25/14", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[10].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[10].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[10].Symbolposition)
	ass.Equal(20200323, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[10].Actiondate.Date)
	ass.Equal("A61P  25/16", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[11].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[11].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[11].Symbolposition)
	ass.Equal(20200323, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[11].Actiondate.Date)
	ass.Equal("A61P  25/28", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[12].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[12].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[12].Symbolposition)
	ass.Equal(20200324, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[12].Actiondate.Date)
	ass.Equal("A61P  27/02", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[13].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[13].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[13].Symbolposition)
	ass.Equal(20200324, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[13].Actiondate.Date)
	ass.Equal("A61P  29/00", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[14].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[14].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[14].Symbolposition)
	ass.Equal(20200325, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[14].Actiondate.Date)
	ass.Equal("A61P  31/18", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[15].Classificationsymbol)
	ass.Equal("I", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[15].Classificationvalue)
	ass.Equal("L", exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[15].Symbolposition)
	ass.Equal(20200327, exchangeObject.ExchBibliographicdata.ExchPatentclassifications.Patentclassification[15].Actiondate.Date)

	// application reference
	ass.Equal("application-reference", exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].XMLName.Local)
	ass.Equal("NO", exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].IsrepresentativeAttr)
	ass.Equal("381353354", exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].DocidAttr)
	ass.Equal("document-id", exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].Documentid.XMLName.Local)
	ass.Equal("YU", exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].Documentid.Country)
	ass.Equal("A", *exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].Documentid.Kind)
	ass.Equal(19990602, exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].Documentid.Date)

	ass.Equal("docdb", exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].DataformatAttr)
	ass.Equal("6701", *exchangeObject.ExchBibliographicdata.ExchApplicationreference[0].Documentid.Docnumber)
	ass.Equal("epodoc", exchangeObject.ExchBibliographicdata.ExchApplicationreference[1].DataformatAttr)
	ass.Equal("YU20010000067", *exchangeObject.ExchBibliographicdata.ExchApplicationreference[1].Documentid.Docnumber)
	ass.Equal("original", exchangeObject.ExchBibliographicdata.ExchApplicationreference[2].DataformatAttr)
	ass.Equal("6701", *exchangeObject.ExchBibliographicdata.ExchApplicationreference[2].Documentid.Docnumber)

	// language of publication
	ass.Equal("language-of-publication", exchangeObject.ExchBibliographicdata.ExchLanguageofpublication.XMLName.Local)
	ass.Equal("sh", exchangeObject.ExchBibliographicdata.ExchLanguageofpublication.Value)

	// priority claims
	ass.Equal("priority-claims", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.XMLName.Local)
	ass.Equal("priority-claim", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[0].XMLName.Local)
	ass.Equal("1", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[0].SequenceAttr)
	ass.Equal("docdb", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[0].DataformatAttr)
	ass.Equal("document-id", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[0].Documentid.XMLName.Local)
	ass.Equal("9500698", *exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[0].Documentid.Docnumber)
	ass.Equal("P", *exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[0].Documentid.Kind)
	ass.Equal("US", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[0].Documentid.Country)
	ass.Equal(19980730, exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[0].Documentid.Date)
	ass.Equal("Y", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[0].ExchPriorityactiveindicator)
	ass.Equal("1", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[1].SequenceAttr)
	ass.Equal("epodoc", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[1].DataformatAttr)
	ass.Equal("US19980095006P", *exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[1].Documentid.Docnumber)
	ass.Equal("1", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[2].SequenceAttr)
	ass.Equal("original", exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[2].DataformatAttr)
	ass.Equal("60095006", *exchangeObject.ExchBibliographicdata.ExchPriorityclaims.ExchPriorityclaim[2].Documentid.Docnumber)

	// parties
	ass.Equal("parties", exchangeObject.ExchBibliographicdata.ExchParties.XMLName.Local)
	// applicants
	ass.Equal("applicants", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.XMLName.Local)
	ass.Equal("applicant", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[0].XMLName.Local)
	ass.Equal("1", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[0].SequenceAttr)
	ass.Equal("docdb", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[0].DataformatAttr)
	ass.Equal("WARNER LAMBERT CO", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[0].ExchApplicantname[0].Name.Value)
	ass.Equal("US", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[0].Residence.Country)
	ass.Equal("1", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[1].SequenceAttr)
	ass.Equal("docdba", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[1].DataformatAttr)
	ass.Equal("WARNER-LAMBERT COMPANY", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[1].ExchApplicantname[0].Name.Value)
	ass.Equal("1", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[2].SequenceAttr)
	ass.Equal("original", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[2].DataformatAttr)
	ass.Equal("WARNER-LAMBERT COMPANY", exchangeObject.ExchBibliographicdata.ExchParties.ExchApplicants.ExchApplicant[2].ExchApplicantname[0].Name.Value)

	// inventors
	ass.Equal("inventors", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.XMLName.Local)
	for i := 0; i <= 2; i++ {
		ass.Equal(strconv.Itoa(i+1), exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[i].SequenceAttr)
		ass.Equal("docdb", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[i].DataformatAttr)
		ass.Equal("US", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[i].Residence.Country)
	}
	for i := 3; i <= 5; i++ {
		ass.Equal(strconv.Itoa(i-2), exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[i].SequenceAttr)
		ass.Equal("docdba", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[i].DataformatAttr)
	}
	for i := 6; i <= 8; i++ {
		ass.Equal(strconv.Itoa(i-5), exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[i].SequenceAttr)
		ass.Equal("original", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[i].DataformatAttr)
	}

	ass.Equal("O'BRIEN PATRICK MICHAEL", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[0].ExchInventorname[0].Name.Value)
	ass.Equal("PICARD JOSEPH ARMAND", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[1].ExchInventorname[0].Name.Value)
	ass.Equal("SLISKOVIC DRAGO ROBERT", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[2].ExchInventorname[0].Name.Value)
	ass.Equal("O'BRIEN, PATRICK MICHAEL", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[3].ExchInventorname[0].Name.Value)
	ass.Equal("PICARD, JOSEPH ARMAND", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[4].ExchInventorname[0].Name.Value)
	ass.Equal("SLISKOVIC, DRAGO ROBERT", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[5].ExchInventorname[0].Name.Value)
	ass.Equal("O'Brien, Patrick Michael", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[6].ExchInventorname[0].Name.Value)
	ass.Equal("Picard, Joseph Armand", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[7].ExchInventorname[0].Name.Value)
	ass.Equal("Sliskovic, Drago Robert", exchangeObject.ExchBibliographicdata.ExchParties.ExchInventors.ExchInventor[8].ExchInventorname[0].Name.Value)

	// invention title
	ass.Equal("invention-title", exchangeObject.ExchBibliographicdata.ExchInventiontitle[0].XMLName.Local)
	ass.Equal("en", exchangeObject.ExchBibliographicdata.ExchInventiontitle[0].LangAttr)
	ass.Equal("docdba", exchangeObject.ExchBibliographicdata.ExchInventiontitle[0].DataformatAttr)
	ass.Equal("TRICYCLIC SULFONAMIDES AND THEIR DERIVATIVES AS INHIBITORS OF MATRIX METALLOPROTEINASES", exchangeObject.ExchBibliographicdata.ExchInventiontitle[0].Value)
	ass.Equal("sh", exchangeObject.ExchBibliographicdata.ExchInventiontitle[1].LangAttr)
	ass.Equal("docdba", exchangeObject.ExchBibliographicdata.ExchInventiontitle[1].DataformatAttr)
	ass.Equal("TRICIKLIČNI SULFONAMIDI I NJIHOVI DERIVATI KAO INHIBITORI MATRIČNIH METALOPROTEINAZA", exchangeObject.ExchBibliographicdata.ExchInventiontitle[1].Value)

	ass.Equal("dates-of-public-availability", exchangeObject.ExchBibliographicdata.ExchDatesofpublicavailability.XMLName.Local)
	// examined/ unexamined (not) printed with(out) grant
	ass.Equal("unexamined-not-printed-without-grant", exchangeObject.ExchBibliographicdata.ExchDatesofpublicavailability.ExchUnexaminednotprintedwithoutgrant.XMLName.Local)
	ass.Equal(20031231, exchangeObject.ExchBibliographicdata.ExchDatesofpublicavailability.ExchUnexaminednotprintedwithoutgrant.Documentid.Date)

	// abstract
	ass.Equal("abstract", exchangeObject.ExchAbstract[0].XMLName.Local)
	ass.Equal("sh", exchangeObject.ExchAbstract[0].LangAttr)
	ass.Equal("docdba", exchangeObject.ExchAbstract[0].DataformatAttr)
	ass.Equal("national office", exchangeObject.ExchAbstract[0].AbstractsourceAttr)
	ass.Equal("p", exchangeObject.ExchAbstract[0].ExchP[0].XMLName.Local)
	ass.Equal(1030, len(exchangeObject.ExchAbstract[0].ExchP[0].Value))

	ass.Equal("abstract", exchangeObject.ExchAbstract[1].XMLName.Local)
	ass.Equal("en", exchangeObject.ExchAbstract[1].LangAttr)
	ass.Equal("docdba", exchangeObject.ExchAbstract[1].DataformatAttr)
	ass.Equal("national office", exchangeObject.ExchAbstract[1].AbstractsourceAttr)
	ass.Equal("p", exchangeObject.ExchAbstract[1].ExchP[0].XMLName.Local)
	ass.Equal(996, len(exchangeObject.ExchAbstract[1].ExchP[0].Value))

	ass.Equal("abstract", exchangeObject.ExchAbstract[2].XMLName.Local)
	ass.Equal("original", exchangeObject.ExchAbstract[2].DataformatAttr)
	ass.Equal("p", exchangeObject.ExchAbstract[2].ExchP[0].XMLName.Local)
	ass.Equal(2125, len(exchangeObject.ExchAbstract[2].ExchP[0].Value))

	fmt.Println(exchangeObject.ExchAbstract[0].ExchP[0].Value)
}

func TestReadFileParsingErr(t *testing.T) {
	// process file
	exchangeObject, err := ParseXmlFileToStruct("./test-data/WO-2023012807-A1_507242069.xml")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(exchangeObject.ExchAbstract[0].ExchP[0].Value)
	return
}

func TestReadFileParsingErrAll(t *testing.T) {

	var testDataDir = "./test-data"

	files, err := os.ReadDir(testDataDir)
	if err != nil {
		log.Fatal(err)
	}
	var filenames []string
	for _, file := range files {
		if !file.IsDir() && strings.Contains(file.Name(), ".xml") {
			filenames = append(filenames, file.Name())
		}
	}
	for _, filename := range filenames {
		fp := filepath.Join(testDataDir, filename)
		exchangeObject, err := ParseXmlFileToStruct(fp)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(exchangeObject.ExchBibliographicdata.ExchInventiontitle[0].Value)
	}
}

func TestReadBigFileProcessExchangeFileContent(t *testing.T) {
	// t.Skip()

	var parserHandler = func(fileName, fileContent string) {
		doc, err := ParseXmlStringToStruct(fileContent)
		if err != nil {
			t.Error(err)
			return
		}
		dateString := strconv.Itoa(doc.DatepublAttr)
		publicationDate, err := time.Parse("20060102", dateString)
		if err != nil {
			slog.With("err", err).Error("can not parse date")
			// set data to 9999-12-31
		} else {
			slog.With("publicationDate", publicationDate.Format("2006-01-02")).Info("publicationDate")
		}
	}

	var testDataPath = "./big-test-data/DOCDB-202407-025-US-0090.xml"

	// open file
	file, err := os.Open(testDataPath)
	if err != nil {
		t.Error(err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			t.Error(err)
		}
	}(file)

	p := NewProcessor()
	p.SetContentHandler(parserHandler)

	// process file
	err = p.ProcessExchangeFileContent(slog.With("test", testDataPath), file)
	if err != nil {
		t.Error(err)
	}
}
