"use client"

import { useEffect, useState } from "react"
import { aiApi, type BusinessResearch } from "@/lib/api/ai-api"
import { T } from "./tokens"
import { Chip, Donut, KV, TPanel } from "./primitives"

interface TickerTabAboutProps {
  ticker: string
}

const DONUT_COLORS = ["#b8751a", "#2563eb", "#7c3aed", "#14804a", "#6b7280", "#b42318", "#0ea5a3", "#d97706"]

const severityMap: Record<string, { label: string; color: string; soft: string }> = {
  critical: { label: "Крит", color: T.neg, soft: T.negSoft },
  high: { label: "Высокий", color: "#c55a18", soft: "#fcebdf" },
  moderate: { label: "Умер", color: T.textDim, soft: T.neuSoft },
}

const depTypeLabel: Record<string, string> = {
  commodity: "Сырьё",
  currency: "Валюта",
  regulation: "Регулирование",
  macro: "Макро",
  technology: "Технологии",
  geopolitics: "Геополитика",
  infrastructure: "Инфраструктура",
  demand: "Спрос",
}

export function TickerTabAbout({ ticker }: TickerTabAboutProps) {
  const [data, setData] = useState<BusinessResearch | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const controller = new AbortController()

    aiApi
      .getBusinessResearch(ticker, controller.signal)
      .then((result) => {
        setData(result)
      })
      .catch(() => {
        setData(null)
      })
      .finally(() => {
        setLoading(false)
      })

    return () => controller.abort()
  }, [ticker])

  const revenueDonut =
    data?.revenue_sources?.map((rs, i) => ({
      label: rs.segment,
      pct: rs.share_pct,
      color: DONUT_COLORS[i % DONUT_COLORS.length],
      description: rs.description,
    })) || []

  const totalRevenuePct = revenueDonut.reduce((s, r) => s + r.pct, 0)

  const profileSections = buildProfileSections(data)
  const dependencies = data?.dependencies || []

  if (loading) {
    return (
      <div
        style={{
          padding: 40,
          textAlign: "center",
          fontFamily: T.mono,
          color: T.textDim,
          fontSize: 11,
          letterSpacing: 1,
          textTransform: "uppercase",
          border: `1px solid ${T.border}`,
          background: T.panel,
        }}
      >
        Загрузка профиля...
      </div>
    )
  }

  if (!data) {
    return (
      <div
        style={{
          padding: 40,
          textAlign: "center",
          fontFamily: T.sans,
          color: T.textDim,
          fontSize: 13,
          border: `1px solid ${T.border}`,
          background: T.panel,
        }}
      >
        Информация о компании пока недоступна.
      </div>
    )
  }

  return (
    <div style={{ display: "grid", gridTemplateColumns: "minmax(0, 1.5fr) minmax(0, 1fr)", gap: 14 }}>
      <TPanel title="Профиль компании" accent={T.accent}>
        <div style={{ padding: "4px 16px 16px" }}>
          {profileSections.map((section, i) => (
            <div key={i} style={{ marginTop: 14 }}>
              <div
                style={{
                  fontFamily: T.mono,
                  fontSize: 10,
                  fontWeight: 700,
                  color: T.accent,
                  textTransform: "uppercase",
                  letterSpacing: 1.2,
                  paddingBottom: 6,
                  marginBottom: 4,
                  borderBottom: `1px solid ${T.borderSoft}`,
                  display: "flex",
                  alignItems: "center",
                  gap: 8,
                }}
              >
                <span style={{ color: T.textFaint }}>§ {String(i + 1).padStart(2, "0")}</span>
                <span>{section.label}</span>
              </div>
              {section.items.map(([k, v], j) => (
                <KV key={j} k={k} v={v} width={120} />
              ))}
            </div>
          ))}
        </div>
      </TPanel>

      <div style={{ display: "flex", flexDirection: "column", gap: 14 }}>
        <TPanel
          title="Структура выручки"
          accent={T.deps}
          right={
            <span
              style={{
                fontFamily: T.mono,
                fontSize: 10,
                color: T.textDim,
                letterSpacing: 0.5,
              }}
            >
              AI research
            </span>
          }
        >
          {revenueDonut.length === 0 ? (
            <div
              style={{
                padding: 24,
                textAlign: "center",
                color: T.textDim,
                fontSize: 12,
                fontFamily: T.sans,
              }}
            >
              Нет данных о выручке
            </div>
          ) : (
            <div style={{ padding: "18px 16px", display: "flex", alignItems: "center", gap: 20 }}>
              <Donut
                data={revenueDonut}
                size={170}
                thickness={26}
                centerLabel="СЕГМЕНТОВ"
                centerValue={String(revenueDonut.length)}
              />
              <div style={{ flex: 1, minWidth: 0 }}>
                {revenueDonut.map((d, i) => (
                  <div
                    key={i}
                    style={{
                      display: "grid",
                      gridTemplateColumns: "10px 1fr auto",
                      gap: 8,
                      alignItems: "center",
                      padding: "5px 0",
                      borderBottom:
                        i < revenueDonut.length - 1 ? `1px dashed ${T.borderSoft}` : "none",
                      fontSize: 11,
                    }}
                  >
                    <span style={{ width: 8, height: 8, background: d.color }} />
                    <span
                      style={{
                        color: T.text2,
                        lineHeight: 1.3,
                        overflow: "hidden",
                        textOverflow: "ellipsis",
                        whiteSpace: "nowrap",
                      }}
                      title={d.label}
                    >
                      {d.label}
                    </span>
                    <span
                      style={{
                        fontFamily: T.mono,
                        fontSize: 11,
                        fontWeight: 700,
                        color: T.text,
                        letterSpacing: 0.2,
                      }}
                    >
                      {String(d.pct).padStart(2, "0")}%
                    </span>
                  </div>
                ))}
                {totalRevenuePct > 0 && totalRevenuePct < 100 && (
                  <div
                    style={{
                      marginTop: 8,
                      paddingTop: 6,
                      borderTop: `1px solid ${T.borderSoft}`,
                      fontFamily: T.mono,
                      fontSize: 9,
                      color: T.textFaint,
                      letterSpacing: 0.3,
                    }}
                  >
                    Сумма: {totalRevenuePct}% · остальное не раскрыто
                  </div>
                )}
              </div>
            </div>
          )}
        </TPanel>

        <TPanel
          title="Ключевые зависимости"
          accent={T.hist}
          count={dependencies.length}
          right={
            <span
              style={{
                fontFamily: T.mono,
                fontSize: 10,
                color: T.textDim,
                letterSpacing: 0.5,
              }}
            >
              AI risk map
            </span>
          }
        >
          {dependencies.length === 0 ? (
            <div
              style={{
                padding: 24,
                textAlign: "center",
                color: T.textDim,
                fontSize: 12,
                fontFamily: T.sans,
              }}
            >
              Зависимости не выявлены
            </div>
          ) : (
            <div style={{ padding: "10px 14px 14px" }}>
              {dependencies.map((dep, i) => {
                const sev = severityMap[dep.severity] || severityMap.moderate
                return (
                  <div
                    key={`${dep.factor}-${i}`}
                    style={{
                      padding: "10px 0",
                      borderBottom:
                        i < dependencies.length - 1 ? `1px dashed ${T.borderSoft}` : "none",
                    }}
                  >
                    <div
                      style={{
                        display: "flex",
                        alignItems: "center",
                        justifyContent: "space-between",
                        gap: 8,
                        marginBottom: 4,
                      }}
                    >
                      <span
                        style={{
                          fontFamily: T.sans,
                          fontSize: 12.5,
                          fontWeight: 600,
                          color: T.text,
                        }}
                      >
                        {dep.factor}
                      </span>
                      <Chip color={sev.color} bg={sev.soft} bd={sev.color}>
                        {sev.label}
                      </Chip>
                    </div>
                    <div
                      style={{
                        display: "flex",
                        alignItems: "center",
                        gap: 8,
                        marginBottom: 6,
                      }}
                    >
                      <Chip>{depTypeLabel[dep.type] || dep.type}</Chip>
                    </div>
                    <div
                      style={{
                        fontFamily: T.sans,
                        fontSize: 12,
                        color: T.text2,
                        lineHeight: 1.45,
                      }}
                    >
                      {dep.description}
                    </div>
                  </div>
                )
              })}
            </div>
          )}
        </TPanel>
      </div>
    </div>
  )
}

function buildProfileSections(data: BusinessResearch | null) {
  if (!data) return []

  const profile = data.profile
  const sections: { label: string; items: [string, string][] }[] = []

  const products = profile.products_and_services || []
  if (products.length > 0) {
    sections.push({
      label: "Продукты и услуги",
      items: products.map((p, i) => [`#${String(i + 1).padStart(2, "0")}`, p] as [string, string]),
    })
  }

  const markets = profile.markets || []
  if (markets.length > 0) {
    sections.push({
      label: "Рынки",
      items: markets.map((m) => [m.market, m.role] as [string, string]),
    })
  }

  if (profile.key_clients) {
    sections.push({
      label: "Ключевые клиенты",
      items: [["Целевая группа", profile.key_clients]],
    })
  }

  if (profile.business_model) {
    sections.push({
      label: "Бизнес-модель",
      items: [["Подход", profile.business_model]],
    })
  }

  return sections
}
