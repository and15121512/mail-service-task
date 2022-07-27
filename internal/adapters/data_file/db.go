package data_file

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"gitlab.com/sukharnikov.aa/mail-service-task/internal/domain/models"
	"gitlab.com/sukharnikov.aa/mail-service-task/internal/utils"
	"go.uber.org/zap"
)

type DataFile struct {
	logger *zap.SugaredLogger
}

func New(logger *zap.SugaredLogger) (*DataFile, error) {
	return &DataFile{
		logger: logger,
	}, nil
}

func (db *DataFile) annotatedLogger(ctx context.Context) *zap.SugaredLogger {
	request_id, _ := ctx.Value(utils.CtxKeyRequestIDGet()).(string)
	method, _ := ctx.Value(utils.CtxKeyMethodGet()).(string)
	url, _ := ctx.Value(utils.CtxKeyURLGet()).(string)

	return db.logger.With(
		"request_id", request_id,
		"method", method,
		"url", url,
	)
}

func (db *DataFile) InsertTask(ctx context.Context, task models.Task) error {
	logger := db.annotatedLogger(ctx)

	tasks, err := db.read(ctx)
	if err != nil {
		logger.Errorf("failed to insert task in data file: cannot read data file")
		return fmt.Errorf("failed to insert task in data file: cannot read data file")
	}

	tasks = append(tasks, task)
	err = db.write(ctx, tasks)
	if err != nil {
		logger.Errorf("failed to insert task in data file: cannot write data file")
		return fmt.Errorf("failed to insert task in data file: cannot write data file")
	}
	return nil
}

func (db *DataFile) read(ctx context.Context) ([]models.Task, error) {
	logger := db.annotatedLogger(ctx)

	if _, err := os.Stat("data_files/tasks.json"); err == nil {
		return []models.Task{}, nil
	}

	file, err := ioutil.ReadFile("data_files/tasks.json")
	if err != nil {
		logger.Errorf("failed to insert task in data file: %s", err.Error())
		return []models.Task{}, fmt.Errorf("failed to insert task in data file: %s", err.Error())
	}

	tasks := []models.Task{}
	err = json.Unmarshal(file, &tasks)
	if err != nil {
		logger.Errorf("failed to insert task in data file: %s", err.Error())
		return []models.Task{}, fmt.Errorf("failed to insert task in data file: %s", err.Error())
	}
	return tasks, nil
}

func (db *DataFile) write(ctx context.Context, tasks []models.Task) error {
	logger := db.annotatedLogger(ctx)

	data, err := json.Marshal(tasks)
	if err != nil {
		logger.Errorf("failed to insert task in data file: %s", err.Error())
		return fmt.Errorf("failed to insert task in data file: %s", err.Error())
	}

	err = ioutil.WriteFile("data_files/tasks.json", data, 0644)
	if err != nil {
		logger.Errorf("failed to insert task in data file: %s", err.Error())
		return fmt.Errorf("failed to insert task in data file: %s", err.Error())
	}
	return nil
}
