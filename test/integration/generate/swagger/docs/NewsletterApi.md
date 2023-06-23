# {{classname}}

All URIs are relative to *https://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AuthorsNewsletters**](NewsletterApi.md#AuthorsNewsletters) | **Get** /api/v1/authors/current/newsletters | 
[**CreateNewsletter**](NewsletterApi.md#CreateNewsletter) | **Post** /api/v1/authors/current/newsletters | 
[**DeleteNewsletter**](NewsletterApi.md#DeleteNewsletter) | **Delete** /api/v1/newsletters/{newsletterId} | 
[**GetNewsletter**](NewsletterApi.md#GetNewsletter) | **Get** /api/v1/newsletters/{newsletterId} | 
[**ListNewsletters**](NewsletterApi.md#ListNewsletters) | **Get** /api/v1/newsletters | 
[**UpdateNewsletter**](NewsletterApi.md#UpdateNewsletter) | **Patch** /api/v1/newsletters/{newsletterId} | 

# **AuthorsNewsletters**
> []Newsletter AuthorsNewsletters(ctx, )


The list of author's newsletters

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**[]Newsletter**](Newsletter.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CreateNewsletter**
> Newsletter CreateNewsletter(ctx, body)


Creates a new newsletter.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**CreateNewsletterReq**](CreateNewsletterReq.md)|  | 

### Return type

[**Newsletter**](Newsletter.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteNewsletter**
> DeleteNewsletter(ctx, newsletterId)


Delete the newsletter

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **newsletterId** | [**string**](.md)| ID of the newsletter | 

### Return type

 (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetNewsletter**
> Newsletter GetNewsletter(ctx, newsletterId)


Get newsletter by ID

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **newsletterId** | [**string**](.md)| ID of the newsletter | 

### Return type

[**Newsletter**](Newsletter.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListNewsletters**
> []Newsletter ListNewsletters(ctx, )


The list of all newsletters

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**[]Newsletter**](Newsletter.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateNewsletter**
> Newsletter UpdateNewsletter(ctx, newsletterId, optional)


Update newsletter

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **newsletterId** | [**string**](.md)| ID of the newsletter | 
 **optional** | ***NewsletterApiUpdateNewsletterOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a NewsletterApiUpdateNewsletterOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **body** | [**optional.Interface of CreateNewsletterReq**](CreateNewsletterReq.md)|  | 

### Return type

[**Newsletter**](Newsletter.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json; charset=utf-8
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

