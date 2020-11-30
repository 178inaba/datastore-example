package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/k0kubun/pp"

	"github.com/178inaba/datastore-example/repository"
)

func main() {
	var desc, text string
	flag.StringVar(&desc, "d", "", "Task description")
	flag.StringVar(&text, "t", "", "Task text")
	flag.Parse()

	ctx := context.Background()

	c, err := datastore.NewClient(ctx, os.Getenv("GCP_PROJECT"))
	if err != nil {
		log.Fatal(err)
	}

	tr := repository.NewTaskRepository(c)

	k, err := tr.AddTask(ctx, desc, text, time.Now())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Key of added task:")
	pp.Println(k)

	ts, err := tr.ListTasks(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("ListTasks:")
	pp.Println(ts)

	keyTs, err := tr.FilterKey(ctx, k)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("FilterKey:")
	pp.Println(keyTs)

	descTs, err := tr.FilterDescription(ctx, "test description")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("FileterDescription:")
	pp.Println(descTs)

	ids, err := tr.GetTaskIDsFilterDescription(ctx, "test description")
	if err != nil {
		log.Fatalf("GetTaskIDsFilterDescription: %v.", err)
	}

	fmt.Println("GetTaskIDsFilterDescription:")
	pp.Println(ids)

	allCnt, err := tr.CountAll(ctx)
	if err != nil {
		log.Fatal(err)
	}

	descNotNullCnt, err := tr.CountDescNotNull(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("All Count: %d\n", allCnt)
	fmt.Printf("Description Not Null Count: %d\n", descNotNullCnt)
}
