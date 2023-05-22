package server

import (
	"context"
	"fmt"
	"strconv"

	api "github.com/studio1767/studio-api/api/v1"
	"github.com/studio1767/studio-api/internal/auth"
)

func (svr *studioServer) CreateProject(ctx context.Context, preq *api.ProjectRequest) (*api.Project, error) {
	fmt.Printf("CreateProject: %s %s\n", preq.Name, preq.Code)

	if err := auth.Authorize(ctx, "/", auth.CREATE); err != nil {
		return nil, err
	}

	result, err := svr.dbClient.Exec("INSERT INTO project (name, code) VALUES (?, ?)", preq.Name, preq.Code)
	if err != nil {
		return nil, fmt.Errorf("create project failed: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("create project id failed: %w", err)
	}

	project := &api.Project{
		Id:   strconv.FormatInt(id, 10),
		Name: preq.Name,
		Code: preq.Code,
	}

	return project, nil
}

func (svr *studioServer) Projects(filter *api.ProjectFilter, stream api.Studio_ProjectsServer) error {
	fmt.Println("Projects")

	ctx := stream.Context()
	if err := auth.Authorize(ctx, "/", auth.READ); err != nil {
		return err
	}

	rows, err := svr.dbClient.Query("SELECT * FROM project")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var project api.Project
		if err = rows.Scan(&project.Id, &project.Name, &project.Code); err != nil {
			return err
		}
		stream.Send(&project)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}
