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

	createdNewsletter, _, err := client.NewsletterApi.CreateNewsletter(user1.Context, createNewsletterReq)
	common.AssertNoError(err)

	readNewsletter, _, err := client.NewsletterApi.GetNewsletter(user1.Context, createdNewsletter.Id)
	common.AssertNoError(err)

	common.Assert(createdNewsletter.Id == readNewsletter.Id, "Newsletter ids must be the same")
	common.Assert(createdNewsletter.Name == readNewsletter.Name, "Newsletter names must be the same")
	common.Assert(createdNewsletter.Description == readNewsletter.Description, "Newsletter descriptions must be the same")

	newsletters, _, err := client.NewsletterApi.AuthorsNewsletters(user1.Context)
	common.AssertNoError(err)
	common.Assert(len(newsletters) == 1, "There should be exactly one newsletter")

	// updateNewsletterReq := swagger.UpdateNewsletterReq{
	// 	Name:        "Newsletter's name",
	// 	Description: "Newsletter' description",
	// }

	// // client.NewsletterApi.UpdateNewsletter(user1.Context, createdNewsletter.Id, updateNewsletterReq)

	list, _, err := client.NewsletterApi.ListNewsletters(user1.Context)
	common.AssertNoError(err)
	common.Assert(len(list) == 1, "There should be exactly one newsletter")

	testlog.EndFlow(name)
}
