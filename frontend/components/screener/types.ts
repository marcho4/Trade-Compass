export interface FilterValues {
  search: string
  sector: string
  ratingMin: string
}

export interface CompanyRating {
  health: number
  growth: number
  moat: number
  dividends: number
  value: number
  total: number
}
