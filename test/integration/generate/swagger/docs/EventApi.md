# {{classname}}

All URIs are relative to *https://api.dev.svc.icebreaker.strv.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CancelEvent**](EventApi.md#CancelEvent) | **Post** /api/v1/events/{eventId}/cancel | Cancel event. Event \&quot;soft\&quot; deletion.
[**CreateEvent**](EventApi.md#CreateEvent) | **Post** /api/v1/events | Create a new event as an organizer
[**DeleteParticipant**](EventApi.md#DeleteParticipant) | **Delete** /api/v1/events/{eventId}/participants/{userId} | Delete a participant from event
[**GetEvent**](EventApi.md#GetEvent) | **Get** /api/v1/events/{eventId} | Read an existing event
[**GetEventByCode**](EventApi.md#GetEventByCode) | **Get** /api/v1/events/code/{eventCode} | Get an event by its code for joining
[**GetEventCategoryImages**](EventApi.md#GetEventCategoryImages) | **Get** /api/v1/event-categories/{eventCategory}/images | (TODO) Get possible images for particular event category
[**GetEventLeaderboard**](EventApi.md#GetEventLeaderboard) | **Get** /api/v1/events/{eventId}/leaderboard | Get the leaderboard of the event&#x27;s game
[**GetEventQuestions**](EventApi.md#GetEventQuestions) | **Get** /api/v1/events/{eventId}/questions | Get questions for an event that organizer has selected
[**GetStaticQuestions**](EventApi.md#GetStaticQuestions) | **Get** /api/v1/events/questions | Get static questions for all event categories
[**GetStaticQuestionsWithOptions**](EventApi.md#GetStaticQuestionsWithOptions) | **Get** /api/v1/events/questions-with-options | Get static questions for all event categories with their options
[**JoinEvent**](EventApi.md#JoinEvent) | **Post** /api/v1/events/{eventId}/join | Join an event as a participant
[**LeaveEvent**](EventApi.md#LeaveEvent) | **Post** /api/v1/events/{eventId}/leave | Leave an event as a participant
[**ListEventParticipants**](EventApi.md#ListEventParticipants) | **Get** /api/v1/events/{eventId}/participants | List of event participants
[**ReportGameRoundResult**](EventApi.md#ReportGameRoundResult) | **Post** /api/v1/events/{eventId}/score | Report a result of one game round
[**StartEventImageUpload**](EventApi.md#StartEventImageUpload) | **Post** /api/v1/events/upload-image | Initialize upload of event image
[**UpdateEvent**](EventApi.md#UpdateEvent) | **Patch** /api/v1/events/{eventId} | Update an existing event

# **CancelEvent**
> CancelEvent(ctx, eventId)
Cancel event. Event \"soft\" deletion.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **eventId** | [**string**](.md)| ID of the event | 

### Return type

 (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CreateEvent**
> Event CreateEvent(ctx, body)
Create a new event as an organizer

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**CreateEventReq**](CreateEventReq.md)|  | 

### Return type

[**Event**](Event.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteParticipant**
> DeleteParticipant(ctx, eventId, userId)
Delete a participant from event

Permanently delete a participant and all related records from event.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **eventId** | [**string**](.md)| ID of the event | 
  **userId** | [**string**](.md)| ID of the participant to delete | 

### Return type

 (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetEvent**
> Event GetEvent(ctx, eventId)
Read an existing event

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **eventId** | [**string**](.md)| ID of the event | 

### Return type

[**Event**](Event.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetEventByCode**
> EventWithStatus GetEventByCode(ctx, eventCode)
Get an event by its code for joining

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **eventCode** | **string**| Code for the event | 

### Return type

[**EventWithStatus**](EventWithStatus.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetEventCategoryImages**
> []EventImage GetEventCategoryImages(ctx, eventCategory)
(TODO) Get possible images for particular event category

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **eventCategory** | [**EventCategory**](.md)| Type of the event | 

### Return type

[**[]EventImage**](array.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetEventLeaderboard**
> []Score GetEventLeaderboard(ctx, eventId)
Get the leaderboard of the event's game

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **eventId** | [**string**](.md)| ID of the event | 

### Return type

[**[]Score**](array.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetEventQuestions**
> []EventQuestionWithOption GetEventQuestions(ctx, eventId)
Get questions for an event that organizer has selected

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **eventId** | [**string**](.md)| ID of the event | 

### Return type

[**[]EventQuestionWithOption**](array.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetStaticQuestions**
> []EventQuestion GetStaticQuestions(ctx, )
Get static questions for all event categories

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**[]EventQuestion**](array.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetStaticQuestionsWithOptions**
> []EventQuestionWithOption GetStaticQuestionsWithOptions(ctx, )
Get static questions for all event categories with their options

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**[]EventQuestionWithOption**](array.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **JoinEvent**
> JoinEvent(ctx, body, eventId)
Join an event as a participant

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**EventJoinReq**](EventJoinReq.md)|  | 
  **eventId** | [**string**](.md)| ID of the event | 

### Return type

 (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **LeaveEvent**
> LeaveEvent(ctx, eventId)
Leave an event as a participant

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **eventId** | [**string**](.md)| ID of the event | 

### Return type

 (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListEventParticipants**
> []ParticipantWithAnswers ListEventParticipants(ctx, eventId)
List of event participants

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **eventId** | [**string**](.md)| ID of the event | 

### Return type

[**[]ParticipantWithAnswers**](array.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ReportGameRoundResult**
> ReportGameRoundResult(ctx, body, eventId)
Report a result of one game round

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**ScoreReq**](ScoreReq.md)|  | 
  **eventId** | [**string**](.md)| ID of the event | 

### Return type

 (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **StartEventImageUpload**
> UploadImageResp StartEventImageUpload(ctx, )
Initialize upload of event image

Create an upload link that the event image can be uploaded to. After uploading an image to returned URL, client should call `PATCH /api/v1/events/{eventId}` with ID of uploaded image. (ID is also returned from this endpoint). It is also possible to upload the event image prior to the event creation. In this case `POST /api/v1/events` also accepts ID of the uploaded image.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**UploadImageResp**](UploadImageResp.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateEvent**
> Event UpdateEvent(ctx, body, eventId)
Update an existing event

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**UpdateEventReq**](UpdateEventReq.md)|  | 
  **eventId** | [**string**](.md)| ID of the event | 

### Return type

[**Event**](Event.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

