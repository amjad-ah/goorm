# GoORM

## What is GoORM?

A simple ORM I made it for learning purpose.

---

## How to use it?

```go
package main

import (
	"database/sql"

	"github.com/amjad-ah/goorm"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	dbConfig()
	qb := goorm.NewQuery("users", db)

	qb = qb.Select("users.id", "users.name", "users.email", "users.age", "other_table.other_data").
		Where("users.id", "=", 1).
		Where("users.age", ">", 26).
		Where("users.name", "=", "Amjad").
		OrWhere("other_table.other_data", "!=", "something").
		Join("LEFT", "other_table", "users.other_table_id", "=", "other_table.id")

	// Read:
	rows, err := qb.Get()

	// Update:
	rows, err := qb.Update(map[string]interface{}{
		"name": "Amjad",
	})

	// Delete:
	_, err := qb.Update(map[string]interface{}{
		"name": "Amjad",
	})

	// Insert:
	rows, err := goorm.NewQuery("users", db).Insert([]string{
		"name",
		"email",
		"age",
	}, []interface{}{
		"Amjad",
		"amjad@example.com",
		26,
	})

	if err != nil {
		panic(err)
	}

	// Do something with your data in `rows` variable
}

func dbConfig() {
	var err error
	db, err = sql.Open("mysql", "{USERNAME}:{PASSWORD}@tcp({DB_HOST}:{DB_PORT})/{DB_NAME}")

	if err != nil {
		panic(err.Error())
	}
}
```

# THAT'S IT!
