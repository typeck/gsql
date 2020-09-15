package driver

type MysqlDialector struct {
}

func(d *MysqlDialector) WritePlaceholder(writer Writer, l int) {
	writer.WriteByte('?')
}

func (d *MysqlDialector) GetPlaceholder(l int) string {
	return "?"
}

var mysqlDialector *PostgresDialector

func init() {
	mysqlDialector = &PostgresDialector{}
	MDriver["mysql"] = mysqlDialector
}
