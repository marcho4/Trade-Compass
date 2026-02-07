"use client"

import * as React from "react"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { LayoutDashboard, Briefcase, User } from "lucide-react"
import Link from "next/link"
import { usePathname } from "next/navigation"

type NavItem = {
  icon: React.ElementType
  label: string
  href: string
  key: string
}

const navItems: NavItem[] = [
  {
    icon: LayoutDashboard,
    label: "Скринер",
    href: "/dashboard/screener",
    key: "screener",
  },
  {
    icon: Briefcase,
    label: "Портфель",
    href: "/dashboard/portfolio",
    key: "portfolio",
  },
]

export const TopNavbar = () => {
  const pathname = usePathname()

  const isActive = (href: string) => {
    if (href === "/dashboard/screener") {
      return pathname === href
    }
    return pathname.startsWith(href)
  }

  return (
    <header
      className={cn(
        "fixed top-0 left-0 right-0 z-50 h-16",
        "bg-background/90 backdrop-blur-xl border-b border-border/50"
      )}
    >
      <nav className="h-full max-w-screen-2xl mx-auto flex items-center justify-between px-6">
        <Link
          href="/dashboard/screener"
          className="text-lg font-bold tracking-tight text-foreground transition-colors hover:text-primary"
        >
          TradeCompass
        </Link>

        <div className="flex items-center gap-1">
          {navItems.map((item) => {
            const Icon = item.icon
            const active = isActive(item.href)

            return (
              <Link key={item.key} href={item.href}>
                <Button
                  variant="ghost"
                  className={cn(
                    "gap-2 px-4 h-9 transition-all duration-200",
                    active &&
                      "bg-primary text-primary-foreground shadow-sm hover:bg-primary/90 hover:text-primary-foreground",
                    !active && "text-muted-foreground hover:text-foreground"
                  )}
                >
                  <Icon className="h-4 w-4" />
                  <span className="text-sm font-medium">{item.label}</span>
                </Button>
              </Link>
            )
          })}
        </div>

        <Link href="/dashboard/account">
          <Button
            variant="ghost"
            size="icon"
            className="h-9 w-9 rounded-full hover:scale-105 transition-all duration-200"
          >
            <div className="h-8 w-8 rounded-full bg-linear-to-br from-primary via-[hsl(var(--chart-1))] to-[hsl(var(--chart-3))] flex items-center justify-center shadow-sm">
              <User className="h-3.5 w-3.5 text-primary-foreground" />
            </div>
          </Button>
        </Link>
      </nav>
    </header>
  )
}
