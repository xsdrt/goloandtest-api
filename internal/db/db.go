package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool" //3rd party lib...
)

type Item struct {
	Task   string
	Status string
}

type DB struct {
	pool *pgxpool.Pool
}

func New(user, password, dbname, host string, port int) (*DB, error) {
	//Define a connection string...
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, dbname)

	//Connect to a database using the connection string...
	pool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the db: %w", err)
	}

	//Pinged the db to make sure we have a connection...
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping the db: %w", err)
	}
	//Return the database with a connection pool available for it...
	return &DB{pool: pool}, nil
}

func (db *DB) InsertItem(ctx context.Context, item Item) error {
	query := `INSERT INTO todo_items (task, status) VALUES ($1, $2)` //Use parameters to avoid SqlInjection attacks...
	_, err := db.pool.Exec(ctx, query, item.Task, item.Status)
	return err
}

func (db *DB) GetAllItems(ctx context.Context) ([]Item, error) {
	query := `SELECT task, status FROM todo_items`
	rows, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.Task, &item.Status)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return items, nil

}

func (db *DB) Close() {
	db.pool.Close()
}
