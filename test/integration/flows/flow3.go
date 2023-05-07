package flows

import (
	"newsletter-manager-go/test/integration/common"
	"newsletter-manager-go/test/integration/generate/swagger"
	"newsletter-manager-go/test/integration/testlog"
	"newsletter-manager-go/test/util"
	"time"
)

func Flow3(client *swagger.APIClient) {
	var name = "flow3"
	var description = "Delete event participant"

	testlog.StartFlow(name, description)

	common.WipePostgres()
	common.MigrateUp()
	common.PopulatePostgres()

	user1 := common.NewUser()
	user2 := common.NewUser()

	signInResp1, _, err := client.SigninApi.SignInFirebase(user1.Context)
	common.AssertNoError(err)
	user1.UpdateWith(signInResp1)

	signInResp2, _, err := client.SigninApi.SignInFirebase(user2.Context)
	common.AssertNoError(err)
	user2.UpdateWith(signInResp2)

	eventRequest1 := swagger.CreateEventReq{
		Title:         "Event1",
		Category:      util.Ptr(swagger.MEETUP_EventCategory),
		Location:      "",
		Size:          util.Ptr(swagger.SIZE_TO10_EventSize),
		ImageId:       "123e4567-e89b-12d3-a456-426614174000",
		StaticImageId: "",
		StartTime:     util.ValueWithoutError(time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")),
		EndTime:       util.ValueWithoutError(time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")),
		GameType:      util.Ptr(swagger.LET_ME_GUESS_GameType),
		QuestionIds:   []string{"76f05f80-7a83-4bfe-a280-13947533757d", "f79ed307-7b42-441c-9ece-eeb0a2d35dc0"},
	}

	event1, _, err := client.EventApi.CreateEvent(user1.Context, eventRequest1)
	common.AssertNoError(err)

	eventJoinReq := swagger.EventJoinReq{OptionIds: []string{"bc874ca8-7124-4000-bd99-b6eb501a9157"}}

	_, err = client.EventApi.JoinEvent(user2.Context, eventJoinReq, event1.Id)
	common.AssertNoError(err)

	_, err = client.EventApi.DeleteParticipant(user1.Context, event1.Id, user2.AuthorID)
	common.AssertNoError(err)

	list, _, err := client.UserApi.ListParticipatingEvents(user2.Context)
	common.AssertNoError(err)

	common.Assert(len(list) == 0, "The list should be empty")

	testlog.EndFlow(name)
}
