package logRushClient

import "errors"

const DefaultBatchSize = 0

var (
	ErrStreamExists    = errors.New("stream already exists")
	ErrStreamNotExists = errors.New("stream already exists")
)

type ClientOptions struct {
	DataSourceUrl string
	BatchSize     int
}

type LogRushLog struct {
	Log       string
	Timestamp int64
}

type ApiStreamResponse struct {
	Id    string `json:"id"`
	Alias string `json:"alias"`
	Key   string `json:"key"`
}

type ApiSuccessOrErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
