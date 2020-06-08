package gsql

type Logger interface {
	Printf(format string, v ...interface{})
}
