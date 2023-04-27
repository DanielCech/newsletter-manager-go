package flows

import (
	"event-facematch-backend/test/integration/common"
	"event-facematch-backend/test/integration/generate/swagger"
	"event-facematch-backend/test/integration/testlog"
)

func Flow1(client *swagger.APIClient) {
	var name = "flow1"
	var description = "Initial DB setup. It can be convenient when called separately."

	testlog.StartFlow(name, description)

	common.WipePostgres()
	common.MigrateUp()
	common.PopulatePostgres()

	testlog.EndFlow(name)
}
