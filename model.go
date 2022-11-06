package go_epo_docdb

import "encoding/xml"

type ExchangeDocuments struct {
	XMLName           xml.Name `xml:"exchange-documents"`
	Text              string   `xml:",chardata"`
	Exch              string   `xml:"exch,attr"`
	Xsi               string   `xml:"xsi,attr"`
	SchemaLocation    string   `xml:"schemaLocation,attr"`
	DateOfExchange    string   `xml:"date-of-exchange,attr"`
	DtdVersion        string   `xml:"dtd-version,attr"`
	File              string   `xml:"file,attr"`
	NoOfDocuments     string   `xml:"no-of-documents,attr"`
	OriginatingOffice string   `xml:"originating-office,attr"`
	ExchangeDocument  []struct {
		Text                   string `xml:",chardata"`
		Country                string `xml:"country,attr"`
		DocNumber              string `xml:"doc-number,attr"`
		Kind                   string `xml:"kind,attr"`
		DocID                  string `xml:"doc-id,attr"`
		DatePubl               string `xml:"date-publ,attr"`
		FamilyID               string `xml:"family-id,attr"` // family-identifier of the DOCDB simple patent family representing the family that the publication is in at the time of 		exchange
		IsRepresentative       string `xml:"is-representative,attr"`
		DateOfLastExchange     string `xml:"date-of-last-exchange,attr"`
		DateOfPreviousExchange string `xml:"date-of-previous-exchange,attr"`
		DateAddedDocdb         string `xml:"date-added-docdb,attr"`
		OriginatingOffice      string `xml:"originating-office,attr"`
		Status                 string `xml:"status,attr"`
		BibliographicData      struct {
			Text                 string `xml:",chardata"`
			PublicationReference []struct {
				Text       string `xml:",chardata"`
				DataFormat string `xml:"data-format,attr"`
				DocumentID struct {
					Text      string `xml:",chardata"`
					Lang      string `xml:"lang,attr"`
					Country   string `xml:"country"`
					DocNumber string `xml:"doc-number"`
					Kind      string `xml:"kind"`
					Date      string `xml:"date"`
				} `xml:"document-id"`
			} `xml:"publication-reference"`
			ClassificationIpc struct {
				Text               string `xml:",chardata"`
				Edition            string `xml:"edition"`
				MainClassification string `xml:"main-classification"`
			} `xml:"classification-ipc"`
			ClassificationsIpcr struct {
				Text               string `xml:",chardata"`
				Status             string `xml:"status,attr"`
				ClassificationIpcr []struct {
					Chardata string `xml:",chardata"`
					Sequence string `xml:"sequence,attr"`
					Text     string `xml:"text"`
				} `xml:"classification-ipcr"`
			} `xml:"classifications-ipcr"`
			PatentClassifications struct {
				Text                 string `xml:",chardata"`
				PatentClassification []struct {
					Text                 string `xml:",chardata"`
					Sequence             string `xml:"sequence,attr"`
					ClassificationScheme struct {
						Text   string `xml:",chardata"`
						Office string `xml:"office,attr"`
						Scheme string `xml:"scheme,attr"`
						Date   string `xml:"date"`
					} `xml:"classification-scheme"`
					ClassificationSymbol     string `xml:"classification-symbol"`
					SymbolPosition           string `xml:"symbol-position"`
					ClassificationValue      string `xml:"classification-value"`
					ClassificationStatus     string `xml:"classification-status"`
					ClassificationDataSource string `xml:"classification-data-source"`
					GeneratingOffice         string `xml:"generating-office"`
					ActionDate               struct {
						Text string `xml:",chardata"`
						Date string `xml:"date"`
					} `xml:"action-date"`
				} `xml:"patent-classification"`
				CombinationSet []struct {
					Text            string `xml:",chardata"`
					Sequence        string `xml:"sequence,attr"`
					GroupNumber     string `xml:"group-number"`
					CombinationRank []struct {
						Text                 string `xml:",chardata"`
						RankNumber           string `xml:"rank-number"`
						PatentClassification struct {
							Text                 string `xml:",chardata"`
							ClassificationScheme struct {
								Text   string `xml:",chardata"`
								Office string `xml:"office,attr"`
								Scheme string `xml:"scheme,attr"`
								Date   string `xml:"date"`
							} `xml:"classification-scheme"`
							ClassificationSymbol     string `xml:"classification-symbol"`
							SymbolPosition           string `xml:"symbol-position"`
							ClassificationValue      string `xml:"classification-value"`
							ClassificationStatus     string `xml:"classification-status"`
							ClassificationDataSource string `xml:"classification-data-source"`
							GeneratingOffice         string `xml:"generating-office"`
							ActionDate               struct {
								Text string `xml:",chardata"`
								Date string `xml:"date"`
							} `xml:"action-date"`
						} `xml:"patent-classification"`
					} `xml:"combination-rank"`
				} `xml:"combination-set"`
			} `xml:"patent-classifications"`
			ApplicationReference []struct {
				Text             string `xml:",chardata"`
				IsRepresentative string `xml:"is-representative,attr"`
				DocID            string `xml:"doc-id,attr"`
				DataFormat       string `xml:"data-format,attr"`
				DocumentID       struct {
					Text      string `xml:",chardata"`
					Country   string `xml:"country"`
					DocNumber string `xml:"doc-number"`
					Kind      string `xml:"kind"`
					Date      string `xml:"date"`
				} `xml:"document-id"`
			} `xml:"application-reference"`
			LanguageOfPublication string `xml:"language-of-publication"`
			PriorityClaims        struct {
				Text          string `xml:",chardata"`
				PriorityClaim []struct {
					Text       string `xml:",chardata"`
					Sequence   string `xml:"sequence,attr"`
					DataFormat string `xml:"data-format,attr"`
					DocumentID struct {
						Text      string `xml:",chardata"`
						DocID     string `xml:"doc-id,attr"`
						Country   string `xml:"country"`
						DocNumber string `xml:"doc-number"`
						Kind      string `xml:"kind"`
						Date      string `xml:"date"`
					} `xml:"document-id"`
					PriorityActiveIndicator string `xml:"priority-active-indicator"`
					PriorityLinkageType     string `xml:"priority-linkage-type"`
				} `xml:"priority-claim"`
			} `xml:"priority-claims"`
			Parties struct {
				Text       string `xml:",chardata"`
				Applicants struct {
					Text      string `xml:",chardata"`
					Applicant []struct {
						Text          string `xml:",chardata"`
						Sequence      string `xml:"sequence,attr"`
						DataFormat    string `xml:"data-format,attr"`
						Status        string `xml:"status,attr"`
						ApplicantName struct {
							Text string `xml:",chardata"`
							Name string `xml:"name"`
						} `xml:"applicant-name"`
						Residence struct {
							Text    string `xml:",chardata"`
							Country string `xml:"country"`
						} `xml:"residence"`
					} `xml:"applicant"`
				} `xml:"applicants"`
				Inventors struct {
					Text     string `xml:",chardata"`
					Inventor []struct {
						Text         string `xml:",chardata"`
						Sequence     string `xml:"sequence,attr"`
						DataFormat   string `xml:"data-format,attr"`
						InventorName struct {
							Text string `xml:",chardata"`
							Name string `xml:"name"`
						} `xml:"inventor-name"`
						Residence struct {
							Text    string `xml:",chardata"`
							Country string `xml:"country"`
						} `xml:"residence"`
					} `xml:"inventor"`
				} `xml:"inventors"`
			} `xml:"parties"`
			InventionTitle struct {
				Text       string `xml:",chardata"`
				Lang       string `xml:"lang,attr"`
				DataFormat string `xml:"data-format,attr"`
			} `xml:"invention-title"`
			DatesOfPublicAvailability struct {
				Text             string `xml:",chardata"`
				PrintedWithGrant struct {
					Text       string `xml:",chardata"`
					DocumentID struct {
						Text string `xml:",chardata"`
						Date string `xml:"date"`
					} `xml:"document-id"`
				} `xml:"printed-with-grant"`
				GazetteReference struct {
					Text string `xml:",chardata"`
					Date string `xml:"date"`
				} `xml:"gazette-reference"`
			} `xml:"dates-of-public-availability"`
			ReferencesCited struct {
				Text     string `xml:",chardata"`
				Citation []struct {
					Text       string `xml:",chardata"`
					CitedPhase string `xml:"cited-phase,attr"`
					Sequence   string `xml:"sequence,attr"`
					Patcit     struct {
						Text       string `xml:",chardata"`
						Num        string `xml:"num,attr"`
						Dnum       string `xml:"dnum,attr"`
						DnumType   string `xml:"dnum-type,attr"`
						DocumentID struct {
							Text      string `xml:",chardata"`
							DocID     string `xml:"doc-id,attr"`
							Country   string `xml:"country"`
							DocNumber string `xml:"doc-number"`
							Kind      string `xml:"kind"`
							Name      string `xml:"name"`
							Date      string `xml:"date"`
						} `xml:"document-id"`
					} `xml:"patcit"`
					Nplcit struct {
						Chardata string `xml:",chardata"`
						Num      string `xml:"num,attr"`
						NplType  string `xml:"npl-type,attr"`
						Text     string `xml:"text"`
					} `xml:"nplcit"`
					RelPassage struct {
						Text    string `xml:",chardata"`
						Passage string `xml:"passage"`
					} `xml:"rel-passage"`
				} `xml:"citation"`
			} `xml:"references-cited"`
		} `xml:"bibliographic-data"`
		Abstract struct {
			Text           string `xml:",chardata"`
			Lang           string `xml:"lang,attr"`
			DataFormat     string `xml:"data-format,attr"`
			AbstractSource string `xml:"abstract-source,attr"`
			P              string `xml:"p"`
		} `xml:"abstract"`
		PatentFamily struct {
			Text         string `xml:",chardata"`
			FamilyMember []struct {
				Text                 string `xml:",chardata"`
				ApplicationReference []struct {
					Text             string `xml:",chardata"`
					DataFormat       string `xml:"data-format,attr"`
					IsRepresentative string `xml:"is-representative,attr"`
					DocumentID       struct {
						Text      string `xml:",chardata"`
						Country   string `xml:"country"`
						DocNumber string `xml:"doc-number"`
						Kind      string `xml:"kind"`
					} `xml:"document-id"`
				} `xml:"application-reference"`
				PublicationReference []struct {
					Text       string `xml:",chardata"`
					DataFormat string `xml:"data-format,attr"`
					Sequence   string `xml:"sequence,attr"`
					DocumentID struct {
						Text      string `xml:",chardata"`
						Country   string `xml:"country"`
						DocNumber string `xml:"doc-number"`
						Kind      string `xml:"kind"`
					} `xml:"document-id"`
				} `xml:"publication-reference"`
			} `xml:"family-member"`
			Abstract struct {
				Text           string `xml:",chardata"`
				Lang           string `xml:"lang,attr"`
				Country        string `xml:"country,attr"`
				DocNumber      string `xml:"doc-number,attr"`
				Kind           string `xml:"kind,attr"`
				DataFormat     string `xml:"data-format,attr"`
				AbstractSource string `xml:"abstract-source,attr"`
				P              string `xml:"p"`
			} `xml:"abstract"`
		} `xml:"patent-family"`
	} `xml:"exchange-document"`
}
