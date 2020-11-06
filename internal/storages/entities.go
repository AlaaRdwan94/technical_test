package storages

// Task reflects tasks in DB
type Task struct {
	ID          string `json:"id"`
	Content     string `json:"content"`
	UserID      string `json:"user_id"`
	CreatedDate string `json:"created_date"`
}

// User reflects users data from DB
type User struct {
	ID       string
	Password string
	Max      int
}

type TaskUpdate struct {
	ID         string
	TaskID     string `json:"task_id"`
	TaskStatus string `json:"status"`
	TaskDue    string `json:"due_date"`
	AssignerID string `json:"assigner_id"`
	AssigneeID string `json:"assignee_id"`
}
