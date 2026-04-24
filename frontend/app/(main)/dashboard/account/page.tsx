"use client"

import { useAuth } from "@/contexts/AuthContext"
import { UserCard } from "@/components/user/UserCard"
import { LimitsCard } from "@/components/user/LimitsCard"

const HEATMAP = [
  1, 2, 0, 3, 4, 2, 1, 0, 1, 3, 2, 4, 3, 2, 1, 0, 2, 3, 2, 1, 0, 2, 3, 4, 3, 2, 1, 0, 1, 2,
]

const LIMITS = {
  plan: "PRO",
  period: "Месячный лимит",
  used: 184,
  total: 300,
  resetsAt: "01 мая 2026, 00:00 МСК",
  daysToReset: 9,
  hoursToReset: 14,
  dailyUsed: 12,
  dailyTotal: 25,
  savedUsed: 47,
  savedTotal: 100,
  screenersUsed: 6,
  screenersTotal: 20,
  peakPerDay: 18,
  peakDate: "14 апр",
  avgPerDay: 6.1,
  heatmap: HEATMAP,
}

const RECENT = [
  { sym: "SBER", name: "Сбербанк", sector: "Финансы", when: "14:22 · сегодня", views: 4 },
  { sym: "GAZP", name: "Газпром", sector: "Нефть/Газ", when: "11:08 · сегодня", views: 2 },
  { sym: "LKOH", name: "Лукойл", sector: "Нефть/Газ", when: "вчера · 18:44", views: 3 },
  { sym: "YDEX", name: "Яндекс", sector: "IT", when: "вчера · 10:12", views: 5 },
  { sym: "GMKN", name: "ГМК Норникель", sector: "Металлы", when: "19 апр · 16:50", views: 1 },
  { sym: "MGNT", name: "Магнит", sector: "Ритейл", when: "19 апр · 09:30", views: 2 },
  { sym: "ROSN", name: "Роснефть", sector: "Нефть/Газ", when: "18 апр · 13:18", views: 1 },
  { sym: "PLZL", name: "Полюс", sector: "Металлы", when: "17 апр · 20:04", views: 3 },
]

const PLANS = [
  { id: "FREE", views: 50, daily: 5, saved: 10, screeners: 3, price: "0 ₽" },
  { id: "PRO", views: 300, daily: 25, saved: 100, screeners: 20, price: "590 ₽" },
  { id: "PRO+", views: 1500, daily: "∞", saved: 500, screeners: 100, price: "1 490 ₽" },
]

function buildInitials(name?: string): string {
  if (!name) return "——"
  const parts = name.trim().split(/\s+/)
  const first = parts[0]?.[0] ?? ""
  const second = parts[1]?.[0] ?? ""
  return (first + second).toUpperCase() || name[0].toUpperCase()
}

function buildHandle(email?: string): string {
  if (!email) return "@user"
  return "@" + email.split("@")[0]
}

export default function AccountPage() {
  const { user, isLoading, logout } = useAuth()

  if (isLoading) {
    return (
      <div className="flex items-center justify-center p-12">
        <p className="text-muted-foreground">Загрузка...</p>
      </div>
    )
  }

  const userData = {
    initials: buildInitials(user?.name),
    name: user?.name || "Гость",
    handle: buildHandle(user?.email),
    email: user?.email || "—",
    phone: "+7 (916) 234-18-09",
    plan: "PRO",
    planSince: "12 авг 2024",
    memberSince: "03 мар 2023",
    tz: "Europe/Moscow · UTC+3",
    lang: "Русский",
    twoFA: true,
    verified: true,
    broker: "Т-Инвестиции",
    brokerLinked: "14 сен 2025",
    portfolioSize: "2 340 000 ₽",
    riskProfile: "Умеренно-агрессивный",
    userId: `USR-${String(user?.id ?? 41882).padStart(6, "0")}`,
    lastLogin: "22.04.2026 · 10:02",
  }

  return (
    <div className="mx-auto max-w-[1280px] px-2 pt-1 pb-10">
      <div className="mb-4">
        <div className="mb-1 font-mono text-[10px] font-bold uppercase tracking-[0.12em] text-primary">
          // Аккаунт
        </div>
        <div className="font-sans text-[22px] font-semibold tracking-[-0.015em] text-foreground">
          Профиль и лимиты просмотра компаний
        </div>
        <div className="mt-1 max-w-[720px] font-sans text-[12.5px] leading-[1.45] text-muted-foreground">
          Личные данные, связанный брокерский счёт, квоты по тарифу PRO и журнал
          последних просмотренных тикеров. Сброс квоты — раз в месяц.
        </div>
      </div>

      <div className="grid gap-4">
        <UserCard user={userData} onLogout={logout} />
        <LimitsCard limits={LIMITS} recent={RECENT} plans={PLANS} />
      </div>

      <div className="mt-7 flex justify-between border-t border-border pt-3.5 font-mono text-[10px] uppercase tracking-[0.05em] text-muted-foreground/70">
        <span>Trade Compass · Account · {userData.userId}</span>
        <span>Build 2026.04.22</span>
      </div>
    </div>
  )
}
