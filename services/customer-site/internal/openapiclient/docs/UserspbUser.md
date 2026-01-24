# UserspbUser

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CreatedAt** | [**TimestamppbTimestamp**](TimestamppbTimestamp.md) |  | 
**Email** | **string** |  | 
**Id** | **string** |  | 
**Phone** | Pointer to **string** |  | [optional] 
**Profile** | [**UserspbProfile**](UserspbProfile.md) |  | 
**UpdatedAt** | [**TimestamppbTimestamp**](TimestamppbTimestamp.md) |  | 

## Methods

### NewUserspbUser

`func NewUserspbUser(createdAt TimestamppbTimestamp, email string, id string, profile UserspbProfile, updatedAt TimestamppbTimestamp, ) *UserspbUser`

NewUserspbUser instantiates a new UserspbUser object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUserspbUserWithDefaults

`func NewUserspbUserWithDefaults() *UserspbUser`

NewUserspbUserWithDefaults instantiates a new UserspbUser object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCreatedAt

`func (o *UserspbUser) GetCreatedAt() TimestamppbTimestamp`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *UserspbUser) GetCreatedAtOk() (*TimestamppbTimestamp, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *UserspbUser) SetCreatedAt(v TimestamppbTimestamp)`

SetCreatedAt sets CreatedAt field to given value.


### GetEmail

`func (o *UserspbUser) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *UserspbUser) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *UserspbUser) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetId

`func (o *UserspbUser) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *UserspbUser) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *UserspbUser) SetId(v string)`

SetId sets Id field to given value.


### GetPhone

`func (o *UserspbUser) GetPhone() string`

GetPhone returns the Phone field if non-nil, zero value otherwise.

### GetPhoneOk

`func (o *UserspbUser) GetPhoneOk() (*string, bool)`

GetPhoneOk returns a tuple with the Phone field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPhone

`func (o *UserspbUser) SetPhone(v string)`

SetPhone sets Phone field to given value.

### HasPhone

`func (o *UserspbUser) HasPhone() bool`

HasPhone returns a boolean if a field has been set.

### GetProfile

`func (o *UserspbUser) GetProfile() UserspbProfile`

GetProfile returns the Profile field if non-nil, zero value otherwise.

### GetProfileOk

`func (o *UserspbUser) GetProfileOk() (*UserspbProfile, bool)`

GetProfileOk returns a tuple with the Profile field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProfile

`func (o *UserspbUser) SetProfile(v UserspbProfile)`

SetProfile sets Profile field to given value.


### GetUpdatedAt

`func (o *UserspbUser) GetUpdatedAt() TimestamppbTimestamp`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *UserspbUser) GetUpdatedAtOk() (*TimestamppbTimestamp, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *UserspbUser) SetUpdatedAt(v TimestamppbTimestamp)`

SetUpdatedAt sets UpdatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


