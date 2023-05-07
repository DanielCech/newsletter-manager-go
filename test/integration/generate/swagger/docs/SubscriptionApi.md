# {{classname}}

All URIs are relative to *https://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**NewsletterSubscriptions**](SubscriptionApi.md#NewsletterSubscriptions) | **Get** /api/v1/newsletters/{newsletterId}/subscriptions | 
[**SubscribeToNewsletter**](SubscriptionApi.md#SubscribeToNewsletter) | **Post** /api/v1/newsletters/{newsletterId}/subscribe | 
[**SubscriptionsByEmail**](SubscriptionApi.md#SubscriptionsByEmail) | **Get** /api/v1/subscriptions | 
[**UnsubscribeFromNewsletter**](SubscriptionApi.md#UnsubscribeFromNewsletter) | **Post** /api/v1/newsletters/{newsletterId}/unsubscribe | 

# **NewsletterSubscriptions**
> CreateAuthorResp NewsletterSubscriptions(ctx, newsletterId)


Newsletter's subscriptions

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

# **SubscribeToNewsletter**
> []Author SubscribeToNewsletter(ctx, newsletterId)


Subscribe to newsletter

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **newsletterId** | [**string**](.md)| ID of the author | 

### Return type

[**[]Author**](Author.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SubscriptionsByEmail**
> CreateAuthorResp SubscriptionsByEmail(ctx, )


All newsletter subscriptions by Email

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**CreateAuthorResp**](CreateAuthorResp.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UnsubscribeFromNewsletter**
> []Author UnsubscribeFromNewsletter(ctx, newsletterId)


Unsubscribe from newsletter

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **newsletterId** | [**string**](.md)| ID of the author | 

### Return type

[**[]Author**](Author.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

