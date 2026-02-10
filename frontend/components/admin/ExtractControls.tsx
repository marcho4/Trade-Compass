'use client';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { PERIOD_OPTIONS } from '@/types/raw-data';
import { Company } from '@/lib/api/financial-data-api';
import { Loader2 } from 'lucide-react';

interface ExtractControlsProps {
  companies: Company[];
  ticker: string;
  period: string;
  year: string;
  onTickerChange: (ticker: string) => void;
  onPeriodChange: (period: string) => void;
  onYearChange: (year: string) => void;
  onExtract: () => void;
  onForceExtract: () => void;
  isExtracting: boolean;
}

export function ExtractControls({
  companies,
  ticker,
  period,
  year,
  onTickerChange,
  onPeriodChange,
  onYearChange,
  onExtract,
  onForceExtract,
  isExtracting,
}: ExtractControlsProps) {
  const canExtract = ticker && period && year && !isExtracting;

  return (
    <div className="flex flex-wrap items-end gap-4">
      <div className="space-y-1.5 min-w-[180px]">
        <Label className="text-xs text-muted-foreground">Тикер</Label>
        <Select value={ticker} onValueChange={onTickerChange}>
          <SelectTrigger className="h-9">
            <SelectValue placeholder="Выберите компанию" />
          </SelectTrigger>
          <SelectContent>
            {companies.map((c) => (
              <SelectItem key={c.ticker} value={c.ticker}>
                {c.ticker} {c.name ? `— ${c.name}` : ''}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      <div className="space-y-1.5 min-w-[150px]">
        <Label className="text-xs text-muted-foreground">Период</Label>
        <Select value={period} onValueChange={onPeriodChange}>
          <SelectTrigger className="h-9">
            <SelectValue placeholder="Период" />
          </SelectTrigger>
          <SelectContent>
            {PERIOD_OPTIONS.map((p) => (
              <SelectItem key={p.value} value={p.value}>
                {p.label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      <div className="space-y-1.5 w-[100px]">
        <Label className="text-xs text-muted-foreground">Год</Label>
        <Input
          type="number"
          value={year}
          onChange={(e) => onYearChange(e.target.value)}
          placeholder="2024"
          className="h-9"
          min={2000}
          max={2100}
        />
      </div>

      <Button onClick={onExtract} disabled={!canExtract} className="h-9">
        {isExtracting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
        Extract
      </Button>

      <Button
        onClick={onForceExtract}
        disabled={!canExtract}
        variant="outline"
        className="h-9"
      >
        {isExtracting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
        Force Re-extract
      </Button>
    </div>
  );
}
