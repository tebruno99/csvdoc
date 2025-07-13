package csvdoc

import (
	"encoding/csv"
	"log"
	"os"
	"reflect"
	"sync"
)

// FileWriter is a generic CSV document writer that maps struct fields to csv headers using reflect.
// It provides functionality to write CSV files line by line, converting each Go struct into an array of csv strings.
// The writer supports overriding converters for specific columns and provides default converters for standard types.
type FileWriter[T any] struct {
	reflectIndexes      map[string]int
	headerIndex         map[string]int
	indexHeader         map[int]string
	defaultConverters   map[reflect.Type]ToStringConversion
	customConverters    map[string]ToStringConversion
	f                   *os.File
	cw                  *csv.Writer
	fp                  string
	headers             []string
	hasWrittenHeaderMux sync.RWMutex
	hasWrittenHeaders   bool
}

// NewFileWriter creates a new CSV FileWriter for the specified file path. If a sort array is not provided, it is assumed
// the header names will come from the struct csv output tags and order will be random.
func NewFileWriter[T any](fp string, sort []string) (*FileWriter[T], error) {
	//nolint:gosec // The purpose of this library is to open user provided files.
	f, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	if err != nil {
		return nil, err
	}

	csvWriter := csv.NewWriter(f)

	reflectIndexes, err := buildReflectTagIndexCache[T](true)
	if err != nil {
		return nil, err
	}
	if sort != nil && len(sort) > len(reflectIndexes) {
		return nil, ErrToFewStructTags
	}

	if sort == nil {
		sort = make([]string, len(reflectIndexes))
		for i, v := range reflectIndexes {
			sort[v] = i
		}
	}

	nameIndex, indexName, err := buildWriteHeaderNameIndexCache(sort, reflectIndexes)
	if err != nil {
		cerr := f.Close()
		if cerr != nil {
			log.Println("Error closeing file: ", cerr)
		}
		return nil, err
	}

	doc := &FileWriter[T]{
		reflectIndexes:    reflectIndexes,
		headerIndex:       nameIndex,
		indexHeader:       indexName,
		defaultConverters: buildWriteDefaultConverters(),
		headers:           sort,
		customConverters:  nil,
		f:                 f,
		cw:                csvWriter,
		fp:                fp,
	}

	return doc, nil
}

// Write converts a Go struct of type T to a []string and writes to the configured io.Writer.
func (doc *FileWriter[T]) Write(tm *T) error {
	var err error
	doc.hasWrittenHeaderMux.Lock()
	if !doc.hasWrittenHeaders {
		err = doc.cw.Write(doc.headers)
		if err != nil {
			return err
		}
		doc.hasWrittenHeaders = true
	}
	doc.hasWrittenHeaderMux.Unlock()
	row := make([]string, len(doc.headers))

	elemVal := reflect.ValueOf(tm).Elem()
	for fieldName, fieldIndex := range doc.reflectIndexes {
		// first get the output position
		if _, ok := doc.headerIndex[fieldName]; !ok {
			continue
		}
		outIndex := doc.headerIndex[fieldName]
		f := elemVal.Field(fieldIndex)
		var columnString string
		if fnc, ok := doc.customConverters[fieldName]; ok {
			columnString, err = fnc(&f)
			if err != nil {
				return err
			}
		} else {
			tp := f.Type()
			columnString, err = doc.defaultConverters[tp](&f)
			if err != nil {
				return err
			}
		}
		row[outIndex] = columnString
	}

	return doc.cw.Write(row)
}

// Close will close the file writer's file and flush the contents.
func (doc *FileWriter[T]) Close() error {
	doc.cw.Flush()
	return doc.f.Close()
}

// AddConverter adds a custom converter function to a specific header/column.
func (doc *FileWriter[T]) AddConverter(header string, handler ToStringConversion) error {
	if doc.customConverters == nil {
		doc.customConverters = make(map[string]ToStringConversion, 1)
	}
	doc.customConverters[header] = handler

	return nil
}

// RemoveConverter removes custom converter for specific header.
func (doc *FileWriter[T]) RemoveConverter(header string) error {
	delete(doc.customConverters, header)
	return nil
}
