package logRushClient

type LogRushClient struct {
	options ClientOptions
}

func NewClient(options ClientOptions) LogRushClient {
	return LogRushClient{options}
}
