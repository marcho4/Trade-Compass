"use client"

import * as React from "react"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { LayoutDashboard, Briefcase, User } from "lucide-react"
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip"
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

export const FloatingSidebar = () => {
  const pathname = usePathname()

  const isActive = (href: string) => {
    if (href === "/dashboard/screener") {
      return pathname === href
    }
    return pathname.startsWith(href)
  }

  return (
    <TooltipProvider delayDuration={0}>
      <aside
        className={cn(
          "fixed left-6 top-1/2 -translate-y-1/2 z-50",
          "bg-background/95 backdrop-blur-xl border border-border/60",
          "rounded-3xl shadow-2xl w-18"
        )}
      >
        <div className="flex flex-col items-center justify-center p-4 gap-4 py-6">
          {/* Navigation Items */}
          {navItems.map((item) => {
            const Icon = item.icon
            const active = isActive(item.href)

            return (
              <Tooltip key={item.key}>
                <TooltipTrigger asChild>
                  <Link href={item.href}>
                    <Button
                      variant={active ? "default" : "ghost"}
                      size="icon"
                      className={cn(
                        "w-11 h-11 transition-all duration-300",
                        active && "shadow-lg scale-105",
                        !active && "hover:scale-105"
                      )}
                    >
                      <Icon className="h-5 w-5" />
                    </Button>
                  </Link>
                </TooltipTrigger>
                <TooltipContent side="right" className="ml-2">
                  {item.label}
                </TooltipContent>
              </Tooltip>
            )
          })}

          {/* Separator */}
          <div className="w-8 h-px bg-border/60 my-2" />

          {/* User Account */}
          <Tooltip>
            <TooltipTrigger asChild>
              <Link href="/dashboard/account">
                <Button
                  variant="ghost"
                  size="icon"
                  className="w-11 h-11 hover:scale-110 transition-all duration-300"
                >
                  <div className="h-9 w-9 rounded-full bg-gradient-to-br from-primary via-[hsl(var(--chart-1))] to-[hsl(var(--chart-3))] flex items-center justify-center shadow-md">
                    <User className="h-4 w-4 text-primary-foreground" />
                  </div>
                </Button>
              </Link>
                
            </TooltipTrigger>
            <TooltipContent side="right" className="ml-2">
              <div className="text-sm">
                <div className="font-medium">Пользователь</div>
                <div className="text-xs text-muted-foreground">
                  user@email.com
                </div>
              </div>
            </TooltipContent>
          </Tooltip>
        </div>
      </aside>
    </TooltipProvider>
  )
}

