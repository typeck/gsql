package driver

type Driver interface {
	WritePlaceholder(writer Writer, l int)
	GetPlaceholder(l int) string
}

type Writer interface {
	WriteString(s string) (int, error)
	WriteByte(c byte) error
}

var MDriver = make(map[string]Driver)