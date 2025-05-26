package domain

type User struct {
	Id       string `json:"id" bson:"-"`
	Login    string `json:"login" bson:"login"`
	Password string `json:"password" bson:"password"`
	Email    string `json:"email" bson:"email"`
}
