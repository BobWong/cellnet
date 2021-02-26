package mysql

import (
	"database/sql"
	"errors"
	"github.com/bobwong89757/cellnet/log"
)

type Wrapper struct {
	drv   *sql.DB
	row   *sql.Rows
	query string

	Err error
}

var ErrDriverNotReady = errors.New("driver not ready")

func (self *Wrapper) Query(query string, args ...interface{}) *Wrapper {

	if self.drv == nil {
		self.Err = ErrDriverNotReady
		return self
	}

	self.query = query
	log.GetLog().Debug("[DB]", query, args)

	self.row, self.Err = self.drv.Query(query, args...)

	if self.Err != nil {
		log.GetLog().Error("[DB] ", self.query, self.Err.Error())
	}

	return self
}

func (self *Wrapper) Execute(query string, args ...interface{}) *Wrapper {
	if self.drv == nil {
		self.Err = ErrDriverNotReady
		return self
	}

	self.query = query
	log.GetLog().Debug("[DB]", query, args)

	_, self.Err = self.drv.Exec(query, args...)

	if self.Err != nil {
		log.GetLog().Error("[DB] ", self.query, self.Err.Error())
	}

	return self
}

func (self *Wrapper) One(data ...interface{}) *Wrapper {

	if self.Err != nil {
		return self
	}

	if self.drv == nil {
		self.Err = ErrDriverNotReady
		return self
	}

	if !self.row.Next() {
		return self
	}

	self.Err = self.row.Scan(data...)

	if self.Err != nil {
		log.GetLog().Error("One.Row.Scan failed", self.query, self.Err)
	}

	self.row.Close()
	self.row = nil

	return self
}

func (self *Wrapper) Scan(dest ...interface{}) {

	self.Err = self.row.Scan(dest...)

	if self.Err != nil {
		log.GetLog().Error("Scan.Scan failed", self.query, self.Err)
	}

}

func (self *Wrapper) Each(callback func(*Wrapper) bool) *Wrapper {

	if self.Err != nil {
		return self
	}

	if self.drv == nil {
		self.Err = ErrDriverNotReady
		return self
	}

	for self.row.Next() {

		if !callback(self) {
			break
		}

		if self.Err != nil {
			return self
		}

	}

	self.row.Close()

	return self
}

func NewWrapper(drv *sql.DB) *Wrapper {

	return &Wrapper{
		drv: drv,
	}
}
