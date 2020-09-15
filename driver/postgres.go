package driver

import "strconv"


type PostgresDialector struct {
}

func(d *PostgresDialector) WritePlaceholder(writer Writer, l int) {
	writer.WriteByte('$')
	writer.WriteString(strconv.Itoa(l))
}

func (d *PostgresDialector) GetPlaceholder(l int) string {
	return "$" + strconv.Itoa(l)
}

var postgresDialector *PostgresDialector

func init() {
	postgresDialector = &PostgresDialector{}
	MDriver["postgres"] = postgresDialector
}
