# {{classname}}

All URIs are relative to *https://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DeleteAuthor**](AuthorApi.md#DeleteAuthor) | **Delete** /api/v1/authors/{authorId} | 
[**GetAuthor**](AuthorApi.md#GetAuthor) | **Get** /api/v1/authors/{authorId} | 
[**ListOfAllAuthors**](AuthorApi.md#ListOfAllAuthors) | **Get** /api/v1/authors | 
[**UpdateAuthor**](AuthorApi.md#UpdateAuthor) | **Patch** /api/v1/authors/{authorId} | 

# **DeleteAuthor**
> []Author DeleteAuthor(ctx, authorId)


(TODO:) Delete the author

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


(TODO:) Get author by ID

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


(TODO:) Returns the list of all authors

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


(TODO:) Update author data

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

