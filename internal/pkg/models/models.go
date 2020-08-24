package models

import (
	"time"
)

type PersonAddress struct {
	Primary       bool
	Recorded      time.Time
	AddressString string
}

type PhoneNumber struct {
	Primary bool
	// Landline     bool
	ReportedDate time.Time
	Number       string
}

type Associate struct {
	BirthDate time.Time
	Name      string
}

type Relative struct {
	BirthDate time.Time
	Spouse    bool
	Name      string
}

type Person struct {
	Deceased     bool
	DeathDate    time.Time
	BirthDate    time.Time
	AlsoKnownAs  []string
	Name         string
	FullName     string
	Emails       []string
	Addresses    []PersonAddress
	PhoneNumbers []PhoneNumber
	Relatives    []Relative
	Associates   []Associate
	Businesses   []string
}

type ElasticHitResult struct {
	Index  string  `json:"_index"`
	Type   string  `json:"_type"`
	Id     string  `json:"_id"`
	Score  float64 `json:"_score"`
	Source Person  `json:"_source"`
}

type ElasticTotalResult struct {
	TotalValue int `json:"value"`
}

type ElasticHitResults struct {
	MaxScore     float64            `json:"max_score"`
	TotalResults ElasticTotalResult `json:"total"`
	Hits         []ElasticHitResult `json:"hits"`
}

type ElasticQueryResult struct {
	Took     int               `json:"took"`
	TimedOut bool              `json:"timed_out"`
	Hits     ElasticHitResults `json:"hits"`
}

type PaginationResult struct {
	Pages              int
	Offset             int
	IsNextPage         bool
	NextPageOffset     int
	IsPreviousPage     bool
	PreviousPageOffset int
}

type SiteSearchResult struct {
	Pagination    PaginationResult
	Query         string
	ElasticResult ElasticQueryResult
}
