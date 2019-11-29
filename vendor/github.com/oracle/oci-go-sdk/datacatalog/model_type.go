// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// DataCatalog API
//
// A description of the DataCatalog API
//

package datacatalog

import (
	"github.com/oracle/oci-go-sdk/common"
)

// ModelType Full Data Catalog Type Definition. Fully defines a type of the data catalog. All types are statically defined
// in the system and are immutable. It isn't possible to create new types or update existing types via the api.
type ModelType struct {

	// Unique Type key that is immutable.
	Key *string `mandatory:"true" json:"key"`

	// The immutable name of the type.
	Name *string `mandatory:"false" json:"name"`

	// Detailed description of the Type.
	Description *string `mandatory:"false" json:"description"`

	// The Catalog's Oracle ID (OCID).
	CatalogId *string `mandatory:"false" json:"catalogId"`

	// A map of arrays which defines the Type specific properties, both required and optional. The map keys are
	// category names and the values are arrays contiaing all property details. Every property is contained inside
	// of a category. Most Types have required properties within the "default" category.
	// Example:
	// `{
	//    "properties": {
	//      "default": {
	//        "attributes:": [
	//          {
	//            "name": "host",
	//            "type": "string",
	//            "isRequired": true,
	//            "isUpdatable": false
	//          },
	//          ...
	//        ]
	//      }
	//    }
	//  }`
	Properties map[string][]PropertyDefinition `mandatory:"false" json:"properties"`

	// The current state of the Type.
	LifecycleState LifecycleStateEnum `mandatory:"false" json:"lifecycleState,omitempty"`

	// Indicates whether the type is internal, making it unavailable for use by metadata elements.
	IsInternal *bool `mandatory:"false" json:"isInternal"`

	// Indicates whether the type can be used for tagging metadata elements.
	IsTag *bool `mandatory:"false" json:"isTag"`

	// Indicates whether the type is approved for use as a classifying object.
	IsApproved *bool `mandatory:"false" json:"isApproved"`

	// Indicates the category this type belongs to. For instance , data assets , connections.
	TypeCategory *string `mandatory:"false" json:"typeCategory"`

	// Mapping type equivalence in the external system.
	ExternalTypeName *string `mandatory:"false" json:"externalTypeName"`

	// URI to the Type instance in the API.
	Uri *string `mandatory:"false" json:"uri"`
}

func (m ModelType) String() string {
	return common.PointerString(m)
}