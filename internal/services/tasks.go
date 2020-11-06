package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go/togo/internal/storages"
	sqllite "github.com/go/togo/internal/storages/sqlite"
	"github.com/google/uuid"
)

// ToDoService implement HTTP server
type ToDoService struct {
	JWTKey string
	Store  *sqllite.LiteDB
}

// task for API request
type Tsk struct {
	ID          string `json:"id"`
	Content     string `json:"content"`
	UserID      string `json:"user_id"`
	CreatedDate string `json:"created_date"`
	TaskID      string `json:"task_id"`
	TaskDue     string `json:"due_date"`
	AssignerID  string `json:"assigner_id"`
	AssigneeID  string `json:"assignee_id"`
	TaskStatus  string `json:"status"`
}

func (s *ToDoService) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, req.URL.Path)
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Allow-Headers", "*")
	resp.Header().Set("Access-Control-Allow-Methods", "*")

	if req.Method == http.MethodOptions {
		resp.WriteHeader(http.StatusOK)
		return
	}

	switch req.URL.Path {
	case "/login":
		s.getAuthToken(resp, req)
		return
	case "/tasks":
		var ok bool
		req, ok = s.validToken(req)
		if !ok {
			resp.WriteHeader(http.StatusUnauthorized)
			return
		}

		switch req.Method {
		case http.MethodGet:
			s.listTasks(resp, req)
		case http.MethodPost:
			s.addTask(resp, req)
		}
		return
	case "/dashboard":
		var ok bool
		req, ok = s.validToken(req)
		if !ok {
			resp.WriteHeader(http.StatusUnauthorized)
			return
		}
		if req.Method == http.MethodGet {
			s.Dashboard(resp, req)
		}
		return
	}

}

//Dashboard
func (s *ToDoService) Dashboard(resp http.ResponseWriter, req *http.Request) {

	tasks, err := s.Store.RetrieveDashboard(
		req.Context(),
		value(req, "created_date"),
	)
	resp.Header().Set("Content-Type", "application/json")

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	json.NewEncoder(resp).Encode(map[string][]*storages.TaskUpdate{
		"data": tasks,
	})
}

func (s *ToDoService) getAuthToken(resp http.ResponseWriter, req *http.Request) {
	id := value(req, "user_id")
	if !s.Store.ValidateUser(req.Context(), id, value(req, "password")) {
		resp.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(resp).Encode(map[string]string{
			"error": "incorrect user_id/pwd",
		})
		return
	}
	resp.Header().Set("Content-Type", "application/json")

	token, err := s.createToken(id.String)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	json.NewEncoder(resp).Encode(map[string]string{
		"data": token,
	})
}

func (s *ToDoService) listTasks(resp http.ResponseWriter, req *http.Request) {
	id, _ := userIDFromCtx(req.Context())
	tasks, err := s.Store.RetrieveTasks(
		req.Context(),
		sql.NullString{
			String: id,
			Valid:  true,
		},
		value(req, "created_date"),
	)

	resp.Header().Set("Content-Type", "application/json")

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	arrtasks := []*Tsk{}
	if tasks != nil {
		for _, task := range tasks {
			arrobj := &Tsk{}
			arrobj.ID = task.ID
			arrobj.Content = task.Content
			arrobj.CreatedDate = task.CreatedDate
			arrobj.UserID = task.UserID
			taskid := arrobj.ID
			taskdata, err := s.Store.RetrieveTasksData(req.Context(), sql.NullString{
				String: taskid,
				Valid:  true,
			})
			if err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(resp).Encode(map[string]string{
					"error": err.Error(),
				})
				return
			}
			arrobj.TaskID = taskdata.TaskID
			arrobj.TaskStatus = taskdata.TaskStatus
			arrobj.TaskDue = taskdata.TaskDue
			arrobj.AssignerID = taskdata.AssignerID
			arrobj.AssigneeID = taskdata.AssigneeID
			arrtasks = append(arrtasks, arrobj)
		}
	}

	json.NewEncoder(resp).Encode(map[string][]*Tsk{
		"data": arrtasks,
	})
}
func (s *ToDoService) addTask(resp http.ResponseWriter, req *http.Request) {
	t := &Tsk{}
	err := json.NewDecoder(req.Body).Decode(t)
	defer req.Body.Close()
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	now := time.Now()
	userID, _ := userIDFromCtx(req.Context())
	//get users tasks
	tasks, err := s.Store.RetrieveTasks(
		req.Context(),
		sql.NullString{
			String: userID,
			Valid:  true,
		},
		sql.NullString{
			String: now.Format("2006-01-02"),
			Valid:  true,
		},
	)
	//get max number of tasks from database
	maxtimes := s.Store.GetMax(req.Context(), sql.NullString{
		String: userID,
		Valid:  true,
	})

	if len(tasks) > maxtimes {
		json.NewEncoder(resp).Encode(map[string]string{
			"message": "you reach max tasks ber day ",
		})
		return
	}

	t.UserID = userID
	t.AssignerID = userID
	t.CreatedDate = now.Format("2006-01-02")
	Task := &storages.Task{}
	Task.ID = uuid.New().String()
	Task.Content = t.Content
	Task.CreatedDate = t.CreatedDate
	Task.UserID = userID
	resp.Header().Set("Content-Type", "application/json")

	err = s.Store.AddTask(req.Context(), Task)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	TaskData := &storages.TaskUpdate{}
	TaskData.ID = uuid.New().String()
	TaskData.TaskID = Task.ID
	TaskData.TaskDue = t.TaskDue
	TaskData.AssigneeID = t.AssigneeID
	TaskData.AssignerID = t.AssignerID
	TaskData.TaskStatus = t.TaskStatus
	err = s.Store.AddTaskDate(req.Context(), TaskData)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	t.ID = Task.ID
	t.TaskID = Task.ID

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	json.NewEncoder(resp).Encode(map[string]*Tsk{
		"data": t,
	})
}

// func (s *ToDoService) addTask(resp http.ResponseWriter, req *http.Request) {
// 	t := &storages.Task{}
// 	err := json.NewDecoder(req.Body).Decode(t)
// 	defer req.Body.Close()
// 	if err != nil {
// 		resp.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	now := time.Now()
// 	userID, _ := userIDFromCtx(req.Context())
// 	t.ID = uuid.New().String()
// 	t.UserID = userID
// 	t.CreatedDate = now.Format("2006-01-02")

// 	resp.Header().Set("Content-Type", "application/json")

// 	err = s.Store.AddTask(req.Context(), t)
// 	if err != nil {
// 		resp.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(resp).Encode(map[string]string{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	json.NewEncoder(resp).Encode(map[string]*storages.Task{
// 		"data": t,
// 	})
// }

func value(req *http.Request, p string) sql.NullString {
	return sql.NullString{
		String: req.FormValue(p),
		Valid:  true,
	}
}

func (s *ToDoService) createToken(id string) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["user_id"] = id
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(s.JWTKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *ToDoService) validToken(req *http.Request) (*http.Request, bool) {
	token := req.Header.Get("Authorization")

	claims := make(jwt.MapClaims)
	t, err := jwt.ParseWithClaims(token, claims, func(*jwt.Token) (interface{}, error) {
		return []byte(s.JWTKey), nil
	})
	if err != nil {
		log.Println(err)
		return req, false
	}

	if !t.Valid {
		return req, false
	}

	id, ok := claims["user_id"].(string)
	if !ok {
		return req, false
	}

	req = req.WithContext(context.WithValue(req.Context(), userAuthKey(0), id))
	return req, true
}

type userAuthKey int8

func userIDFromCtx(ctx context.Context) (string, bool) {
	v := ctx.Value(userAuthKey(0))
	id, ok := v.(string)
	return id, ok
}
