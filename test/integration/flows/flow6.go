package flows

import (
	"newsletter-manager-go/test/integration/common"
	"newsletter-manager-go/test/integration/generate/swagger"
	"newsletter-manager-go/test/integration/testlog"
)

func Flow6(client *swagger.APIClient) {
	var name = "flow6"
	var description = "Create newsletter"

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

	createNewsletterReq := swagger.CreateNewsletterReq{
		Name:        "Newsletter's name",
		Description: "Newsletter' description",
	}

	_, _, err = client.NewsletterApi.CreateNewsletter(user1.Context, createNewsletterReq)
	common.AssertNoError(err)

	// Read logged user
	_, _, err = client.SessionApi.GetCurrentAuthor(user1.Context)
	common.Assert(err != nil, "Author should be deleted")

	testlog.EndFlow(name)
}
