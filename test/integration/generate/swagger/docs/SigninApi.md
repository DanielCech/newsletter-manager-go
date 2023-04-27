# {{classname}}

All URIs are relative to *https://api.dev.svc.icebreaker.strv.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**SignInFirebase**](SigninApi.md#SignInFirebase) | **Post** /api/v1/signin/firebase | Sign in via Firebase

# **SignInFirebase**
> SignInResp SignInFirebase(ctx, )
Sign in via Firebase

Sign in using Firebase ID token. Can use Google or Apple OAuth provider, or anonymous provider. If there is no user account associated with the OAuth provider account, the user is created and their name, email and image url are prefilled from the provider. Returns the existing user if it (associated OAuth provider account) exists.  The ID token should be refreshed after this for further api calls, as this endpoint adds custom claims to it that other endpoints need. Use `auth.currentUser.getIdToken(true)` to force refreshing the ID token.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**SignInResp**](SignInResp.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json; charset=utf-8

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

