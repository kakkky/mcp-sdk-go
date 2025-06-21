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
				Roots: &schema.Roots{
					ListChanged: true,
				},
			},
			ProtocolOptions: protocol.ProtocolOptions{
				EnforceStrictCapabilities: true,
			},
		},
	)
	c.SetRequestHandler(&schema.ListRootsRequestSchema{MethodName: "roots/list"}, func(jrr schema.JsonRpcRequest) (schema.Result, error) {
		return &schema.ListRootsResultSchema{
			Roots: []schema.RootSchema{
				{
					Uri:  "file:///example/root1",
					Name: "Root 1",
				},
			},
		}, nil
	})
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
	fmt.Println("Initialization complete 🎉 Client is ready to send commands.")
	c.ListTools()
	c.CallTool(schema.CallToolRequestParams{
		Name: "calculate",
		Arguments: map[string]any{
			"first":  5,
			"second": []float64{10, 20},
		},
	})
	// コマンド入力のためのループ
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter method :  ")
	for scanner.Scan() {
		switch scanner.Text() {
		case "ping":
			result, err := c.Ping()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Ping", result)
		case "resources/list":
			result, err := c.ListResources()
			if err != nil {
				fmt.Println(err)
			}
			resourceList := result.(*schema.ListResourcesResultSchema)
			var resources []string
			for _, resource := range resourceList.Resources {
				resources = append(resources, fmt.Sprintf("Name: %s, URI: %s ,Metadata:%v", resource.Name, resource.Uri, *resource.ResourceMetadata))
			}
			fmt.Println("Resources:", resources)
		}
		fmt.Println("Enter method :  ")
	}

}
