package csvdoc

import (
	"database/sql"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Conversion is a function type that is used to provide converters from csv string values to a specific Go type.
// default Conversion functions are provided for Reader and custom Conversion functions may be added to a Reader to
// override the defaults of a specific column identified by the header name.
// The function is responsible for parsing the string and setting the appropriate value in the reflect.Value.
type Conversion func(string, *reflect.Value) error

// Reader is an interface for reading and converting CSV data into Go types.
type Reader[T any] interface {
	Read() (*T, error)                                    // Reads a row from the csv and converts it to type *T
	Close() error                                         // Closes the io.Reader
	Reset() error                                         // Resets the io.Reader back to the beginning of the file.
	AddConverter(header string, handler Conversion) error // AddConverter registers a custom conversion function for the specified header.
	RemoveConverter(header string) error                  // Removes a custom conversion function for the specified header.
}

// buildReadDefaultConverters produces a map[reflect.Type]Conversion for each default handled type. Reader implementations can use
// this map to aid in building default csv string values into Go types.
func buildReadDefaultConverters() map[reflect.Type]Conversion {
	converts := make(map[reflect.Type]Conversion)
	intConversion := Conversion(func(s string, field *reflect.Value) error {
		if s != "" {
			val, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return err
			}
			if field.OverflowInt(val) {
				return errors.New("overflow convert")
			}
			field.SetInt(val)
			return nil
		}

		return errors.New("cannot convert empty string to int*")
	})
	sqlNullInt64Conversion := Conversion(func(a string, field *reflect.Value) error {
		if a != "" {
			val, err := strconv.ParseInt(a, 10, 64)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(sql.NullInt64{Int64: val, Valid: true}))
		}

		return nil
	})
	sqlNullInt32Conversion := Conversion(func(a string, field *reflect.Value) error {
		if a != "" {
			val, err := strconv.ParseInt(a, 10, 32)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(sql.NullInt32{Int32: int32(val), Valid: true}))
		}

		return nil
	})
	sqlNullInt16Conversion := Conversion(func(a string, field *reflect.Value) error {
		if a != "" {
			val, err := strconv.ParseInt(a, 10, 16)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(sql.NullInt16{Int16: int16(val), Valid: true}))
		}

		return nil
	})
	uintConversion := Conversion(func(s string, field *reflect.Value) error {
		if s != "" {
			val, err := strconv.ParseUint(s, 10, 64)
			if err != nil {
				return err
			}
			if field.OverflowUint(val) {
				return ErrTypeOverflow
			}
			field.SetUint(val)
			return nil
		}
		return errors.New("cannot convert empty string to uint*")
	})
	floatConversion := Conversion(func(s string, field *reflect.Value) error {
		if s != "" {
			val, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return err
			}
			if field.OverflowFloat(val) {
				return ErrTypeOverflow
			}
			field.SetFloat(val)
			return nil
		}

		return errors.New("cannot convert empty string to float*")
	})
	sqlNullFloat64Conversion := Conversion(func(s string, field *reflect.Value) error {
		if s != "" {
			val, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(sql.NullFloat64{Float64: val, Valid: true}))
		}

		return nil
	})
	stringConversion := Conversion(func(s string, field *reflect.Value) error {
		field.SetString(s)
		return nil
	})
	timeConversion := Conversion(func(s string, field *reflect.Value) error {
		formats := []string{time.DateTime, time.DateOnly, time.RFC3339, time.RFC3339Nano, "01/02/2006 15:04:05 PM", "1/2/2006 15:04:05 PM"}
		if s != "" {
			for _, format := range formats {
				val, err := time.Parse(format, s)
				if err == nil {
					field.Set(reflect.ValueOf(val))
					return nil
				}
			}
			return errors.New("cannot convert string to time")
		}
		return errors.New("cannot convert empty string to time")
	})
	boolConversion := Conversion(func(a string, field *reflect.Value) error {
		if a != "" {
			cmp := strings.ToLower(a)
			if cmp == "true" || cmp == "1" || cmp == "on" || cmp == "yes" || cmp == "y" {
				field.SetBool(true)
			} else {
				field.SetBool(false)
			}
			return nil
		}

		return errors.New("cannot convert empty string to bool")
	})
	sqlNullTimeConversion := Conversion(func(s string, field *reflect.Value) error {
		formats := []string{time.DateTime, time.DateOnly, time.RFC3339, time.RFC3339Nano, "01/02/2006 15:04:05 PM", "1/2/2006 15:04:05 PM"}
		if s != "" {
			for _, format := range formats {
				val, err := time.Parse(format, s)
				if err == nil {
					field.Set(reflect.ValueOf(sql.NullTime{Time: val, Valid: true}))
					return nil
				}
			}
			return errors.New("cannot convert string to time")
		}

		return nil
	})
	sqlNullStringConversion := Conversion(func(a string, field *reflect.Value) error {
		if a != "" {
			field.Set(reflect.ValueOf(sql.NullString{String: a, Valid: true}))
		}

		return nil
	})
	sqlNullBoolConversion := Conversion(func(a string, field *reflect.Value) error {
		if a != "" {
			cmp := strings.ToLower(a)
			if cmp == "true" || cmp == "1" || cmp == "on" || cmp == "yes" || cmp == "y" {
				field.Set(reflect.ValueOf(sql.NullBool{Bool: true, Valid: true}))
			} else {
				field.Set(reflect.ValueOf(sql.NullBool{Bool: false, Valid: true}))
			}
		}
		return nil
	})

	converts[reflect.TypeOf(int64(1))] = intConversion
	converts[reflect.TypeOf(int32(1))] = intConversion
	converts[reflect.TypeOf(int16(1))] = intConversion
	converts[reflect.TypeOf(int(1))] = intConversion
	converts[reflect.TypeOf(uint64(1))] = uintConversion
	converts[reflect.TypeOf(uint32(1))] = uintConversion
	converts[reflect.TypeOf(uint16(1))] = uintConversion
	converts[reflect.TypeOf(uint(1))] = uintConversion
	converts[reflect.TypeOf("")] = stringConversion
	converts[reflect.TypeOf(sql.NullString{})] = sqlNullStringConversion
	converts[reflect.TypeOf(sql.NullInt64{})] = sqlNullInt64Conversion
	converts[reflect.TypeOf(sql.NullInt32{})] = sqlNullInt32Conversion
	converts[reflect.TypeOf(sql.NullInt16{})] = sqlNullInt16Conversion
	converts[reflect.TypeOf(sql.NullTime{})] = sqlNullTimeConversion
	converts[reflect.TypeOf(time.Time{})] = timeConversion
	converts[reflect.TypeOf(float64(0))] = floatConversion
	converts[reflect.TypeOf(float32(0))] = floatConversion
	converts[reflect.TypeOf(sql.NullFloat64{})] = sqlNullFloat64Conversion
	converts[reflect.TypeOf(sql.NullBool{})] = sqlNullBoolConversion
	converts[reflect.TypeOf(true)] = boolConversion

	return converts
}

// buildHeaderNameIndexCache creates two maps representing the csv header and the column number of the header.
// The first map links column names to their index positions, the second maps indices back to names.
// It validates that all headers exist in struct tags and checks for duplicate headers.
func buildReadHeaderNameIndexCache(headerLine []string, tFieldsIndexes map[string]int) (map[string]int, map[int]string, error) {
	nameIndex := make(map[string]int, len(tFieldsIndexes))
	indexName := make(map[int]string, len(tFieldsIndexes))

	// Collect indexes for headers in struct tags.
	for i, col := range headerLine {
		if _, ok := nameIndex[col]; ok {
			return nil, nil, ErrDuplicateHeaderInCSV
		}
		if _, ok := tFieldsIndexes[col]; ok {
			nameIndex[col] = i
			indexName[i] = col
		}
	}

	// check that all struct tags were in the header.
	for k := range tFieldsIndexes {
		if _, ok := nameIndex[k]; !ok {
			return nil, nil, ErrStructTagNotInCSV
		}
	}

	return nameIndex, indexName, nil
}
