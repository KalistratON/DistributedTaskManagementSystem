package repository

import (
	"context"
	helper "dtms/pkg/database"
	"dtms/pkg/domain"
	user_errors "dtms/pkg/errors"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"os"
	"strconv"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func prepareResources() (*mongo.Collection, func(context.Context) error, error) {
	mngCln, err := helper.ConnectMongo()
	if err != nil {
		return nil, nil, err
	}

	mngCol := mngCln.Database("user").Collection("list")
	if mngCol == nil {
		return nil, nil, fmt.Errorf("nil mongo collection")
	}
	mngCol.Drop(context.Background())

	return mngCol, mngCln.Disconnect, nil
}

func TestMain(m *testing.M) {
	log.SetOutput(io.Discard)

	os.Exit(m.Run())
}

func TestCreate(t *testing.T) {
	validTests := []domain.User{
		// valid cases
		domain.User{Login: "arthur", Password: "1234567890", Email: "template1@index.ru"},
		domain.User{Login: "oleg", Password: "1234567890", Email: "template2@index.ru"},
		domain.User{Login: "anton", Password: "1234567890", Email: "template3@index.ru"},
		domain.User{Login: "akakiy", Password: "1234567890", Email: "template4@index.ru"},
		domain.User{Login: "armavir", Password: "1234567890", Email: "template5@index.ru"},
		domain.User{Login: "gangut", Password: "1234567890", Email: "template6@index.ru"},
	}

	alreadyExistTests := []domain.User{
		// already exist cases
		domain.User{Login: "arthur", Password: "1234567890", Email: "template7@index.ru"},
		domain.User{Login: "oleg", Password: "1234567890", Email: "template8@index.ru"},
		domain.User{Login: "arthur1", Password: "1234567890", Email: "template1@index.ru"},
		domain.User{Login: "oleg2", Password: "1234567890", Email: "template2@index.ru"},
	}

	invalidPasswordTests := []domain.User{
		// invalid password cases
		domain.User{Login: "arthur1_password", Password: "123456789", Email: "template1_password@index.ru"},
		domain.User{Login: "arthur2_password", Password: "", Email: "template2_password@index.ru"},
	}

	mngCln, closeMongoFunc, err := prepareResources()
	if err != nil {
		t.Fatalf("%v", err)
	}

	t.Cleanup(func() {
		mngCln.Drop(context.Background())
		closeMongoFunc(context.Background())
	})

	rep, err := NewUserRepository(mngCln, context.Background())
	if err != nil {
		t.Fatalf("%v", err)
	}

	// valid cases
	for _, test := range validTests {
		t.Run("User-Create | "+test.Login+" | "+test.Password+" | "+test.Email, func(t *testing.T) {
			result, err := rep.Create(&test)
			if err != nil {
				t.Errorf("%v", err)
			}

			if result.Login != test.Login || result.Email != test.Email || result.Password != test.Password {
				t.Errorf("non-equal result to origin: %v", result)
			}

			_, err = primitive.ObjectIDFromHex(result.Id)
			if err != nil {
				t.Errorf("invalid id of user: %s", result.Id)
			}
		})
	}

	// already exist cases
	for _, test := range alreadyExistTests {
		t.Run("User-Create | "+test.Login+" | "+test.Password+" | "+test.Email, func(t *testing.T) {
			_, err := rep.Create(&test)
			if !errors.Is(err, user_errors.LoginOrEmailAlreadtExist{}) {
				t.Errorf("no expected error %v", err)
			}
		})
	}

	// already exist cases
	for _, test := range invalidPasswordTests {
		t.Run("User-Create | "+test.Login+" | "+test.Password+" | "+test.Email, func(t *testing.T) {
			_, err := rep.Create(&test)
			if !errors.Is(err, user_errors.PasswordValidateFail{}) {
				t.Errorf("no expected error %v", err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	validTests := [][2]domain.User{
		// valid cases
		{
			domain.User{Login: "user_create_1", Password: "1234567890", Email: "user_create_1"},
			domain.User{Login: "user_update_1", Password: "1234567890", Email: "user_create_1"},
		},
		{
			domain.User{Login: "user_create_2", Password: "1234567890", Email: "user_create_2"},
			domain.User{Login: "user_create_2", Password: "update_password_2", Email: "user_create_2"},
		},
		{
			domain.User{Login: "user_create_3", Password: "1234567890", Email: "user_create_3"},
			domain.User{Login: "user_create_3", Password: "1234567890", Email: "update_email_3"},
		},
		{
			domain.User{Login: "user_create_4", Password: "1234567890", Email: "user_create_4"},
			domain.User{Login: "user_update_4", Password: "update_password_4", Email: "user_create_4"},
		},
		{
			domain.User{Login: "user_create_5", Password: "1234567890", Email: "user_create_5"},
			domain.User{Login: "user_update_5", Password: "1234567890", Email: "update_email_5"},
		},
		{
			domain.User{Login: "user_create_6", Password: "1234567890", Email: "user_create_6"},
			domain.User{Login: "user_create_6", Password: "update_password_6", Email: "update_email_6"},
		},
	}

	// empty update cases
	emptyUpdateTests := [][2]domain.User{
		{
			domain.User{Login: "empty_update_1", Password: "empty_update_1", Email: "empty_update_1"},
			domain.User{Login: "empty_update_1", Password: "empty_update_1", Email: "empty_update_1"},
		},
	}

	mngCln, closeMongoFunc, err := prepareResources()
	if err != nil {
		t.Fatalf("%v", err)
	}

	t.Cleanup(func() {
		mngCln.Drop(context.Background())
		closeMongoFunc(context.Background())
	})

	rep, err := NewUserRepository(mngCln, context.Background())
	if err != nil {
		t.Fatalf("%v", err)
	}

	// valid cases
	for _, test := range validTests {
		userCreate := test[0]
		userUpdate := test[1]
		t.Run("User-Update | "+userUpdate.Login+" | "+userUpdate.Password+" | "+userUpdate.Email, func(t *testing.T) {
			result, err := rep.Create(&test[0])
			if err != nil {
				t.Errorf("error while create: %v", err)
			}

			userUpdate.Id = result.Id
			result, err = rep.Update(&userUpdate)
			if err != nil {
				t.Errorf("%v", err)
			}

			if result.Login != userUpdate.Login ||
				result.Email != userUpdate.Email ||
				result.Password != userUpdate.Password ||
				result.Id != userUpdate.Id {
				t.Errorf("non-equal result to origin: %v -> %v", userCreate, result)
			}

			_, err = primitive.ObjectIDFromHex(result.Id)
			if err != nil {
				t.Errorf("invalid id of user: %s", result.Id)
			}
		})
	}

	// empty update cases
	for _, test := range emptyUpdateTests {
		userCreate := test[0]
		userUpdate := test[1]
		t.Run("User-Update | "+userUpdate.Login+" | "+userUpdate.Password+" | "+userUpdate.Email, func(t *testing.T) {
			result, err := rep.Create(&userCreate)
			if err != nil {
				t.Errorf("error while create: %v", err)
			}

			userUpdate.Id = result.Id
			result, err = rep.Update(&userUpdate)
			if !errors.Is(err, user_errors.EmptyUpdate{}) {
				t.Errorf("no expected error: %v", err)
			}
		})
	}
}

func TestGet(t *testing.T) {
	validTests := []domain.User{
		// valid cases
		domain.User{Login: "arthur", Password: "1234567890", Email: "template1@index.ru"},
		domain.User{Login: "oleg", Password: "1234567890", Email: "template2@index.ru"},
		domain.User{Login: "anton", Password: "1234567890", Email: "template3@index.ru"},
		domain.User{Login: "akakiy", Password: "1234567890", Email: "template4@index.ru"},
		domain.User{Login: "armavir", Password: "1234567890", Email: "template5@index.ru"},
		domain.User{Login: "gangut", Password: "1234567890", Email: "template6@index.ru"},
	}

	mngCln, closeMongoFunc, err := prepareResources()
	if err != nil {
		t.Fatalf("%v", err)
	}

	t.Cleanup(func() {
		mngCln.Drop(context.Background())
		closeMongoFunc(context.Background())
	})

	rep, err := NewUserRepository(mngCln, context.Background())
	if err != nil {
		t.Fatalf("%v", err)
	}

	for _, test := range validTests {
		t.Run("User-Get | "+test.Login+" | "+test.Password+" | "+test.Email, func(t *testing.T) {
			result, err := rep.Create(&test)
			if err != nil {
				t.Errorf("%v", err)
			}

			getResult, err := rep.Get(result.Id)
			if err != nil {
				t.Errorf("%v", err)
			}

			if getResult.Login != result.Login ||
				getResult.Email != result.Email ||
				getResult.Password != result.Password ||
				getResult.Id != result.Id {
				t.Errorf("non-equal result to origin: %v", result)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	validTests := []domain.User{
		// valid cases
		domain.User{Login: "arthur", Password: "1234567890", Email: "template1@index.ru"},
		domain.User{Login: "oleg", Password: "1234567890", Email: "template2@index.ru"},
		domain.User{Login: "anton", Password: "1234567890", Email: "template3@index.ru"},
		domain.User{Login: "akakiy", Password: "1234567890", Email: "template4@index.ru"},
		domain.User{Login: "armavir", Password: "1234567890", Email: "template5@index.ru"},
		domain.User{Login: "gangut", Password: "1234567890", Email: "template6@index.ru"},
	}

	mngCln, closeMongoFunc, err := prepareResources()
	if err != nil {
		t.Fatalf("%v", err)
	}

	t.Cleanup(func() {
		mngCln.Drop(context.Background())
		closeMongoFunc(context.Background())
	})

	rep, err := NewUserRepository(mngCln, context.Background())
	if err != nil {
		t.Fatalf("%v", err)
	}

	for _, test := range validTests {
		t.Run("User-Get | "+test.Login+" | "+test.Password+" | "+test.Email, func(t *testing.T) {
			result, err := rep.Create(&test)
			if err != nil {
				t.Errorf("%v", err)
			}

			delResult, err := rep.Delete(result.Id)
			if err != nil {
				t.Errorf("%v", err)
			}

			if delResult.Login != result.Login ||
				delResult.Email != result.Email ||
				delResult.Password != result.Password ||
				delResult.Id != result.Id {
				t.Errorf("non-equal result to origin: %v", result)
			}

			getResult, err := rep.Get(result.Id)
			if err == nil {
				t.Errorf("user still can be reached: %v", getResult)
			}
		})
	}
}

func BenchmarkCreate(b *testing.B) {
	mngCln, closeMongoFunc, err := prepareResources()
	if err != nil {
		b.Fatalf("%v", err)
	}

	defer closeMongoFunc(context.Background())
	defer mngCln.Drop(context.Background())

	rep, err := NewUserRepository(mngCln, context.Background())
	if err != nil {
		b.Fatalf("%v", err)
	}

	for i := 0; i < b.N; i++ {
		b.Run("User-Create", func(b *testing.B) {
			user := domain.User{
				Login:    "user_" + strconv.Itoa(rand.Int()),
				Password: "1234567890",
				Email:    "user_" + strconv.Itoa(rand.Int()),
			}

			_, err := rep.Create(&user)
			if err != nil {
				b.Fatalf("%v", err)
			}
		})
	}
}
