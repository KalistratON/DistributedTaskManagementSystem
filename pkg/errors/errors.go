package errors

type NilMdb struct{}

func (err NilMdb) Error() string {
	return `variable mongo.Collection is nil pointer`
}

type NilRdb struct{}

func (err NilRdb) Error() string {
	return `variable redis.Client is nil pointer`
}

type NilTask struct{}

func (err NilTask) Error() string {
	return `variable "task" is nil pointer`
}

type InvalidIdTask struct{}

func (err InvalidIdTask) Error() string {
	return `task has invalid id`
}

type EmptyNameTask struct{}

func (err EmptyNameTask) Error() string {
	return `task has empty name`
}

type EmptyAuthorIdTask struct{}

func (err EmptyAuthorIdTask) Error() string {
	return `task has empty author id`
}

type EmptyStatusTask struct{}

func (err EmptyStatusTask) Error() string {
	return `task has empty status`
}

type EmptyUpdateResultTask struct{}

func (err EmptyUpdateResultTask) Error() string {
	return `no task was updated`
}

type EmptyDeleteResultTask struct{}

func (err EmptyDeleteResultTask) Error() string {
	return `no task was deleted`
}

type NilSqlDb struct{}

func (err NilSqlDb) Error() string {
	return `variable sql.db is nil pointer`
}

type UnconnectSqlDb struct{}

func (err UnconnectSqlDb) Error() string {
	return `can't connect to database`
}

type NoStatusCreate struct{}

func (err NoStatusCreate) Error() string {
	return `no status was created`
}

type NilUser struct{}

func (err NilUser) Error() string {
	return `nil user`
}

type LoginOrEmailAlreadtExist struct{}

func (err LoginOrEmailAlreadtExist) Error() string {
	return `login or email alreaty exist`
}

type PasswordValidateFail struct{}

func (err PasswordValidateFail) Error() string {
	return `password validate fail`
}

type InvalidObjectId struct{}

func (err InvalidObjectId) Error() string {
	return `can't convert string to objectId`
}

type NoResult struct{}

func (err NoResult) Error() string {
	return `no result on query`
}

type EmptyUpdate struct{}

func (err EmptyUpdate) Error() string {
	return `no data was updated`
}

type EmptyDelete struct{}

func (err EmptyDelete) Error() string {
	return `no data was deleted`
}

type EmptyId struct{}

func (err EmptyId) Error() string {
	return `empty id`
}

type EmptyToken struct{}

func (err EmptyToken) Error() string {
	return `empty token`
}

type NoTokenExpirationTime struct{}

func (err NoTokenExpirationTime) Error() string {
	return `no TOKEN_EXPIRATION_TIME`
}

type NilRepo struct{}

func (err NilRepo) Error() string {
	return `repo is nil`
}

type ExpiredToken struct{}

func (err ExpiredToken) Error() string {
	return `token is expired`
}
