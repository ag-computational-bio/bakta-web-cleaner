package cleaner

import (
	"github.com/ag-computational-bio/bakta-web-cleaner/util"
	"github.com/golang/glog"
	"golang.org/x/sync/errgroup"
)

type Handler struct {
	Database *util.DatabaseHandler
	S3       *util.S3ObjectStorageHandler
	K8s      *util.K8sHandler
}

func Init() (*Handler, error) {
	db, err := util.InitDatabaseHandler()
	if err != nil {
		glog.Errorln(err.Error())
		return nil, err
	}

	s3, err := util.InitS3ObjectStorageHandler()
	if err != nil {
		glog.Errorln(err.Error())
		return nil, err
	}

	return &Handler{
		Database: db,
		S3:       s3,
	}, nil
}

func (handler *Handler) RemoveExpired() error {
	expiredJobs, err := handler.Database.ListExpiredJobs()
	if err != nil {
		glog.Errorln(err.Error())
		return err
	}

	jobsChan := make(chan *util.Job, 500)
	errgrp := errgroup.Group{}

	for i := 1; i <= 100; i++ {
		errgrp.Go(func() error {
			for job := range jobsChan {
				err := handler.deleteExpiredJob(job)
				// Log errors but dont stop execution if one occurs to clean up as many jobs as possible
				// Straggler can then be investigated further
				if err != nil {
					glog.Errorln(err.Error())
				}
			}

			return nil
		})
	}

	go func() {
		for _, job := range expiredJobs {
			jobsChan <- job
		}
	}()

	err = errgrp.Wait()
	if err != nil {
		glog.Errorln(err.Error())
		return err
	}

	return nil

}

func (handler *Handler) deleteExpiredJob(job *util.Job) error {
	err := handler.S3.DeleteObject(job.DataBucket, job.FastaKey)
	if err != nil {
		glog.Errorln(err.Error())
		return err
	}

	err = handler.S3.DeleteObject(job.DataBucket, job.ResultKey)
	if err != nil {
		glog.Errorln(err.Error())
		return err
	}

	err = handler.S3.DeleteObject(job.DataBucket, job.ProdigalKey)
	if err != nil {
		glog.Errorln(err.Error())
		return err
	}

	err = handler.S3.DeleteObject(job.DataBucket, job.RepliconKey)
	if err != nil {
		glog.Errorln(err.Error())
		return err
	}

	if !job.IsDeleted {
		err = handler.K8s.DeleteK8sJob(job.JobID)
		if err != nil {
			glog.Errorln(err.Error())
			return err
		}
	}

	err = handler.Database.DeleteJob(job.JobID)
	if err != nil {
		glog.Errorln(err.Error())
		return err
	}

	return nil
}
