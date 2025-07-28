package infra

import (
	"database/sql"
	"github.com/cristianorosa/eSimulate/backend/domain"
)

// ThemeRepositoryPG implementa ThemeRepository usando PostgreSQL

type ThemeRepositoryPG struct {
	DB *sql.DB
}

func (r *ThemeRepositoryPG) Create(theme *domain.Theme) error {
	query := `INSERT INTO themes (name, parent_id) VALUES ($1, $2) RETURNING id`
	return r.DB.QueryRow(query, theme.Name, theme.ParentID).Scan(&theme.ID)
}

func (r *ThemeRepositoryPG) Update(theme *domain.Theme) error {
	query := `UPDATE themes SET name = $1, parent_id = $2 WHERE id = $3`
	_, err := r.DB.Exec(query, theme.Name, theme.ParentID, theme.ID)
	return err
}

func (r *ThemeRepositoryPG) Delete(id int) error {
	query := `DELETE FROM themes WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *ThemeRepositoryPG) FindByID(id int) (*domain.Theme, error) {
	query := `SELECT id, name, parent_id FROM themes WHERE id = $1`
	row := r.DB.QueryRow(query, id)
	var t domain.Theme
	var parentID sql.NullInt64
	if err := row.Scan(&t.ID, &t.Name, &parentID); err != nil {
		return nil, err
	}
	if parentID.Valid {
		pid := int(parentID.Int64)
		t.ParentID = &pid
	}
	return &t, nil
}

func (r *ThemeRepositoryPG) ListAll() ([]*domain.Theme, error) {
	query := `SELECT id, name, parent_id FROM themes`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var themes []*domain.Theme
	for rows.Next() {
		var t domain.Theme
		var parentID sql.NullInt64
		if err := rows.Scan(&t.ID, &t.Name, &parentID); err != nil {
			return nil, err
		}
		if parentID.Valid {
			pid := int(parentID.Int64)
			t.ParentID = &pid
		} else {
			t.ParentID = nil
		}
		themes = append(themes, &t)
	}
	return themes, nil
}
