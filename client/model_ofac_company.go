/*
 * OFAC API
 *
 * OFAC (Office of Foreign Assets Control) API is designed to facilitate the enforcement of US government economic sanctions programs required by federal law. This project implements a modern REST HTTP API for companies and organizations to obey federal law and use OFAC data in their applications.
 *
 * API version: v1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

// OFAC Company and metadata
type OfacCompany struct {
	// OFAC Company ID
	Id        string            `json:"id,omitempty"`
	Sdn       Sdn               `json:"sdn,omitempty"`
	Addresses []Address         `json:"addresses,omitempty"`
	Alts      []Alt             `json:"alts,omitempty"`
	Status    OfacCompanyStatus `json:"status,omitempty"`
}
