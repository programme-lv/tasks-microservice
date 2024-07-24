package handlers

import "net/http"

type ListTasksResponse struct {
	Tasks []Task `json:"tasks"`
}

func (c *Controller) ListTasks(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, ListTasksResponse{}, http.StatusOK)
}

// func (c *Controller) ListUsers(w http.ResponseWriter, r *http.Request) {
// 	user, err := c.UserService.ListUsers()
// 	if err != nil {
// 		respondWithBadRequest(w, "user not found")
// 		return
// 	}

// 	users := []Task{}
// 	for _, u := range user {
// 		users = append(users, Task{})
// 	}

// 	response := ListTasksResponse{
// 		Users: users,
// 	}

// 	respondWithJSON(w, response, http.StatusOK)
// }
