### Overview
This is a simple backend for a good old todo service, right now this service can handle login/list/create simple tasks.  
To make it run:
- `go run main.go`
- Import Postman collection from `docs` to check example
 
### Requirements
Design database to serve addition requirements:
- Tasks should have status (open/doing/resolved/overdue).
- Tasks should have an assigner, assignee.
- Tasks should have a due date.

Right now a user can add as many tasks as they want, we want the ability to limit N task per day.
Example: users are limited to create only 5 tasks only per day, if the daily limit is reached, return 4xx code to the client and ignore the create request.

Write API to create a dashboard that will show the top of assigners, assignees, resolvers by time range, or the current time.
Feel free to add more features to the dashboard.

#### Backend requirements
- A nice README on how to run, what is missing, what else you want to improve but don't have enough time
- Fork this repo and show us your development progess by a PR.
- Write integration tests for this project
- Make this code DRY
- Write unit test for `services` layer
- Write migration to migrate data from SQLite to MySQL and alter tables.
- This project include many issues from code to DB strucutre, feel free to optimize them.
#### Optional requirements
- Write unit test for `storages` layer
- Split `services` layer to `use case` and `transport` layer


### DB Schema
```sql
-- users definition

CREATE TABLE users (
	id TEXT NOT NULL,
	password TEXT NOT NULL,
	max_todo INTEGER DEFAULT 5 NOT NULL,
	CONSTRAINT users_PK PRIMARY KEY (id)
);

INSERT INTO users (id, password, max_todo) VALUES('firstUser', 'example', 5);

-- tasks definition

CREATE TABLE tasks (
	id TEXT NOT NULL,
	content TEXT NOT NULL,
	user_id TEXT NOT NULL,
    created_date TEXT NOT NULL,
	CONSTRAINT tasks_PK PRIMARY KEY (id),
	CONSTRAINT tasks_FK FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### Sequence diagram
