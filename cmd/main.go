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

	cd.AddConvertor("birthDate", csvdoc.Conversion(func(s string, field *reflect.Value) error {
		formats := []string{"1/2/2006 15:04:05 AM"}
		if s != "" {
			for _, format := range formats {
				val, err := time.Parse(format, s)
				if err == nil {
					field.Set(reflect.ValueOf(val))
					return nil
				}
			}
		}
		return errors.New("cannot convert string to time")
	}))

	for {
		m, err := cd.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatalf("Error reading CSV file: %v", err)
			}
		}

		fmt.Printf("%#v\n", m)
		ct += 1
	}

	seconds := time.Since(startTime).Seconds()
	fmt.Printf("Processed %d total rows in %.2fs: %.2f/s\n", ct, seconds, float64(ct)/seconds)
	fmt.Printf("Done\n")
}

type Example struct {
	BirthDate time.Time     `csv:"birthDate"`
	SystemId  string        `csv:"systemId"`
	UserId    string        `csv:"userId"`
	Gender    string        `csv:"gender"`
	Maximum   string        `csv:"Maximum"`
	GovId     sql.NullInt64 `csv:"govId"`
	Id        int64         `csv:"Id"`
	Year      uint          `csv:"year"`
	Minium    float64       `csv:"Minimum"`
}
