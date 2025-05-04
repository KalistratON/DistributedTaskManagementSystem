package main

import (
	"context"
	helper "dtms/helper/pkg"
	"dtms/notification/internal"
	"log"
)

func main() {
	mngCln, err := helper.ConnectMongo()
	if err != nil {
		log.Fatalf("task service down with fatal error related to mongo connection= %v", err)
	}
	defer mngCln.Disconnect(context.Background())

	service := internal.NewNotificationService(mngCln)
	service.Run()
}
