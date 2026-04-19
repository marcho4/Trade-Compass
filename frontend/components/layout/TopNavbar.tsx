"use client"

import { Fragment, useMemo } from "react"
import Link from "next/link"
import { usePathname } from "next/navigation"
import { cn } from "@/lib/utils"

function isMoexOpen(): boolean {
  const now = new Date()
  const h = Number(
    new Intl.DateTimeFormat("en-GB", {
      timeZone: "Europe/Moscow",
      hour: "2-digit",
      hour12: false,
    }).format(now),
  )
  const day = now.getUTCDay()
  return day !== 0 && day !== 6 && h >= 10 && h < 24
}

interface Crumb {
  label: string
  href?: string
}

function parseCrumbs(pathname: string): Crumb[] {
  const parts = pathname.replace(/^\/dashboard\/?/, "").split("/").filter(Boolean)

  if (!parts.length || parts[0] === "screener") {
    return [{ label: "СКРИНЕР" }]
  }
  if (parts[0] === "portfolio") {
    if (parts[1]) {
      return [
        { label: "ПОРТФЕЛЬ", href: "/dashboard/portfolio" },
        { label: parts[1].toUpperCase() },
      ]
    }
    return [{ label: "ПОРТФЕЛЬ" }]
  }
  if (parts[0] === "account") {
    return [{ label: "АККАУНТ" }]
  }

  return [
    { label: "СКРИНЕР", href: "/dashboard/screener" },
    { label: parts[0].toUpperCase() },
  ]
}

function NavLink({ href, label }: { href: string; label: string }) {
  return (
    <Link
      href={href}
      className="font-mono text-[11px] tracking-[0.5px] text-muted-foreground hover:text-primary transition-colors duration-[120ms] no-underline"
    >
      {label}
    </Link>
  )
}

function AccountButton() {
  return (
    <Link
      href="/dashboard/account"
      className="inline-flex items-center gap-[7px] font-mono text-[11px] font-semibold tracking-[0.8px] text-muted-foreground hover:text-primary no-underline px-[11px] py-[5px] border border-border hover:border-primary rounded-[2px] bg-transparent hover:bg-primary/10 transition-all duration-[120ms]"
    >
      <svg
        width="12"
        height="12"
        viewBox="0 0 12 12"
        fill="none"
        className="shrink-0"
      >
        <circle cx="6" cy="4" r="2.5" stroke="currentColor" strokeWidth="1.2" />
        <path
          d="M1.5 10.5C1.5 8.567 3.567 7 6 7s4.5 1.567 4.5 3.5"
          stroke="currentColor"
          strokeWidth="1.2"
          strokeLinecap="round"
        />
      </svg>
      АККАУНТ
    </Link>
  )
}

export const TopNavbar = () => {
  const pathname = usePathname()
  const crumbs = useMemo(() => parseCrumbs(pathname), [pathname])
  const marketOpen = useMemo(isMoexOpen, [])

  return (
    <header className="fixed top-0 left-0 right-0 z-50 flex items-center justify-between px-7 h-[52px] border-b border-border bg-background">
      <div className="flex items-center gap-4">
        <Link
          href="/dashboard/screener"
          className="font-mono text-[13px] font-bold text-primary tracking-[2px] no-underline"
        >
          Trade Compass
        </Link>

        <div className="flex items-center gap-2.5 font-mono text-[11px] text-muted-foreground tracking-[0.5px]">
          <span className="text-muted-foreground/50">&gt;</span>
          {crumbs.map((crumb, i) => (
            <Fragment key={i}>
              {i > 0 && (
                <span className="text-muted-foreground/50 mx-px">/</span>
              )}
              {crumb.href ? (
                <NavLink href={crumb.href} label={crumb.label} />
              ) : (
                <span className="text-foreground">{crumb.label}</span>
              )}
            </Fragment>
          ))}
        </div>
      </div>

      <div className="flex items-center gap-[18px]">
        <span className="hidden md:inline-flex items-center gap-[6px] font-mono text-[11px] text-muted-foreground tracking-[0.5px]">
          <span
            className={cn(
              "w-1.5 h-1.5 rounded-full shrink-0",
              marketOpen ? "bg-positive" : "bg-muted-foreground/50",
            )}
          />
          MOEX {marketOpen ? "OPEN" : "CLOSED"}
        </span>

        <AccountButton />
      </div>
    </header>
  )
}
