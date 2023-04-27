# User

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | User ID | [default to null]
**Name** | **string** | User full name | [default to null]
**Instagram** | **string** | User&#x27;s instagram | [optional] [default to null]
**Linkedin** | **string** | User&#x27;s linkedin | [optional] [default to null]
**ImageUrl** | **string** | link to user profile image | [default to null]
**Finalized** | **bool** | Whether the user finished registration. False for newly created users, true if update &#x60;PUT /api/v1/users/current&#x60; endpoint was called for existing user (i.e. they confirmed their profile info). | [default to null]
**StrongAuth** | **bool** | Whether the user is using useful authentication provider, which enables them to be organizer of events. False for anonymous (guest) users. | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

