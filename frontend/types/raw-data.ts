export type MetricsStatus = 'draft' | 'confirmed';

export interface RawData {
  ticker: string;
  year: number;
  period: string;
  status: MetricsStatus;

  revenue?: number | null;
  costOfRevenue?: number | null;
  grossProfit?: number | null;
  operatingExpenses?: number | null;
  ebit?: number | null;
  ebitda?: number | null;
  interestExpense?: number | null;
  taxExpense?: number | null;
  netProfit?: number | null;

  totalAssets?: number | null;
  currentAssets?: number | null;
  cashAndEquivalents?: number | null;
  inventories?: number | null;
  receivables?: number | null;

  totalLiabilities?: number | null;
  currentLiabilities?: number | null;
  debt?: number | null;
  longTermDebt?: number | null;
  shortTermDebt?: number | null;
  equity?: number | null;
  retainedEarnings?: number | null;

  operatingCashFlow?: number | null;
  investingCashFlow?: number | null;
  financingCashFlow?: number | null;
  capex?: number | null;
  freeCashFlow?: number | null;

  sharesOutstanding?: number | null;
  marketCap?: number | null;

  workingCapital?: number | null;
  capitalEmployed?: number | null;
  enterpriseValue?: number | null;
  netDebt?: number | null;
}

export interface MetricFieldConfig {
  key: keyof RawData;
  label: string;
}

export const PNL_FIELDS: MetricFieldConfig[] = [
  { key: 'revenue', label: 'Выручка' },
  { key: 'costOfRevenue', label: 'Себестоимость' },
  { key: 'grossProfit', label: 'Валовая прибыль' },
  { key: 'operatingExpenses', label: 'Операционные расходы' },
  { key: 'ebit', label: 'EBIT' },
  { key: 'ebitda', label: 'EBITDA' },
  { key: 'interestExpense', label: 'Проценты к уплате' },
  { key: 'taxExpense', label: 'Налог на прибыль' },
  { key: 'netProfit', label: 'Чистая прибыль' },
];

export const BALANCE_SHEET_FIELDS: MetricFieldConfig[] = [
  { key: 'totalAssets', label: 'Всего активов' },
  { key: 'currentAssets', label: 'Оборотные активы' },
  { key: 'cashAndEquivalents', label: 'Денежные средства' },
  { key: 'inventories', label: 'Запасы' },
  { key: 'receivables', label: 'Дебиторская задолженность' },
  { key: 'totalLiabilities', label: 'Всего обязательств' },
  { key: 'currentLiabilities', label: 'Краткосрочные обязательства' },
  { key: 'debt', label: 'Долг' },
  { key: 'longTermDebt', label: 'Долгосрочный долг' },
  { key: 'shortTermDebt', label: 'Краткосрочный долг' },
  { key: 'equity', label: 'Собственный капитал' },
  { key: 'retainedEarnings', label: 'Нераспределённая прибыль' },
];

export const CASH_FLOW_FIELDS: MetricFieldConfig[] = [
  { key: 'operatingCashFlow', label: 'Операционный ДП' },
  { key: 'investingCashFlow', label: 'Инвестиционный ДП' },
  { key: 'financingCashFlow', label: 'Финансовый ДП' },
  { key: 'capex', label: 'CAPEX' },
  { key: 'freeCashFlow', label: 'Свободный ДП' },
];

export const MARKET_DATA_FIELDS: MetricFieldConfig[] = [
  { key: 'sharesOutstanding', label: 'Акций в обращении' },
  { key: 'marketCap', label: 'Капитализация' },
];

export const CALCULATED_FIELDS: MetricFieldConfig[] = [
  { key: 'workingCapital', label: 'Оборотный капитал' },
  { key: 'capitalEmployed', label: 'Задействованный капитал' },
  { key: 'enterpriseValue', label: 'Enterprise Value' },
  { key: 'netDebt', label: 'Чистый долг' },
];

export const PERIOD_OPTIONS = [
  { value: '3', label: '3 мес. (Q1)' },
  { value: '6', label: '6 мес. (Q2)' },
  { value: '9', label: '9 мес. (Q3)' },
  { value: '12', label: '12 мес. (Год)' },
];

export const PERIOD_TO_FD: Record<string, string> = {
  '3': 'Q1',
  '6': 'Q2',
  '9': 'Q3',
  '12': 'YEAR',
};
