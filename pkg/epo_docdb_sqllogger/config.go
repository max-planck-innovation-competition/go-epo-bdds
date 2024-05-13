package epo_docdb_sqllogger

import "gorm.io/gorm"

type SqlLogger struct {
	//initialize these
	DatabaseName  string //e.g. log.db, for the initializer
	DatabaseDir   string //path of the .db, e.g. C:\docdb\ or .\ for relative path
	ProcessingDir string //directory containing the downloaded zip files
	//these are initialized in NewSqlLogger(...)
	ProcessingDirSQL   ProcessDirectorySQL //SQL Struct of current processing directory for faster access
	currentBulkFileSQL BulkFileSQL
	currentXMLFileSQL  XMLFileSQL
	DatabasePath       string //Database Dir + Database Name
	db                 *gorm.DB
}

// NewProcessor creates a new processor
// the default handler is PrintLineHandler
func NewSqlLogger(DatabaseName string, DatabaseDir string, ProcessingDir string) *SqlLogger {
	p := SqlLogger{
		DatabaseName:  DatabaseName,
		DatabaseDir:   DatabaseDir,
		ProcessingDir: ProcessingDir,
		DatabasePath:  DatabaseDir + DatabaseName,
	}
	p.Inizialize() //Initializes the other fields, gets the last known state
	return &p
}
