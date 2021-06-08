package util

import (
	"context"
	"log"
	"testing"
	"time"

	testutil "github.com/ag-computational-bio/bakta-web-cleaner/test_util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDatabase(t *testing.T) {
	testutil.InitTestConf()

	db, err := InitDatabaseHandler()
	if err != nil {
		log.Fatalln(err.Error())
	}

	base_model := Job{
		JobID:       "test",
		Secret:      "test",
		K8sID:       "test",
		Created:     primitive.Timestamp{T: uint32(time.Now().Unix())},
		Updated:     primitive.Timestamp{T: uint32(time.Now().Unix())},
		FastaKey:    "test",
		ProdigalKey: "test",
		RepliconKey: "test",
		DataBucket:  "test",
		ResultKey:   "test",
		Jobname:     "test",
		IsDeleted:   true,
		ConfString:  "fgoo",
		Status:      "RUNNING",
		ExpiryDate:  primitive.Timestamp{T: uint32(time.Now().Unix())},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err = db.Collection.InsertOne(ctx, base_model)
	if err != nil {
		log.Fatalln(err.Error())
	}

	base_model.JobID = "test1"
	base_model.Created.T = uint32(time.Now().AddDate(0, 0, -100).Unix())

	_, err = db.Collection.InsertOne(ctx, base_model)
	if err != nil {
		log.Fatalln(err.Error())
	}

	base_model.JobID = "test2"
	_, err = db.Collection.InsertOne(ctx, base_model)
	if err != nil {
		log.Fatalln(err.Error())
	}

	names := make(map[string]struct{})
	names["test1"] = struct{}{}
	names["test2"] = struct{}{}

	jobs, err := db.ListExpiredJobs()
	if err != nil {
		log.Fatalln(err.Error())
	}

	for _, job := range jobs {
		if _, ok := names[job.JobID]; !ok {
			log.Fatalln("expired job found that should not have been expired")
		}
	}

}
