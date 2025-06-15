package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/kakkky/mcp-sdk-go/client"
	"github.com/kakkky/mcp-sdk-go/client/transport"
	"github.com/kakkky/mcp-sdk-go/shared/protocol"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

func main() {
	c := client.NewClient(
		schema.Implementation{
			Name:    "example-client",
			Version: "1.0.0",
		},
		&client.ClientOptions{
			Capabilities: schema.ClientCapabilities{
				Experimental: map[string]any{
					"exampleFeature": true,
				},
			},
			ProtocolOptions: protocol.ProtocolOptions{
				EnforceStrictCapabilities: true,
			},
		},
	)
	t := transport.NewStdioClientTransport(
		transport.StdioServerParameters{
			Command: "go",
			Args:    []string{"run", "./examples/server/with-stdio/main.go"}, // サーバープログラムの実行コマンド
		},
	)
	go func() {
		err := c.Connect(t)
		if err != nil {
			panic(err)
		}
	}()
	<-client.OperationPhaseStartNotify

	// コマンド入力のためのループ
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(">")
	for scanner.Scan() {
		switch scanner.Text() {
		case "ping":
			result, err := c.Ping()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("<-", result)
		case "resources/list":
			result, err := c.ListResources()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("<-", result)
		}
		fmt.Print(">")
	}

}
