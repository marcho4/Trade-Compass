# Структура базы данных

```dbml
Table reports {
  id serial pk
  report_year integer
  report_period enum
  report_storage_url text
  company_id integer
}

Table companies {
  id serial pk
  inn integer
  ticker varchar(4)
  owner text
  sector_id integer
  lot_size integer
  ceo varchar(100)
  employees int
}

Table sectors {
  id serial pk
  name varchar(50)
}

Table metrics {
  report_id int
  
  // P&L (Отчёт о прибылях и убытках)
  revenue bigint // Выручка
  cost_of_revenue bigint // Себестоимость
  gross_profit bigint // Валовая прибыль
  operating_expenses bigint // Операционные расходы
  ebit bigint // Прибыль до вычета процентов и налогов
  ebitda bigint // EBITDA
  interest_expense bigint // Проценты к уплате
  tax_expense bigint // Налоги
  net_profit bigint // Чистая прибыль
  
  // Balance Sheet (Баланс)
  total_assets bigint // Всего активов
  current_assets bigint // Оборотные активы
  cash_and_equivalents bigint // Денежные средства и эквиваленты
  inventories bigint // Запасы
  receivables bigint // Дебиторская задолженность
  
  total_liabilities bigint // Всего обязательств
  current_liabilities bigint // Краткосрочные обязательства
  debt bigint // Долг (краткосрочный + долгосрочный)
  long_term_debt bigint // Долгосрочный долг
  short_term_debt bigint // Краткосрочный долг
  
  equity bigint // Собственный капитал
  retained_earnings bigint // Нераспределённая прибыль
  
  // Cash Flow Statement (Отчёт о движении денежных средств)
  operating_cash_flow bigint // Операционный денежный поток
  investing_cash_flow bigint // Инвестиционный денежный поток
  financing_cash_flow bigint // Финансовый денежный поток
  capex bigint // Капитальные затраты
  free_cash_flow bigint // Свободный денежный поток (OCF - CapEx)
  
  // Market Data (для мультипликаторов)
  shares_outstanding bigint // Количество акций в обращении
  market_cap bigint // Рыночная капитализация на дату отчёта
  
  // Дополнительные расчётные поля
  working_capital bigint // Оборотный капитал (current_assets - current_liabilities)
  capital_employed bigint // Задействованный капитал (total_assets - current_liabilities)
  enterprise_value bigint // EV = market_cap + debt - cash
  net_debt bigint // Чистый долг (debt - cash)
}

Table indicators {
  report_id int 
  PE decimal(10, 2)
  PB decimal(10, 2)
  PS decimal(10, 2)
  PV decimal(10, 2)
  EPS decimal(10, 2)
  ROE decimal(10, 2)
  ROA decimal(10, 2)
  ROCE decimal(10, 2)
  margin decimal(10, 2)
  ICR decimal(10, 2)
  COR decimal(10, 2)
  debt_to_equity decimal(10, 2)
  ev_ebitda decimal(10,2)

  dividend_yield decimal(5,2)
  fcf_yield decimal(5,2) // Free Cash Flow Yield
  current_ratio decimal(5,2) // ликвидность
  quick_ratio decimal(5,2)
  debt_to_ebitda decimal(5,2)
  interest_coverage decimal(5,2)
  working_capital bigint
  capex bigint
  fcf bigint
  
  revenue_growth_yoy decimal(5,2)
  profit_growth_yoy decimal(5,2)
  
  // Качество прибыли
  accruals_ratio decimal(5,2)
}

Table dividends {
  id serial pk
  company_id int
  ex_dividend_date date
  payment_date date
  amount_per_share decimal(10,2)
  dividend_yield decimal(5,2)
  payout_ratio decimal(5,2)
  currency varchar(3)
}

Table ai_analyses {
  id serial pk
  company_id int
  report_id int
  analysis_type enum // 'full', 'quick', 'term_explanation'
  prompt_hash varchar(64) // для дедупликации
  response_text text
  tokens_used int
  cost_usd decimal(8,4)
  created_at timestamp
  cache_until timestamp
}

Table portfolio {
  id serial pk
  name varchar(100)
  user_id int [ref: > users.id, not null]
  description text

  created_at timestamptz [default: `now()`]
  updated_at timestamptz

  indexes {
    user_id
    (user_id, name) [unique]
  }
}

Table position {
  id serial pk
  portfolio_id int
  company_id int
  avg_price decimal(10, 2)
  last_buy_date timestamptz // для рассчета налогов при ребалансировке.
  quantity int // кол-во лотов

  created_at timestamptz [default: `now()`]
  updated_at timestamptz

  indexes {
    portfolio_id
    (portfolio_id, company_id) [unique]
  }
}

Table users {
  id serial pk
  name varchar(50)
  last_login_at timestamptz
  status varchar(20) [default: 'active'] // active, blocked, deleted

  created_at timestamptz [default: `now()`]
  updated_at timestamptz
  
  indexes {
    created_at
  }
}

Table auth {
  user_id int pk [ref: - users.id]
  email varchar(100) [unique, not null]
  hashed_password varchar(100)

  indexes {
    email
  }
}

Table provider_auth {
  id serial pk
  user_id int [ref: > users.id, not null]
  provider_user_id varchar(255) [not null]
  provider_type varchar(20) [not null]
  
  indexes {
    (provider_type, provider_user_id) [unique]
  }
}

Table subscriptions {
  id serial pk
  user_id int [ref: - users.id, not null]
  start_date timestamptz
  end_date timestamptz
  level varchar(20)

  indexes {
    user_id
    level
  }
}

Ref: reports.company_id > companies.id
Ref: companies.sector_id > sectors.id
Ref: metrics.report_id > reports.id
Ref: "reports"."id" < "indicators"."report_id"
Ref: "dividends"."company_id" > "companies"."id"
Ref: "portfolio"."id" < "position"."portfolio_id"
Ref: "position"."company_id" - "companies"."id"
```
