# UpdateEventReq

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Title** | **string** | Title name of the event | [optional] [default to null]
**Category** | [***EventCategory**](EventCategory.md) |  | [optional] [default to null]
**Location** | **string** | The location of an event or the link | [optional] [default to null]
**Size** | [***EventSize**](EventSize.md) |  | [optional] [default to null]
**ImageId** | **string** | ID of uploaded image (only one of imageId and staticImageId should be non-null). Use it only if you already have ID of uploaded event image. | [optional] [default to null]
**StaticImageId** | **string** | ID of static image on frontend (only one of imageId and staticImageId should be non-null) | [optional] [default to null]
**StartTime** | [**time.Time**](time.Time.md) | Start time of the real-life event (not the game) | [optional] [default to null]
**EndTime** | [**time.Time**](time.Time.md) | End time of the real-life event (not the game) | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
