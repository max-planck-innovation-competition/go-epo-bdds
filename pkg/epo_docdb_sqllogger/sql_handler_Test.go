package epo_docdb_sqllogger

import (
	"fmt"
	"testing"
)

func TestSQLHandlesr(t *testing.T) {
	sqllogger := NewSqlLogger("log.db", "./", "C:\\docdb")
	fmt.Println(sqllogger)
}
