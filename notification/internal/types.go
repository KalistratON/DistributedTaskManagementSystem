package internal

type NotificationTaskInfo struct {
	TaskId   string `json:"task_id"`
	AuthorId string `json:"author_id"`
	Status   string `json:"status"`
}

type NotificationTaskLogInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Deadline    string `json:"deadline"`
	Author      string `json:"author"`
	Status      string `json:"status"`
}
