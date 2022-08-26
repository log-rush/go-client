package logRushClient

type LogRushClient struct {
	options ClientOptions
	streams map[string]LogRushStream
}

func NewClient(options ClientOptions) LogRushClient {
	return LogRushClient{
		options: options,
		streams: map[string]LogRushStream{},
	}
}

func (c *LogRushClient) CreateStream(name string) (LogRushStream, error) {
	stream := NewLogStream(c.options, name, "", "")
	if _, ok := c.streams[stream.id]; ok {
		return LogRushStream{}, ErrStreamExists
	}
	c.streams[stream.id] = stream
	return stream, nil
}

func (c *LogRushClient) ResumeStream(name, id, key string) (LogRushStream, error) {
	stream := NewLogStream(c.options, name, id, key)
	if _, ok := c.streams[stream.id]; ok {
		return LogRushStream{}, ErrStreamExists
	}
	c.streams[stream.id] = stream
	return stream, nil
}

func (c *LogRushClient) DeleteStream(id string, sendRemainingLogs bool) error {
	stream, ok := c.streams[id]
	if !ok {
		return ErrStreamNotExists
	}

	if sendRemainingLogs {
		stream.FlushLogs()
	}
	delete(c.streams, id)

	return nil
}

func (c *LogRushClient) Disconnect(sendRemainingLogs bool) error {
	var err error
	for _, stream := range c.streams {
		err = c.DeleteStream(stream.id, sendRemainingLogs)
	}
	return err
}
