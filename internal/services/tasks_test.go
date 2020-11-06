package services

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go/togo/internal/storages"
	sqllite "github.com/go/togo/internal/storages/sqlite"
	"gopkg.in/go-playground/assert.v1"
)

func TestDashboard(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "task_id", "status", "assigner_id", "assignee_id", "due_date"}).
		AddRow("dcbd2592-f469-4433-b569-05638104b514", "ed5e80eb-383e-461a-b491-e86f23d170ad", "resolved", "secondUser", "thirdUser", "3 days")

	mock.ExpectQuery("^SELECT id , task_id , status , assigner_id , assignee_id , due_date FROM tasksdata WHERE task_id IN (SELECT id FROM tasks)*").
		WithArgs("2020-11-03").
		WillReturnRows(rows)

	ctx := context.TODO()
	var l sqllite.LiteDB
	l.DB = db
	tasks, err := l.RetrieveDashboard(
		ctx,
		sql.NullString{String: "2020-11-03", Valid: true},
	)
	if err != nil {
		t.Errorf("this is the error getting the data: %v\n", err)
		return
	}
	assert.NotEqual(t, tasks, nil)
}

func TestValidateUser(t *testing.T) {

	user_id := "firstUser"
	password := "example"

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(user_id)

	mock.ExpectQuery("^SSELECT id FROM users*").
		WithArgs(user_id, password).
		WillReturnRows(rows)

	ctx := context.TODO()
	var l sqllite.LiteDB
	l.DB = db
	valid := l.ValidateUser(ctx, sql.NullString{String: user_id, Valid: true}, sql.NullString{String: password, Valid: true})
	assert.Equal(t, valid, true)
}

func TestGetMax(t *testing.T) {
	user_id := "secondUser"
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow("4")

	mock.ExpectQuery("^SELECT max_todo FROM users*").
		WithArgs(user_id).
		WillReturnRows(rows)

	ctx := context.TODO()
	var l sqllite.LiteDB
	l.DB = db
	max := l.GetMax(ctx, sql.NullString{String: user_id, Valid: true})

	assert.Equal(t, max, 4)
}

func TestAddTaskDate(t *testing.T) {
	id := "dcbd2592-f469-4444-b569-05638104b514"
	task_id := "ed5e80eb-383e-461a-b491-e86f23d170ad"
	status := "resolved"
	due_date := "3 days"
	assigner_id := "secondUser"
	assignee_id := "thirdUser"
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("^INSERT INTO tasksdata*").
		WithArgs(id, task_id, status, due_date, assigner_id, assignee_id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.TODO()
	var l sqllite.LiteDB
	l.DB = db

	err = l.AddTaskDate(ctx, &storages.TaskUpdate{
		ID:         id,
		TaskID:     task_id,
		TaskStatus: status,
		TaskDue:    due_date,
		AssignerID: assigner_id,
		AssigneeID: assignee_id,
	})
	if err != nil {
		t.Errorf("this is the error getting the data: %v\n", err)
		return
	}
	assert.Equal(t, err, nil)
}

func TestAddTask(t *testing.T) {
	id := "dcbd2592-f469-4444-b569-05638104b514"
	content := "content 1"
	created_date := "2020-11-4"
	user_id := "secondUser"
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("^INSERT INTO tasks*").
		WithArgs(id, content, user_id, created_date).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.TODO()
	var l sqllite.LiteDB
	l.DB = db

	err = l.AddTask(ctx, &storages.Task{
		ID:          id,
		Content:     content,
		CreatedDate: created_date,
		UserID:      user_id,
	})
	if err != nil {
		t.Errorf("this is the error getting the data: %v\n", err)
		return
	}
	assert.Equal(t, err, nil)
}

func TestRetrieveTasksData(t *testing.T) {
	task_id := "ed5e80eb-383e-461a-b491-e86f23d170ad"
	status := "resolved"
	due_date := "3 days"
	assigner_id := "secondUser"
	assignee_id := "thirdUser"

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a database connection", err)
	}
	rows := sqlmock.NewRows([]string{"status", "due_date", "assigner_id", "assignee_id"}).
		AddRow(status, due_date, assigner_id, assignee_id)
	defer db.Close()
	mock.ExpectQuery(`SELECT status, due_date, assigner_id, assignee_id FROM tasksdata WHERE*`).
		WithArgs(task_id).
		WillReturnRows(rows)

	ctx := context.TODO()
	var l sqllite.LiteDB
	l.DB = db

	update, err := l.RetrieveTasksData(ctx, sql.NullString{
		String: task_id,
		Valid:  true,
	})
	if err != nil {
		t.Errorf("this is the error getting the data: %v\n", err)
		return
	}
	assert.NotEqual(t, update, nil)
}

func TestRetrieveTasks(t *testing.T) {

	user_id := "secondUser"
	created_date := "2020-11-03"

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a database connection", err)
	}
	rows := sqlmock.NewRows([]string{"id", "content", "user_id", "created_date"}).
		AddRow("b9490d6d-4599-4f17-b455-5ddfce2a6c73", "third content", "secondUser", "2020-11-03")
	defer db.Close()
	mock.ExpectQuery(`SELECT id, content, user_id, created_date FROM tasks*`).
		WithArgs(user_id, created_date).
		WillReturnRows(rows)

	ctx := context.TODO()
	var l sqllite.LiteDB
	l.DB = db

	update, err := l.RetrieveTasks(ctx, sql.NullString{
		String: user_id,
		Valid:  true,
	}, sql.NullString{
		String: created_date,
		Valid:  true,
	})
	if err != nil {
		t.Errorf("this is the error getting the data: %v\n", err)
		return
	}
	assert.NotEqual(t, update, nil)
}
