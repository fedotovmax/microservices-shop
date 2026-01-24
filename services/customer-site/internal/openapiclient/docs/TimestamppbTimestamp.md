# TimestamppbTimestamp

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Nanos** | Pointer to **int32** | Non-negative fractions of a second at nanosecond resolution. This field is the nanosecond portion of the duration, not an alternative to seconds. Negative second values with fractions must still have non-negative nanos values that count forward in time. Must be between 0 and 999,999,999 inclusive. | [optional] 
**Seconds** | Pointer to **int32** | Represents seconds of UTC time since Unix epoch 1970-01-01T00:00:00Z. Must be between -315576000000 and 315576000000 inclusive (which corresponds to 0001-01-01T00:00:00Z to 9999-12-31T23:59:59Z). | [optional] 

## Methods

### NewTimestamppbTimestamp

`func NewTimestamppbTimestamp() *TimestamppbTimestamp`

NewTimestamppbTimestamp instantiates a new TimestamppbTimestamp object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTimestamppbTimestampWithDefaults

`func NewTimestamppbTimestampWithDefaults() *TimestamppbTimestamp`

NewTimestamppbTimestampWithDefaults instantiates a new TimestamppbTimestamp object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNanos

`func (o *TimestamppbTimestamp) GetNanos() int32`

GetNanos returns the Nanos field if non-nil, zero value otherwise.

### GetNanosOk

`func (o *TimestamppbTimestamp) GetNanosOk() (*int32, bool)`

GetNanosOk returns a tuple with the Nanos field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNanos

`func (o *TimestamppbTimestamp) SetNanos(v int32)`

SetNanos sets Nanos field to given value.

### HasNanos

`func (o *TimestamppbTimestamp) HasNanos() bool`

HasNanos returns a boolean if a field has been set.

### GetSeconds

`func (o *TimestamppbTimestamp) GetSeconds() int32`

GetSeconds returns the Seconds field if non-nil, zero value otherwise.

### GetSecondsOk

`func (o *TimestamppbTimestamp) GetSecondsOk() (*int32, bool)`

GetSecondsOk returns a tuple with the Seconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSeconds

`func (o *TimestamppbTimestamp) SetSeconds(v int32)`

SetSeconds sets Seconds field to given value.

### HasSeconds

`func (o *TimestamppbTimestamp) HasSeconds() bool`

HasSeconds returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


