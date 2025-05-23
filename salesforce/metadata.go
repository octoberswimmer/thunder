package salesforce

import (
	"embed"
)

// SalesforceMetadataFS embeds Apex class definitions and static resource metadata
//
//go:embed classes/*.cls classes/*.cls-meta.xml
//go:embed lwc/go/*
//go:embed lwc/thunder/*
var SalesforceMetadataFS embed.FS
