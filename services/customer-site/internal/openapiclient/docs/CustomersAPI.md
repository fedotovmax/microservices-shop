# \CustomersAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CustomersSessionLoginPost**](CustomersAPI.md#CustomersSessionLoginPost) | **Post** /customers/session/login | Login in account
[**CustomersSessionRefreshSessionPost**](CustomersAPI.md#CustomersSessionRefreshSessionPost) | **Post** /customers/session/refresh-session | Refresh session
[**CustomersUsersPost**](CustomersAPI.md#CustomersUsersPost) | **Post** /customers/users | Create user account
[**CustomersUsersProfileGet**](CustomersAPI.md#CustomersUsersProfileGet) | **Get** /customers/users/profile | Get user profile
[**CustomersUsersProfilePatch**](CustomersAPI.md#CustomersUsersProfilePatch) | **Patch** /customers/users/profile | Update user profile



## CustomersSessionLoginPost

> SessionspbCreateSessionResponse CustomersSessionLoginPost(ctx).Dto(dto).XRequestLocale(xRequestLocale).Execute()

Login in account



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	dto := *openapiclient.NewGithubComFedotovmaxMicroservicesShopApiGatewayInternalDomainLoginInput("makc-dgek@mail.ru", "19.56.186.122", "Password_example", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36") // GithubComFedotovmaxMicroservicesShopApiGatewayInternalDomainLoginInput | Dto for login in account
	xRequestLocale := "xRequestLocale_example" // string | Locale (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CustomersAPI.CustomersSessionLoginPost(context.Background()).Dto(dto).XRequestLocale(xRequestLocale).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CustomersAPI.CustomersSessionLoginPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CustomersSessionLoginPost`: SessionspbCreateSessionResponse
	fmt.Fprintf(os.Stdout, "Response from `CustomersAPI.CustomersSessionLoginPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCustomersSessionLoginPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **dto** | [**GithubComFedotovmaxMicroservicesShopApiGatewayInternalDomainLoginInput**](GithubComFedotovmaxMicroservicesShopApiGatewayInternalDomainLoginInput.md) | Dto for login in account | 
 **xRequestLocale** | **string** | Locale | 

### Return type

[**SessionspbCreateSessionResponse**](SessionspbCreateSessionResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CustomersSessionRefreshSessionPost

> SessionspbCreateSessionResponse CustomersSessionRefreshSessionPost(ctx).Dto(dto).XRequestLocale(xRequestLocale).Execute()

Refresh session



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	dto := *openapiclient.NewSessionspbRefreshSessionRequest("19.56.186.122", "RefreshToken_example", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36") // SessionspbRefreshSessionRequest | Refresh session dto
	xRequestLocale := "xRequestLocale_example" // string | Locale (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CustomersAPI.CustomersSessionRefreshSessionPost(context.Background()).Dto(dto).XRequestLocale(xRequestLocale).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CustomersAPI.CustomersSessionRefreshSessionPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CustomersSessionRefreshSessionPost`: SessionspbCreateSessionResponse
	fmt.Fprintf(os.Stdout, "Response from `CustomersAPI.CustomersSessionRefreshSessionPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCustomersSessionRefreshSessionPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **dto** | [**SessionspbRefreshSessionRequest**](SessionspbRefreshSessionRequest.md) | Refresh session dto | 
 **xRequestLocale** | **string** | Locale | 

### Return type

[**SessionspbCreateSessionResponse**](SessionspbCreateSessionResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CustomersUsersPost

> UserspbCreateUserResponse CustomersUsersPost(ctx).Dto(dto).XRequestLocale(xRequestLocale).Execute()

Create user account



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	dto := *openapiclient.NewUserspbCreateUserRequest("makc-dgek@mail.ru", "Password_example") // UserspbCreateUserRequest | Create user account with body dto
	xRequestLocale := "xRequestLocale_example" // string | Locale (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CustomersAPI.CustomersUsersPost(context.Background()).Dto(dto).XRequestLocale(xRequestLocale).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CustomersAPI.CustomersUsersPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CustomersUsersPost`: UserspbCreateUserResponse
	fmt.Fprintf(os.Stdout, "Response from `CustomersAPI.CustomersUsersPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCustomersUsersPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **dto** | [**UserspbCreateUserRequest**](UserspbCreateUserRequest.md) | Create user account with body dto | 
 **xRequestLocale** | **string** | Locale | 

### Return type

[**UserspbCreateUserResponse**](UserspbCreateUserResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CustomersUsersProfileGet

> UserspbUser CustomersUsersProfileGet(ctx).XRequestLocale(xRequestLocale).Execute()

Get user profile



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	xRequestLocale := "xRequestLocale_example" // string | Locale (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CustomersAPI.CustomersUsersProfileGet(context.Background()).XRequestLocale(xRequestLocale).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CustomersAPI.CustomersUsersProfileGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CustomersUsersProfileGet`: UserspbUser
	fmt.Fprintf(os.Stdout, "Response from `CustomersAPI.CustomersUsersProfileGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCustomersUsersProfileGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xRequestLocale** | **string** | Locale | 

### Return type

[**UserspbUser**](UserspbUser.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CustomersUsersProfilePatch

> HttputilsErrorResponse CustomersUsersProfilePatch(ctx).Dto(dto).XRequestLocale(xRequestLocale).Execute()

Update user profile



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	dto := *openapiclient.NewUserspbUpdateUserProfileData() // UserspbUpdateUserProfileData | Update user profile with body dto
	xRequestLocale := "xRequestLocale_example" // string | Locale (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CustomersAPI.CustomersUsersProfilePatch(context.Background()).Dto(dto).XRequestLocale(xRequestLocale).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CustomersAPI.CustomersUsersProfilePatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CustomersUsersProfilePatch`: HttputilsErrorResponse
	fmt.Fprintf(os.Stdout, "Response from `CustomersAPI.CustomersUsersProfilePatch`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCustomersUsersProfilePatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **dto** | [**UserspbUpdateUserProfileData**](UserspbUpdateUserProfileData.md) | Update user profile with body dto | 
 **xRequestLocale** | **string** | Locale | 

### Return type

[**HttputilsErrorResponse**](HttputilsErrorResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

