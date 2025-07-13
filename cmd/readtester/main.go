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

	td "github.com/tebruno99/csvdoc/test-data"

	"github.com/tebruno99/csvdoc"
)

func main() {
	startTime := time.Now()
	ct := 1

	cd, err := csvdoc.NewFileReader[td.Example]("test-data/example.csv")
	if err != nil {
		log.Fatal(err)
	}

	err = cd.AddConverter("MonYear", csvdoc.Conversion(func(s string, field *reflect.Value) error {
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
