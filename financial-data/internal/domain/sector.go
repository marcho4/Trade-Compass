package domain

type Sector int

const (
	Oils Sector = iota + 1
	Finance
	Technology
	Telecom
	Metals
	Mining
	Utilities
	RealEstate
	ConsumerStaples
	ConsumerDiscretionary
	Healthcare
	Industrial
	Energy
	Materials
	Transportation
	Agriculture
	Chemicals
	Construction
	Retail
)

func (s Sector) IsValid() bool {
	return s >= 1 && s <= 19
}
