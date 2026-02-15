import { RawData } from '@/types/raw-data';

const AI_BASE_URL = '/api/ai';

export interface AnalysisReport {
  id: number;
  ticker: string;
  year: number;
  period: number;
  analysis: string;
  created_at: string;
  updated_at: string;
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

  async getAnalysesByTicker(
    ticker: string,
    signal?: AbortSignal
  ): Promise<AnalysisReport[]> {
    const params = new URLSearchParams({ ticker });

    const response = await fetch(`${AI_BASE_URL}/analyses?${params}`, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
      signal,
    });

    if (!response.ok) {
      const body = await response.json().catch(() => ({}));
      throw new Error(body.error || `Failed to fetch analyses (${response.status})`);
    }

    const json = await response.json();
    return json.data || [];
  },
};
