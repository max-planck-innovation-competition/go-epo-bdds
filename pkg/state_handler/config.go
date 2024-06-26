package state_handler

import (
	"gorm.io/gorm"
	"path/filepath"
)

// StateHandler contains the config for the state handler
type StateHandler struct {
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

// New creates a new state handler
func New(databaseName string, databaseDir string, processingDir string) *StateHandler {
	stateHandler := StateHandler{
		DatabaseName:   databaseName,
		DatabaseDir:    databaseDir,
		ProcessingDir:  processingDir,
		DatabasePath:   filepath.Join(databaseDir, databaseName),
		SafeDeleteOnly: false,
	}
	stateHandler.Initialize() //Initializes the other fields
	return &stateHandler
}
