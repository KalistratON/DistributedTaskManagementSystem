package domain

type TaskHistoryAction string

const (
	CREATE TaskHistoryAction = "create"
	UPDATE TaskHistoryAction = "update"
	DELETE TaskHistoryAction = "delete"
)

type TaskHistory struct {
	Common Task `json:",inline"`

	Action string `json:"action"` // action of user
}
