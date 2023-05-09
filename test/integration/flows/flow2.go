package flows

import (
	"newsletter-manager-go/test/integration/common"
	"newsletter-manager-go/test/integration/generate/swagger"
	"newsletter-manager-go/test/integration/testlog"
)

func Flow2(client *swagger.APIClient) {
	var name = "flow2"
	var description = "Initial DB setup. It can be convenient when called separately."

	testlog.StartFlow(name, description)

	common.WipePostgres()
	common.MigrateUp()

	user1 := common.NewUser("dummy1@test.com", "TheSecretPassword5")

	createAuthorInput1 := swagger.CreateAuthorInput{
		Name:     "John Doe",
		Email:    "john.doe@dummy.com",
		Password: "TheSecretPassword5",
	}

	//param := swagger.SessionApiAuthorSignUpOpts{Body: createAuthorInput1}

	signUpResp1, _, err := client.SessionApi.AuthorSignUp(user1.Context, createAuthorInput1)
	common.AssertNoError(err)
	user1.UpdateWith(signUpResp1)

	testlog.EndFlow(name)
}
