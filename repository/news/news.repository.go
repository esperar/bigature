package news

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
)

type NewsDto struct {
	Title   string
	Author  string
	Content string
}

type NewsRaw struct {
	Id      string
	Author  string
	Title   string
	Content string
}

type NewsRepository struct {
	DB *sql.DB
}

var Repository NewsRepository

func (r *NewsRepository) AssignDB(db *sql.DB) {
	r.DB = db
}

func (r *NewsRepository) InsertNews(n NewsDto) (sql.Result, error) {
	id, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	query := `
    INSERT INTO news
    (id, title, author, content)
    VALUES (?, ?, ?, ?)
  `

	result, err := r.DB.Exec(query, id.String(), n.Title, n.Author, n.Content)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *NewsRepository) GetAllNews() (*[]NewsRaw, error) {
	var raws []NewsRaw

	query := `SELECT * FROM news`
	rows, err := r.DB.Query(query)

	for rows.Next() {
		var raw NewsRaw
		rows.Scan(&raw.Id, &raw.Title, &raw.Author, &raw.Content)
		raws = append(raws, raw)
	}

	if err != nil {
		return nil, err
	} else {
		return &raws, nil
	}
}

func (r *NewsRepository) GetOneNews(id string) (*NewsRaw, error) {
	var raw NewsRaw

	query := `SELECT * FROM news WHERE id = ?`
	err := r.DB.QueryRow(query, id).Scan(&raw.Id, &raw.Title, &raw.Author, &raw.Content)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New("NOT FOUND")
		} else {
			return nil, err
		}
	} else {
		return &raw, nil
	}
}

func (r *NewsRepository) DeleteOneNews(id string) (sql.Result, error) {
	query := `DELETE FROM news WHERE id = ?`
	result, err := r.DB.Exec(query, id)

	if err != nil {
		return nil, err
	}

	affected, err := result.RowsAffected()

	if err != nil {
		return nil, err
	}

	if affected == 0 {
		return nil, errors.New("NOT FOUND")
	}

	return result, nil
}

func (r *NewsRepository) UpdateOneNews(id string, n NewsDto) (sql.Result, error) {
	query := `UPDATE news SET title = IFNULL(?, title), author = IFNULL(?, author), content = IFNULL(?, content) WHERE id = ?`
	var title, author, content *string

	if n.Title != "" {
		title = &n.Title
	}

	if n.Author != "" {
		author = &n.Author
	}

	if n.Content != "" {
		content = &n.Content
	}

	result, err := r.DB.Exec(query, title, author, content, id)

	if err != nil {
		return nil, err
	}

	affected, err := result.RowsAffected()

	if err != nil {
		return nil, err
	}

	if affected == 0 {
		return nil, errors.New("NOT FOUND")
	}

	return result, nil
}
