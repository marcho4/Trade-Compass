"use client"

import { Subscription, SubscriptionLevel, SUBSCRIPTION_PLANS } from "@/types"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Progress } from "@/components/ui/progress"
import { Crown, Calendar, TrendingUp, AlertCircle } from "lucide-react"

interface SubscriptionCardProps {
  subscription: Subscription
  onUpgrade?: () => void
  onCancel?: () => void
}

export const SubscriptionCard = ({ subscription, onUpgrade, onCancel }: SubscriptionCardProps) => {
  const plan = SUBSCRIPTION_PLANS.find((p) => p.level === subscription.level)

  const formatDate = (date: Date) => {
    return new Intl.DateTimeFormat("ru-RU", {
      day: "numeric",
      month: "long",
      year: "numeric",
    }).format(date)
  }

  const getLevelColor = (level: SubscriptionLevel) => {
    switch (level) {
      case "pro":
        return "bg-gradient-to-br from-[hsl(var(--chart-3))] via-[hsl(var(--chart-4))] to-[hsl(var(--chart-5))]"
      case "premium":
        return "bg-gradient-to-br from-primary via-[hsl(var(--chart-1))] to-[hsl(var(--chart-2))]"
      default:
        return "bg-gradient-to-br from-muted to-muted-foreground/20"
    }
  }

  const getLevelBadgeVariant = (level: SubscriptionLevel) => {
    switch (level) {
      case "pro":
        return "default"
      case "premium":
        return "secondary"
      default:
        return "outline"
    }
  }

  const daysUntilRenewal = subscription.renewsAt
    ? Math.ceil((subscription.renewsAt.getTime() - new Date().getTime()) / (1000 * 60 * 60 * 24))
    : null

  const isExpiringSoon = daysUntilRenewal !== null && daysUntilRenewal <= 7

  return (
    <Card>
      <CardHeader className="pb-3">
        <CardTitle className="text-lg">Подписка</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Текущий план */}
        <div className={`rounded-lg p-4 ${getLevelColor(subscription.level)}`}>
          <div className="flex items-center justify-between mb-3">
            <div className="flex items-center gap-2">
              <Crown className="h-5 w-5 text-white" />
              <div>
                <h3 className="text-lg font-bold text-white">{plan?.name || "Бесплатный"}</h3>
                <p className="text-white/80 text-xs">
                  {subscription.isActive ? "Активна" : "Неактивна"}
                </p>
              </div>
            </div>
            <Badge variant={getLevelBadgeVariant(subscription.level)} className="bg-white/20 text-white border-white/30 text-xs">
              {subscription.level.toUpperCase()}
            </Badge>
          </div>

          {subscription.price && (
            <div className="text-white">
              <span className="text-2xl font-bold">{subscription.price} ₽</span>
              <span className="text-white/80 text-sm"> / месяц</span>
            </div>
          )}
        </div>

        {/* Информация о подписке */}
        <div className="space-y-2">
          <div className="flex items-center justify-between text-xs">
            <span className="text-muted-foreground flex items-center gap-1.5">
              <Calendar className="h-3.5 w-3.5" />
              Начало:
            </span>
            <span className="font-medium">{formatDate(subscription.startDate)}</span>
          </div>

          {subscription.endDate && (
            <div className="flex items-center justify-between text-xs">
              <span className="text-muted-foreground flex items-center gap-1.5">
                <Calendar className="h-3.5 w-3.5" />
                Окончание:
              </span>
              <span className="font-medium">{formatDate(subscription.endDate)}</span>
            </div>
          )}

          {subscription.renewsAt && (
            <div className="flex items-center justify-between text-xs">
              <span className="text-muted-foreground flex items-center gap-1.5">
                <TrendingUp className="h-3.5 w-3.5" />
                Продление:
              </span>
              <span className="font-medium">{formatDate(subscription.renewsAt)}</span>
            </div>
          )}
        </div>

        {/* Предупреждение о скором окончании */}
        {isExpiringSoon && daysUntilRenewal !== null && (
          <div className="flex items-start gap-2 p-3 rounded-lg bg-amber-500/10 border border-amber-500/20">
            <AlertCircle className="h-4 w-4 text-amber-600 shrink-0 mt-0.5" />
            <div>
              <p className="font-medium text-xs text-amber-900 dark:text-amber-100">
                Истекает через {daysUntilRenewal} {daysUntilRenewal === 1 ? "день" : "дней"}
              </p>
            </div>
          </div>
        )}

        {/* Действия */}
        <div className="flex gap-2">
          {subscription.level !== "pro" && onUpgrade && (
            <Button onClick={onUpgrade} size="sm" className="flex-1">
              <Crown className="h-3.5 w-3.5 mr-1.5" />
              Улучшить
            </Button>
          )}
          {subscription.level !== "free" && onCancel && (
            <Button variant="outline" size="sm" onClick={onCancel} className="flex-1">
              Отменить
            </Button>
          )}
        </div>

        {/* Функции плана */}
        {plan && (
          <div className="pt-3 border-t">
            <h4 className="font-semibold mb-2 text-xs">Возможности:</h4>
            <div className="space-y-1.5">
              {plan.features.filter(f => f.included).slice(0, 4).map((feature, idx) => (
                <div key={idx} className="flex items-center gap-1.5 text-xs">
                  <div className="h-1 w-1 rounded-full bg-green-500" />
                  <span>{feature.name}</span>
                </div>
              ))}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  )
}

