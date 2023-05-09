package flows

import (
	"newsletter-manager-go/test/integration/common"
	"newsletter-manager-go/test/integration/generate/swagger"
	"newsletter-manager-go/test/integration/testlog"
)

func Flow2(client *swagger.APIClient) {
	var name = "flow2"
	var description = "User sign up."

	testlog.StartFlow(name, description)

	common.WipePostgres()
	common.MigrateUp()

	user1 := common.NewUser()

	createAuthorInput1 := swagger.CreateAuthorInput{
		Name:     "John Doe",
		Email:    "john.doe@dummy.com",
		Password: "TheSecretPassword5",
	}

	signUpResp1, _, err := client.SessionApi.AuthorSignUp(user1.Context, createAuthorInput1)
	common.AssertNoError(err)
	user1.UpdateWithResponse(signUpResp1.Author.Id, signUpResp1.Session)

	testlog.EndFlow(name)
}
