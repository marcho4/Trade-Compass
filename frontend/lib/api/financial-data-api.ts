export interface Sector {
  id: number;
  name: string;
}

export interface Company {
  id: number;
  ticker: string;
  name?: string;
  sectorId: number;
  lotSize?: number;
  ceo?: string;
}

export interface Candle {
  open: number;
  close: number;
  high: number;
  low: number;
  value: number;
  volume: number;
  begin: string;
  end: string;
}

interface ApiResponse<T> {
  status: string;
  data: T;
  message?: string;
}

const FINANCIAL_DATA_BASE_URL = 'https://trade-compass.ru/api/financial-data';

export const financialDataApi = {
  async getSectors(): Promise<Sector[]> {
    const response = await fetch(`${FINANCIAL_DATA_BASE_URL}/sectors`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error('Failed to fetch sectors');
    }

    const result: ApiResponse<Sector[]> = await response.json();
    return result.data;
  },

  async getCompanies(): Promise<Company[]> {
    const response = await fetch(`${FINANCIAL_DATA_BASE_URL}/companies`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error('Failed to fetch companies');
    }

    const result: ApiResponse<Company[]> = await response.json();
    return result.data;
  },

  async getCompanyByTicker(ticker: string): Promise<Company> {
    const response = await fetch(`${FINANCIAL_DATA_BASE_URL}/companies/${ticker}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch company ${ticker}`);
    }

    const result: ApiResponse<Company> = await response.json();
    return result.data;
  },

  async getCompaniesBySector(sectorId: number): Promise<Company[]> {
    const response = await fetch(`${FINANCIAL_DATA_BASE_URL}/companies/sector/${sectorId}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch companies for sector ${sectorId}`);
    }

    const result: ApiResponse<Company[]> = await response.json();
    return result.data;
  },

  async getPriceCandles(ticker: string, days: number, interval: number, signal?: AbortSignal): Promise<Candle[]> {
    const response = await fetch(
      `${FINANCIAL_DATA_BASE_URL}/price?ticker=${ticker}&days=${days}&interval=${interval}`,
      {
        method: 'GET',
        headers: { 'Content-Type': 'application/json' },
        signal,
      }
    );

    if (!response.ok) {
      throw new Error(`Failed to fetch price candles for ${ticker}`);
    }

    const result: ApiResponse<Candle[]> = await response.json();
    return result.data;
  },

  async getLatestPrice(ticker: string, signal?: AbortSignal): Promise<number> {
    const response = await fetch(`${FINANCIAL_DATA_BASE_URL}/price/latest?ticker=${ticker}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      signal,
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch latest price for ${ticker}`);
    }

    const result: ApiResponse<number> = await response.json();
    return result.data;
  },

  async getMarketCap(ticker: string, signal?: AbortSignal): Promise<number> {
    const response = await fetch(`${FINANCIAL_DATA_BASE_URL}/market-cap?ticker=${ticker}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      signal,
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch market cap for ${ticker}`);
    }

    const result: ApiResponse<number> = await response.json();
    return result.data;
  },
};
