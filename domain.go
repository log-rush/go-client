package logRushClient

const DefaultBatchSize = 0

type ClientOptions struct {
	DataSourceUrl string
	BatchSize     int
}

type LogRushLog struct {
	Log       string
	Timestamp int64
}

type LogRushApiStreamResponse struct {
	Id    string `json:"id"`
	Alias string `json:"alias"`
	Key   string `json:"key"`
}

type LogRushApiSuccessResponse struct {
	Success bool `json:"success"`
}

type LogRushApiErrorResponse struct {
	Message string `json:"message"`
}

type LogRushApiSuccessOrErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
