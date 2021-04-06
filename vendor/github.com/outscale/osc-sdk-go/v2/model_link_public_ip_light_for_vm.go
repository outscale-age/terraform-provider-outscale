/*
 * 3DS OUTSCALE API
 *
 * Welcome to the 3DS OUTSCALE's API documentation.<br /><br />  The 3DS OUTSCALE API enables you to manage your resources in the 3DS OUTSCALE Cloud. This documentation describes the different actions available along with code examples.<br /><br />  Note that the 3DS OUTSCALE Cloud is compatible with Amazon Web Services (AWS) APIs, but some resources have different names in AWS than in the 3DS OUTSCALE API. You can find a list of the differences [here](https://wiki.outscale.net/display/EN/3DS+OUTSCALE+APIs+Reference).<br /><br />  You can also manage your resources using the [Cockpit](https://wiki.outscale.net/display/EN/About+Cockpit) web interface.
 *
 * API version: 1.7
 * Contact: support@outscale.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package osc

import (
	"encoding/json"
)

// LinkPublicIpLightForVm Information about the EIP associated with the NIC.
type LinkPublicIpLightForVm struct {
	// The name of the public DNS.
	PublicDnsName *string `json:"PublicDnsName,omitempty"`
	// The External IP address (EIP) associated with the NIC.
	PublicIp *string `json:"PublicIp,omitempty"`
	// The account ID of the owner of the EIP.
	PublicIpAccountId *string `json:"PublicIpAccountId,omitempty"`
}

// NewLinkPublicIpLightForVm instantiates a new LinkPublicIpLightForVm object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewLinkPublicIpLightForVm() *LinkPublicIpLightForVm {
	this := LinkPublicIpLightForVm{}
	return &this
}

// NewLinkPublicIpLightForVmWithDefaults instantiates a new LinkPublicIpLightForVm object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewLinkPublicIpLightForVmWithDefaults() *LinkPublicIpLightForVm {
	this := LinkPublicIpLightForVm{}
	return &this
}

// GetPublicDnsName returns the PublicDnsName field value if set, zero value otherwise.
func (o *LinkPublicIpLightForVm) GetPublicDnsName() string {
	if o == nil || o.PublicDnsName == nil {
		var ret string
		return ret
	}
	return *o.PublicDnsName
}

// GetPublicDnsNameOk returns a tuple with the PublicDnsName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LinkPublicIpLightForVm) GetPublicDnsNameOk() (*string, bool) {
	if o == nil || o.PublicDnsName == nil {
		return nil, false
	}
	return o.PublicDnsName, true
}

// HasPublicDnsName returns a boolean if a field has been set.
func (o *LinkPublicIpLightForVm) HasPublicDnsName() bool {
	if o != nil && o.PublicDnsName != nil {
		return true
	}

	return false
}

// SetPublicDnsName gets a reference to the given string and assigns it to the PublicDnsName field.
func (o *LinkPublicIpLightForVm) SetPublicDnsName(v string) {
	o.PublicDnsName = &v
}

// GetPublicIp returns the PublicIp field value if set, zero value otherwise.
func (o *LinkPublicIpLightForVm) GetPublicIp() string {
	if o == nil || o.PublicIp == nil {
		var ret string
		return ret
	}
	return *o.PublicIp
}

// GetPublicIpOk returns a tuple with the PublicIp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LinkPublicIpLightForVm) GetPublicIpOk() (*string, bool) {
	if o == nil || o.PublicIp == nil {
		return nil, false
	}
	return o.PublicIp, true
}

// HasPublicIp returns a boolean if a field has been set.
func (o *LinkPublicIpLightForVm) HasPublicIp() bool {
	if o != nil && o.PublicIp != nil {
		return true
	}

	return false
}

// SetPublicIp gets a reference to the given string and assigns it to the PublicIp field.
func (o *LinkPublicIpLightForVm) SetPublicIp(v string) {
	o.PublicIp = &v
}

// GetPublicIpAccountId returns the PublicIpAccountId field value if set, zero value otherwise.
func (o *LinkPublicIpLightForVm) GetPublicIpAccountId() string {
	if o == nil || o.PublicIpAccountId == nil {
		var ret string
		return ret
	}
	return *o.PublicIpAccountId
}

// GetPublicIpAccountIdOk returns a tuple with the PublicIpAccountId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LinkPublicIpLightForVm) GetPublicIpAccountIdOk() (*string, bool) {
	if o == nil || o.PublicIpAccountId == nil {
		return nil, false
	}
	return o.PublicIpAccountId, true
}

// HasPublicIpAccountId returns a boolean if a field has been set.
func (o *LinkPublicIpLightForVm) HasPublicIpAccountId() bool {
	if o != nil && o.PublicIpAccountId != nil {
		return true
	}

	return false
}

// SetPublicIpAccountId gets a reference to the given string and assigns it to the PublicIpAccountId field.
func (o *LinkPublicIpLightForVm) SetPublicIpAccountId(v string) {
	o.PublicIpAccountId = &v
}

func (o LinkPublicIpLightForVm) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.PublicDnsName != nil {
		toSerialize["PublicDnsName"] = o.PublicDnsName
	}
	if o.PublicIp != nil {
		toSerialize["PublicIp"] = o.PublicIp
	}
	if o.PublicIpAccountId != nil {
		toSerialize["PublicIpAccountId"] = o.PublicIpAccountId
	}
	return json.Marshal(toSerialize)
}

type NullableLinkPublicIpLightForVm struct {
	value *LinkPublicIpLightForVm
	isSet bool
}

func (v NullableLinkPublicIpLightForVm) Get() *LinkPublicIpLightForVm {
	return v.value
}

func (v *NullableLinkPublicIpLightForVm) Set(val *LinkPublicIpLightForVm) {
	v.value = val
	v.isSet = true
}

func (v NullableLinkPublicIpLightForVm) IsSet() bool {
	return v.isSet
}

func (v *NullableLinkPublicIpLightForVm) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableLinkPublicIpLightForVm(val *LinkPublicIpLightForVm) *NullableLinkPublicIpLightForVm {
	return &NullableLinkPublicIpLightForVm{value: val, isSet: true}
}

func (v NullableLinkPublicIpLightForVm) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableLinkPublicIpLightForVm) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}