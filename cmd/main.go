// Package main is the main executable of this example application.
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"reflect"
	"time"

	"github.com/tebruno99/csvdoc"
)

func main() {
	startTime := time.Now()
	ct := 1

	cd, err := csvdoc.NewFileReader[Example]("test-data/example.csv")
	if err != nil {
		log.Fatal(err)
	}

	err = cd.AddConvertor("MonYear", csvdoc.Conversion(func(s string, field *reflect.Value) error {
		formats := []string{"1/2006"}
		if s != "" {
			for _, format := range formats {
				val, serr := time.Parse(format, s)
				if serr == nil {
					field.Set(reflect.ValueOf(sql.NullTime{Time: val, Valid: true}))
					return nil
				}
			}
			return errors.New("cannot convert string to time")
		}
		return nil
	}))
	if err != nil {
		log.Fatalf("Failed to install converter for MonYear: %s", err)
	}

	for {
		m, rerr := cd.Read()
		if rerr != nil {
			if errors.Is(rerr, io.EOF) {
				break
			}
			log.Fatalf("Error reading CSV file: %v", rerr)
		}

		fmt.Printf("%#v\n", m)
		ct++
	}

	seconds := time.Since(startTime).Seconds()
	fmt.Printf("Processed %d total rows in %.2fs: %.2f/s\n", ct, seconds, float64(ct)/seconds)
	fmt.Printf("Done\n")
}

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
}

// ExampleMods is a struct that defines a subset of columns from the test-data/example.csv.
type ExampleMods struct {
	BirthDate time.Time `csv:"birthDate,birthDateTime"`
	SystemID  string    `csv:"systemId,systemId"`
	UserID    string    `csv:"userId,userId"`
	Gender    string    `csv:"gender,gender"`
	Errors    string    `csv:"-,errors"`
}
