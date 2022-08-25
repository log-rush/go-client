package logRushClient

import (
	"time"
)

type LogRushStream struct {
	options   ClientOptions
	logsQueue chan LogRushLog
	name      string
	id        string
	key       string
}

func NewLogStream(options ClientOptions, name, id, key string) LogRushStream {
	stream := LogRushStream{
		options:   options,
		logsQueue: make(chan LogRushLog, options.BatchSize*3),
		name:      name,
		id:        id,
		key:       key,
	}
	if options.BatchSize == 0 {
		stream.options.BatchSize = 1
	}
	return stream
}

func (s *LogRushStream) Id() string {
	return s.id
}

func (s *LogRushStream) SecretKey() string {
	return s.key
}

func (s *LogRushStream) Name() string {
	return s.name
}

func (s *LogRushStream) Register() error {
	stream, err := logRushHttpApi.RegisterStream(s.options.DataSourceUrl, s.name, s.id, s.key)
	s.id = stream.Id
	s.key = stream.Key
	return err
}

func (s *LogRushStream) Log(msg string) error {
	s.logsQueue <- LogRushLog{
		Log:       msg,
		Timestamp: time.Now().UnixMilli(),
	}

	if len(s.logsQueue) >= s.options.BatchSize {
		// send logs
		if len(s.logsQueue) == 1 {
			_, err := logRushHttpApi.Log(s.options.DataSourceUrl, s.id, <-s.logsQueue)
			return err
		} else {
			logs := []LogRushLog{}
			for i := 0; i < s.options.BatchSize; i++ {
				logs = append(logs, <-s.logsQueue)
			}
			_, err := logRushHttpApi.Batch(s.options.DataSourceUrl, s.id, logs)
			return err

		}
	}
	return nil
}

func (s *LogRushStream) Destroy() error {
	close(s.logsQueue)
	_, err := logRushHttpApi.UnregisterStream(s.options.DataSourceUrl, s.id, s.key)

	return err
}
