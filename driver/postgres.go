package driver

import "strconv"


type PostgresDriver struct {
}

func(d *PostgresDriver) WritePlaceholder(writer Writer, l int) {
	writer.WriteByte('$')
	writer.WriteString(strconv.Itoa(l))
}

func (d *PostgresDriver) GetPlaceholder(l int) string {
	return "$" + strconv.Itoa(l)
}


func init() {
	postgresDriver := &PostgresDriver{}
	MDriver["postgres"] = postgresDriver
}
