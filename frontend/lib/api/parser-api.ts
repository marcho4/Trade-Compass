import type { Report, ReportsResponse } from '@/types';

const PARSER_BASE_URL = 'https://trade-compass.ru/api/parser';

export const parserApi = {
  async getReportsByTicker(ticker: string, signal?: AbortSignal): Promise<Report[]> {
    const response = await fetch(`${PARSER_BASE_URL}/reports/${ticker}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      signal,
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch reports for ${ticker}`);
    }

    const result: ReportsResponse = await response.json();
    return result.reports;
  },
};
