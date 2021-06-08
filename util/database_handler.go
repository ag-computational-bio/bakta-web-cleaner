package util

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"k8s.io/client-go/kubernetes"
)

const COLLECTIONNAME = "jobs"

type DatabaseHandler struct {
	DB             *mongo.Client
	Collection     *mongo.Collection
	K8sClient      *kubernetes.Clientset
	BaseKey        string
	UserDataBucket string
	DBBucket       string
	Namespace      string
	ExpiryTime     int64
}

func InitDatabaseHandler() (*DatabaseHandler, error) {
	host := viper.GetString("Database.MongoHost")
	dbName := viper.GetString("Database.MongoDBName")
	dbUser := viper.GetString("Database.MongoUser")
	dbAuthSource := viper.GetString("Database.MongoAuthSource")
	dbPassword := os.Getenv("MongoPassword")
	dbPort := viper.GetString("Database.MongoPort")

	if dbPassword == "" {
		return nil, fmt.Errorf("password for mongodb required, can be set with env var MongoPassword")
	}

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%v:%v", host, dbPort)).SetAuth(
		options.Credential{
			AuthSource: dbAuthSource,
			Username:   dbUser,
			Password:   dbPassword,
		},
	))
	if err != nil {
		glog.Errorln(err.Error())
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		glog.Errorln(err.Error())
		return nil, err
	}

	collection := client.Database(dbName).Collection(COLLECTIONNAME)

	userBucket := viper.GetString("Objectstorage.S3.UserBucket")
	baseKey := viper.GetString("Objectstorage.S3.BaseKey")
	expiryTime := viper.GetInt64("ExpiryTime")

	handler := DatabaseHandler{
		DB:             client,
		Collection:     collection,
		UserDataBucket: userBucket,
		BaseKey:        baseKey,
		ExpiryTime:     expiryTime,
	}

	return &handler, nil
}

func (handler *DatabaseHandler) ListRunningJobs() ([]*Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	jobs_running_filter := bson.M{
		"Status": "RUNNING",
	}

	csr, err := handler.Collection.Find(ctx, jobs_running_filter)
	if err != nil {
		glog.Errorln(err.Error())
		return nil, err
	}

	var running_jobs []*Job
	err = csr.All(ctx, running_jobs)
	if err != nil {
		glog.Errorln(err.Error())
		return nil, err
	}

	return running_jobs, nil
}

func (handler *DatabaseHandler) ListExpiredJobs() ([]*Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	maxTime := time.Now().AddDate(0, 0, -10)

	query := bson.M{
		"Created": bson.M{
			"$lt": maxTime,
		},
	}

	var jobs []*Job

	csr, err := handler.Collection.Find(ctx, query)
	if err != nil {
		glog.Errorln(err.Error())
		return nil, err
	}

	err = csr.All(context.Background(), &jobs)
	if err != nil {
		glog.Errorln(err.Error())
		return nil, err
	}

	return jobs, nil
}

func (handler *DatabaseHandler) DeleteJob(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := bson.M{
		"JobID": id,
	}

	_, err := handler.Collection.DeleteOne(ctx, query)
	if err != nil {
		glog.Errorln(err.Error())
		return err
	}

	return nil
}
