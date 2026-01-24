# UserspbUserSessionActionResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Payload** | Pointer to **map[string]interface{}** | Types that are valid to be assigned to Payload:   *UserSessionActionResponse_Deleted  *UserSessionActionResponse_EmailNotVerified  *UserSessionActionResponse_BadCredentials  *UserSessionActionResponse_Ok | [optional] 

## Methods

### NewUserspbUserSessionActionResponse

`func NewUserspbUserSessionActionResponse() *UserspbUserSessionActionResponse`

NewUserspbUserSessionActionResponse instantiates a new UserspbUserSessionActionResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUserspbUserSessionActionResponseWithDefaults

`func NewUserspbUserSessionActionResponseWithDefaults() *UserspbUserSessionActionResponse`

NewUserspbUserSessionActionResponseWithDefaults instantiates a new UserspbUserSessionActionResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPayload

`func (o *UserspbUserSessionActionResponse) GetPayload() map[string]interface{}`

GetPayload returns the Payload field if non-nil, zero value otherwise.

### GetPayloadOk

`func (o *UserspbUserSessionActionResponse) GetPayloadOk() (*map[string]interface{}, bool)`

GetPayloadOk returns a tuple with the Payload field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPayload

`func (o *UserspbUserSessionActionResponse) SetPayload(v map[string]interface{})`

SetPayload sets Payload field to given value.

### HasPayload

`func (o *UserspbUserSessionActionResponse) HasPayload() bool`

HasPayload returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


