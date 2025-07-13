package csvdoc

import (
	"database/sql"
	"errors"
	"reflect"
	"strconv"
	"time"
)

// ToStringConversion is a function type that is used to provide converters from Go types to csv string values.
// default Conversion functions are provided for Writer and custom Conversion functions may be added to a Writer to
// override the defaults of a specific column identified by the header name.
// The function is responsible for converting struct field value and outputting the appropriate string.
type ToStringConversion func(*reflect.Value) (string, error)

// Writer is an interface for converting Go types to csv string arrays and writing to io.Writer.
type Writer[T any] interface {
	Write(tm *T) error                                            // Writes a row by converting it to string and writing
	Close() error                                                 // Closes the io.Writer
	AddConverter(header string, handler ToStringConversion) error // AddConverter registers a custom conversion function for the specified header.
	RemoveConverter(header string) error                          // Removes a custom conversion function for the specified header.
}

// buildHeaderNameIndexCache creates two maps representing the csv header and the column number of the header.
// The first map links column names to their index positions, the second maps indices back to names.
// It validates that all headers exist in struct tags and checks for duplicate headers.
func buildWriteHeaderNameIndexCache(headerLine []string, tFieldsIndexes map[string]int) (map[string]int, map[int]string, error) {
	nameIndex := make(map[string]int, len(headerLine))
	indexName := make(map[int]string, len(headerLine))

	// Collect indexes for headers in struct tags.
	for i, col := range headerLine {
		if _, ok := nameIndex[col]; ok {
			return nil, nil, ErrDuplicateHeaderInCSV
		}
		if _, ok := tFieldsIndexes[col]; ok {
			nameIndex[col] = i
			indexName[i] = col
		} else {
			return nil, nil, ErrStructTagNotInCSV
		}
	}

	return nameIndex, indexName, nil
}

// buildWriteDefaultConverters produces a map[reflect.Type]Conversion for each default handled type. Writer implementations can use
// this map to aid in building Go types into csv string values.
func buildWriteDefaultConverters() map[reflect.Type]ToStringConversion {
	converts := make(map[reflect.Type]ToStringConversion)
	intConversion := ToStringConversion(func(v *reflect.Value) (string, error) {
		return strconv.FormatInt(v.Int(), 10), nil
	})
	uintConversion := ToStringConversion(func(v *reflect.Value) (string, error) {
		return strconv.FormatUint(v.Uint(), 10), nil
	})
	floatConversion := ToStringConversion(func(v *reflect.Value) (string, error) {
		return strconv.FormatFloat(v.Float(), 'f', -1, 64), nil
	})
	stringConversion := ToStringConversion(func(v *reflect.Value) (string, error) {
		return v.String(), nil
	})
	boolConversion := ToStringConversion(func(v *reflect.Value) (string, error) {
		return strconv.FormatBool(v.Bool()), nil
	})
	sqlNullStringConversion := ToStringConversion(func(v *reflect.Value) (string, error) {
		ns, ok := v.Interface().(sql.NullString)
		if !ok {
			return "", errors.New("cannot convert to sql.NullString")
		}
		if ns.Valid {
			return ns.String, nil
		}

		return "", nil
	})
	sqlNullInt64Conversion := ToStringConversion(func(v *reflect.Value) (string, error) {
		ns, ok := v.Interface().(sql.NullInt64)
		if !ok {
			return "", errors.New("cannot convert to sql.NullInt64")
		}
		if ns.Valid {
			return strconv.FormatInt(ns.Int64, 10), nil
		}

		return "", nil
	})
	sqlNullInt32Conversion := ToStringConversion(func(v *reflect.Value) (string, error) {
		ns, ok := v.Interface().(sql.NullInt32)
		if !ok {
			return "", errors.New("cannot convert to sql.NullInt32")
		}
		if ns.Valid {
			return strconv.FormatInt(int64(ns.Int32), 10), nil
		}

		return "", nil
	})
	sqlNullInt16Conversion := ToStringConversion(func(v *reflect.Value) (string, error) {
		ns, ok := v.Interface().(sql.NullInt16)
		if !ok {
			return "", errors.New("cannot convert to sql.NullInt16")
		}

		if ns.Valid {
			return strconv.FormatInt(int64(ns.Int16), 10), nil
		}

		return "", nil
	})
	sqlNullTimeConversion := ToStringConversion(func(v *reflect.Value) (string, error) {
		ns, ok := v.Interface().(sql.NullTime)
		if !ok {
			return "", errors.New("cannot convert to sql.NullTime")
		}
		if ns.Valid {
			return ns.Time.Format(time.DateTime), nil
		}
		return "", nil
	})
	timeConversion := ToStringConversion(func(v *reflect.Value) (string, error) {
		tm, ok := v.Interface().(time.Time)
		if !ok {
			return "", errors.New("cannot convert to time.Time{}")
		}
		return tm.Format(time.DateTime), nil
	})
	sqlNullFloat64Conversion := ToStringConversion(func(v *reflect.Value) (string, error) {
		ns, ok := v.Interface().(sql.NullFloat64)
		if !ok {
			return "", errors.New("cannot convert to sql.NullFloat64")
		}
		if ns.Valid {
			return strconv.FormatFloat(float64(ns.Float64), 'f', -1, 64), nil
		}

		return "", nil
	})
	sqlNullBoolConversion := ToStringConversion(func(v *reflect.Value) (string, error) {
		ns, ok := v.Interface().(sql.NullBool)
		if !ok {
			return "", errors.New("cannot convert to sql.NullBool")
		}

		if ns.Valid {
			return strconv.FormatBool(ns.Bool), nil
		}
		return "", nil
	})
	//
	converts[reflect.TypeOf(int64(1))] = intConversion
	converts[reflect.TypeOf(int32(1))] = intConversion
	converts[reflect.TypeOf(int16(1))] = intConversion
	converts[reflect.TypeOf(int(1))] = intConversion
	converts[reflect.TypeOf(uint64(1))] = uintConversion
	converts[reflect.TypeOf(uint32(1))] = uintConversion
	converts[reflect.TypeOf(uint16(1))] = uintConversion
	converts[reflect.TypeOf(uint(1))] = uintConversion
	converts[reflect.TypeOf("")] = stringConversion
	converts[reflect.TypeOf(true)] = boolConversion
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

	return converts
}
