package flows

import (
	"newsletter-manager-go/test/integration/common"
	"newsletter-manager-go/test/integration/generate/swagger"
	"newsletter-manager-go/test/integration/testlog"
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
