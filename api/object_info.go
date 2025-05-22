package api

import "encoding/json"

// ObjectInfo represents SObject metadata from the UI API's object-info endpoint.
type ObjectInfo struct {
	APIName               string                            `json:"apiName"`
	AssociateEntityType   *string                           `json:"associateEntityType"`
	AssociateParentEntity *string                           `json:"associateParentEntity"`
	ChildRelationships    []ChildRelationshipInfo           `json:"childRelationships"`
	CompactLayoutable     bool                              `json:"compactLayoutable"`
	Createable            bool                              `json:"createable"`
	Custom                bool                              `json:"custom"`
	DefaultRecordTypeID   string                            `json:"defaultRecordTypeId"`
	Deletable             bool                              `json:"deletable"`
	DependentFields       map[string]map[string]interface{} `json:"dependentFields"`
	ETag                  string                            `json:"eTag"`
	FeedEnabled           bool                              `json:"feedEnabled"`
	Fields                map[string]FieldInfo              `json:"fields"`
	KeyPrefix             string                            `json:"keyPrefix"`
	Label                 string                            `json:"label"`
	LabelPlural           string                            `json:"labelPlural"`
	Layoutable            bool                              `json:"layoutable"`
	MRUEnabled            bool                              `json:"mruEnabled"`
	NameFields            []string                          `json:"nameFields"`
	Queryable             bool                              `json:"queryable"`
	RecordTypeInfos       map[string]RecordTypeInfo         `json:"recordTypeInfos"`
	SearchLayoutable      bool                              `json:"searchLayoutable"`
	Searchable            bool                              `json:"searchable"`
	ThemeInfo             ThemeInfo                         `json:"themeInfo"`
	Updateable            bool                              `json:"updateable"`
}

// ChildRelationshipInfo holds metadata for a child relationship.
type ChildRelationshipInfo struct {
	ChildObjectAPIName  string   `json:"childObjectApiName"`
	FieldName           string   `json:"fieldName"`
	JunctionIDListNames []string `json:"junctionIdListNames"`
	JunctionReferenceTo []string `json:"junctionReferenceTo"`
	RelationshipName    string   `json:"relationshipName"`
}

// ThemeInfo holds theme metadata for an SObject.
type ThemeInfo struct {
	Color   string `json:"color"`
	IconURL string `json:"iconUrl"`
}

// RecordTypeInfo holds metadata for a record type.
type RecordTypeInfo struct {
	Available                bool   `json:"available"`
	DefaultRecordTypeMapping bool   `json:"defaultRecordTypeMapping"`
	Master                   bool   `json:"master"`
	Name                     string `json:"name"`
	RecordTypeID             string `json:"recordTypeId"`
}

// FieldInfo holds metadata for a single field.
type FieldInfo struct {
	APIName               string            `json:"apiName"`
	Calculated            bool              `json:"calculated"`
	Compound              bool              `json:"compound"`
	CompoundComponentName *string           `json:"compoundComponentName"`
	CompoundFieldName     *string           `json:"compoundFieldName"`
	ControllerName        *string           `json:"controllerName"`
	ControllingFields     []string          `json:"controllingFields"`
	Createable            bool              `json:"createable"`
	Custom                bool              `json:"custom"`
	DataType              string            `json:"dataType"`
	ExternalID            bool              `json:"externalId"`
	ExtraTypeInfo         *string           `json:"extraTypeInfo"`
	Filterable            bool              `json:"filterable"`
	FilteredLookupInfo    json.RawMessage   `json:"filteredLookupInfo"`
	HighScaleNumber       bool              `json:"highScaleNumber"`
	HTMLFormatted         bool              `json:"htmlFormatted"`
	InlineHelpText        *string           `json:"inlineHelpText"`
	Label                 string            `json:"label"`
	Length                *int              `json:"length"`
	MaskType              *string           `json:"maskType"`
	NameField             bool              `json:"nameField"`
	PolymorphicForeignKey interface{}       `json:"polymorphicForeignKey"`
	Precision             *int              `json:"precision"`
	Reference             json.RawMessage   `json:"reference"`
	ReferenceTargetField  *string           `json:"referenceTargetField"`
	ReferenceToInfos      []json.RawMessage `json:"referenceToInfos"`
	RelationshipName      *string           `json:"relationshipName"`
	Required              bool              `json:"required"`
	Scale                 *int              `json:"scale"`
	SearchPrefilterable   bool              `json:"searchPrefilterable"`
	Sortable              bool              `json:"sortable"`
	Unique                bool              `json:"unique"`
	Updateable            bool              `json:"updateable"`
}

// UnmarshalObjectInfo parses JSON data into an ObjectInfo struct.
func UnmarshalObjectInfo(data []byte) (ObjectInfo, error) {
	var info ObjectInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return ObjectInfo{}, err
	}
	return info, nil
}
