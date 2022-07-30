package mongodb

import (
	"context"
	"fmt"

	"gitlab.com/sukharnikov.aa/mail-service-task/internal/domain/models"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type Database struct {
	mClient  *mongo.Client
	dbName   string
	collName string
	logger   *zap.SugaredLogger
}

func New(ctx context.Context, conn string, dbName string, collName string, logger *zap.SugaredLogger) (*Database, func(), error) {
	cOpts := options.Client().ApplyURI(conn)
	mClient, err := mongo.Connect(ctx, cOpts)
	if err != nil {
		logger.Errorf("failed to init mongo connection: %s", err.Error())
		return &Database{}, nil, fmt.Errorf("failed to init mongo connection: %s", err.Error())
	}
	return &Database{
			mClient:  mClient,
			dbName:   dbName,
			collName: collName,
			logger:   logger,
		}, func() {
			mClient.Disconnect(ctx)
		}, nil
}

func (db *Database) annotatedLogger(ctx context.Context) *zap.SugaredLogger {
	request_id, _ := ctx.Value(utils.CtxKeyRequestIDGet()).(string)
	method, _ := ctx.Value(utils.CtxKeyMethodGet()).(string)
	url, _ := ctx.Value(utils.CtxKeyURLGet()).(string)

	return db.logger.With(
		"request_id", request_id,
		"method", method,
		"url", url,
	)
}

func (db *Database) InsertTask(ctx context.Context, task models.Task) error {
	logger := db.annotatedLogger(ctx)

	mCollection := db.mClient.Database(db.dbName).Collection(db.collName)
	_, err := mCollection.InsertOne(ctx, task)
	if err != nil {
		logger.Errorf("failed to insert task %s: %s", task.ID, err.Error())
		return fmt.Errorf("failed to insert task %s: %s", task.ID, err.Error())
	}

	return nil
}

func (db *Database) GetTask(ctx context.Context, task_id string) (models.Task, error) {
	logger := db.annotatedLogger(ctx)

	mCollection := db.mClient.Database(db.dbName).Collection(db.collName)
	task := models.Task{}
	err := mCollection.FindOne(ctx, bson.M{"_id": task_id}).Decode(&task)
	if err != nil {
		logger.Errorf("failed to find task %s: %s", task.ID, err.Error())
		return models.Task{}, fmt.Errorf("failed to find task %s: %s", task.ID, err.Error())
	}

	return task, nil
}

func (db *Database) UpdateTask(ctx context.Context, task models.Task) error {
	logger := db.annotatedLogger(ctx)

	mCollection := db.mClient.Database(db.dbName).Collection(db.collName)
	mUpdateResult, err := mCollection.ReplaceOne(ctx, bson.M{"_id": task.ID}, task)
	if err != nil {
		logger.Errorf("failed to update task %s: %s", task.ID, err.Error())
		return fmt.Errorf("failed to update task %s: %s", task.ID, err.Error())
	}
	if mUpdateResult.MatchedCount == 0 {
		logger.Errorf("failed to update task %s: task not found", task.ID)
		return fmt.Errorf("failed to update task %s: task not found", task.ID)
	}

	return nil
}

func (db *Database) DeleteTask(ctx context.Context, task_id string) error {
	logger := db.annotatedLogger(ctx)

	mCollection := db.mClient.Database(db.dbName).Collection(db.collName)
	mDeleteResult, err := mCollection.DeleteOne(ctx, bson.M{"_id": task_id})
	if err != nil {
		logger.Errorf("failed to delete task %s: %s", task_id, err.Error())
		return fmt.Errorf("failed to delete task %s: %s", task_id, err.Error())
	}
	if mDeleteResult.DeletedCount == 0 {
		logger.Errorf("failed to delete task %s: task not found", task_id)
		return fmt.Errorf("failed to delete task %s: task not found", task_id)
	}

	return nil
}
