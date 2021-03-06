/*
 * OFAC API
 *
 * OFAC (Office of Foreign Assets Control) API is designed to facilitate the enforcement of US government economic sanctions programs required by federal law. This project implements a modern REST HTTP API for companies and organizations to obey federal law and use OFAC data in their applications.
 *
 * API version: v1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

// Entity List (EL) - Bureau of Industry and Security
type El struct {
	// The name of the entity
	Name string `json:"name,omitempty"`
	// Addresses associated with the entity
	Addresses []string `json:"addresses,omitempty"`
	// Known aliases associated with the entity
	AlternateNames []string `json:"alternateNames,omitempty"`
	// Date when the restriction came into effect
	StartDate string `json:"startDate,omitempty"`
	// Specifies the license requirement imposed on the named entity
	LicenseRequirement string `json:"licenseRequirement,omitempty"`
	// Identifies the policy BIS uses to review the licenseRequirements
	LicensePolicy string `json:"licensePolicy,omitempty"`
	// Identifies the corresponding Notice in the Federal Register
	FrNotice string `json:"frNotice,omitempty"`
	// The link to the official SSI list
	SourceListURL string `json:"sourceListURL,omitempty"`
	// The link for information regarding the source
	SourceInfoURL string `json:"sourceInfoURL,omitempty"`
}
