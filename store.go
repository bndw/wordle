package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gliderlabs/ssh"
	_ "github.com/mattn/go-sqlite3"
)

func newRepo(dbFile string) (*sqliteRepo, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	r := sqliteRepo{
		dbFile: dbFile,
		DB:     db,
	}

	if err := r.createSchema(); err != nil {
		return nil, err
	}

	return &r, nil
}

type sqliteRepo struct {
	dbFile string
	DB     *sql.DB
}

func (r *sqliteRepo) SaveGame(ctx context.Context, userID string, game *Game) error {
	const insert = `insert into game(user, data) values(?, ?)`

	stmt, err := r.DB.Prepare(insert)
	if err != nil {
		return err
	}

	data, err := json.Marshal(game)
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(userID, data); err != nil {
		return err
	}

	return nil
}

func (r *sqliteRepo) ListGames(ctx context.Context, user string) ([]Game, error) {
	var (
		query string
		rows  *sql.Rows
		err   error
	)
	if user != "" {
		query = `SELECT id, data FROM game WHERE user=? ORDER BY id DESC;`
		rows, err = r.DB.Query(query, user)
	} else {
		query = `SELECT id, data FROM game ORDER BY id DESC;`
		rows, err = r.DB.Query(query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	games := make([]Game, 0)
	for rows.Next() {
		var (
			id   int64
			data []byte
		)

		rows.Scan(&id, &data)

		var game Game
		if err := json.Unmarshal(data, &game); err != nil {
			return nil, fmt.Errorf("failed to decode game")
		}

		games = append(games, game)
	}

	return games, nil
}

// Close closes the SQLite database connection.
func (r *sqliteRepo) Close() error {
	return r.DB.Close()
}

// createSchema initializes the SQLite database schema.
func (r *sqliteRepo) createSchema() error {
	const schema = `
	CREATE TABLE IF NOT EXISTS game(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
		user TEXT NOT NULL,
		data BLOB NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_game_user ON game(user);`

	if _, err := r.DB.Exec(schema); err != nil {
		return err
	}

	return nil
}

func userKey(s ssh.Session) string {
	parts := strings.Split(s.RemoteAddr().String(), ":")
	ip := parts[0]
	return fmt.Sprintf("%s|%s", s.User(), ip)
}
