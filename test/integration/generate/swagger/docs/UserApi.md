# {{classname}}

All URIs are relative to *https://api.dev.svc.icebreaker.strv.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CurrentUserInfo**](UserApi.md#CurrentUserInfo) | **Get** /api/v1/users/current | Get information about the current user
[**DeleteUser**](UserApi.md#DeleteUser) | **Delete** /api/v1/users/current | Delete logged in user
[**ListOrganizingEvents**](UserApi.md#ListOrganizingEvents) | **Get** /api/v1/users/current/organizing-events | List events the logged in user is organizing
[**ListParticipatingEvents**](UserApi.md#ListParticipatingEvents) | **Get** /api/v1/users/current/participating-events | List events the logged in user is participating in
[**StartUserImageUpload**](UserApi.md#StartUserImageUpload) | **Post** /api/v1/users/current/upload-image | Initialize image upload
[**UpdateUser**](UserApi.md#UpdateUser) | **Patch** /api/v1/users/current | Update info of logged in user

# **CurrentUserInfo**
> User CurrentUserInfo(ctx, )
Get information about the current user

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**User**](User.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteUser**
> DeleteUser(ctx, )
Delete logged in user

Permanently delete currently logged in user. Removes them from all events.

### Required Parameters
This endpoint does not need any parameter.

### Return type

 (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListOrganizingEvents**
> []Event ListOrganizingEvents(ctx, )
List events the logged in user is organizing

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**[]Event**](array.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListParticipatingEvents**
> []Event ListParticipatingEvents(ctx, )
List events the logged in user is participating in

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**[]Event**](array.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **StartUserImageUpload**
> UploadImageResp StartUserImageUpload(ctx, )
Initialize image upload

Create an upload link that the user profile image can be uploaded to. After uploading an image to returned URL, client should call `PATCH /api/v1/users/current` with ID of uploaded image (ID is also returned from this endpoint).

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

# **UpdateUser**
> UserUpdateResp UpdateUser(ctx, body)
Update info of logged in user

Updates the user's info, also setting their `finalized` flag.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**UserUpdateReq**](UserUpdateReq.md)|  | 

### Return type

[**UserUpdateResp**](UserUpdateResp.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

