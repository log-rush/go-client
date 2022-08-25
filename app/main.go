package main

import (
	"fmt"

	logRushClient "github.com/log-rush/go-client"
)

func main() {

	stream := logRushClient.NewLogStream(logRushClient.ClientOptions{
		DataSourceUrl: "http://localhost:7001/",
		BatchSize:     1,
	}, "abc", "id", "key")

	err := stream.Register()
	fmt.Println(err)

	fmt.Println(stream.Log("1"))
	fmt.Println(stream.Log("2"))
	fmt.Println(stream.Log("3"))
}
