package state_handler

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"

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
	DatabasePath  string
	Status        ProcessStatus
	Info          string
	//Finished   *time.Time //changed to pointer
	ZipFilesSQL []ZipFileSQL `gorm:"foreignkey:ProcessDirID"`
}

// One Process Directory has many Bulk Files (zip)
// XMLs löschen nach processing
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

// One Zip File has many XML Files
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

// One XML File has many Exchange Lines
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
func (p *StateHandler) Initialize() {
	mydb, err := gorm.Open(sqlite.Open(p.DatabasePath), &gorm.Config{})
	if err != nil {
		panic("failed to open " + p.DatabasePath)
	}

	p.db = mydb
	//this will create the database, or simply use it if it does not exist
	err = p.db.AutoMigrate(&ProcessDirectorySQL{}, &ZipFileSQL{}, &XMLFileSQL{}, &ExchangeLineSQL{})
	if err != nil {
		slog.With("err", err).Error("could not migrate")
		return
	}

	//Load the Process Dir Struct if it doesnt exist
	var processDirSQL ProcessDirectorySQL

	dirResult := p.db.Where("processing_dir = ?", p.ProcessingDir).First(&processDirSQL)
	if dirResult.Error != nil {
		if errors.Is(dirResult.Error, gorm.ErrRecordNotFound) {
			fmt.Println("No record for this processing path, creating one")
			processDir := ProcessDirectorySQL{
				ProcessingDir: p.ProcessingDir,
				DatabasePath:  p.DatabasePath,
				Status:        NotProcessed,
				Info:          "New Directory Process Started",
			}
			p.db.Create(&processDir)
			p.ProcessingDirSQL = processDir
			return
		} else {
			panic(dirResult.Error)
		}
	}

	//loaded successfully, so cache it in the SqlLogger struct
	p.ProcessingDirSQL = processDirSQL
}

// SetSafeDelete no if directory.done = false or no entry exists
func (p *StateHandler) SetSafeDelete(status bool) {
	p.SafeDeleteOnly = status
}

// GetDirectoryProcessStatus no if directory.done = false or no entry exists
func (p *StateHandler) GetDirectoryProcessStatus() (ProcessStatus, error) {
	return p.ProcessingDirSQL.Status, nil
}

// RegisterOrSkipZipFile returns Done if the Bulk file is already processed
// If the Bulk file entry does not exist,
// creates a new one (using the current processDir as foreign key)
// or loads the existing bulk file information if the entry exists but is not done
func (p *StateHandler) RegisterOrSkipZipFile(fileName string) (ProcessStatus, error) {
	// So a directory has started processing, but never finished, get the last known ZIP File
	// The Processor starts at the last known ZIP File, not at the specific Exchange Document
	// find the last unfinished zip file
	var zipFile ZipFileSQL
	errBulkFile := p.db.Where("zip_name = ?", fileName).First(&zipFile)
	if errBulkFile.Error != nil {
		if errors.Is(errBulkFile.Error, gorm.ErrRecordNotFound) {
			fmt.Println("No record for this bulk file, creating")
			zipFileNew := ZipFileSQL{
				ZipName:      fileName,
				Status:       NotProcessed,
				FullPath:     p.ProcessingDir + "\\" + fileName,
				ProcessDirID: p.ProcessingDirSQL.ID,
				Info:         "New Zip Process Started",
			}
			p.db.Create(&zipFileNew)
			p.currentZipFileSQL = zipFileNew
			return NotProcessed, nil //new processing project
		} else {
			return NotProcessed, errBulkFile.Error
		}
	}

	//wont be registered if already done
	if zipFile.Status != Done {
		p.currentZipFileSQL = zipFile
	}

	return zipFile.Status, nil
}

// RegisterOrSkipXMLFile returns Done if the Bulk file is already processed
// If the XML file entry does not exist,
// creates a new one (using the current Zip file as foreign key)
// or loads the existing bulk file information if the entry exists but is not done
func (p *StateHandler) RegisterOrSkipXMLFile(fileName string, innerZipPath string) (ProcessStatus, error) {
	// So a directory has started processing, but never finished, get the last known ZIP File
	// The Processor starts at the last known ZIP File, not at the specific Exchange Document
	// find the last unfinished zip file
	var xmlFile XMLFileSQL
	errXMLFile := p.db.Where("xml_name = ?", fileName).First(&xmlFile)
	if errXMLFile.Error != nil {
		if errors.Is(errXMLFile.Error, gorm.ErrRecordNotFound) {
			fmt.Println("No record for this xml file, creating")
			xmlFile := XMLFileSQL{
				XmlName:   fileName,
				Status:    NotProcessed,
				FullPath:  p.currentZipFileSQL.FullPath + "::" + innerZipPath + fileName,
				ZipFileID: p.currentZipFileSQL.ID,
				Info:      "New XML Process Started",
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

// RegisterOrSkipExchangeLine returns Done if the Bulk file is already processed
// If the XML file entry does not exist,
// creates a new one (using the current Zip file as foreign key)
// or loads the existing bulk file information if the entry exists but is not done
func (p *StateHandler) RegisterOrSkipExchangeLine(exchangeID string, lineNumber int) (ProcessStatus, error) {
	//So a directory has started processing, but never finished, get the last known ZIP File
	//The Processor starts at the last known ZIP File, not at the specific Exchange Document
	//find the last unfinished zip file
	var exchangeLine ExchangeLineSQL
	errExchangeFile := p.db.Where("exchange_name = ?", exchangeID).First(&exchangeLine)
	if errExchangeFile.Error != nil {
		if errors.Is(errExchangeFile.Error, gorm.ErrRecordNotFound) {
			fmt.Println("No record for this xml file, creating")
			exchangeLineNew := ExchangeLineSQL{
				XMLFileID:    p.currentXMLFileSQL.ID,
				ExchangeName: exchangeID,
				Status:       NotProcessed,
				Info:         "new exchange line",
				FullPath:     p.currentXMLFileSQL.FullPath + "::" + exchangeID + " (line: " + strconv.Itoa(lineNumber) + ")",
			}
			p.db.Create(&exchangeLineNew)
			p.currentExchangeLineSQL = exchangeLineNew
			return NotProcessed, nil //new processing project
		} else {
			return NotProcessed, errExchangeFile.Error
		}
	}

	// won't be registered if already done
	if exchangeLine.Status != Done {
		p.currentExchangeLineSQL = exchangeLine
	}

	return exchangeLine.Status, nil
}

// MarkProcessingDirectoryAsFinished sets the status of directory as finished, no deleting downwards
func (p *StateHandler) MarkProcessingDirectoryAsFinished() {
	resultStatus := p.db.Model(&p.ProcessingDirSQL).Update("status", Done)

	if resultStatus.Error != nil {
		panic(resultStatus.Error)
	}

	resultInfo := p.db.Model(&p.ProcessingDirSQL).Update("info", "finished")
	if resultInfo.Error != nil {
		panic(resultInfo.Error)
	}
}

// MarkZipFileAsFinished if a finishes, delete all recorded XML lines
// And mark the Zip as finished, but always keep it
func (p *StateHandler) MarkZipFileAsFinished() {
	//Delete All Exchange Files belonging to the current XML File
	var resultXMLDelete *gorm.DB

	if p.SafeDeleteOnly {
		resultXMLDelete = p.db.Where("zip_file_id = ?", p.currentZipFileSQL.ID).Delete(&XMLFileSQL{})
	} else {
		resultXMLDelete = p.db.Unscoped().Where("zip_file_id = ?", p.currentZipFileSQL.ID).Delete(&XMLFileSQL{})
	}

	if resultXMLDelete.Error != nil {
		panic(resultXMLDelete.Error)
	}

	// Check deleted records
	fmt.Printf("Deleted %v record(s) with name '%s'", resultXMLDelete.RowsAffected, p.currentXMLFileSQL.XmlName)
	// set current Exchange File to empty
	p.currentXMLFileSQL = XMLFileSQL{}

	//We're keeping the Zip Entry for now
	resultInfo := p.db.Model(&p.currentZipFileSQL).Update("info", "finished")
	if resultInfo.Error != nil {
		panic(resultInfo.Error)
	}

	resultStatus := p.db.Model(&p.currentZipFileSQL).Update("status", Done)
	if resultStatus.Error != nil {
		panic(resultStatus.Error)
	}
}

// MarkXMLAsFinished if an XML finishes, delete all recorded exchange lines
// And mark the XML as finished, but keep it
// the xml gets deleted when the whole Zip is finished
func (p *StateHandler) MarkXMLAsFinished() {
	//Delete All Exchange Files belonging to the current XML File
	var resultExchangeDelete *gorm.DB

	if p.SafeDeleteOnly {
		resultExchangeDelete = p.db.Where("xml_file_id = ?", p.currentXMLFileSQL.ID).Delete(&ExchangeLineSQL{})
	} else {
		resultExchangeDelete = p.db.Unscoped().Where("xml_file_id = ?", p.currentXMLFileSQL.ID).Delete(&ExchangeLineSQL{})
	}

	if resultExchangeDelete.Error != nil {
		panic(resultExchangeDelete.Error)
	}

	// Check deleted records
	fmt.Printf("Deleted %v record(s) with ID '%v'", resultExchangeDelete.RowsAffected, p.currentXMLFileSQL.ID)
	// set current Exchange File to empty
	p.currentExchangeLineSQL = ExchangeLineSQL{}

	//We're keeping the XML Entry for now
	resultInfo := p.db.Model(&p.currentXMLFileSQL).Update("info", "finished")
	if resultInfo.Error != nil {
		panic(resultInfo.Error)
	}

	resultStatus := p.db.Model(&p.currentXMLFileSQL).Update("status", Done)
	if resultStatus.Error != nil {
		panic(resultStatus.Error)
	}
}

// MarkExchangeAsFinished exchange Records only get deleted when the XML is done
func (p *StateHandler) MarkExchangeAsFinished() {
	resultStatus := p.db.Model(&p.currentExchangeLineSQL).Update("status", Done)

	if resultStatus.Error != nil {
		panic(resultStatus.Error)
	}

	resultInfo := p.db.Model(&p.currentExchangeLineSQL).Update("info", "finished")
	if resultInfo.Error != nil {
		panic(resultInfo.Error)
	}
}

// IsDirectoryFinished returns true if the directory is finished
func (p *StateHandler) IsDirectoryFinished() bool {
	return p.ProcessingDirSQL.Status == Done
}
