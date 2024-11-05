# EPO DocDB

This package provides the ETL process to process the EPO DocDB data.
It can be used to ingest the data into any database or any file format.

## Usage 

The `Processor` struct can be used to process the data.

## Environment Variables

```env
EPO_USERNAME=hello@world.com
EPO_PASSWORD=******
DOCDB_FRONTFILES_PATH=/media/oem/Volume/docdb/frontfiles
DOCDB_BACKFILES_PATH=/media/oem/Volume/docdb/backfiles_2024_02_27
```

## Usage

```go

// your custom handler
var parserHandler = func(fileName, fileContent string) {
	// converts the docdb xml string to a golang struct
    doc, err := epo_docdb.ParseXmlStringToStruct(fileContent)
    if err != nil {
		slog.With("err", err).Error("can not parse xml")
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


p := epo_docdb.NewProcessor()
// Include the file types you want to process
p.IncludeFileTypes("CreateDelete", "bck")
// Include the authorities you want to process
p.IncludeAuthorities("EP", "US", "WO")
// Set the content handler
p.SetContentHandler(parserHandler)

err := p.ProcessDirectory("/docdb/backfiles")
```


