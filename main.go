package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/178inaba/datastore-example/repository"
	"github.com/k0kubun/pp"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()

	c, err := datastore.NewClient(ctx, os.Getenv("GCP_PROJECT"), option.WithCredentialsFile(os.Getenv("CRED_FILEPATH")))
	if err != nil {
		log.Fatal(err)
	}

	tr := repository.NewTaskRepository(c)

	k, err := tr.AddTask(ctx, "test description", time.Now())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Key of added task:")
	pp.Println(k)

	ts, err := tr.ListTasks(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Tasks:")
	pp.Println(ts)
}
