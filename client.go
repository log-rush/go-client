package logRushClient

type Client struct {
	options ClientOptions
	streams map[string]Stream
}

func NewClient(options ClientOptions) Client {
	return Client{
		options: options,
		streams: map[string]Stream{},
	}
}

func (c *Client) CreateStream(name string) (Stream, error) {
	stream := NewLogStream(c.options, name, "", "")
	if _, ok := c.streams[stream.id]; ok {
		return Stream{}, ErrStreamExists
	}
	c.streams[stream.id] = stream
	return stream, nil
}

func (c *Client) ResumeStream(name, id, key string) (Stream, error) {
	stream := NewLogStream(c.options, name, id, key)
	if _, ok := c.streams[stream.id]; ok {
		return Stream{}, ErrStreamExists
	}
	c.streams[stream.id] = stream
	return stream, nil
}

func (c *Client) DeleteStream(id string, sendRemainingLogs bool) error {
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

func (c *Client) Disconnect(sendRemainingLogs bool) error {
	var err error
	for _, stream := range c.streams {
		err = c.DeleteStream(stream.id, sendRemainingLogs)
	}
	return err
}
