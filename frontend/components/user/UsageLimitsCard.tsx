"use client"

import { UsageLimits } from "@/types"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Progress } from "@/components/ui/progress"
import { Badge } from "@/components/ui/badge"
import { MessageSquare, Building2, Briefcase, Bell } from "lucide-react"

interface UsageLimitsCardProps {
  limits: UsageLimits
}

export const UsageLimitsCard = ({ limits }: UsageLimitsCardProps) => {
  const formatResetDate = (date: Date) => {
    return new Intl.DateTimeFormat("ru-RU", {
      day: "numeric",
      month: "long",
    }).format(date)
  }

  const getProgressColor = (percentage: number) => {
    if (percentage >= 90) return "bg-red-500"
    if (percentage >= 70) return "bg-amber-500"
    return "bg-green-500"
  }

  const calculatePercentage = (used: number, limit: number) => {
    if (limit === -1) return 0 // Безлимит
    return Math.min((used / limit) * 100, 100)
  }

  const formatLimit = (limit: number) => {
    return limit === -1 ? "∞" : limit.toString()
  }

  const limits_data = [
    {
      icon: MessageSquare,
      label: "AI запросы",
      used: limits.aiQueries.used,
      limit: limits.aiQueries.limit,
      resetsAt: limits.aiQueries.resetsAt,
      color: "text-blue-600",
    },
    {
      icon: Building2,
      label: "Анализы компаний",
      used: limits.companyAnalyses.used,
      limit: limits.companyAnalyses.limit,
      resetsAt: limits.companyAnalyses.resetsAt,
      color: "text-purple-600",
    },
    {
      icon: Briefcase,
      label: "Портфели",
      used: limits.portfolios.used,
      limit: limits.portfolios.limit,
      color: "text-green-600",
    },
    {
      icon: Bell,
      label: "Алерты",
      used: limits.alerts.used,
      limit: limits.alerts.limit,
      color: "text-amber-600",
    },
  ]

  return (
    <Card>
      <CardHeader className="pb-3">
        <CardTitle className="text-lg">Лимиты</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {limits_data.map((item, idx) => {
          const Icon = item.icon
          const percentage = calculatePercentage(item.used, item.limit)
          const isUnlimited = item.limit === -1

          return (
            <div key={idx} className="space-y-1.5">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-1.5">
                  <Icon className={`h-4 w-4 ${item.color}`} />
                  <span className="font-medium text-xs">{item.label}</span>
                </div>
                <div className="flex items-center gap-1.5">
                  <span className="text-xs font-semibold">
                    {item.used} / {formatLimit(item.limit)}
                  </span>
                  {isUnlimited && (
                    <Badge variant="secondary" className="text-[10px] h-4 px-1.5">
                      ∞
                    </Badge>
                  )}
                </div>
              </div>

              {!isUnlimited && (
                <>
                  <Progress value={percentage} className="h-1.5" />
                  {item.resetsAt && (
                    <p className="text-[10px] text-muted-foreground">
                      Обновится {formatResetDate(item.resetsAt)}
                    </p>
                  )}
                </>
              )}

              {!isUnlimited && percentage >= 90 && (
                <p className="text-[10px] text-red-600">
                  Почти исчерпан
                </p>
              )}
            </div>
          )
        })}

        <div className="pt-3 border-t">
          <p className="text-[10px] text-muted-foreground">
            Лимиты обновляются ежемесячно
          </p>
        </div>
      </CardContent>
    </Card>
  )
}

