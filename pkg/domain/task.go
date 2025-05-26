package domain

type Task struct {
	Id          string `json:"id" bson:"-"`
	AuthorId    string `json:"author_id" bson:"author_id"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	Deadline    string `json:"deadline,omitempty" bson:"deadline,omitempty"`
	Status      string `json:"status" bson:"status"`
}
