### CsvDoc

Personal project used in tooling for several csv processing tasks such as converting CSVs with hundreds millions of lines into SQL tables.

Using struct tags paired with column headers auto converts csv columns into Go types. The custom csv tag can be split read and output headers

Example:

```
type MyCsv struct {
	BirthDateTime `csv:"birthDate:birthDateTime"` // birthdate is the read column name, birthDateTime is the write column name
	IncrementalId `csv:"-,id"` // during read this field is ignored. During write the value is written in the id column
	FirstName `csv:"firstName" // Read and Write both use this column
	MiddleName `csv:"-"` // Ignored by read and write
```


### License
see LICENSE file.
