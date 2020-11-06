package sqllite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go/togo/internal/storages"
)

// LiteDB for working with sqllite
type LiteDB struct {
	DB *sql.DB
}

// RetrieveTasks returns tasks if match userID AND createDate.
func (l *LiteDB) RetrieveTasks(ctx context.Context, userID, createdDate sql.NullString) ([]*storages.Task, error) {

	stmt := `SELECT id, content, user_id, created_date FROM tasks WHERE user_id = ? AND created_date = ?`
	rows, err := l.DB.QueryContext(ctx, stmt, userID, createdDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*storages.Task
	for rows.Next() {
		t := &storages.Task{}
		err := rows.Scan(&t.ID, &t.Content, &t.UserID, &t.CreatedDate)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

//gt dashboard
// RetrieveTasks returns tasks if match userID AND createDate.
func (l *LiteDB) RetrieveDashboard(ctx context.Context, createdDate sql.NullString) ([]*storages.TaskUpdate, error) {

	stmt := `SELECT id , task_id , status , assigner_id , assignee_id , due_date FROM tasksdata WHERE task_id IN (SELECT id FROM tasks WHERE created_date = ? )`
	fmt.Println("stmt", stmt)
	rows, err := l.DB.QueryContext(ctx, stmt, createdDate)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []*storages.TaskUpdate
	for rows.Next() {
		t := &storages.TaskUpdate{}
		err := rows.Scan(&t.ID, &t.TaskID, &t.TaskStatus, &t.AssignerID, &t.AssigneeID, &t.TaskDue)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// RetrieveTasks returns tasks if match userID AND createDate.
func (l *LiteDB) RetrieveTasksData(ctx context.Context, taskID sql.NullString) (*storages.TaskUpdate, error) {

	stmt := `SELECT status, due_date, assigner_id, assignee_id FROM tasksdata WHERE task_id = ? `
	row := l.DB.QueryRowContext(ctx, stmt, taskID)
	task := &storages.TaskUpdate{}
	err := row.Scan(&task.TaskStatus, &task.TaskDue, &task.AssignerID, &task.AssigneeID)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// AddTask adds a new task to DB
func (l *LiteDB) AddTask(ctx context.Context, t *storages.Task) error {
	stmt := `INSERT INTO tasks (id, content, user_id, created_date) VALUES (?, ?, ?, ?)`
	_, err := l.DB.ExecContext(ctx, stmt, &t.ID, &t.Content, &t.UserID, &t.CreatedDate)
	if err != nil {
		return err
	}

	return nil
}

// AddTask adds a new task to DB
func (l *LiteDB) AddTaskDate(ctx context.Context, t *storages.TaskUpdate) error {
	stmt := `INSERT INTO tasksdata (id, task_id, status, due_date , assigner_id ,assignee_id) VALUES (?, ?, ?, ? , ? , ?)`
	_, err := l.DB.ExecContext(ctx, stmt, &t.ID, &t.TaskID, &t.TaskStatus, &t.TaskDue, &t.AssignerID, &t.AssigneeID)
	if err != nil {
		return err
	}

	return nil
}

// ValidateUser returns true if match userID AND password
func (l *LiteDB) ValidateUser(ctx context.Context, userID, pwd sql.NullString) bool {
	stmt := `SELECT id FROM users WHERE id = ? AND password = ?`
	row := l.DB.QueryRowContext(ctx, stmt, userID, pwd)
	u := &storages.User{}
	err := row.Scan(&u.ID)
	if err != nil {
		return false
	}

	return true
}

// ValidateUser returns true if match userID AND password
func (l *LiteDB) GetMax(ctx context.Context, userID sql.NullString) int {
	stmt := `SELECT max_todo FROM users WHERE id = ?`
	row := l.DB.QueryRowContext(ctx, stmt, userID)
	u := &storages.User{}
	err := row.Scan(&u.Max)
	if err != nil {
		return -1
	}

	return u.Max
}
