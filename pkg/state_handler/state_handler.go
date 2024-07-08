package state_handler

import (
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"strconv"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ProcessStatus string

const (
	// Todo is the status of a process that has not started yet
	Todo ProcessStatus = "todo"
	// PartlyProcessed is the status of a process that has started but not finished
	PartlyProcessed ProcessStatus = "partly"
	// Done is the status of a process that has finished
	Done ProcessStatus = "done"
)

// ProcessDirectorySQL represents the directory that is being processed
type ProcessDirectorySQL struct {
	gorm.Model
	ProcessingDir string `gorm:"unique"`
	DatabasePath  string
	Status        ProcessStatus
	Info          string
	//Finished   *time.Time //changed to pointer
	ZipFilesSQL []ZipFileSQL `gorm:"foreignkey:ProcessDirID"`
}

// ZipFileSQL represents a the bulk files (zip)
type ZipFileSQL struct {
	gorm.Model
	ProcessDirID uint `gorm:"index"` //foreign key
	ZipName      string
	FullPath     string
	Status       ProcessStatus
	Info         string
	//Finished   *time.Time //changed to pointer
	XMLFilesSQL []XMLFileSQL `gorm:"foreignkey:ZipFileID"`
}

// XMLFileSQL represents single XML files inside the bulk file
type XMLFileSQL struct {
	gorm.Model
	ZipFileID uint `gorm:"index"` // Foreign key to BulkFileSQL
	XmlName   string
	Status    ProcessStatus
	Info      string
	FullPath  string
	//Finished   *time.Time //changed to pointer
	ExchangeFilesSQL []ExchangeLineSQL `gorm:"foreignkey:XMLFileID"`
}

// ExchangeLineSQL represents a single exchange document inside the XML file
type ExchangeLineSQL struct {
	gorm.Model
	XMLFileID    uint `gorm:"index"` // Foreign key to BulkFileSQL
	ExchangeName string
	Status       ProcessStatus
	Info         string
	FullPath     string
	//Finished   *time.Time //changed to pointer
}

// Initialize loads Last Known State
// Creates DB if there is none
// returns false if the processing is already finished
// returns true if there is some processing left to be done
func (sh *StateHandler) Initialize() {
	db, err := gorm.Open(sqlite.Open(sh.DatabasePath), &gorm.Config{})
	if err != nil {
		panic("failed to open " + sh.DatabasePath)
	}

	sh.db = db
	// this will create the tables in the database, or migrate them if they already exist
	err = sh.db.AutoMigrate(&ProcessDirectorySQL{}, &ZipFileSQL{}, &XMLFileSQL{}, &ExchangeLineSQL{})
	if err != nil {
		slog.With("err", err).Error("could not migrate")
		return
	}

	// Load the Process Dir Struct if it doesn't exist
	var processDirSQL ProcessDirectorySQL

	dirResult := sh.db.Where("processing_dir = ?", sh.ProcessingDir).First(&processDirSQL)
	if dirResult.Error != nil {
		if errors.Is(dirResult.Error, gorm.ErrRecordNotFound) {
			fmt.Println("No record for this processing path, creating one")
			processDir := ProcessDirectorySQL{
				ProcessingDir: sh.ProcessingDir,
				DatabasePath:  sh.DatabasePath,
				Status:        Todo,
				Info:          "New Directory Process Started",
			}
			sh.db.Create(&processDir)
			sh.ProcessingDirSQL = processDir
			return
		} else {
			panic(dirResult.Error)
		}
	}

	// loaded successfully, so cache it in the SqlLogger struct
	sh.ProcessingDirSQL = processDirSQL
}

// SetSafeDelete no if directory.done = false or no entry exists
func (sh *StateHandler) SetSafeDelete(status bool) {
	sh.SafeDeleteOnly = status
}

// GetDirectoryProcessStatus no if directory.done = false or no entry exists
func (sh *StateHandler) GetDirectoryProcessStatus() (ProcessStatus, error) {
	return sh.ProcessingDirSQL.Status, nil
}

// RegisterOrSkipZipFile returns Done if the Bulk file is already processed
// If the Bulk file entry does not exist,
// creates a new one (using the current processDir as foreign key)
// or loads the existing bulk file information if the entry exists but is not done
func (sh *StateHandler) RegisterOrSkipZipFile(fileName string) (ProcessStatus, error) {
	// So a directory has started processing, but never finished, get the last known ZIP File
	// The Processor starts at the last known ZIP File, not at the specific Exchange Document
	// find the last unfinished zip file
	var zipFile ZipFileSQL
	errBulkFile := sh.db.Where("zip_name = ?", fileName).First(&zipFile).Error
	if errBulkFile != nil {
		if errors.Is(errBulkFile, gorm.ErrRecordNotFound) {
			slog.With("fileName", fileName).Info("No record for this zip file, creating")
			newZipFile := ZipFileSQL{
				ZipName:      fileName,
				Status:       Todo,
				FullPath:     filepath.Join(sh.ProcessingDir, fileName),
				ProcessDirID: sh.ProcessingDirSQL.ID,
				Info:         "New Zip Process Started",
			}
			errCreate := sh.db.Create(&newZipFile).Error
			if errCreate != nil {
				slog.With("err", errCreate).Error("failed to create zip file")
			}
			sh.currentZipFileSQL = newZipFile
			return Todo, nil //new processing project
		} else {
			return Todo, errBulkFile
		}
	}

	//won't be registered if already done
	if zipFile.Status != Done {
		sh.currentZipFileSQL = zipFile
	}

	return zipFile.Status, nil
}

// RegisterOrSkipXMLFile returns Done if the Bulk file is already processed
// If the XML file entry does not exist,
// creates a new one (using the current Zip file as foreign key)
// or loads the existing bulk file information if the entry exists but is not done
func (sh *StateHandler) RegisterOrSkipXMLFile(fileName string, innerZipPath string) (ProcessStatus, error) {
	// So a directory has started processing, but never finished, get the last known ZIP File
	// The Processor starts at the last known ZIP File, not at the specific Exchange Document
	// find the last unfinished zip file
	var xmlFile XMLFileSQL
	errXMLFile := sh.db.Where("xml_name = ?", fileName).First(&xmlFile).Error
	if errXMLFile != nil {
		if errors.Is(errXMLFile, gorm.ErrRecordNotFound) {
			fmt.Println("No record for this xml file, creating")
			newXmlFile := XMLFileSQL{
				XmlName:   fileName,
				Status:    Todo,
				FullPath:  sh.currentZipFileSQL.FullPath + "::" + innerZipPath + fileName,
				ZipFileID: sh.currentZipFileSQL.ID,
				Info:      "New XML Process Started",
			}
			errCreate := sh.db.Create(&newXmlFile).Error
			if errCreate != nil {
				slog.With("err", errCreate).Error("failed to create xml file")
			}
			sh.currentXMLFileSQL = xmlFile
			return Todo, nil
		} else {
			return Todo, errXMLFile
		}
	}

	// won't be registered if already done
	if xmlFile.Status != Done {
		sh.currentXMLFileSQL = xmlFile
	}

	return xmlFile.Status, nil
}

// RegisterOrSkipExchangeLine returns Done if the Bulk file is already processed
// If the XML file entry does not exist,
// creates a new one (using the current Zip file as foreign key)
// or loads the existing bulk file information if the entry exists but is not done
func (sh *StateHandler) RegisterOrSkipExchangeLine(exchangeID string, lineNumber int) (ProcessStatus, error) {
	//So a directory has started processing, but never finished, get the last known ZIP File
	//The Processor starts at the last known ZIP File, not at the specific Exchange Document
	//find the last unfinished zip file
	var exchangeLine ExchangeLineSQL
	errExchangeFile := sh.db.Where("exchange_name = ?", exchangeID).First(&exchangeLine).Error
	if errExchangeFile != nil {
		if errors.Is(errExchangeFile, gorm.ErrRecordNotFound) {
			fmt.Println("No record for this xml file, creating")
			newExchangeLine := ExchangeLineSQL{
				XMLFileID:    sh.currentXMLFileSQL.ID,
				ExchangeName: exchangeID,
				Status:       Todo,
				Info:         "new exchange line",
				FullPath:     sh.currentXMLFileSQL.FullPath + "::" + exchangeID + " (line: " + strconv.Itoa(lineNumber) + ")",
			}
			errCreate := sh.db.Create(&newExchangeLine).Error
			if errCreate != nil {
				slog.With("err", errCreate).Error("failed to create exchange line")
			}
			sh.currentExchangeLineSQL = newExchangeLine
			return Todo, nil //new processing project
		} else {
			return Todo, errExchangeFile
		}
	}

	// won't be registered if already done
	if exchangeLine.Status != Done {
		sh.currentExchangeLineSQL = exchangeLine
	}

	return exchangeLine.Status, nil
}

// MarkProcessingDirectoryAsFinished sets the status of directory as finished, no deleting downwards
func (sh *StateHandler) MarkProcessingDirectoryAsFinished() {
	// set the status of the directory as finished
	err := sh.db.Model(&sh.ProcessingDirSQL).Update("status", Done).Error
	if err != nil {
		panic(err)
	}
	// set the info of the directory as finished
	err = sh.db.Model(&sh.ProcessingDirSQL).Update("info", "finished").Error
	if err != nil {
		panic(err)
	}
}

// MarkZipFileAsFinished if a finishes, delete all recorded XML lines
// And mark the Zip as finished, but always keep it
func (sh *StateHandler) MarkZipFileAsFinished() {
	//Delete All Exchange Files belonging to the current XML File
	var resultXMLDelete *gorm.DB

	if sh.SafeDeleteOnly {
		resultXMLDelete = sh.db.Where("zip_file_id = ?", sh.currentZipFileSQL.ID).Delete(&XMLFileSQL{})
	} else {
		resultXMLDelete = sh.db.Unscoped().Where("zip_file_id = ?", sh.currentZipFileSQL.ID).Delete(&XMLFileSQL{})
	}

	if resultXMLDelete.Error != nil {
		panic(resultXMLDelete.Error)
	}

	// Check deleted records
	fmt.Printf("Deleted %v record(s) with name '%s'", resultXMLDelete.RowsAffected, sh.currentXMLFileSQL.XmlName)
	// set current Exchange File to empty
	sh.currentXMLFileSQL = XMLFileSQL{}

	//We're keeping the Zip Entry for now
	resultInfo := sh.db.Model(&sh.currentZipFileSQL).Update("info", "finished")
	if resultInfo.Error != nil {
		panic(resultInfo.Error)
	}

	resultStatus := sh.db.Model(&sh.currentZipFileSQL).Update("status", Done)
	if resultStatus.Error != nil {
		panic(resultStatus.Error)
	}
}

// MarkXMLAsFinished if an XML finishes, delete all recorded exchange lines
// And mark the XML as finished, but keep it
// the xml gets deleted when the whole Zip is finished
func (sh *StateHandler) MarkXMLAsFinished() {
	//Delete All Exchange Files belonging to the current XML File
	var resultExchangeDelete *gorm.DB

	if sh.SafeDeleteOnly {
		resultExchangeDelete = sh.db.Where("xml_file_id = ?", sh.currentXMLFileSQL.ID).Delete(&ExchangeLineSQL{})
	} else {
		resultExchangeDelete = sh.db.Unscoped().Where("xml_file_id = ?", sh.currentXMLFileSQL.ID).Delete(&ExchangeLineSQL{})
	}

	if resultExchangeDelete.Error != nil {
		panic(resultExchangeDelete.Error)
	}

	// Check deleted records
	fmt.Printf("Deleted %v record(s) with ID '%v'", resultExchangeDelete.RowsAffected, sh.currentXMLFileSQL.ID)
	// set current Exchange File to empty
	sh.currentExchangeLineSQL = ExchangeLineSQL{}

	//We're keeping the XML Entry for now
	resultInfo := sh.db.Model(&sh.currentXMLFileSQL).Update("info", "finished")
	if resultInfo.Error != nil {
		panic(resultInfo.Error)
	}

	resultStatus := sh.db.Model(&sh.currentXMLFileSQL).Update("status", Done)
	if resultStatus.Error != nil {
		panic(resultStatus.Error)
	}
}

// MarkExchangeFileAsFinished exchange Records only get deleted when the XML is done
func (sh *StateHandler) MarkExchangeFileAsFinished() {
	resultStatus := sh.db.Model(&sh.currentExchangeLineSQL).Update("status", Done)

	if resultStatus.Error != nil {
		panic(resultStatus.Error)
	}

	resultInfo := sh.db.Model(&sh.currentExchangeLineSQL).Update("info", "finished")
	if resultInfo.Error != nil {
		panic(resultInfo.Error)
	}
}

// IsDirectoryFinished returns true if the directory is finished
func (sh *StateHandler) IsDirectoryFinished() bool {
	return sh.ProcessingDirSQL.Status == Done
}
