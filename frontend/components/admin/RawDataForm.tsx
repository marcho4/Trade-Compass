'use client';

import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  RawData,
  MetricFieldConfig,
  BANK_FIELDS,
  PNL_FIELDS,
  BALANCE_SHEET_FIELDS,
  CASH_FLOW_FIELDS,
  MARKET_DATA_FIELDS,
  CALCULATED_FIELDS,
} from '@/types/raw-data';

interface RawDataFormProps {
  data: RawData;
  onChange: (data: RawData) => void;
  disabled?: boolean;
}

function MetricGroup({
  title,
  fields,
  data,
  onChange,
  disabled,
}: {
  title: string;
  fields: MetricFieldConfig[];
  data: RawData;
  onChange: (key: keyof RawData, value: number | null) => void;
  disabled?: boolean;
}) {
  return (
    <div className="space-y-3">
      <h3 className="text-sm font-semibold text-muted-foreground uppercase tracking-wider">
        {title}
      </h3>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
        {fields.map((field) => {
          const value = data[field.key] as number | null | undefined;
          const isEmpty = value === null || value === undefined;

          return (
            <div key={field.key} className="space-y-1">
              <Label
                htmlFor={field.key}
                className={`text-xs ${isEmpty ? 'text-orange-500' : 'text-muted-foreground'}`}
              >
                {field.label}
              </Label>
              <Input
                id={field.key}
                type="number"
                value={value ?? ''}
                onChange={(e) => {
                  const val = e.target.value;
                  onChange(field.key, val === '' ? null : parseInt(val, 10));
                }}
                disabled={disabled}
                className={`h-9 text-sm ${isEmpty ? 'border-orange-500/50' : ''}`}
                placeholder="null"
              />
            </div>
          );
        })}
      </div>
    </div>
  );
}

export function RawDataForm({ data, onChange, disabled }: RawDataFormProps) {
  const handleFieldChange = (key: keyof RawData, value: number | null) => {
    onChange({ ...data, [key]: value });
  };

  const isBank = data.companyType === 'bank';

  const groups = [
    { title: 'P&L (Прибыли и убытки)', fields: PNL_FIELDS },
    { title: 'Баланс', fields: BALANCE_SHEET_FIELDS },
    { title: 'Денежный поток', fields: CASH_FLOW_FIELDS },
    { title: 'Рыночные данные', fields: MARKET_DATA_FIELDS },
    { title: 'Расчётные показатели', fields: CALCULATED_FIELDS },
  ];

  return (
    <div className="space-y-6">
      <div className="space-y-2">
        <Label htmlFor="companyType" className="text-xs text-muted-foreground uppercase tracking-wider font-semibold">
          Тип компании
        </Label>
        <select
          id="companyType"
          value={data.companyType ?? 'corporate'}
          onChange={(e) => onChange({ ...data, companyType: e.target.value })}
          disabled={disabled}
          className="h-9 w-48 rounded-md border border-input bg-background px-3 text-sm focus:outline-none focus:ring-1 focus:ring-ring"
        >
          <option value="corporate">Corporate</option>
          <option value="bank">Bank</option>
        </select>
      </div>

      {isBank && (
        <MetricGroup
          title="Банковские показатели"
          fields={BANK_FIELDS}
          data={data}
          onChange={handleFieldChange}
          disabled={disabled}
        />
      )}

      {groups.map((group) => (
        <MetricGroup
          key={group.title}
          title={group.title}
          fields={group.fields}
          data={data}
          onChange={handleFieldChange}
          disabled={disabled}
        />
      ))}
    </div>
  );
}
