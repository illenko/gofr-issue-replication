package main

import (
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/migration"
)

type transaction struct {
	Id    int
	Title string
}

func main() {
	app := gofr.New()

	app.Migrate(map[int64]migration.Migrate{
		20240426170000: {
			UP: func(d migration.Datasource) error {
				d.SQL.Exec(`create table transaction (id integer, title varchar);`)
				d.SQL.Exec(`insert into transaction (id, title) values (1, 'transaction 1');`)
				return nil
			},
		},
	})

	app.Subscribe("transaction-topic", func(c *gofr.Context) error {

		rows, _ := c.SQL.QueryContext(c, "SELECT id, title FROM transaction")

		for rows.Next() {
			var t transaction
			rows.Scan(&t.Id, &t.Title)

			c.Logger.Info("Read transaction", t)
		}

		defer rows.Close()

		return nil
	})

	app.GET("/transactions", func(c *gofr.Context) (interface{}, error) {
		var transactions []transaction

		rows, _ := c.SQL.QueryContext(c, "SELECT id, title FROM transaction")

		for rows.Next() {
			var t transaction
			rows.Scan(&t.Id, &t.Title)

			c.Logger.Info("Read transaction", t)

			transactions = append(transactions, t)
		}

		defer rows.Close()

		return transactions, nil
	})
	app.Run()
}
