package flows

import (
	"event-facematch-backend/test/integration/common"
	"event-facematch-backend/test/integration/generate/swagger"
	"event-facematch-backend/test/integration/testlog"
	"event-facematch-backend/test/util"
	"time"
)

func Flow2(client *swagger.APIClient) {
	var name = "flow2"
	var description = "Delete user"

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

	eventRequest2 := swagger.CreateEventReq{
		Title:         "Event2",
		Category:      util.Ptr(swagger.SPORT_EventCategory),
		Location:      "",
		Size:          util.Ptr(swagger.SIZE10_TO30_EventSize),
		ImageId:       "123e4567-e89b-12d3-a456-426614174000",
		StaticImageId: "",
		StartTime:     util.ValueWithoutError(time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")),
		EndTime:       util.ValueWithoutError(time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")),
		GameType:      util.Ptr(swagger.LET_ME_GUESS_GameType),
		QuestionIds:   []string{"76f05f80-7a83-4bfe-a280-13947533757d", "f79ed307-7b42-441c-9ece-eeb0a2d35dc0"},
	}

	event1, _, err := client.EventApi.CreateEvent(user1.Context, eventRequest1)
	common.AssertNoError(err)

	event2, _, err := client.EventApi.CreateEvent(user2.Context, eventRequest2)
	common.AssertNoError(err)

	_, err = client.UserApi.DeleteUser(user1.Context)
	common.AssertNoError(err)

	event1, _, err = client.EventApi.GetEvent(user2.Context, event1.Id)
	common.AssertNoError(err)

	common.Assert(event1.IsCanceled, "Event1 should be canceled")

	event2, _, err = client.EventApi.GetEvent(user2.Context, event2.Id)
	common.AssertNoError(err)

	common.Assert(!event2.IsCanceled, "Event2 should not be canceled")

	testlog.EndFlow(name)
}
