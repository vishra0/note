package mysql

import (
	"database/sql"
	"vis/note/pkg/models"
)

type Snippetmodel struct {
	DB *sql.DB
}

func (m *Snippetmodel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO Snippet(title, content, created ,expires)
	values (?,?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL  ? DAY))`
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, nil
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return int(id), nil
}

func (m *Snippetmodel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT id,title,content,created ,expires FROM snippet
	where expires > UTC_TIMESTAMP() AND id = ? `
	row := m.DB.QueryRow(stmt, id)
	s := &models.Snippet{}
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return s, nil
}

func (m *Snippetmodel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT id,title,content,created,expires FROM snippet
	where expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	snippets := []*models.Snippet{}
	for rows.Next() {
		s := &models.Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
