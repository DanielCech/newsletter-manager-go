package main

import (
	"newsletter-manager-go/test/integration/flows"
	"newsletter-manager-go/test/integration/generate/swagger"
	"os"
)

const (
	localhostAPI                 = "http://localhost:8080"
	localhostAPIWithCharlesProxy = "http://localhost.charlesproxy.com:8080"

	allFlows = "all"
)

func main() {
	configuration := swagger.NewConfiguration()
	configuration.BasePath = localhostAPI

	client := swagger.NewAPIClient(configuration)

	flow := "all"
	if len(os.Args) > 1 {
		flow = os.Args[1]
	}

	runFlows(flow, client)
}

func runFlows(flow string, client *swagger.APIClient) {
	runFlow(flow, client, "flow1", flows.Flow1)
}

func runFlow(currentFlow string, client *swagger.APIClient, flowName string, run func(client *swagger.APIClient)) {
	if currentFlow == flowName || currentFlow == allFlows {
		run(client)
	}
}
