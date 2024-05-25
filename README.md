# GO Bulk Data Sets

Go API Client for the European Patent Office Bulk Data Sets and DocDB Data.

## Status

Alpha Version

**⚠️ Experimental - Not ready for production.**

## Author
Sebastian Erhardt

## DocDB

The structs for the DocDB data are generated from the xsd files provided by the EPO.

The [xgen](https://github.com/xuri/xgen) library was used to generate the structs.
```
xgen -i /path/to/your/xsd -o /path/to/your/output -l Go
```

## Environment Variables

```
EPO_USERNAME=XYZ
EPO_PASSWORD=XXXXXX
```
## Installation

```shell
go get -u github.com/max-planck-innovation-competition/go-epo-bdds
```

## Usage

There are separate packages for the bulk data service and the DocDB data.

With the bulk data service package, you can download the data from the EPO.
Among other digital data sets, the EPO provides the DocDB data.

With the DocDB package, you can process the DocDB data.

### Bulk Data Service

To interact with the bulk data service, you can use the `epo_bdds` package.

```go
package main
import (
	"github.com/max-planck-innovation-competition/go-epo-bdds/pkg/epo_bdds"
)

func main() {
    // Get the authorization token
    token, err := epo_bdds.GetAuthorizationToken()
	if err != nil {
        log.Fatal(err)
    }
	// get the products 
    products, err := epo_bdds.GetProducts(token)
	...
}
```

### DocDB

The DocDB data can be processed with the `Processor` struct.
Depending on the data you want to process, you can include or exclude authorities.

You can also define your own content handler to process the data.
For example, you can write the data to a database or a file.

```go
p := NewProcessor()
p.IncludeAuthorities("EP")
err := p.ProcessDirectory("/docdb/backfiles")
if err != nil {
    t.Error(err)
}
p.SetContentHandler(yourContentHandler)
```