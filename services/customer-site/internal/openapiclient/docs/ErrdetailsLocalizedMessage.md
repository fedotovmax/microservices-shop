# ErrdetailsLocalizedMessage

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Locale** | Pointer to **string** | The locale used following the specification defined at https://www.rfc-editor.org/rfc/bcp/bcp47.txt. Examples are: \&quot;en-US\&quot;, \&quot;fr-CH\&quot;, \&quot;es-MX\&quot; | [optional] 
**Message** | Pointer to **string** | The localized error message in the above locale. | [optional] 

## Methods

### NewErrdetailsLocalizedMessage

`func NewErrdetailsLocalizedMessage() *ErrdetailsLocalizedMessage`

NewErrdetailsLocalizedMessage instantiates a new ErrdetailsLocalizedMessage object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewErrdetailsLocalizedMessageWithDefaults

`func NewErrdetailsLocalizedMessageWithDefaults() *ErrdetailsLocalizedMessage`

NewErrdetailsLocalizedMessageWithDefaults instantiates a new ErrdetailsLocalizedMessage object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLocale

`func (o *ErrdetailsLocalizedMessage) GetLocale() string`

GetLocale returns the Locale field if non-nil, zero value otherwise.

### GetLocaleOk

`func (o *ErrdetailsLocalizedMessage) GetLocaleOk() (*string, bool)`

GetLocaleOk returns a tuple with the Locale field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocale

`func (o *ErrdetailsLocalizedMessage) SetLocale(v string)`

SetLocale sets Locale field to given value.

### HasLocale

`func (o *ErrdetailsLocalizedMessage) HasLocale() bool`

HasLocale returns a boolean if a field has been set.

### GetMessage

`func (o *ErrdetailsLocalizedMessage) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *ErrdetailsLocalizedMessage) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *ErrdetailsLocalizedMessage) SetMessage(v string)`

SetMessage sets Message field to given value.

### HasMessage

`func (o *ErrdetailsLocalizedMessage) HasMessage() bool`

HasMessage returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


