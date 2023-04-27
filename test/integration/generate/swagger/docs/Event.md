# Event

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | Event ID | [default to null]
**Title** | **string** | Title name of the event | [default to null]
**Category** | [***EventCategory**](EventCategory.md) |  | [default to null]
**Size** | [***EventSize**](EventSize.md) |  | [default to null]
**Location** | **string** | The optional location or link | [default to null]
**ImageURL** | **string** | The event image URL | [optional] [default to null]
**StaticImageId** | **string** | ID of static image on frontend (only one of imageId and staticImageId should be non-null) | [optional] [default to null]
**StartTime** | [**time.Time**](time.Time.md) | Start time of the real-life event (not the game) | [default to null]
**EndTime** | [**time.Time**](time.Time.md) | End time of the real-life event (not the game) | [default to null]
**Code** | **string** | Code for joining and generating share dynamic links. Only available if the requesting user is an organizer | [optional] [default to null]
**IsCanceled** | **bool** | The event has been canceled. | [optional] [default to null]
**ShareLink** | **string** | Dynamic links for sharing. Only available if the requesting user is an organizer | [optional] [default to null]
**GameType** | [***GameType**](GameType.md) |  | [default to null]
**Organizer** | [***Organizer**](Organizer.md) |  | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

