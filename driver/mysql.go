package driver

type MysqlDriver struct {
}

func(d *MysqlDriver) WritePlaceholder(writer Writer, l int) {
	writer.WriteByte('?')
}

func (d *MysqlDriver) GetPlaceholder(l int) string {
	return "?"
}

func init() {
	mysqlDriver := &MysqlDriver{}
	MDriver["mysql"] = mysqlDriver
}
