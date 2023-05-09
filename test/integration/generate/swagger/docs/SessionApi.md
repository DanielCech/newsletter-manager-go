# {{classname}}

All URIs are relative to *https://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AuthorSignUp**](SessionApi.md#AuthorSignUp) | **Post** /api/v1/authors/sign-up | 
[**ChangeAuthorsPassword**](SessionApi.md#ChangeAuthorsPassword) | **Post** /api/v1/authors/current/change-password | 
[**CreateSession**](SessionApi.md#CreateSession) | **Post** /api/v1/authors/sign-in | 
[**DestroySession**](SessionApi.md#DestroySession) | **Post** /api/v1/authors/current/logout | 
[**RefreshSession**](SessionApi.md#RefreshSession) | **Post** /api/v1/authors/current/refresh-token | 

# **AuthorSignUp**
> CreateAuthorResp AuthorSignUp(ctx, body)


Author sign-up

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**CreateAuthorInput**](CreateAuthorInput.md)|  | 

### Return type

[**CreateAuthorResp**](CreateAuthorResp.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ChangeAuthorsPassword**
> []Author ChangeAuthorsPassword(ctx, )


(TODO:) Change author's password

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**[]Author**](Author.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CreateSession**
> CreateSessionResp CreateSession(ctx, optional)


User sign-in

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***SessionApiCreateSessionOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a SessionApiCreateSessionOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**optional.Interface of CreateSessionInput**](CreateSessionInput.md)|  | 

### Return type

[**CreateSessionResp**](CreateSessionResp.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json; charset=utf-8
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DestroySession**
> DestroySession(ctx, optional)


(TODO:) Logout author completely

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***SessionApiDestroySessionOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a SessionApiDestroySessionOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**optional.Interface of DestroySessionInput**](DestroySessionInput.md)|  | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json; charset=utf-8
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **RefreshSession**
> Session RefreshSession(ctx, optional)


(TODO:) Refresh access token

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***SessionApiRefreshSessionOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a SessionApiRefreshSessionOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**optional.Interface of RefreshSessionInput**](RefreshSessionInput.md)|  | 

### Return type

[**Session**](Session.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json; charset=utf-8
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

