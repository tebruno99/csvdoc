package csvdoc

import "errors"

var (
	// ErrStructTagNotInCSV struct tags must be found in the csv file.
	ErrStructTagNotInCSV = errors.New("struct tag not in csv")

	// ErrDuplicateHeaderInCSV csv had 2 columns with the same header.
	ErrDuplicateHeaderInCSV = errors.New("duplicate header in csv")

	// ErrStructTagDuplicate struct has 2 tags with the same csv name.
	ErrStructTagDuplicate = errors.New("struct tag duplicate")

	// ErrConverterNotFoundForType struct field type doesn't have a default converter (need to provide a customer column converter).
	ErrConverterNotFoundForType = errors.New("converter not found for type")

	// ErrNotFoundHeaderInCSV provided header name was not in the csv headers.
	ErrNotFoundHeaderInCSV = errors.New("header not found in csv")

	// ErrTypeOverflow conversion to the field's type resulted in overflow.
	ErrTypeOverflow = errors.New("conversion type overlow")

	// ErrToFewStructTags more headers were provided than struct tags available.
	ErrToFewStructTags = errors.New("to few struct tags")
)
