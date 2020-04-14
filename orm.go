package goorm

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// Qb : QueryBuilder object
type Qb struct {
	t          string
	err        error
	joins      string
	slct       []string
	db         *sql.DB
	conditions string
	params     []interface{}
}

// NewQuery creates a new instance of the query builder
func NewQuery(t string, db *sql.DB) *Qb {
	return &Qb{
		t:  t,
		db: db,
	}
}

// Select columns
func (qb *Qb) Select(fields ...string) *Qb {
	qb.slct = fields

	return qb
}

func (qb Qb) String() string {
	return fmt.Sprintf("table: %s", qb.t)
}

// Where condition
func (qb *Qb) Where(key string, operator string, val interface{}) *Qb {
	if qb.conditions != "" {
		qb.conditions = fmt.Sprintf("%s AND %s %s ?", qb.conditions, key, operator)
	} else {
		qb.conditions = fmt.Sprintf("WHERE %s %s ?", key, operator)
	}

	qb.params = append(qb.params, val)

	return qb
}

// OrWhere condition
func (qb *Qb) OrWhere(key string, operator string, val interface{}) *Qb {
	if qb.conditions != "" {
		qb.conditions = fmt.Sprintf("%s OR %s %s ?", qb.conditions, key, operator)
	} else {
		qb.conditions = fmt.Sprintf("WHERE %s %s ?", key, operator)
	}

	qb.params = append(qb.params, val)

	return qb
}

// Get the query you have just built...
func (qb *Qb) Get() (*sql.Rows, error) {
	var slct string
	if len(qb.slct) > 0 {
		slct = strings.Join(qb.slct, ",")
	} else {
		slct = "*"
	}

	q := fmt.Sprintf("SELECT %s FROM %s %s %s", slct, qb.t, qb.joins, qb.conditions)
	return qb.run(q)
}

func (qb *Qb) run(q string) (*sql.Rows, error) {
	return qb.db.Query(q, qb.params...)
}

// Insert a new record
func (qb *Qb) Insert(keys []string, values []interface{}) (*sql.Rows, error) {
	if len(values) != len(keys) && len(keys) == 0 {
		return nil, errors.New("keys and values don't match")
	}

	k := fmt.Sprintf("%s", strings.Join(keys, ","))
	for i := 0; i < len(keys); i++ {
		keys[i] = "?"
	}
	v := fmt.Sprintf("%s", strings.Join(keys, ","))

	q := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", qb.t, k, v)
	qb.params = values

	return qb.run(q)
}

// Update an exists record
func (qb *Qb) Update(data map[string]interface{}) (*sql.Rows, error) {
	var values string

	params := []interface{}{}
	for k, v := range data {
		values = fmt.Sprintf("%s%s = ?,", values, k)
		params = append(params, v)
	}

	qb.params = append(params, qb.params...)
	values = strings.Trim(values, ",")
	q := fmt.Sprintf("UPDATE %s SET %s %s", qb.t, values, qb.conditions)

	return qb.run(q)
}

// Delete an exists record
func (qb *Qb) Delete() (*sql.Rows, error) {
	q := fmt.Sprintf("DELETE FROM %s %s", qb.t, qb.conditions)

	return qb.run(q)
}

// Join two tables together
func (qb *Qb) Join(dir, tbl, local, operator, foreign string) *Qb {
	qb.joins = fmt.Sprintf("%s %s JOIN %s ON %s %s %s ", qb.joins, dir, tbl, local, operator, foreign)

	return qb
}
