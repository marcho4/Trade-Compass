package entity

import (
	"fmt"
	"strings"
)

type CompanyProfile struct {
	Ticker              string   `json:"ticker"`
	CompanyName         string   `json:"company_name"`
	Description         string   `json:"description"`
	ProductsAndServices []string `json:"products_and_services"`
	Markets             []Market `json:"markets"`
	KeyClients          string   `json:"key_clients"`
	BusinessModel       string   `json:"business_model"`
}

func (c *CompanyProfile) String() string {
	var s strings.Builder
	s.WriteString("Профиль компании\n")
	fmt.Fprintf(&s, "Тикер: %s\n", c.Ticker)
	fmt.Fprintf(&s, "Название: %s\n", c.CompanyName)
	fmt.Fprintf(&s, "Описание: %s\n", c.Description)
	fmt.Fprintf(&s, "Бизнес-модель: %s\n", c.BusinessModel)
	fmt.Fprintf(&s, "Ключевые клиенты: %s\n", c.KeyClients)
	if len(c.ProductsAndServices) > 0 {
		fmt.Fprintf(&s, "Продукты и услуги: %s\n", strings.Join(c.ProductsAndServices, ", "))
	}
	if len(c.Markets) > 0 {
		s.WriteString("Рынки:\n")
		for _, m := range c.Markets {
			fmt.Fprintf(&s, "  - %s (роль: %s)\n", m.Market, m.Role)
		}
	}
	return s.String()
}

type Market struct {
	Market string `json:"market"`
	Role   string `json:"role"`
}

type RevenueSource struct {
	Ticker      string  `json:"ticker"`
	Segment     string  `json:"segment"`
	SharePct    float64 `json:"share_pct"`
	Approximate bool    `json:"approximate"`
	Description string  `json:"description"`
	Trend       string  `json:"trend"`
}

func (r *RevenueSource) String() string {
	approx := ""
	if r.Approximate {
		approx = " (приблизительно)"
	}
	return fmt.Sprintf("  - %s: %.1f%%%s — %s [тренд: %s]",
		r.Segment, r.SharePct, approx, r.Description, r.Trend)
}

type CompanyDependency struct {
	Ticker      string `json:"ticker"`
	Factor      string `json:"factor"`
	Type        string `json:"type"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
}

func (d *CompanyDependency) String() string {
	return fmt.Sprintf("  - %s [тип: %s, критичность: %s]: %s",
		d.Factor, d.Type, d.Severity, d.Description)
}

type BusinessResearchResponse struct {
	Ticker      string `json:"ticker"`
	CompanyName string `json:"company_name"`
	Profile     struct {
		Description         string   `json:"description"`
		ProductsAndServices []string `json:"products_and_services"`
		Markets             []Market `json:"markets"`
		KeyClients          string   `json:"key_clients"`
		BusinessModel       string   `json:"business_model"`
	} `json:"profile"`
	RevenueSources []struct {
		Segment     string  `json:"segment"`
		SharePct    float64 `json:"share_pct"`
		Approximate bool    `json:"approximate"`
		Description string  `json:"description"`
		Trend       string  `json:"trend"`
	} `json:"revenue_sources"`
	Dependencies []struct {
		Factor      string `json:"factor"`
		Type        string `json:"type"`
		Severity    string `json:"severity"`
		Description string `json:"description"`
	} `json:"dependencies"`
}

type BusinessResearchResult struct {
	Profile      CompanyProfile      `json:"profile"`
	Revenue      []RevenueSource     `json:"revenue_sources"`
	Dependencies []CompanyDependency `json:"dependencies"`
}

func (b *BusinessResearchResult) String() string {
	var s strings.Builder

	s.WriteString(b.Profile.String())

	if len(b.Revenue) > 0 {
		s.WriteString("\n=== Источники выручки ===\n")
		for _, r := range b.Revenue {
			s.WriteString(r.String())
			s.WriteByte('\n')
		}
	}

	if len(b.Dependencies) > 0 {
		s.WriteString("\n=== Зависимости ===\n")
		for _, d := range b.Dependencies {
			s.WriteString(d.String())
			s.WriteByte('\n')
		}
	}

	return s.String()
}
