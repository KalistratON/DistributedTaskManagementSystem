package types

type TaskMessage struct {
	Id          string `json:"task_id"`
	AuthorId    string `json:"author_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Deadline    string `json:"deadline"`
	Status      string `json:"status"`
}

type TaskResponse struct {
	TaskId string `json:"task_id"`
	Status string `json:"status"`
	ErrMsg string `json:"err_msg"`
}
