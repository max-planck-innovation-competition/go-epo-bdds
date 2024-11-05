package state_handler

import (
	"errors"
	"gorm.io/gorm/logger"
	"log/slog"
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

// Initialize loads Last Known State
// Creates DB if there is none
// returns false if the processing is already finished
// returns true if there is some processing left to be done
func (sh *StateHandler) Initialize() {
	db, err := gorm.Open(sqlite.Open(sh.DatabasePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to open " + sh.DatabasePath)
	}
	sh.db = db
	// this will create the tables in the database, or migrate them if they already exist
	err = sh.db.AutoMigrate(
		&Object{},
	)
	if err != nil {
		slog.With("err", err).Error("could not migrate")
		return
	}
}
