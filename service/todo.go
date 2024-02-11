package service

import (
	"context"
	"database/sql"
	"strings"
	"fmt"

	"github.com/sungyo4869/go-basic/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	var todo model.TODO

	result, err := s.db.ExecContext(ctx, insert, subject, description)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	
	row := s.db.QueryRowContext(ctx, confirm, id)

	err = row.Scan(
		&todo.Subject,
		&todo.Description,
		&todo.CreatedAt,
		&todo.UpdatedAt,	
	)
	
	todo.ID = id

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	todos:= []*model.TODO{}

	var rows *sql.Rows
	var err error

	if prevID == 0 {
		rows, err = s.db.QueryContext(ctx, read, size)
	}else{
		rows, err = s.db.QueryContext(ctx, readWithID, prevID, size)
	}
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		todo := model.TODO{}
	
		if err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}

	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	var todo model.TODO

	result, err := s.db.ExecContext(ctx, update, subject, description, id)
	if err != nil {
		return nil, err
	}

	sum, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if sum == 0 {
		return nil, &model.ErrNotFound{}
	}
	row := s.db.QueryRowContext(ctx, confirm, id)

	err = row.Scan(
		&todo.Subject,
		&todo.Description,
		&todo.CreatedAt,
		&todo.UpdatedAt,	
	)
	
	todo.ID = id

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	var idSlice []interface{}

	if len(ids) == 0 {
		return nil
	}

	placeholders := ""
	if len(ids) > 1 {
		placeholders = strings.Repeat(", ?", len(ids)-1)
	}
	delete := fmt.Sprintf(deleteFmt, placeholders)

    for _, id := range ids {
        idSlice = append(idSlice, id)
    }

	result, err := s.db.ExecContext(ctx, delete, idSlice...)
	if err != nil {
		return err
	}

	sum, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if sum == 0 {
		return &model.ErrNotFound{}
	}

	return nil
}
