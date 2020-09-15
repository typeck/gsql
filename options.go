package gsql

import "github.com/typeck/gsql/driver"

type Option interface {
	apply(*gsql)
}

type optionFunc func(*gsql)

func (f optionFunc) apply(db *gsql) {
	f(db)
}

func SetLogger(logger Logger) Option {
	return optionFunc(func (db *gsql) {
		db.logger = logger
	})
}

func SetDriver(driver driver.Driver) Option {
	return optionFunc(func(db *gsql) {
		db.Driver = driver
		//panic("don't support placeholder")
	})
}

func SetTag(name string) Option {
	return optionFunc(func(db *gsql) {
		db.orm.Tag = name
	})
}