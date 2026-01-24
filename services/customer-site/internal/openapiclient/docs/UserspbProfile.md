# UserspbProfile

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AvatarUrl** | Pointer to **string** |  | [optional] 
**BirthDate** | Pointer to **string** |  | [optional] 
**FirstName** | Pointer to **string** |  | [optional] 
**Gender** | [**UserspbGenderValue**](UserspbGenderValue.md) | GENDER_UNSPECIFIED &#x3D; 0 Reserved for Proto, not a valid value GENDER_UNSELECTED &#x3D; 1 User has not selected a gender GENDER_MALE &#x3D; 2 Represents male gender GENDER_FEMALE &#x3D; 3 Represents female gender | 
**LastName** | Pointer to **string** |  | [optional] 
**MiddleName** | Pointer to **string** |  | [optional] 
**UpdatedAt** | [**TimestamppbTimestamp**](TimestamppbTimestamp.md) |  | 

## Methods

### NewUserspbProfile

`func NewUserspbProfile(gender UserspbGenderValue, updatedAt TimestamppbTimestamp, ) *UserspbProfile`

NewUserspbProfile instantiates a new UserspbProfile object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUserspbProfileWithDefaults

`func NewUserspbProfileWithDefaults() *UserspbProfile`

NewUserspbProfileWithDefaults instantiates a new UserspbProfile object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAvatarUrl

`func (o *UserspbProfile) GetAvatarUrl() string`

GetAvatarUrl returns the AvatarUrl field if non-nil, zero value otherwise.

### GetAvatarUrlOk

`func (o *UserspbProfile) GetAvatarUrlOk() (*string, bool)`

GetAvatarUrlOk returns a tuple with the AvatarUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvatarUrl

`func (o *UserspbProfile) SetAvatarUrl(v string)`

SetAvatarUrl sets AvatarUrl field to given value.

### HasAvatarUrl

`func (o *UserspbProfile) HasAvatarUrl() bool`

HasAvatarUrl returns a boolean if a field has been set.

### GetBirthDate

`func (o *UserspbProfile) GetBirthDate() string`

GetBirthDate returns the BirthDate field if non-nil, zero value otherwise.

### GetBirthDateOk

`func (o *UserspbProfile) GetBirthDateOk() (*string, bool)`

GetBirthDateOk returns a tuple with the BirthDate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBirthDate

`func (o *UserspbProfile) SetBirthDate(v string)`

SetBirthDate sets BirthDate field to given value.

### HasBirthDate

`func (o *UserspbProfile) HasBirthDate() bool`

HasBirthDate returns a boolean if a field has been set.

### GetFirstName

`func (o *UserspbProfile) GetFirstName() string`

GetFirstName returns the FirstName field if non-nil, zero value otherwise.

### GetFirstNameOk

`func (o *UserspbProfile) GetFirstNameOk() (*string, bool)`

GetFirstNameOk returns a tuple with the FirstName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFirstName

`func (o *UserspbProfile) SetFirstName(v string)`

SetFirstName sets FirstName field to given value.

### HasFirstName

`func (o *UserspbProfile) HasFirstName() bool`

HasFirstName returns a boolean if a field has been set.

### GetGender

`func (o *UserspbProfile) GetGender() UserspbGenderValue`

GetGender returns the Gender field if non-nil, zero value otherwise.

### GetGenderOk

`func (o *UserspbProfile) GetGenderOk() (*UserspbGenderValue, bool)`

GetGenderOk returns a tuple with the Gender field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGender

`func (o *UserspbProfile) SetGender(v UserspbGenderValue)`

SetGender sets Gender field to given value.


### GetLastName

`func (o *UserspbProfile) GetLastName() string`

GetLastName returns the LastName field if non-nil, zero value otherwise.

### GetLastNameOk

`func (o *UserspbProfile) GetLastNameOk() (*string, bool)`

GetLastNameOk returns a tuple with the LastName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastName

`func (o *UserspbProfile) SetLastName(v string)`

SetLastName sets LastName field to given value.

### HasLastName

`func (o *UserspbProfile) HasLastName() bool`

HasLastName returns a boolean if a field has been set.

### GetMiddleName

`func (o *UserspbProfile) GetMiddleName() string`

GetMiddleName returns the MiddleName field if non-nil, zero value otherwise.

### GetMiddleNameOk

`func (o *UserspbProfile) GetMiddleNameOk() (*string, bool)`

GetMiddleNameOk returns a tuple with the MiddleName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMiddleName

`func (o *UserspbProfile) SetMiddleName(v string)`

SetMiddleName sets MiddleName field to given value.

### HasMiddleName

`func (o *UserspbProfile) HasMiddleName() bool`

HasMiddleName returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *UserspbProfile) GetUpdatedAt() TimestamppbTimestamp`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *UserspbProfile) GetUpdatedAtOk() (*TimestamppbTimestamp, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *UserspbProfile) SetUpdatedAt(v TimestamppbTimestamp)`

SetUpdatedAt sets UpdatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


