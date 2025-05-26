package repository

import (
	"context"
	"dtms/pkg/domain"
	"dtms/pkg/errors"
	"fmt"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	mdb *mongo.Collection
	ctx context.Context
}

func NewUserRepository(mdb *mongo.Collection, ctx context.Context) (UserRepository, error) {
	if mdb == nil {
		return nil, errors.NilMdb{}
	}

	return &userRepository{
		mdb: mdb,
		ctx: ctx,
	}, nil
}

var _ UserRepository = &userRepository{}

func (r *userRepository) fitUser(user *domain.User) {
	user.Id = strings.TrimSpace(user.Id)
	user.Login = strings.TrimSpace(user.Login)
	user.Password = strings.TrimSpace(user.Password)
	user.Email = strings.TrimSpace(user.Email)
}

func (r *userRepository) isExist(filter *[]bson.M) (bool, error) {
	fitFilter := bson.M{
		"$or": *filter,
	}

	var result bson.M
	err := r.mdb.FindOne(r.ctx, fitFilter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		log.Printf("no result was found: %v", err)
		return false, nil
	} else if err != nil {
		log.Printf("error while trying to get result from mongo: %v", err)
		return false, err
	}
	return true, nil
}

func (r *userRepository) validatePassword(password string) error {
	if len(strings.TrimSpace(password)) < 10 {
		return errors.PasswordValidateFail{}
	}
	return nil
}

func (r *userRepository) Create(user *domain.User) (*domain.User, error) {
	if user == nil {
		return nil, errors.NilUser{}
	}
	r.fitUser(user)

	filter := []bson.M{
		bson.M{"login": user.Login},
		bson.M{"email": user.Email},
	}
	exist, err := r.isExist(&filter)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, errors.LoginOrEmailAlreadtExist{}
	}

	if err := r.validatePassword(user.Password); err != nil {
		return nil, err
	}

	insRes, err := r.mdb.InsertOne(r.ctx, *user)
	if err != nil {
		log.Printf("error during insert on task creation: %v", err)
		return nil, err
	}
	log.Printf("was inserterd by id = %v", insRes.InsertedID)

	id := insRes.InsertedID.(primitive.ObjectID).Hex()
	result, err := r.Get(id)
	if err != nil {
		log.Printf("error while checking inserting in mdb: %v", err)
		return nil, err
	}

	return result, nil
}

func (r *userRepository) Update(user *domain.User) (*domain.User, error) {
	if user == nil {
		return nil, errors.NilUser{}
	}
	r.fitUser(user)

	id, err := primitive.ObjectIDFromHex(user.Id)
	if err != nil {
		log.Printf("can't convert id to objectId")
		return nil, fmt.Errorf("%v: %v", errors.InvalidObjectId{}, err)
	}

	filter := []bson.M{
		bson.M{"_id": id},
	}
	exist, err := r.isExist(&filter)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.NoResult{}
	}

	updRes, err := r.mdb.UpdateByID(r.ctx, id, bson.M{
		"$set": bson.M{
			"email":    user.Email,
			"login":    user.Login,
			"password": user.Password,
		},
	})
	if err != nil {
		log.Printf("error while updating on user update")
		return nil, err
	}

	if updRes.MatchedCount < 1 || updRes.ModifiedCount < 1 {
		log.Println("no data was updated")
		return nil, errors.EmptyUpdate{}
	}

	result, err := r.Get(user.Id)
	if err != nil {
		log.Printf("error while checking updating result")
		return nil, err
	}
	return result, nil
}

func (r *userRepository) Get(id string) (*domain.User, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("can't convert id to objectID")
		return nil, err
	}

	var result domain.User
	filter := bson.M{
		"_id": objId,
	}
	err = r.mdb.FindOne(r.ctx, filter).Decode(&result)
	if err != nil {
		log.Printf("error while getting user: %v", err)
		return nil, err
	}
	result.Id = id
	return &result, nil
}

func (r *userRepository) GetAll(filter map[string]interface{}) (*[]domain.User, error) {
	cur, err := r.mdb.Find(r.ctx, filter)
	if err != nil {
		log.Printf("error while trying get data with filter: %v", err)
		return nil, err
	}
	defer cur.Close(r.ctx)

	type mgdbBufferUser struct {
		Id primitive.ObjectID `bson:"_id"`

		domain.User `bson:",inline"`
	}

	buffer := make([]mgdbBufferUser, 0)
	if err = cur.All(r.ctx, &buffer); err != nil {
		log.Printf("error while decoding data: %v", err)
		return nil, err
	}

	result := make([]domain.User, 0, len(buffer))
	for _, v := range buffer {
		result = append(result, domain.User{
			Id:       v.Id.Hex(),
			Login:    v.Login,
			Email:    v.Email,
			Password: v.Password,
		})
	}
	return &result, nil
}

func (r *userRepository) Delete(id string) (*domain.User, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("user has invalid id: %v", err)
		return nil, fmt.Errorf("%v: %v", errors.InvalidObjectId{}, err)
	}

	task, err := r.Get(id)
	if err != nil {
		log.Printf("user can't be found: %v", err)
		return nil, err
	}

	delRes, err := r.mdb.DeleteOne(r.ctx, bson.M{"_id": objId})
	if err != nil {
		log.Printf("task can not be deleted: %v", err)
		return nil, err
	}

	if delRes.DeletedCount < 1 {
		log.Printf("no data was deleted")
		return nil, errors.EmptyDelete{}
	}

	return task, nil
}
