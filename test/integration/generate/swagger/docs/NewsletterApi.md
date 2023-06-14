# {{classname}}

All URIs are relative to *https://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AuthorsNewsletters**](NewsletterApi.md#AuthorsNewsletters) | **Get** /api/v1/authors/current/newsletters | 
[**CreateNewsletter**](NewsletterApi.md#CreateNewsletter) | **Post** /api/v1/authors/current/newsletters | 
[**DeleteNewsletter**](NewsletterApi.md#DeleteNewsletter) | **Delete** /api/v1/newsletters/{newsletterId} | 
[**GetNewsletter**](NewsletterApi.md#GetNewsletter) | **Get** /api/v1/newsletters/{newsletterId} | 
[**UpdateNewsletter**](NewsletterApi.md#UpdateNewsletter) | **Patch** /api/v1/newsletters/{newsletterId} | 

# **AuthorsNewsletters**
> []Author AuthorsNewsletters(ctx, authorId)


(TODO:) The list of author's newsletters

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **authorId** | [**string**](.md)| ID of the author | 

### Return type

[**[]Author**](Author.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CreateNewsletter**
> []Author CreateNewsletter(ctx, authorId)


(TODO:) Creates a new newsletter.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **authorId** | [**string**](.md)| ID of the author | 

### Return type

[**[]Author**](Author.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteNewsletter**
> CreateAuthorResp DeleteNewsletter(ctx, newsletterId)


(TODO:) Delete the newsletter

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **newsletterId** | [**string**](.md)| ID of the author | 

### Return type

[**CreateAuthorResp**](CreateAuthorResp.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetNewsletter**
> CreateAuthorResp GetNewsletter(ctx, newsletterId)


(TODO:) Get newsletter by ID

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **newsletterId** | [**string**](.md)| ID of the author | 

### Return type

[**CreateAuthorResp**](CreateAuthorResp.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateNewsletter**
> CreateAuthorResp UpdateNewsletter(ctx, newsletterId, optional)


(TODO:) Update newsletter

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **newsletterId** | [**string**](.md)| ID of the author | 
 **optional** | ***NewsletterApiUpdateNewsletterOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a NewsletterApiUpdateNewsletterOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **body** | [**optional.Interface of CreateAuthorInput**](CreateAuthorInput.md)|  | 

### Return type

[**CreateAuthorResp**](CreateAuthorResp.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json; charset=utf-8
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

