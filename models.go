package endoflife

import "time"

// Uri is a link to a resource.
type Uri struct {
	// Name of the URI, e.g. "tags".
	Name string `json:"name"`
	// URI value, e.g. "https://endoflife.date/tags/".
	Uri string `json:"uri"`
}

// Identifier is a product identifier, such as a purl, repology or cpe identifier.
type Identifier struct {
	// ID is the identifier value, e.g. "cpe:/o:canonical:ubuntu_linux".
	ID string `json:"id"`
	// Type of the identifier, e.g. "cpe".
	Type string `json:"type"`
}

// ProductVersion holds information about a specific product version.
type ProductVersion struct {
	// Name of the version, e.g. "22.04.2".
	Name string `json:"name"`
	// Date is the release date of the version. Nil when the information is not known.
	Date *Date `json:"date"`
	// Link to the changelog or release notes. Nil when no public link is available.
	Link *string `json:"link"`
}

// ProductRelease holds full information about a product release cycle.
type ProductRelease struct {
	// Name of the product release cycle, e.g. "22.04".
	Name string `json:"name"`
	// Codename of the release cycle. Nil when there is no codename or it is unknown.
	Codename *string `json:"codename"`
	// Label of the release cycle, e.g. "22.04 'Jammy Jellyfish' (LTS)".
	Label string `json:"label"`
	// ReleaseDate is the release date of the release cycle.
	ReleaseDate Date `json:"releaseDate"`

	// IsLts reports whether the release cycle receives long-term support.
	IsLts bool `json:"isLts"`
	// LtsFrom is the start date of the LTS phase. Nil when not applicable or unknown.
	LtsFrom *Date `json:"ltsFrom"`

	// IsEoas reports whether the active support phase is over. Nil when the product
	// does not have an active support phase (field absent in the API response).
	IsEoas *bool `json:"isEoas,omitempty"`
	// EoasFrom is the end of active support date. Nil when not applicable or unknown.
	EoasFrom *Date `json:"eoasFrom,omitempty"`

	// IsEol reports whether the release cycle has reached end of life.
	IsEol bool `json:"isEol"`
	// EolFrom is the end of life date. Nil when the date is not known.
	EolFrom *Date `json:"eolFrom"`

	// IsDiscontinued reports whether the release cycle is discontinued. Mainly used
	// for hardware; nil when not applicable (field absent in the API response).
	IsDiscontinued *bool `json:"isDiscontinued,omitempty"`
	// DiscontinuedFrom is the discontinuation date. Nil when not applicable or unknown.
	DiscontinuedFrom *Date `json:"discontinuedFrom,omitempty"`

	// IsEoes reports whether the extended support phase is over. Nil when the product
	// does not have an extended support phase, or the cycle is not eligible for it.
	IsEoes *bool `json:"isEoes,omitempty"`
	// EoesFrom is the end of extended support date. Nil when not applicable or unknown.
	EoesFrom *Date `json:"eoesFrom,omitempty"`

	// IsMaintained reports whether the release cycle still has some level of support
	// (including extended support).
	IsMaintained bool `json:"isMaintained"`

	// Latest is the latest version for this release cycle. Nil when none is documented.
	Latest *ProductVersion `json:"latest"`

	// Custom holds custom fields for the release cycle. Nil when the product does not
	// declare any custom fields. Values may be null (represented as nil entries).
	Custom map[string]*string `json:"custom,omitempty"`
}

// ProductSummary is a summary of a product, as returned by list endpoints.
type ProductSummary struct {
	// Name of the product, e.g. "ubuntu".
	Name string `json:"name"`
	// Label of the product, e.g. "Ubuntu".
	Label string `json:"label"`
	// Aliases declared for the product. Empty when none is declared.
	Aliases []string `json:"aliases"`
	// Category of the product, e.g. "os".
	Category string `json:"category"`
	// Tags associated with the product. Always contains at least one tag.
	Tags []string `json:"tags"`
	// Uri is a link to the full product details.
	Uri string `json:"uri"`
}

// ProductLabels holds the labels used to denote a product's lifecycle phases.
type ProductLabels struct {
	// Eoas labels the phase before end of active support. Nil when not applicable.
	Eoas *string `json:"eoas"`
	// Discontinued labels the discontinuation phase. Nil for software.
	Discontinued *string `json:"discontinued"`
	// Eol labels the phase before end of life.
	Eol string `json:"eol"`
	// Eoes labels the phase before end of extended support. Nil when not applicable.
	Eoes *string `json:"eoes"`
}

// ProductLinks holds links related to a product.
type ProductLinks struct {
	// Icon links to the product icon on simpleicons.org. Nil when none exists.
	Icon *string `json:"icon"`
	// HTML links to the product page on endoflife.date.
	HTML string `json:"html"`
	// ReleasePolicy links to the product release policy. Nil when none is available.
	ReleasePolicy *string `json:"releasePolicy"`
}

// ProductDetails holds the full details of a product.
type ProductDetails struct {
	// Name of the product, e.g. "ubuntu".
	Name string `json:"name"`
	// Label of the product, e.g. "Ubuntu".
	Label string `json:"label"`
	// Aliases declared for the product. Empty when none is declared.
	Aliases []string `json:"aliases"`
	// Category of the product, e.g. "os".
	Category string `json:"category"`
	// Tags associated with the product. Always contains at least one tag.
	Tags []string `json:"tags"`
	// VersionCommand is a command to check the current product version. Nil when unknown.
	VersionCommand *string `json:"versionCommand"`
	// Identifiers known for the product (purl, repology, cpe...). Empty when none.
	Identifiers []Identifier `json:"identifiers"`
	// Labels used for the product's lifecycle phases.
	Labels ProductLabels `json:"labels"`
	// Links related to the product.
	Links ProductLinks `json:"links"`
	// Releases is the list of all product releases.
	Releases []ProductRelease `json:"releases"`
}

// IdentifierEntry associates an identifier with its related product.
type IdentifierEntry struct {
	// Identifier value, e.g. "cpe:/o:canonical:ubuntu_linux".
	Identifier string `json:"identifier"`
	// Product references the product this identifier relates to.
	Product Uri `json:"product"`
}

// UriListResponse is a response containing a list of URIs.
type UriListResponse struct {
	SchemaVersion string    `json:"schema_version"`
	GeneratedAt   time.Time `json:"generated_at"`
	Total         int       `json:"total"`
	Result        []Uri     `json:"result"`
}

// ProductListResponse is a response containing a list of product summaries.
type ProductListResponse struct {
	SchemaVersion string           `json:"schema_version"`
	GeneratedAt   time.Time        `json:"generated_at"`
	Total         int              `json:"total"`
	Result        []ProductSummary `json:"result"`
}

// FullProductListResponse is a response containing a list of full product details.
type FullProductListResponse struct {
	SchemaVersion string           `json:"schema_version"`
	GeneratedAt   time.Time        `json:"generated_at"`
	Total         int              `json:"total"`
	Result        []ProductDetails `json:"result"`
}

// ProductResponse is a response containing a single product's full details.
type ProductResponse struct {
	SchemaVersion string         `json:"schema_version"`
	GeneratedAt   time.Time      `json:"generated_at"`
	LastModified  time.Time      `json:"last_modified"`
	Result        ProductDetails `json:"result"`
}

// ProductReleaseResponse is a response containing a single release cycle.
type ProductReleaseResponse struct {
	SchemaVersion string         `json:"schema_version"`
	GeneratedAt   time.Time      `json:"generated_at"`
	Result        ProductRelease `json:"result"`
}

// IdentifierListResponse is a response containing all identifiers for a type.
type IdentifierListResponse struct {
	SchemaVersion string            `json:"schema_version"`
	GeneratedAt   time.Time         `json:"generated_at"`
	Total         int               `json:"total"`
	Result        []IdentifierEntry `json:"result"`
}
