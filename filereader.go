package csvdoc

import (
	"encoding/csv"
	"log"
	"os"
	"reflect"
)

// FileReader is a generic CSV document reader that maps csv headers to struct fields using reflect.
// It provides functionality to read CSV files line by line, converting each line into a struct of type T.
// The reader supports overriding with custom converters for specific columns and provides default converters for standard types.
type FileReader[T any] struct {
	reflectIndexes    map[string]int
	headerIndex       map[string]int
	defaultConverters map[reflect.Type]Conversion
	customConverters  map[string]Conversion
	indexHeader       map[int]string
	f                 *os.File
	cr                *csv.Reader
	fp                string
}

// Close the underlaying file. Close should be called when done reading.
func (fr *FileReader[T]) Close() error {
	return fr.f.Close()
}

// Reset resets the csv reader back to the row after the header (2nd row).
func (fr *FileReader[T]) Reset() error {
	_, err := fr.f.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = fr.cr.Read()
	if err != nil {
		return err
	}
	return nil
}

// NewFileReader creates a new CSV FileReader for the specified file path. This reader assumes the csv file has a header
// and all the header values match the struct tags of type T.
func NewFileReader[T any](fp string) (*FileReader[T], error) {
	fieldIndexes, err := buildReflectTagIndexCache[T]()
	if err != nil {
		return nil, err
	}

	//nolint:gosec // The purpose of this library is to open user provided files.
	f, err := os.Open(fp)
	if err != nil {
		return nil, err
	}

	cr := csv.NewReader(f)
	headerLine, err := cr.Read()
	if err != nil {
		cerr := f.Close()
		if cerr != nil {
			log.Println("Error closeing file: ", cerr)
		}
		return nil, err
	}

	nameIndex, indexName, err := buildHeaderNameIndexCache(headerLine, fieldIndexes)
	if err != nil {
		cerr := f.Close()
		if cerr != nil {
			log.Println("Error closeing file: ", cerr)
		}
		return nil, err
	}

	fileReader := &FileReader[T]{
		fp:                fp,
		f:                 f,
		reflectIndexes:    fieldIndexes,
		cr:                cr,
		headerIndex:       nameIndex,
		indexHeader:       indexName,
		defaultConverters: buildDefaultConverts(),
		customConverters:  make(map[string]Conversion),
	}

	return fileReader, nil
}

// Read uses csv.Reader to obtain the next line as a []string and then builds a struct of type *T from the []string. Returns EOF and closes the open file automatically.
func (fr *FileReader[T]) Read() (*T, error) {
	line, err := fr.cr.Read()
	if err != nil {
		cerr := fr.f.Close()
		if cerr != nil {
			log.Println("Error closeing file: ", cerr)
		}
		return nil, err
	}
	t := new(T)
	elemVal := reflect.ValueOf(t).Elem()
	for i, v := range line {
		if _, ok := fr.indexHeader[i]; !ok {
			continue
		}
		hrName := fr.indexHeader[i]
		tagFieldIndex := fr.reflectIndexes[hrName]
		f := elemVal.Field(tagFieldIndex)
		tp := f.Type()

		if cv, ok := fr.customConverters[hrName]; ok {
			err = cv(v, &f)
			if err != nil {
				return nil, err
			}

			continue
		}

		if cv, ok := fr.defaultConverters[tp]; ok {
			err = cv(v, &f)
			if err != nil {
				return nil, err
			}

			continue
		}

		return nil, ErrConverterNotFoundForType
	}

	return t, nil
}

// AddConvertor adds a customer Conversion func to handle a specific CSV header/struct tag.
func (fr *FileReader[T]) AddConvertor(header string, handler Conversion) error {
	if _, ok := fr.headerIndex[header]; !ok {
		return ErrNotFoundHeaderInCSV
	}

	fr.customConverters[header] = handler
	return nil
}

// RemoveConvertor removes a customer Conversion func for a specific CSV header/struct tag.
func (fr *FileReader[T]) RemoveConvertor(header string) error {
	delete(fr.customConverters, header)
	return nil
}
