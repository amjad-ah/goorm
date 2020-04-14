package orm

import (
	"database/sql"
	"errors"
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
	db     *sql.DB
}

// NewQuery creates a new instance of the query builder
func NewQuery(t string, db *sql.DB) *Qb {
	return &Qb{
		t:  t,
		db: db,
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

// Get ...
func (qb *Qb) Get() (*sql.Rows, error) {
	var slct string
	if len(qb.slct) > 0 {
		slct = strings.Join(qb.slct, ",")
	} else {
		slct = "*"
	}

	qb.q = fmt.Sprintf("SELECT %s FROM %s %s", slct, qb.t, qb.q)
	return qb.run()
}

func (qb *Qb) run() (*sql.Rows, error) {
	return qb.db.Query(qb.q, qb.params...)
}

// Insert ...
func (qb *Qb) Insert(keys []string, values []interface{}) (*sql.Rows, error) {
	if len(values) != len(keys) && len(keys) == 0 {
		return nil, errors.New("keys and values don't match")
	}

	k := fmt.Sprintf("%s", strings.Join(keys, ","))
	for i := 0; i < len(keys); i++ {
		keys[i] = "?"
	}
	v := fmt.Sprintf("%s", strings.Join(keys, ","))

	qb.q = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", qb.t, k, v)
	qb.params = values

	return qb.run()
}

// Update ..
func (qb *Qb) Update(data map[string]interface{}) (*sql.Rows, error) {
	var values string
	for k, v := range data {
		values = fmt.Sprintf("%s%s = '%s',", values, k, v)
	}
	values = strings.Trim(values, ",")
	qb.q = fmt.Sprintf("UPDATE %s SET %s %s", qb.t, values, qb.q)

	return qb.run()
}
