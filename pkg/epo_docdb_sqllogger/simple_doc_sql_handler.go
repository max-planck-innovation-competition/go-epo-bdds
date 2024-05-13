package epo_docdb_sqllogger

import (
	"errors"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ProcessStatus string

const (
	NotProcessed    ProcessStatus = "newproc"
	PartlyProcessed ProcessStatus = "partly"
	Done            ProcessStatus = "done"
)

type ProcessDirectorySQL struct {
	gorm.Model
	ProcessingDir string `gorm:"unique"`
	Status        ProcessStatus
	Info          string
	//Finished   *time.Time //changed to pointer
	BulkFilesSQL []BulkFileSQL `gorm:"foreignkey:ProcessDirID"`
}

// One Process Directory has many Bulk Files (zip)
type BulkFileSQL struct {
	gorm.Model
	ProcessDirID uint `gorm:"index"` //foreign key
	ZipName      string
	Status       ProcessStatus
	Info         string
	//Finished   *time.Time //changed to pointer
	ExchangeFilesSQL []XMLFileSQL `gorm:"foreignkey:BulkFileID"`
}

// One Bulk File has many XML Files
type XMLFileSQL struct {
	gorm.Model
	BulkFileID uint `gorm:"index"` // Foreign key to BulkFileSQL
	XMLName    string
	Status     ProcessStatus
	Info       string
	Path       string
	//Finished   *time.Time //changed to pointer
}

// Loads Last Known State
// Creates DB if there is none
// returns false if the processing is already finished
// returns true if there is some processing left to be done
func (p *SqlLogger) Inizialize() {
	mydb, err := gorm.Open(sqlite.Open(p.DatabasePath), &gorm.Config{})
	if err != nil {
		panic("failed to open " + p.DatabasePath)
	}

	p.db = mydb
	//this will create the database, or simply use it if it does not exist
	p.db.AutoMigrate(&ProcessDirectorySQL{}, &BulkFileSQL{}, &XMLFileSQL{})

	//Load the Process Dir Struct if it doesnt exist
	var processDirSQL ProcessDirectorySQL

	dirResult := p.db.Where("processing_dir = ?", p.ProcessingDir).First(&processDirSQL)
	if dirResult.Error != nil {
		if errors.Is(dirResult.Error, gorm.ErrRecordNotFound) {
			fmt.Println("No record for this processing path, creating one")
			processDir := ProcessDirectorySQL{
				ProcessingDir: p.DatabasePath,
				Status:        NotProcessed,
				Info:          "New Directory Process Started",
			}
			p.db.Create(&processDir)
			p.ProcessingDirSQL = processDir
		} else {
			panic(dirResult.Error)
		}
	}

	//loaded successfully, so cache it in the SqlLogger struct
	p.ProcessingDirSQL = processDirSQL
}

// no if directory.done = false or no entry exists
func (p *SqlLogger) GetDirectoryProcessStatus() (ProcessStatus, error) {
	return p.ProcessingDirSQL.Status, nil
}

// RegisterOrSkip: Returns Done if the Bulk file is already processed
// _______
// If the Bulk file entry does not exist,
// creates a new one (using the current processDir as foreign key)
// _______
// or loads the existing bulk file information if the entry exists but is not done
func (p *SqlLogger) RegisterOrSkipBulkFile(fileName string) (ProcessStatus, error) {
	//So a directory has started processing, but never finished, get the last known ZIP File
	//The Processor starts at the last known ZIP File, not at the specific Exchange Document
	//find the last unfinished zip file
	var bulkFile BulkFileSQL
	errBulkFile := p.db.Where("ZipName = ?", fileName).First(&bulkFile)
	if errBulkFile.Error != nil {
		if errors.Is(errBulkFile.Error, gorm.ErrRecordNotFound) {
			fmt.Println("No record for this bulk file, creating")
			bulkFile := BulkFileSQL{
				ZipName:      fileName,
				Status:       NotProcessed,
				ProcessDirID: p.ProcessingDirSQL.ID,
				Info:         "New Bulk Process Started",
			}
			p.db.Create(&bulkFile)
			p.currentBulkFileSQL = bulkFile
			return NotProcessed, nil //new processing project
		} else {
			return NotProcessed, errBulkFile.Error
		}
	}

	//wont be registered if already done
	if bulkFile.Status != Done {
		p.currentBulkFileSQL = bulkFile
	}

	return bulkFile.Status, nil
}

// RegisterOrSkip: Returns Done if the Bulk file is already processed
// _______
// If the Bulk file entry does not exist,
// creates a new one (using the current processDir as foreign key)
// _______
// or loads the existing bulk file information if the entry exists but is not done
func (p *SqlLogger) RegisterOrSkipXMLFile(fileName string) (ProcessStatus, error) {
	//So a directory has started processing, but never finished, get the last known ZIP File
	//The Processor starts at the last known ZIP File, not at the specific Exchange Document
	//find the last unfinished zip file
	var xmlFile XMLFileSQL
	errXMLFile := p.db.Where("XMLName = ?", fileName).First(&xmlFile)
	if errXMLFile.Error != nil {
		if errors.Is(errXMLFile.Error, gorm.ErrRecordNotFound) {
			fmt.Println("No record for this xml file, creating")
			xmlFile := XMLFileSQL{
				XMLName:    fileName,
				Status:     NotProcessed,
				BulkFileID: p.currentBulkFileSQL.ID,
				Info:       "New Bulk Process Started",
			}
			p.db.Create(&xmlFile)
			p.currentXMLFileSQL = xmlFile
			return NotProcessed, nil //new processing project
		} else {
			return NotProcessed, errXMLFile.Error
		}
	}

	//wont be registered if already done
	if xmlFile.Status != Done {
		p.currentXMLFileSQL = xmlFile
	}

	return xmlFile.Status, nil
}

func (p *SqlLogger) MarkProcessingDirectoryAsFinished() {
	resultStatus := p.db.Model(&p.ProcessingDirSQL).Update("status", Done)

	if resultStatus.Error != nil {
		panic(resultStatus.Error)
	}

	resultInfo := p.db.Model(&p.ProcessingDirSQL).Update("info", "finished")
	if resultInfo.Error != nil {
		panic(resultInfo.Error)
	}
}

func (p *SqlLogger) MarkBulkFileAsFinished() {
	resultStatus := p.db.Model(&p.currentBulkFileSQL).Update("status", Done)

	if resultStatus.Error != nil {
		panic(resultStatus.Error)
	}

	resultInfo := p.db.Model(&p.currentBulkFileSQL).Update("info", "finished")
	if resultInfo.Error != nil {
		panic(resultInfo.Error)
	}
}

func (p *SqlLogger) MarkXMLAsFinished() {
	resultStatus := p.db.Model(&p.currentXMLFileSQL).Update("status", Done)

	if resultStatus.Error != nil {
		panic(resultStatus.Error)
	}

	resultInfo := p.db.Model(&p.currentXMLFileSQL).Update("info", "finished")
	if resultInfo.Error != nil {
		panic(resultInfo.Error)
	}
}
