package epo_docdb_sqllogger

import (
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ProcessDirectorySQL struct {
	gorm.Model
	ProcessingDir string `gorm:"unique"`
	Done          bool
	BulkFilesSQL  []BulkFileSQL `gorm:"foreignkey:ProcessDirectorySQLID"`
}

// One Process Directory has many Bulk Files (zip)
type BulkFileSQL struct {
	gorm.Model
	ZipName          string
	Done             bool
	ExchangeFilesSQL []ExchangeFileSQL `gorm:"foreignkey:BulkFileSQLID"`
}

// One Bulk File has many Exchange Docs (one doc = one line in XML)
type ExchangeFileSQL struct {
	gorm.Model
	Name          string
	Done          bool
	Info          string
	Path          string
	StatusAttr    string
	Finished      *time.Time //changed to pointer
	BulkFileSQLID uint       `gorm:"index"` // Foreign key to BulkFileSQL
}

func (p *SqlLogger) Inizialize() {
	mydb, err := gorm.Open(sqlite.Open(p.DatabasePath), &gorm.Config{})
	if err != nil {
		panic("failed to open " + p.DatabasePath)
	}

	p.db = mydb
	//this will create the database, or simply use it if it does not exist
	p.db.AutoMigrate(&ProcessDirectorySQL{}, &BulkFileSQL{}, &ExchangeFileSQL{})

	var processDirSQL ProcessDirectorySQL

	result := p.db.Where("processing_dir = ?", p.ProcessingDir).First(&processDirSQL)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			fmt.Println("No record for this processing path, so its a new process")
		} else {
			// Handle other possible errors
			log.Printf("Failed to query database: %v", result.Error)
		}
	}

}
