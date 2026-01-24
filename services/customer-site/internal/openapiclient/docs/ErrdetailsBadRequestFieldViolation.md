# ErrdetailsBadRequestFieldViolation

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Description** | Pointer to **string** | A description of why the request element is bad. | [optional] 
**Field** | Pointer to **string** | A path that leads to a field in the request body. The value will be a sequence of dot-separated identifiers that identify a protocol buffer field.  Consider the following:   message CreateContactRequest {    message EmailAddress {      enum Type {        TYPE_UNSPECIFIED &#x3D; 0;        HOME &#x3D; 1;        WORK &#x3D; 2;      }       optional string email &#x3D; 1;      repeated EmailType type &#x3D; 2;    }     string full_name &#x3D; 1;    repeated EmailAddress email_addresses &#x3D; 2;  }  In this example, in proto &#x60;field&#x60; could take one of the following values:    - &#x60;full_name&#x60; for a violation in the &#x60;full_name&#x60; value   - &#x60;email_addresses[1].email&#x60; for a violation in the &#x60;email&#x60; field of the     first &#x60;email_addresses&#x60; message   - &#x60;email_addresses[3].type[2]&#x60; for a violation in the second &#x60;type&#x60;     value in the third &#x60;email_addresses&#x60; message.  In JSON, the same values are represented as:    - &#x60;fullName&#x60; for a violation in the &#x60;fullName&#x60; value   - &#x60;emailAddresses[1].email&#x60; for a violation in the &#x60;email&#x60; field of the     first &#x60;emailAddresses&#x60; message   - &#x60;emailAddresses[3].type[2]&#x60; for a violation in the second &#x60;type&#x60;     value in the third &#x60;emailAddresses&#x60; message. | [optional] 
**LocalizedMessage** | Pointer to [**ErrdetailsLocalizedMessage**](ErrdetailsLocalizedMessage.md) | Provides a localized error message for field-level errors that is safe to return to the API consumer. | [optional] 
**Reason** | Pointer to **string** | The reason of the field-level error. This is a constant value that identifies the proximate cause of the field-level error. It should uniquely identify the type of the FieldViolation within the scope of the google.rpc.ErrorInfo.domain. This should be at most 63 characters and match a regular expression of &#x60;[A-Z][A-Z0-9_]+[A-Z0-9]&#x60;, which represents UPPER_SNAKE_CASE. | [optional] 

## Methods

### NewErrdetailsBadRequestFieldViolation

`func NewErrdetailsBadRequestFieldViolation() *ErrdetailsBadRequestFieldViolation`

NewErrdetailsBadRequestFieldViolation instantiates a new ErrdetailsBadRequestFieldViolation object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewErrdetailsBadRequestFieldViolationWithDefaults

`func NewErrdetailsBadRequestFieldViolationWithDefaults() *ErrdetailsBadRequestFieldViolation`

NewErrdetailsBadRequestFieldViolationWithDefaults instantiates a new ErrdetailsBadRequestFieldViolation object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDescription

`func (o *ErrdetailsBadRequestFieldViolation) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *ErrdetailsBadRequestFieldViolation) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *ErrdetailsBadRequestFieldViolation) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *ErrdetailsBadRequestFieldViolation) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetField

`func (o *ErrdetailsBadRequestFieldViolation) GetField() string`

GetField returns the Field field if non-nil, zero value otherwise.

### GetFieldOk

`func (o *ErrdetailsBadRequestFieldViolation) GetFieldOk() (*string, bool)`

GetFieldOk returns a tuple with the Field field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetField

`func (o *ErrdetailsBadRequestFieldViolation) SetField(v string)`

SetField sets Field field to given value.

### HasField

`func (o *ErrdetailsBadRequestFieldViolation) HasField() bool`

HasField returns a boolean if a field has been set.

### GetLocalizedMessage

`func (o *ErrdetailsBadRequestFieldViolation) GetLocalizedMessage() ErrdetailsLocalizedMessage`

GetLocalizedMessage returns the LocalizedMessage field if non-nil, zero value otherwise.

### GetLocalizedMessageOk

`func (o *ErrdetailsBadRequestFieldViolation) GetLocalizedMessageOk() (*ErrdetailsLocalizedMessage, bool)`

GetLocalizedMessageOk returns a tuple with the LocalizedMessage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocalizedMessage

`func (o *ErrdetailsBadRequestFieldViolation) SetLocalizedMessage(v ErrdetailsLocalizedMessage)`

SetLocalizedMessage sets LocalizedMessage field to given value.

### HasLocalizedMessage

`func (o *ErrdetailsBadRequestFieldViolation) HasLocalizedMessage() bool`

HasLocalizedMessage returns a boolean if a field has been set.

### GetReason

`func (o *ErrdetailsBadRequestFieldViolation) GetReason() string`

GetReason returns the Reason field if non-nil, zero value otherwise.

### GetReasonOk

`func (o *ErrdetailsBadRequestFieldViolation) GetReasonOk() (*string, bool)`

GetReasonOk returns a tuple with the Reason field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReason

`func (o *ErrdetailsBadRequestFieldViolation) SetReason(v string)`

SetReason sets Reason field to given value.

### HasReason

`func (o *ErrdetailsBadRequestFieldViolation) HasReason() bool`

HasReason returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


