"use client"

import { ReactNode } from "react"
import { cn } from "@/lib/utils"

interface UserData {
  initials: string
  name: string
  handle: string
  email: string
  phone: string
  plan: string
  planSince: string
  memberSince: string
  tz: string
  lang: string
  twoFA: boolean
  verified: boolean
  broker: string
  brokerLinked: string
  portfolioSize: string
  riskProfile: string
  userId: string
  lastLogin: string
}

interface UserCardProps {
  user: UserData
  onEdit?: () => void
  onLogout?: () => void
  onChangePassword?: () => void
  onNotifications?: () => void
  onIntegrations?: () => void
}

function SectionLabel({ children, right }: { children: ReactNode; right?: ReactNode }) {
  return (
    <div className="flex items-center justify-between border-b border-border bg-muted/40 px-3.5 py-2.5">
      <div className="flex items-center gap-2.5">
        <span className="h-3 w-[3px] bg-primary" />
        <span className="font-mono text-[11px] font-bold uppercase tracking-[0.1em] text-foreground">
          {children}
        </span>
      </div>
      {right}
    </div>
  )
}

function DashBtn({
  children,
  primary,
  onClick,
  className,
}: {
  children: ReactNode
  primary?: boolean
  onClick?: () => void
  className?: string
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={cn(
        "cursor-pointer rounded-[2px] border px-2.5 py-1.5 font-mono text-[10.5px] font-bold uppercase tracking-[0.08em] transition-colors",
        primary
          ? "border-primary bg-primary text-primary-foreground hover:bg-primary/90"
          : "border-border bg-card text-foreground hover:bg-muted/50",
        className,
      )}
    >
      {children}
    </button>
  )
}

function KV({ k, v }: { k: string; v: ReactNode }) {
  return (
    <div className="grid grid-cols-[120px_1fr] items-baseline gap-3.5 border-b border-dashed border-border/50 py-2 text-[12.5px]">
      <span className="font-mono text-[10px] font-semibold uppercase tracking-[0.07em] text-muted-foreground">
        {k}
      </span>
      <span className="leading-[1.45] text-foreground">{v}</span>
    </div>
  )
}

function Chip({
  children,
  tone = "default",
}: {
  children: ReactNode
  tone?: "default" | "primary" | "positive"
}) {
  const toneClass =
    tone === "primary"
      ? "border-primary bg-primary text-primary-foreground"
      : tone === "positive"
        ? "border-positive/30 bg-positive/10 text-positive"
        : "border-border bg-muted/40 text-foreground"
  return (
    <span
      className={cn(
        "inline-flex items-center rounded-[2px] border px-1.5 py-0.5 font-mono text-[10px] font-semibold uppercase leading-[1.4] tracking-[0.07em]",
        toneClass,
      )}
    >
      {children}
    </span>
  )
}

export const UserCard = ({
  user,
  onEdit,
  onLogout,
  onChangePassword,
  onNotifications,
  onIntegrations,
}: UserCardProps) => {
  return (
    <div className="flex flex-col rounded-[2px] border border-border bg-card">
      <SectionLabel
        right={
          <div className="flex gap-1.5">
            <DashBtn onClick={onEdit}>Редактировать</DashBtn>
            <DashBtn onClick={onLogout}>Выйти</DashBtn>
          </div>
        }
      >
        Профиль · {user.userId}
      </SectionLabel>

      {/* Identity strip */}
      <div className="grid grid-cols-[auto_1fr_auto] items-center gap-5 px-[22px] pt-[22px] pb-[18px]">
        <div className="relative flex h-[84px] w-[84px] items-center justify-center border border-primary/30 bg-primary/10">
          <span className="font-sans text-[30px] font-semibold tracking-[-0.02em] text-primary">
            {user.initials}
          </span>
          {user.verified && (
            <span className="absolute -right-px -bottom-px bg-positive px-[5px] py-[2px] font-mono text-[8.5px] font-bold uppercase tracking-[0.07em] text-white">
              ✓ Verif
            </span>
          )}
        </div>

        <div>
          <div className="font-sans text-[22px] font-semibold leading-[1.15] tracking-[-0.014em] text-foreground">
            {user.name}
          </div>
          <div className="mt-1 font-mono text-[11.5px] tracking-[0.03em] text-muted-foreground">
            {user.handle} · {user.email}
          </div>
          <div className="mt-2.5 flex flex-wrap gap-1.5">
            <Chip tone="primary">{user.plan} · Активен</Chip>
            <Chip>С {user.memberSince}</Chip>
            {user.twoFA && <Chip tone="positive">2FA вкл</Chip>}
          </div>
        </div>

        <div className="text-right">
          <div className="font-mono text-[9.5px] font-bold uppercase tracking-[0.1em] text-muted-foreground">
            Портфель (связан)
          </div>
          <div className="mt-1.5 font-mono text-[20px] font-bold tracking-[-0.015em] text-foreground">
            {user.portfolioSize}
          </div>
          <div className="mt-1 font-mono text-[10px] tracking-[0.03em] text-muted-foreground">
            {user.broker} · с {user.brokerLinked}
          </div>
        </div>
      </div>

      {/* KV details */}
      <div className="grid grid-cols-2 gap-x-7 px-[22px] pt-0.5 pb-2.5">
        <div>
          <KV k="Телефон" v={user.phone} />
          <KV k="Часовой пояс" v={user.tz} />
          <KV k="Язык" v={user.lang} />
        </div>
        <div>
          <KV k="Риск-профиль" v={user.riskProfile} />
          <KV k="Брокер" v={user.broker} />
          <KV k="Тариф с" v={user.planSince} />
        </div>
      </div>

      {/* Footer */}
      <div className="flex items-center justify-between border-t border-dashed border-border bg-muted/40 px-3.5 py-2.5">
        <span className="font-mono text-[10px] tracking-[0.05em] text-muted-foreground/80">
          {user.userId} · Last login {user.lastLogin}
        </span>
        <div className="flex gap-1.5">
          <DashBtn onClick={onChangePassword}>Сменить пароль</DashBtn>
          <DashBtn onClick={onNotifications}>Уведомления</DashBtn>
          <DashBtn onClick={onIntegrations}>Интеграции</DashBtn>
        </div>
      </div>
    </div>
  )
}
