"use client"

import { useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Separator } from "@/components/ui/separator"
import { Switch } from "@/components/ui/switch"
import { 
  Shield, 
  Bell, 
  Mail, 
  Lock, 
  Trash2,
  Save 
} from "lucide-react"

interface AccountSettingsProps {
  onSave?: (settings: any) => void
  onDeleteAccount?: () => void
}

export const AccountSettings = ({ onSave, onDeleteAccount }: AccountSettingsProps) => {
  const [emailNotifications, setEmailNotifications] = useState(true)
  const [portfolioAlerts, setPortfolioAlerts] = useState(true)
  const [marketNews, setMarketNews] = useState(false)
  const [weeklyReport, setWeeklyReport] = useState(true)

  const handleSave = () => {
    const settings = {
      emailNotifications,
      portfolioAlerts,
      marketNews,
      weeklyReport,
    }
    onSave?.(settings)
  }

  return (
    <Card>
      <CardHeader className="pb-3">
        <CardTitle className="text-lg">Настройки</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Безопасность и Email в двух колонках на больших экранах */}
        <div className="grid gap-4 lg:grid-cols-2">
          {/* Безопасность */}
          <div className="space-y-3">
            <div className="flex items-center gap-2">
              <Shield className="h-4 w-4 text-primary" />
              <h3 className="font-semibold text-sm">Безопасность</h3>
            </div>

            <div className="space-y-2 pl-6">
              <div>
                <Label htmlFor="current-password" className="text-xs">Текущий пароль</Label>
                <Input
                  id="current-password"
                  type="password"
                  placeholder="••••••••"
                  className="mt-1 h-8 text-sm"
                />
              </div>
              <div>
                <Label htmlFor="new-password" className="text-xs">Новый пароль</Label>
                <Input
                  id="new-password"
                  type="password"
                  placeholder="••••••••"
                  className="mt-1 h-8 text-sm"
                />
              </div>
              <Button variant="outline" size="sm" className="mt-2 h-7 text-xs">
                <Lock className="h-3 w-3 mr-1.5" />
                Изменить
              </Button>
            </div>
          </div>

          {/* Email */}
          <div className="space-y-3">
            <div className="flex items-center gap-2">
              <Mail className="h-4 w-4 text-primary" />
              <h3 className="font-semibold text-sm">Email</h3>
            </div>

            <div className="space-y-2 pl-6">
              <div>
                <Label htmlFor="email" className="text-xs">Текущий email</Label>
                <Input
                  id="email"
                  type="email"
                  defaultValue="user@example.com"
                  className="mt-1 h-8 text-sm"
                />
              </div>
              <Button variant="outline" size="sm" className="h-7 text-xs">
                Изменить email
              </Button>
            </div>
          </div>
        </div>

        <Separator />

        {/* Уведомления */}
        <div className="space-y-3">
          <div className="flex items-center gap-2">
            <Bell className="h-4 w-4 text-primary" />
            <h3 className="font-semibold text-sm">Уведомления</h3>
          </div>

          <div className="space-y-2 pl-6">
            <div className="flex items-center justify-between">
              <Label htmlFor="email-notif" className="text-xs">Email уведомления</Label>
              <Switch
                id="email-notif"
                checked={emailNotifications}
                onCheckedChange={setEmailNotifications}
              />
            </div>

            <div className="flex items-center justify-between">
              <Label htmlFor="portfolio-alerts" className="text-xs">Алерты по портфелю</Label>
              <Switch
                id="portfolio-alerts"
                checked={portfolioAlerts}
                onCheckedChange={setPortfolioAlerts}
              />
            </div>

            <div className="flex items-center justify-between">
              <Label htmlFor="market-news" className="text-xs">Новости рынка</Label>
              <Switch
                id="market-news"
                checked={marketNews}
                onCheckedChange={setMarketNews}
              />
            </div>

            <div className="flex items-center justify-between">
              <Label htmlFor="weekly-report" className="text-xs">Еженедельный отчет</Label>
              <Switch
                id="weekly-report"
                checked={weeklyReport}
                onCheckedChange={setWeeklyReport}
              />
            </div>
          </div>
        </div>

        <Separator />

        {/* Опасная зона */}
        <div className="space-y-3">
          <h3 className="font-semibold text-red-600 flex items-center gap-2 text-sm">
            <Trash2 className="h-4 w-4" />
            Опасная зона
          </h3>

          <div className="pl-6 space-y-2">
            <p className="text-xs text-muted-foreground">
              Удаление аккаунта приведет к безвозвратной потере всех данных.
            </p>
            <Button
              variant="destructive"
              size="sm"
              onClick={onDeleteAccount}
              className="h-7 text-xs"
            >
              <Trash2 className="h-3 w-3 mr-1.5" />
              Удалить аккаунт
            </Button>
          </div>
        </div>

        <Separator />

        {/* Сохранить изменения */}
        <div className="flex justify-end">
          <Button onClick={handleSave} size="sm" className="h-8 text-xs">
            <Save className="h-3.5 w-3.5 mr-1.5" />
            Сохранить
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}

