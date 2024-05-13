package epo_docdb_sqllogger

import "gorm.io/gorm"

type SqlLogger struct {
	//initialize these
	DatabaseName  string //e.g. log.db, for the initializer
	DatabaseDir   string //absolute db path, e.g. C:\docdb\
	ProcessingDir string //directory containing the downloaded zip files
	//these are implied and generated in NewSqlLogger(...)
	CurrentZip         string //currentZipFile
	CurrentExchangeDoc string //current Exchange Doc in ZipFile
	CurrentStatusAttr  string //current Status Attribute (e.g. A, D etc.)
	DatabasePath       string //Database Dir + Database Name
	Done               bool   //The SQL Logger is done when every zip file in dir is processed
	db                 *gorm.DB
}

func (p *SqlLogger) HasDirectoryBeenProcessed() bool {
	return false
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
