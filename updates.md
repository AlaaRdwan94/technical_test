### DB Schema
```sql
-- taskdatas definition 
CREATE TABLE tasksdata (
	id TEXT NOT NULL,
	task_id TEXT NOT NULL,
	status TEXT DEFAULT 'open' NOT NULL,
	due_date TEXT NOT NULL,
	assigner_id TEXT NOT NULL,
	assignee_id TEXT NOT NULL,
	CONSTRAINT tasksdata_PK PRIMARY KEY (id),
	CONSTRAINT tasksdata_FK FOREIGN KEY (task_id) REFERENCES tasks(id)
);

-- insert another user
INSERT INTO users (id, password, max_todo) VALUES('secondUser', 'example', 4);
INSERT INTO users (id, password, max_todo) VALUES('thirdUser', 'example', 4);
```
-------------------------------
## code 
-- add TaskUpdate struct 
-- add fields 
    --------------------------------------
    name        |  type     |  json name 
    --------------------------------------
    TaskID      |  string  |`json:"task_id"`
    TaskDue     |  string  |`json:"due_date"` 
    AssignerID  |  string  |`json:"assigner_id"` 
    AssigneeID  |  string  |`json:"assignee_id"` 
    TaskStatus  |  string  |`json:"status"` 

-- add AddTaskDate function in internal/storages/sqlite/db.go
   insert into `tasksdata` table 

-- update addTask function in internal/services/tasks.go
   to call AddTaskDate and return new tsk struct 

-- add Tsk struct in internal/services/tasks.go
-- add fields 
    --------------------------------------
    name        |  type     |  json name 
    --------------------------------------
    ID          |  string  |`json:"id"`
    Content     |  string  |`json:"content"` 
    UserID      |  string  |`json:"user_id"` 
    CreatedDate |  string  |`json:"created_date"` 
    TaskID      |  string  |`json:"task_id"` 
    TaskDue     |  string  |`json:"due_date"` 
    AssignerID  |  string  |`json:"assigner_id"` 
    AssigneeID  |  string  |`json:"assignee_id"` 
    TaskStatus  |  string  |`json:"status"` 

--add function RetrieveTasksData in internal/storages/sqlite/db.go
to get task data by task id 

--udate function listTasks in internal/services/tasks.go
to call RetrieveTasksData and return full data of tasks

--add function GetMax in internal/storages/sqlite/db.go
to get mux taskes of user per day 

--update user struct 
  -- add field `Max`  of type  `int` to get data into it from users table 

--update function addTask in internal/services/tasks.go
  -- call gettasks and getmax to prevent users from addng more tasks per day according to his max  

--create function RetrieveDashboard in internal/storages/sqlite/db.go
  -- get all data of tasks with assigners and assignees and the status of the task if open or resolved or any thing else 

--create new API rout /dashboard with mehod GET
  -- get data with all assigners , assignees , status of all tasks and the duration for every task 

-------------------------------
## Test 
--create tasks_test.go file in internal/services/tasks_test.go
  -- create the unit test for all functions in test.go
-------------------------------
## Database migration
-- create directory migrations 
  -- this will contains the `.sql` files used for databse migrations

-- create 1_tables.up.sql in migrations
  -- this will contain database schema
 
-- create 1_tables.down.sql in migrations
  -- this will contain undo sql commands to downgrade the version of the database

-- create migration_test.go in internal/storages/mysql/migration_test.go
  -- to test migration from sqlite to mysql database 

------------------------------
## Whish to do 
-- complete integration test
-- docker-ze the application 
-- create repository between APIs and database model 
-- uses interfaces instead call function directly 

------------------------------
## To run 
-- import postman `togo.postman_collectionv2.json` from `docs` directory  
-- run `go run main.go` 
