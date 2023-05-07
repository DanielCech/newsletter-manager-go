# {{classname}}

All URIs are relative to *https://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AuthorSignIn**](AuthorApi.md#AuthorSignIn) | **Post** /api/v1/authors/sign-in | 
[**AuthorSignUp**](AuthorApi.md#AuthorSignUp) | **Post** /api/v1/authors/sign-up | 
[**ChangeAuthorsPassword**](AuthorApi.md#ChangeAuthorsPassword) | **Post** /api/v1/authors/change-password | 
[**DeleteAuthor**](AuthorApi.md#DeleteAuthor) | **Delete** /api/v1/authors/{authorId} | 
[**GetAuthor**](AuthorApi.md#GetAuthor) | **Get** /api/v1/authors/{authorId} | 
[**ListOfAllAuthors**](AuthorApi.md#ListOfAllAuthors) | **Get** /api/v1/authors | 
[**UpdateAuthor**](AuthorApi.md#UpdateAuthor) | **Patch** /api/v1/authors/{authorId} | 

# **AuthorSignIn**
> []Author AuthorSignIn(ctx, )


Author sign-in

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

# **AuthorSignUp**
> CreateAuthorResp AuthorSignUp(ctx, optional)


Author sign-up

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***AuthorApiAuthorSignUpOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a AuthorApiAuthorSignUpOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**optional.Interface of CreateAuthorInput**](CreateAuthorInput.md)|  | 

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


Change author's password

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

# **DeleteAuthor**
> []Author DeleteAuthor(ctx, authorId)


Delete the author

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

# **GetAuthor**
> []Author GetAuthor(ctx, authorId)


Get author by ID

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

# **ListOfAllAuthors**
> []Author ListOfAllAuthors(ctx, )


Returns the list of all authors

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

# **UpdateAuthor**
> []Author UpdateAuthor(ctx, authorId)


Update author data

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

