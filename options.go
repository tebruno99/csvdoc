package csvdoc

const (
	defaultEnableCLRF  = false
	defaultWriteHeader = true
	defaultEscapeRune  = ','
)

type Option[T WriterOption] func(*T)

type WriterOption struct {
	crlfEnable   bool
	escapeRune   rune
	outputHeader []string
	writeHeader  bool
}

func DefaultWriterOption() *WriterOption {
	return &WriterOption{
		crlfEnable:   defaultEnableCLRF,
		writeHeader:  defaultWriteHeader,
		escapeRune:   defaultEscapeRune,
		outputHeader: nil,
	}
}

func WithEscapeRune[T WriterOption](escapeRune rune) Option[T] {
	return func(o *T) {
		switch x := any(o).(type) {
		case *WriterOption:
			x.escapeRune = escapeRune
		}
	}
}

func WithUseCRLF[T WriterOption](crlfEnable bool) Option[T] {
	return func(o *T) {
		switch x := any(o).(type) {
		case *WriterOption:
			x.crlfEnable = crlfEnable
		}
	}
}

func WithFormatHeaders[T WriterOption](outputHeader []string) Option[T] {
	return func(o *T) {
		switch x := any(o).(type) {
		case *WriterOption:
			cpy := make([]string, len(outputHeader))
			copy(cpy, outputHeader)
			x.outputHeader = cpy
		}
	}
}

func WithWriteHeader[T WriterOption](enable bool) Option[T] {
	return func(o *T) {
		switch x := any(o).(type) {
		case *WriterOption:
			x.writeHeader = enable
		}
	}
}
