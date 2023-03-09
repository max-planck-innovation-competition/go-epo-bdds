# GO EPO-DOCDB

Go API Client for the European Patent Office DocDB Data

## Status

Alpha Version

**⚠️ Experimental - Not ready for production.**

## Author
Sebastian Erhardt

## DocDB

The structs for the docdb data are generated from the xsd files provided by the EPO.

The [xgen](https://github.com/xuri/xgen) library was used to generate the structs.
```
xgen -i /path/to/your/xsd -o /path/to/your/output -l Go
```

## Environment Variables

```
EPO_USERNAME=XYZ
EPO_PASSWORD=XXXXXX
```