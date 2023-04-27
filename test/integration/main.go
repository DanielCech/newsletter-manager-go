package main

import (
	"event-facematch-backend/test/integration/flows"
	"event-facematch-backend/test/integration/generate/swagger"
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
	// Run flows
	if flow == "flow1" || flow == allFlows {
		flows.Flow1(client)
	}
	if flow == "flow2" || flow == allFlows {
		flows.Flow2(client)
	}
	if flow == "flow3" || flow == allFlows {
		flows.Flow3(client)
	}
}
