package service

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/parlaynu/studio1767-api/api/graph/model"
)

func (s *service) CreateProject(ctx context.Context, np *model.NewProject) (*model.Project, error) {

	result, err := s.db.Exec("INSERT INTO project (name, code) VALUES (?, ?)", np.Name, np.Code)
	if err != nil {
		return nil, fmt.Errorf("create project failed: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("create project id failed: %w", err)
	}

	p := model.Project{
		ID:   strconv.FormatInt(id, 10),
		Name: np.Name,
		Code: np.Code,
	}

	return &p, nil
}

func (s *service) Projects(ctx context.Context) ([]*model.Project, error) {

	ps := make([]*model.Project, 0, 5)

	rows, err := s.db.Query("SELECT * FROM project")
	if err != nil {
		return nil, fmt.Errorf("query project table failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var proj model.Project

		if err = rows.Scan(&proj.ID, &proj.Name, &proj.Code); err != nil {
			return nil, fmt.Errorf("scanning project row failed: %w", err)
		}

		ps = append(ps, &proj)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("looping over rows failed: %w", err)
	}

	return ps, nil
}

func (s *service) ProjectById(ctx context.Context, id string) (*model.Project, error) {
	var project model.Project

	row := s.db.QueryRow("SELECT * FROM project WHERE id = ?", id)
	if err := row.Scan(&project.ID, &project.Name, &project.Code); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project by id %s: no such project", id)
		}
		return nil, fmt.Errorf("project by id %s: %w", id, err)
	}
	return &project, nil
}
