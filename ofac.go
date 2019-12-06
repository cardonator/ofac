// Copyright 2018 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package ofac

// ToDo: NON-SDN List, Consolidated List - They appear to have the same format. Other list?

// SDN is a specially Designated National
type SDN struct {
	// EntityID (ent_num) is the unique record identifier/unique listing identifier
	EntityID string `json:"entityID"`
	// SDNName (SDN_name)  is the name of the specially designated national
	SDNName string `json:"sdnName"`
	// SDNType (SDN_Type) is the type of SDN
	SDNType string `json:"sdnType"`
	// Program is the sanctions program name
	Program string `json:"program"`
	// Title is the title of an individual
	Title string `json:"title"`
	// CallSign (Call_Sign) is vessel call sign
	CallSign string `json:"callSign"`
	// VesselType (Vess_type) is the vessel type
	VesselType string `json:"vesselType"`
	// Tonnage is the vessel tonnage
	Tonnage string `json:"tonnage"`
	// GrossRegisteredTonnage (GRT) is gross registered tonnage
	GrossRegisteredTonnage string `json:"grossRegisteredTonnage"`
	// VesselFlag (Vess_flag) is vessel flag
	VesselFlag string `json:"vesselFlag"`
	// VesselOwner  (Vess_owner) is vessel owner
	VesselOwner string `json:"vesselOwner"`
	//  Remarks is remarks on specially designated national
	Remarks string `json:"remarks"`
}

// Address is OFAC SDN Addresses
type Address struct {
	// EntityID (ent_num) is the unique record identifier/unique listing identifier
	EntityID string `json:"entityID"`
	// AddressID (add_num) is the unique record identifier for the address
	AddressID string `json:"addressID"`
	// Address is the street address of the specially designated national
	Address string `json:"address"`
	// CityStateProvincePostalCode is the city, state/province, zip/postal code for the address of the
	// specially designated national
	CityStateProvincePostalCode string `json:"cityStateProvincePostalCode"`
	// Country is the country for the address of the specially designated national
	Country string `json:"country"`
	//AddressRemarks (Add_remarks) is remarks on the address
	AddressRemarks string `json:"addressRemarks"`
}

// AlternateIdentity is OFAC SDN Alternate Identity object
type AlternateIdentity struct {
	// EntityID (ent_num) is the unique record identifier/unique listing identifier
	EntityID string `json:"entityID"`
	// AlternateID (alt_num) is the unique record identifier for the alternate identity
	AlternateID string `json:"alternateID"`
	// AlternateIdentityType (alt_type) is the type of alternate identity (aka, fka, nka)
	AlternateType string `json:"alternateType"`
	// AlternateIdentityName (alt_name) is the alternate identity name of the specially designated national
	AlternateName string `json:"alternateName"`
	// AlternateIdentityRemarks (alt_remarks) is remarks on alternate identity of the specially designated national
	AlternateRemarks string `json:"alternateRemarks"`
}

// SDNComments is OFAC SDN Additional Comments
type SDNComments struct {
	// EntityID (ent_num) is the unique record identifier/unique listing identifier
	EntityID string `json:"entityID"`
	// RemarksExtended is remarks extended on a Specially Designated National
	RemarksExtended string `json:"remarksExtended"`
}

// DPL is the BIS Denied Persons List
type DPL struct {
	// Name is the name of the Denied Person
	Name string `json:"name"`
	// StreetAddress is the Denied Person's street address
	StreetAddress string `json:"streetAddress"`
	// City is the Denied Person's city
	City string `json:"city"`
	// State is the Denied Person's state
	State string `json:"state"`
	// Country is the Denied Person's country
	Country string `json:"country"`
	// PostalCode is the Denied Person's postal code
	PostalCode string `json:"postalCode"`
	// EffectiveDate is the date the denial came into effect
	EffectiveDate string `json:"effectiveDate"`
	// ExpirationDate is the date the denial expires. If blank, the denial has no expiration
	ExpirationDate string `json:"expirationDate"`
	// StandardOrder denotes whether or not the Person was added to the list by a "standard" order
	StandardOrder string `json:"standardOrder"`
	// LastUpdate is the date of the most recent change to the denial
	LastUpdate string `json:"lastUpdate"`
	// Action is the most recent action taken regarding the denial
	Action string `json:"action"`
	// FRCitation is the reference to the order's citation in the Federal Register
	FRCitation string `json:"frCitation"`
}

// SSI is the Sectoral Sanctions Identifications List - Treasury Department
type SSI struct {
	// EntityID (ent_num) is the unique record identifier/unique listing identifier
	EntityID string `json:"entityID"`
	// Type is the entity type (e.g. individual, vessel, aircraft, etc)
	Type string `json:"type"`
	// Programs is the list of sanctions program for which the entity is flagged
	Programs []string `json:"programs"`
	// Name is the entity's name (e.g. given name for individual, company name, etc.)
	Name string `json:"name"`
	// Addresses is a list of known addresses associated with the entity
	Addresses []string `json:"addresses"`
	// Remarks is used to provide additional details for the entity
	Remarks []string `json:"remarks"`
	// AlternateNames is a list of aliases associated with the entity
	AlternateNames []string `json:"alternateNames"`
	// IDsOnRecord is a list of the forms of identification on file for the entity
	IDsOnRecord []string `json:"ids"`
	// SourceListURL is a link to the official SSI list
	SourceListURL string `json:"sourceListURL"`
	// SourceInfoURL is a link to information about the list
	SourceInfoURL string `json:"sourceInfoURL"`
}

// EL is the Entity List (EL) - Bureau of Industry and Security
type EL struct {
	// Name is the primary name of the entity
	Name string `json:"name"`
	// AlternateNames is a list of aliases associated with the entity
	AlternateNames []string `json:"alternateNames"`
	// Addresses is a list of known addresses associated with the entity
	Addresses []string `json:"addresses"`
	// StartDate is the effective date
	StartDate string `json:"startDate"`
	// LicenceRequirement specifies the license requirements that it imposes on each listed person
	LicenceRequirement string `json:"licenseRequirement"`
	// LicensePolicy is the policy with which BIS reviews the requirements set forth in Licence Requirements
	LicensePolicy string `json:"licensePolicy"`
	// FRNotice identifies the notice in the Federal Register
	FRNotice string `json:"FRNotice"`
	// SourceListURL is a link to the official SSI list
	SourceListURL string `json:"sourceListURL"`
	// SourceInfoURL is a link to information about the list
	SourceInfoURL string `json:"sourceInfoURL"`
}
