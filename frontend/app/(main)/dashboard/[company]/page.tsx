"use client"

import { use, useEffect, useMemo, useState } from "react"
import { notFound } from "next/navigation"
import {
  CompanyReports,
  CompanyAnalyses,
  CompanyNews,
} from "@/components/company"
import {
  T,
  TickerTopBar,
  TickerHero,
  TickerTabStrip,
  TickerFooter,
  TickerTabAbout,
  TickerTabMetrics,
  TPanel,
  type TickerHeroData,
  type TickerTab,
} from "@/components/company/ticker"
import { aiApi } from "@/lib/api/ai-api"
import { financialDataApi, type Sector } from "@/lib/api"
import { Company as CompanyType } from "@/types"
import { useTickerScreenData, formatShortNumber } from "@/hooks/use-ticker-screen-data"

type PageProps = {
  params: Promise<{
    company: string
  }>
}

function formatMoscowTime(d: Date): string {
  return new Intl.DateTimeFormat("ru-RU", {
    timeZone: "Europe/Moscow",
    day: "2-digit",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
  })
    .format(d)
    .replace(",", " ·")
    .toUpperCase()
}

function isMoexOpen(d: Date): boolean {
  const moscowHour = Number(
    new Intl.DateTimeFormat("en-GB", {
      timeZone: "Europe/Moscow",
      hour: "2-digit",
      hour12: false,
    }).format(d),
  )
  const isWeekend = d.getUTCDay() === 0 || d.getUTCDay() === 6
  return !isWeekend && moscowHour >= 10 && moscowHour < 24
}

const CompanyDashboardPage = ({ params }: PageProps) => {
  const { company: ticker } = use(params)
  const decodedTicker = decodeURIComponent(ticker).toUpperCase()

  const [company, setCompany] = useState<CompanyType | null>(null)
  const [companyName, setCompanyName] = useState<string>("")
  const [description, setDescription] = useState<string>("")
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [tab, setTab] = useState<string>("about")

  const screen = useTickerScreenData(decodedTicker)

  useEffect(() => {
    const controller = new AbortController()

    const loadCompany = async () => {
      try {
        const [companyData, sectorsData] = await Promise.all([
          financialDataApi.getCompanyByTicker(decodedTicker),
          financialDataApi.getSectors(),
        ])

        const sector = sectorsData.find((s: Sector) => s.id === companyData.sectorId)

        setCompany({
          id: companyData.id,
          ticker: companyData.ticker,
          sectorId: companyData.sectorId,
          sector: sector?.name,
          lotSize: companyData.lotSize,
          ceo: companyData.ceo,
        })
      } catch (err) {
        console.error("Failed to load company:", err)
        setError("Компания не найдена")
      } finally {
        setLoading(false)
      }
    }

    loadCompany()

    aiApi
      .getBusinessResearch(decodedTicker, controller.signal)
      .then((data) => {
        if (data) {
          setCompanyName(data.profile.company_name || "")
          setDescription(data.profile.description || "")
        }
      })
      .catch(() => {})

    return () => controller.abort()
  }, [decodedTicker])

  const now = useMemo(() => new Date(), [])
  const marketOpen = useMemo(() => isMoexOpen(now), [now])
  const marketLabel = useMemo(
    () => `MOEX ${marketOpen ? "OPEN" : "CLOSED"} · ${formatMoscowTime(now)} МСК`,
    [now, marketOpen],
  )

  const lastUpdate = useMemo(
    () =>
      new Intl.DateTimeFormat("ru-RU", {
        timeZone: "Europe/Moscow",
        hour: "2-digit",
        minute: "2-digit",
        hour12: false,
      }).format(now) + " МСК",
    [now],
  )

  const buildTag = useMemo(() => {
    const y = now.getFullYear()
    const m = String(now.getMonth() + 1).padStart(2, "0")
    const d = String(now.getDate()).padStart(2, "0")
    return `${y}.${m}.${d}`
  }, [now])

  const heroData: TickerHeroData = useMemo(
    () => ({
      symbol: decodedTicker,
      exchange: "MOEX",
      currency: "RUB",
      isin: "—",
      name: companyName,
      sector: company?.sector || "",
      industry: "",
      description,
      price: screen.price,
      change: screen.change,
      changePct: screen.changePct,
      lastUpdate,
      w52Low: screen.w52Low,
      w52High: screen.w52High,
      open: screen.open,
      prevClose: screen.prevClose,
      dayLow: screen.dayLow,
      dayHigh: screen.dayHigh,
      volumeShares: screen.volumeShares,
      turnoverRub: screen.turnoverRub,
      marketCap: screen.marketCapRub ? formatShortNumber(screen.marketCapRub, true) : "—",
      freeFloat: "—",
      pe: null,
      pb: null,
      divYield: "—",
      beta: null,
      loading: screen.loading,
    }),
    [decodedTicker, company, companyName, description, screen, lastUpdate],
  )

  const tabs: TickerTab[] = useMemo(
    () => [
      { id: "about", label: "О компании" },
      { id: "metrics", label: "Показатели" },
      { id: "ai", label: "AI Анализ" },
      { id: "news", label: "Новости" },
      { id: "reports", label: "Отчёты" },
    ],
    [],
  )

  if (loading) {
    return (
      <div
        style={{
          padding: 40,
          textAlign: "center",
          fontFamily: T.mono,
          color: T.textDim,
          letterSpacing: 1,
          textTransform: "uppercase",
          fontSize: 11,
        }}
      >
        Загрузка...
      </div>
    )
  }

  if (error || !company) {
    notFound()
  }

  return (
    <div
      style={{
        background: T.bg,
        fontFamily: T.sans,
        color: T.text,
        margin: "-16px -16px 0",
        minHeight: "calc(100vh - 64px)",
      }}
    >
      <TickerTopBar ticker={decodedTicker} marketOpen={marketOpen} marketLabel={marketLabel} />

      <div
        style={{
          padding: "18px 28px",
          maxWidth: 1600,
          margin: "0 auto",
        }}
      >
        <TickerHero data={heroData} />

        <div style={{ marginTop: 18 }}>
          <TickerTabStrip tabs={tabs} active={tab} onChange={setTab} />
        </div>

        <div style={{ marginTop: 16 }}>
          {tab === "about" && <TickerTabAbout ticker={decodedTicker} />}
          {tab === "metrics" && <TickerTabMetrics ticker={decodedTicker} />}
          {tab === "ai" && (
            <TPanel title="AI Анализ" accent={T.accent}>
              <div style={{ padding: 16, background: T.panel }}>
                <CompanyAnalyses ticker={decodedTicker} />
              </div>
            </TPanel>
          )}
          {tab === "news" && (
            <TPanel title="Новости и события" accent={T.deps}>
              <div style={{ padding: 16, background: T.panel }}>
                <CompanyNews ticker={decodedTicker} />
              </div>
            </TPanel>
          )}
          {tab === "reports" && (
            <TPanel title="Отчёты компании" accent={T.hist}>
              <div style={{ padding: 16, background: T.panel }}>
                <CompanyReports ticker={decodedTicker} />
              </div>
            </TPanel>
          )}
        </div>

        <TickerFooter ticker={decodedTicker} buildTag={buildTag} />
      </div>
    </div>
  )
}

export default CompanyDashboardPage
