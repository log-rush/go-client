package logRushClient

import (
	"errors"
	"time"
)

type Stream struct {
	options   ClientOptions
	logsQueue chan Log
	name      string
	id        string
	key       string
}

func NewLogStream(options ClientOptions, name, id, key string) Stream {
	stream := Stream{
		options:   options,
		logsQueue: make(chan Log, options.BatchSize*3),
		name:      name,
		id:        id,
		key:       key,
	}
	if options.BatchSize == 0 {
		stream.options.BatchSize = 1
	}
	return stream
}

func (s *Stream) Id() string {
	return s.id
}

func (s *Stream) SecretKey() string {
	return s.key
}

func (s *Stream) Name() string {
	return s.name
}

func (s *Stream) Register() error {
	stream, err := logRushHttpApi.RegisterStream(s.options.DataSourceUrl, s.name, s.id, s.key)
	s.id = stream.Id
	s.key = stream.Key
	return err
}

func (s *Stream) Log(msg string) error {
	s.logsQueue <- Log{
		Log:       msg,
		Timestamp: time.Now().UnixMilli(),
	}

	if len(s.logsQueue) >= s.options.BatchSize {
		// send logs
		if len(s.logsQueue) == 1 {
			_, err := logRushHttpApi.Log(s.options.DataSourceUrl, s.id, <-s.logsQueue)
			return err
		} else {
			logs := []Log{}
			for i := 0; i < s.options.BatchSize; i++ {
				logs = append(logs, <-s.logsQueue)
			}
			_, err := logRushHttpApi.Batch(s.options.DataSourceUrl, s.id, logs)
			return err
		}
	}
	return nil
}

func (s *Stream) FlushLogs() error {
	logs := []Log{}
	for i := 0; i < s.options.BatchSize; i++ {
		logs = append(logs, <-s.logsQueue)
	}
	_, err := logRushHttpApi.Batch(s.options.DataSourceUrl, s.id, logs)
	return err
}

func (s *Stream) Destroy() error {
	close(s.logsQueue)
	response, err := logRushHttpApi.UnregisterStream(s.options.DataSourceUrl, s.id, s.key)

	if response.Message != "" {
		return errors.New(response.Message)
	}

	return err
}
