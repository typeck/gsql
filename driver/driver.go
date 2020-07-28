package driver

type Dialector interface {
	WritePlaceholder(writer Writer, l int)
	GetPlaceholder(l int) string
}

type Writer interface {
	WriteString(s string) (int, error)
	WriteByte(c byte) error
}

var MDialector = make(map[string]Dialector)