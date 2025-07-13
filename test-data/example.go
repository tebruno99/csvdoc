// Package testdata is for holding common types used in examples and testing
package testdata

import (
	"database/sql"
	"time"
)

// Example is a struct that defines each column of the test-data/example.csv.
type Example struct {
	BirthDate time.Time     `csv:"birthDate"`
	MonYear   sql.NullTime  `csv:"MonYear"`
	SystemID  string        `csv:"systemId"`
	UserID    string        `csv:"userId"`
	Gender    string        `csv:"gender"`
	Maximum   string        `csv:"Maximum"`
	GovID     sql.NullInt64 `csv:"govId"`
	ID        int64         `csv:"Id"`
	Year      uint          `csv:"year"`
	Minium    float64       `csv:"Minimum"`
	DoProcess bool          `csv:"DoProcess"`
	Validated sql.NullBool  `csv:"Validated"`
}

// ExampleMods is a struct that defines a subset of columns from the test-data/example.csv.
type ExampleMods struct {
	BirthDate time.Time `csv:"birthDate,birthDateTime"`
	SystemID  string    `csv:"systemId,systemId"`
	UserID    string    `csv:"userId,userId"`
	Gender    string    `csv:"gender,gender"`
	Errors    string    `csv:"-,errors"`
}
