CREATE TABLE users (
	id TEXT NOT NULL,
	password TEXT NOT NULL,
	max_todo INTEGER DEFAULT 5 NOT NULL,
	CONSTRAINT users_PK PRIMARY KEY (id)
);

INSERT INTO users (id, password, max_todo) VALUES('firstUser', 'example', 5);

CREATE TABLE tasks (
	id TEXT NOT NULL,
	content TEXT NOT NULL,
	user_id TEXT NOT NULL,
    created_date TEXT NOT NULL,
	CONSTRAINT tasks_PK PRIMARY KEY (id),
	CONSTRAINT tasks_FK FOREIGN KEY (user_id) REFERENCES users(id)
);

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

INSERT INTO users (id, password, max_todo) VALUES('secondUser', 'example', 4);
INSERT INTO users (id, password, max_todo) VALUES('thirdUser', 'example', 4);