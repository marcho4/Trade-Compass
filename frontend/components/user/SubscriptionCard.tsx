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
      <CardHeader>
        <CardTitle>Подписка</CardTitle>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* Текущий план */}
        <div className={`rounded-lg p-6 ${getLevelColor(subscription.level)}`}>
          <div className="flex items-center justify-between mb-4">
            <div className="flex items-center gap-3">
              <Crown className="h-6 w-6 text-white" />
              <div>
                <h3 className="text-2xl font-bold text-white">{plan?.name || "Бесплатный"}</h3>
                <p className="text-white/80 text-sm">
                  {subscription.isActive ? "Активна" : "Неактивна"}
                </p>
              </div>
            </div>
            <Badge variant={getLevelBadgeVariant(subscription.level)} className="bg-white/20 text-white border-white/30">
              {subscription.level.toUpperCase()}
            </Badge>
          </div>

          {subscription.price && (
            <div className="text-white">
              <span className="text-3xl font-bold">{subscription.price} ₽</span>
              <span className="text-white/80"> / месяц</span>
            </div>
          )}
        </div>

        {/* Информация о подписке */}
        <div className="space-y-3">
          <div className="flex items-center justify-between text-sm">
            <span className="text-muted-foreground flex items-center gap-2">
              <Calendar className="h-4 w-4" />
              Дата начала:
            </span>
            <span className="font-medium">{formatDate(subscription.startDate)}</span>
          </div>

          {subscription.endDate && (
            <div className="flex items-center justify-between text-sm">
              <span className="text-muted-foreground flex items-center gap-2">
                <Calendar className="h-4 w-4" />
                Дата окончания:
              </span>
              <span className="font-medium">{formatDate(subscription.endDate)}</span>
            </div>
          )}

          {subscription.renewsAt && (
            <div className="flex items-center justify-between text-sm">
              <span className="text-muted-foreground flex items-center gap-2">
                <TrendingUp className="h-4 w-4" />
                Следующее продление:
              </span>
              <span className="font-medium">{formatDate(subscription.renewsAt)}</span>
            </div>
          )}
        </div>

        {/* Предупреждение о скором окончании */}
        {isExpiringSoon && daysUntilRenewal !== null && (
          <div className="flex items-start gap-3 p-4 rounded-lg bg-amber-500/10 border border-amber-500/20">
            <AlertCircle className="h-5 w-5 text-amber-600 shrink-0 mt-0.5" />
            <div>
              <p className="font-medium text-sm text-amber-900 dark:text-amber-100">
                Подписка истекает через {daysUntilRenewal} {daysUntilRenewal === 1 ? "день" : "дней"}
              </p>
              <p className="text-xs text-amber-700 dark:text-amber-200 mt-1">
                Продлите подписку, чтобы не потерять доступ к функциям
              </p>
            </div>
          </div>
        )}

        {/* Действия */}
        <div className="flex gap-3">
          {subscription.level !== "pro" && onUpgrade && (
            <Button onClick={onUpgrade} className="flex-1">
              <Crown className="h-4 w-4 mr-2" />
              Улучшить план
            </Button>
          )}
          {subscription.level !== "free" && onCancel && (
            <Button variant="outline" onClick={onCancel} className="flex-1">
              Отменить подписку
            </Button>
          )}
        </div>

        {/* Функции плана */}
        {plan && (
          <div className="pt-4 border-t">
            <h4 className="font-semibold mb-3 text-sm">Возможности вашего плана:</h4>
            <div className="space-y-2">
              {plan.features.filter(f => f.included).map((feature, idx) => (
                <div key={idx} className="flex items-center gap-2 text-sm">
                  <div className="h-1.5 w-1.5 rounded-full bg-green-500" />
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

