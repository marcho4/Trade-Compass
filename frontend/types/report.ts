export interface Report {
  id: number
  ticker: string
  year: number
  period: string
  s3_path: string
}

export interface ReportsResponse {
  ticker: string
  reports: Report[]
  total: number
}
