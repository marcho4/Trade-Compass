import { RawData } from '@/types/raw-data';

const AI_BASE_URL = '/api/ai';

export interface AvailablePeriod {
  year: number;
  period: number;
}

export interface ReportResults {
  health: number;
  growth: number;
  moat: number;
  dividends: number;
  value: number;
  total: number;
}

export interface BusinessResearch {
  profile: {
    ticker: string;
    company_name: string;
    description: string;
    products_and_services: string[];
    markets: { market: string; role: string }[];
    key_clients: string;
    business_model: string;
  };
  revenue_sources: {
    ticker: string;
    segment: string;
    share_pct: number;
    approximate: boolean;
    description: string;
    trend: string;
  }[];
  dependencies: {
    ticker: string;
    factor: string;
    type: string;
    severity: string;
    description: string;
  }[];
}

export const aiApi = {
  async extractData(
    ticker: string,
    period: string,
    year?: number,
    force?: boolean,
    signal?: AbortSignal
  ): Promise<RawData> {
    const params = new URLSearchParams({ ticker, period });
    if (year) params.set('year', year.toString());
    if (force) params.set('force', 'true');

    const response = await fetch(`${AI_BASE_URL}/extract?${params}`, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
      signal,
    });

    if (!response.ok) {
      const body = await response.json().catch(() => ({}));
      throw new Error(body.error || `Extraction failed (${response.status})`);
    }

    return response.json();
  },

  async getAvailablePeriods(
    ticker: string,
    signal?: AbortSignal
  ): Promise<AvailablePeriod[]> {
    const params = new URLSearchParams({ ticker });

    const response = await fetch(`${AI_BASE_URL}/analyses?${params}`, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
      signal,
    });

    if (!response.ok) {
      const body = await response.json().catch(() => ({}));
      throw new Error(body.error || `Failed to fetch available periods (${response.status})`);
    }

    const json = await response.json();
    return json.data || [];
  },

  async getAnalysis(
    ticker: string,
    year: number,
    period: number,
    signal?: AbortSignal
  ): Promise<string> {
    const params = new URLSearchParams({
      ticker,
      year: year.toString(),
      period: period.toString(),
    });

    const response = await fetch(`${AI_BASE_URL}/analysis?${params}`, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
      signal,
    });

    if (!response.ok) {
      const body = await response.json().catch(() => ({}));
      throw new Error(body.error || `Failed to fetch analysis (${response.status})`);
    }

    const json = await response.json();
    return json.data;
  },

  async getBusinessResearch(
    ticker: string,
    signal?: AbortSignal
  ): Promise<BusinessResearch | null> {
    const params = new URLSearchParams({ ticker });

    const response = await fetch(`${AI_BASE_URL}/business-research?${params}`, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
      signal,
    });

    if (response.status === 404) return null;
    if (!response.ok) return null;

    const json = await response.json();
    return json.data || null;
  },

  async getReportResults(
    ticker: string,
    signal?: AbortSignal
  ): Promise<ReportResults | null> {
    const params = new URLSearchParams({ ticker });

    const response = await fetch(`${AI_BASE_URL}/report-results/latest?${params}`, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
      signal,
    });

    if (response.status === 404) {
      return null;
    }

    if (!response.ok) {
      return null;
    }

    const json = await response.json();
    return json.data || null;
  },
};
