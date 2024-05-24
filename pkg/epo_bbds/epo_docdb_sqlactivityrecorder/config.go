package epo_docdb_sqlactivityrecorder

import "gorm.io/gorm"

// Fitting Name
// Exchange Lines
// Delete Down
type SQLActivityRecorder struct {
	//initialize these
	DatabaseName   string //e.g. log.db, for the initializer
	DatabaseDir    string //path of the .db, e.g. C:\docdb\ or .\ for relative path
	ProcessingDir  string //directory containing the downloaded zip files
	SafeDeleteOnly bool
	//for the state
	//these are initialized in NewSqlLogger(...)
	ProcessingDirSQL       ProcessDirectorySQL //SQL Struct of current processing directory for faster access
	currentZipFileSQL      ZipFileSQL
	currentXMLFileSQL      XMLFileSQL
	currentExchangeLineSQL ExchangeLineSQL
	DatabasePath           string //Database Dir + Database Name
	db                     *gorm.DB
}

// NewProcessor creates a new processor
// the default handler is PrintLineHandler
func NewSqlActivityRecorder(DatabaseName string, DatabaseDir string, ProcessingDir string) *SQLActivityRecorder {
	p := SQLActivityRecorder{
		DatabaseName:   DatabaseName,
		DatabaseDir:    DatabaseDir,
		ProcessingDir:  ProcessingDir,
		DatabasePath:   DatabaseDir + DatabaseName,
		SafeDeleteOnly: false,
	}
	p.Initialize() //Initializes the other fields
	return &p
}
