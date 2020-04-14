package orm

import (
	"database/sql"
	"fmt"
	"strings"
)

// Qb stands for query builder
type Qb struct {
	t      string
	q      string
	params []interface{}
	err    error
	slct   []string
}

// NewQuery creates a new instance of the query builder
func NewQuery(t string) *Qb {
	return &Qb{
		t:   t,
		err: nil,
	}
}

// Select ...
func (qb *Qb) Select(fields ...string) *Qb {
	qb.slct = fields

	return qb
}

func (qb Qb) String() string {
	return fmt.Sprintf("SELECT %s FROM %s %s", strings.Join(qb.slct, ","), qb.t, qb.q)
}

// Where ...
func (qb *Qb) Where(key string, operator string, val interface{}) *Qb {
	if qb.q != "" {
		qb.q = fmt.Sprintf("%s AND %s %s ?", qb.q, key, operator)
	} else {
		qb.q = fmt.Sprintf("WHERE %s %s ?", key, operator)
	}

	qb.params = append(qb.params, val)

	return qb
}

// OrWhere ...
func (qb *Qb) OrWhere(key string, operator string, val interface{}) *Qb {
	if qb.q != "" {
		qb.q = fmt.Sprintf("%s OR %s %s ?", qb.q, key, operator)
	} else {
		qb.q = fmt.Sprintf("WHERE %s %s ?", key, operator)
	}

	qb.params = append(qb.params, val)

	return qb
}

func (qb *Qb) build() (string, []interface{}) {
	var slct string
	if len(qb.slct) > 0 {
		slct = strings.Join(qb.slct, ",")
	} else {
		slct = "*"
	}
	return fmt.Sprintf("SELECT %s FROM %s %s", slct, qb.t, qb.q), qb.params
}

// Get ...
func (qb *Qb) Get(db *sql.DB) (*sql.Rows, error) {
	q, args := qb.build()
	return db.Query(q, args...)
}
