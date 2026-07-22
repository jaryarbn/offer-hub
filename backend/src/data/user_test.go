package data

import (
	"errors"
	"fmt"
	"testing"

	mysqlDriver "github.com/go-sql-driver/mysql"
)

func TestIsDuplicateKeyError(t *testing.T) {
	duplicate := &mysqlDriver.MySQLError{Number: 1062, Message: "duplicate entry"}
	if !isDuplicateKeyError(duplicate) {
		t.Fatal("MySQL error 1062 must be recognized as duplicate key")
	}
	if !isDuplicateKeyError(fmt.Errorf("wrapped: %w", duplicate)) {
		t.Fatal("wrapped MySQL error 1062 must be recognized as duplicate key")
	}
	if isDuplicateKeyError(&mysqlDriver.MySQLError{Number: 1048}) {
		t.Fatal("non-duplicate MySQL error was recognized as duplicate key")
	}
	if isDuplicateKeyError(errors.New("ordinary error")) {
		t.Fatal("ordinary error was recognized as duplicate key")
	}
}

func TestUserRecordTableName(t *testing.T) {
	if got := (UserRecord{}).TableName(); got != "t_user_info" {
		t.Fatalf("UserRecord.TableName() = %q, want t_user_info", got)
	}
}

func TestCreateUserRecordTableName(t *testing.T) {
	if got := (CreateUserRecord{}).TableName(); got != "t_user_info" {
		t.Fatalf("CreateUserRecord.TableName() = %q, want t_user_info", got)
	}
}

func TestUserInfoRecordTableName(t *testing.T) {
	if got := (UserInfoRecord{}).TableName(); got != "t_user_info" {
		t.Fatalf("UserInfoRecord.TableName() = %q, want t_user_info", got)
	}
}
