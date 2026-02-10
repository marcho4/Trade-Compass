'use client';

import { useState, useEffect, useCallback } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { ExtractControls } from '@/components/admin/ExtractControls';
import { RawDataForm } from '@/components/admin/RawDataForm';
import { aiApi } from '@/lib/api/ai-api';
import { financialDataApi, Company } from '@/lib/api/financial-data-api';
import { RawData, PERIOD_TO_FD } from '@/types/raw-data';
import { Loader2, Save, CheckCircle, Download } from 'lucide-react';

export default function AdminPage() {
  const [companies, setCompanies] = useState<Company[]>([]);
  const [ticker, setTicker] = useState('');
  const [period, setPeriod] = useState('12');
  const [year, setYear] = useState(new Date().getFullYear().toString());

  const [rawData, setRawData] = useState<RawData | null>(null);
  const [isExtracting, setIsExtracting] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const [isConfirming, setIsConfirming] = useState(false);
  const [isLoadingConfirmed, setIsLoadingConfirmed] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  useEffect(() => {
    financialDataApi
      .getCompanies()
      .then(setCompanies)
      .catch((err) => setError(`Failed to load companies: ${err.message}`));
  }, []);

  const clearMessages = () => {
    setError(null);
    setSuccess(null);
  };

  const handleExtract = useCallback(
    async (force = false) => {
      clearMessages();
      setIsExtracting(true);

      try {
        const data = await aiApi.extractData(
          ticker,
          period,
          parseInt(year, 10),
          force
        );
        setRawData(data);
        setSuccess(
          force ? 'Данные извлечены заново' : 'Данные загружены'
        );
      } catch (err: unknown) {
        const msg = err instanceof Error ? err.message : 'Unknown error';
        setError(msg);
      } finally {
        setIsExtracting(false);
      }
    },
    [ticker, period, year]
  );

  const handleSave = async () => {
    if (!rawData) return;
    clearMessages();
    setIsSaving(true);

    try {
      const fdPeriod = PERIOD_TO_FD[period] || 'YEAR';
      await financialDataApi.updateRawData(
        ticker,
        parseInt(year, 10),
        fdPeriod,
        rawData
      );
      setSuccess('Данные сохранены');
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Unknown error';
      setError(`Ошибка сохранения: ${msg}`);
    } finally {
      setIsSaving(false);
    }
  };

  const handleConfirm = async () => {
    if (!rawData) return;
    clearMessages();
    setIsConfirming(true);

    try {
      const fdPeriod = PERIOD_TO_FD[period] || 'YEAR';
      await financialDataApi.confirmDraft(ticker, parseInt(year, 10), fdPeriod);
      setRawData({ ...rawData, status: 'confirmed' });
      setSuccess('Данные подтверждены и доступны в основном API');
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Unknown error';
      setError(`Ошибка подтверждения: ${msg}`);
    } finally {
      setIsConfirming(false);
    }
  };

  const handleLoadConfirmed = async () => {
    clearMessages();
    setIsLoadingConfirmed(true);

    try {
      const fdPeriod = PERIOD_TO_FD[period] || 'YEAR';
      const data = await financialDataApi.getRawData(
        ticker,
        parseInt(year, 10),
        fdPeriod
      );
      setRawData(data);
      setSuccess('Загружены подтверждённые данные');
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Unknown error';
      setError(`Подтверждённые данные не найдены: ${msg}`);
    } finally {
      setIsLoadingConfirmed(false);
    }
  };

  return (
    <div className="max-w-6xl mx-auto space-y-6">
      <div>
        <h1 className="text-2xl font-bold">AI Data Extractor</h1>
        <p className="text-sm text-muted-foreground mt-1">
          Извлечение финансовых данных из отчетов с помощью AI
        </p>
      </div>

      <Card>
        <CardHeader className="pb-4">
          <CardTitle className="text-base">Параметры</CardTitle>
        </CardHeader>
        <CardContent>
          <ExtractControls
            companies={companies}
            ticker={ticker}
            period={period}
            year={year}
            onTickerChange={setTicker}
            onPeriodChange={setPeriod}
            onYearChange={setYear}
            onExtract={() => handleExtract(false)}
            onForceExtract={() => handleExtract(true)}
            isExtracting={isExtracting}
          />
        </CardContent>
      </Card>

      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      {success && (
        <Alert>
          <AlertDescription className="text-green-600">{success}</AlertDescription>
        </Alert>
      )}

      {rawData && (
        <Card>
          <CardHeader className="pb-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <CardTitle className="text-base">
                  {rawData.ticker} / {rawData.year} / {rawData.period}
                </CardTitle>
                <Badge
                  variant={rawData.status === 'confirmed' ? 'default' : 'secondary'}
                  className={
                    rawData.status === 'draft'
                      ? 'bg-yellow-500/15 text-yellow-600 border-yellow-500/30'
                      : 'bg-green-500/15 text-green-600 border-green-500/30'
                  }
                >
                  {rawData.status === 'draft' ? 'Draft' : 'Confirmed'}
                </Badge>
              </div>

              <div className="flex items-center gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handleLoadConfirmed}
                  disabled={isLoadingConfirmed}
                >
                  {isLoadingConfirmed ? (
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  ) : (
                    <Download className="mr-2 h-4 w-4" />
                  )}
                  Load confirmed
                </Button>

                <Button
                  variant="outline"
                  size="sm"
                  onClick={handleSave}
                  disabled={isSaving}
                >
                  {isSaving ? (
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  ) : (
                    <Save className="mr-2 h-4 w-4" />
                  )}
                  Save
                </Button>

                {rawData.status === 'draft' && (
                  <Button
                    size="sm"
                    onClick={handleConfirm}
                    disabled={isConfirming}
                  >
                    {isConfirming ? (
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    ) : (
                      <CheckCircle className="mr-2 h-4 w-4" />
                    )}
                    Confirm
                  </Button>
                )}
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <RawDataForm
              data={rawData}
              onChange={setRawData}
              disabled={isExtracting}
            />
          </CardContent>
        </Card>
      )}
    </div>
  );
}
