package state_handler

import (
	"errors"
	"log/slog"
	"path/filepath"
	"strconv"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// PathDelimiter is the delimiter used to separate the path of the objects
const PathDelimiter = "/"

// ProcessStatus represents the status of a process
type ProcessStatus string

const (
	// Todo is the status of a process that has not started yet
	Todo ProcessStatus = "todo"
	// Done is the status of a process that has finished
	Done ProcessStatus = "done"
	// Error is the status of a process that has finished with an error
	Error ProcessStatus = "error"
)

// Object represents a single object in the database
// this will be used to store the last known state of the processing
type Object struct {
	Path   string        `gorm:"primary_key"` // just the path of the object
	Status ProcessStatus // the status of the object
}

// RegisterOrSkip returns Done if the object is already processed
func (sh *StateHandler) RegisterOrSkip(path string) (skip bool, err error) {
	// check if the parent is done
	// if the parent is not done, the object is not done
	// if the parent is done, the object is done
	// if the parent does not exist, the object is not done
	pathParts := strings.Split(path, PathDelimiter)
	if len(pathParts) > 1 {
		for i := 1; i < len(pathParts); i++ {
			parentPath := strings.Join(pathParts[:i], PathDelimiter)
			parentObj, err := sh.Get(parentPath)
			if err != nil {
				slog.With("err", err).Error("could not get parent object")
				return true, err
			}
			// if the parent is done skip the object
			if parentObj.Status == Done {
				return true, nil
			}
		}
	}
	// get the object
	// if its already in the database, return if it's done
	// if it's not in the database, create it and return not done
	obj, err := sh.Get(path)
	if err != nil {
		slog.With("err", err).Error("could not get object")
		return true, err
	}
	return obj.Status == Done, nil
}

// Set sets the status of the object with the given path
func (sh *StateHandler) Set(path string, status ProcessStatus) error {
	return sh.db.Save(&Object{Path: path, Status: status}).Error
}

// Get returns the object with the given path
// if the object does not exist, it will be created with the status Todo
func (sh *StateHandler) Get(path string) (obj *Object, err error) {
	var object Object
	err = sh.db.Where("path = ?", path).First(&object).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = sh.Set(path, Todo)
			if err != nil {
				slog.With("err", err).Error("could not set object")
				return nil, err
			}
			return &Object{Path: path, Status: Todo}, nil
		}
		return nil, err
	}
	return &object, nil
}

// MarkAsDone marks the object as done
// and deletes all objects that are children of the object
func (sh *StateHandler) MarkAsDone(path string) error {
	tx := sh.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	// mark object as done
	err := tx.Model(&Object{}).Where("path = ?", path).Update("status", Done).Error
	if err != nil {
		slog.With("err", err).Error("could not mark as done")
		tx.Rollback()
	}
	// delete all objects that are children of the object
	err = tx.Where("path LIKE ?", path+"/%").Unscoped().Delete(&Object{}).Error
	if err != nil {
		slog.With("err", err).Error("could not delete children")
		tx.Rollback()
	}
	return tx.Commit().Error
}

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
	db, err := gorm.Open(sqlite.Open(sh.DatabasePath), &gorm.Config{
		Logger: nil,
	})
	if err != nil {
		panic("failed to open " + sh.DatabasePath)
	}

	sh.db = db
	// this will create the tables in the database, or migrate them if they already exist
	err = sh.db.AutoMigrate(
		&Object{},
		&ProcessDirectorySQL{},
		&ZipFileSQL{},
		&XMLFileSQL{},
		&ExchangeLineSQL{},
	)
	if err != nil {
		slog.With("err", err).Error("could not migrate")
		return
	}

	// Load the Process Dir Struct if it doesn't exist
	var processDirSQL ProcessDirectorySQL

	dirResult := sh.db.Where("processing_dir = ?", sh.ProcessingDir).First(&processDirSQL)
	if dirResult.Error != nil {
		if errors.Is(dirResult.Error, gorm.ErrRecordNotFound) {
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
	// set current Exchange File to empty
	sh.currentExchangeLineSQL = ExchangeLineSQL{}

	//We're keeping the XML Entry for now
	resultInfo := sh.db.Model(&sh.currentXMLFileSQL).
		Where("id", &sh.currentZipFileSQL.ID).
		Update("info", "finished")
	if resultInfo.Error != nil {
		panic(resultInfo.Error)
	}

	resultStatus := sh.db.Model(&sh.currentXMLFileSQL).
		Where("id", &sh.currentZipFileSQL.ID).
		Update("status", Done)
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
