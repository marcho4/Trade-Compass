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
      <CardHeader>
        <CardTitle>Настройки аккаунта</CardTitle>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* Безопасность */}
        <div className="space-y-4">
          <div className="flex items-center gap-2">
            <Shield className="h-5 w-5 text-primary" />
            <h3 className="font-semibold">Безопасность</h3>
          </div>
          
          <div className="space-y-3 pl-7">
            <div>
              <Label htmlFor="current-password">Текущий пароль</Label>
              <Input 
                id="current-password" 
                type="password" 
                placeholder="Введите текущий пароль"
                className="mt-1.5"
              />
            </div>
            <div>
              <Label htmlFor="new-password">Новый пароль</Label>
              <Input 
                id="new-password" 
                type="password" 
                placeholder="Введите новый пароль"
                className="mt-1.5"
              />
            </div>
            <div>
              <Label htmlFor="confirm-password">Подтвердите пароль</Label>
              <Input 
                id="confirm-password" 
                type="password" 
                placeholder="Повторите новый пароль"
                className="mt-1.5"
              />
            </div>
            <Button variant="outline" size="sm" className="mt-2">
              <Lock className="h-4 w-4 mr-2" />
              Изменить пароль
            </Button>
          </div>
        </div>

        <Separator />

        {/* Уведомления */}
        <div className="space-y-4">
          <div className="flex items-center gap-2">
            <Bell className="h-5 w-5 text-primary" />
            <h3 className="font-semibold">Уведомления</h3>
          </div>

          <div className="space-y-4 pl-7">
            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label htmlFor="email-notif">Email уведомления</Label>
                <p className="text-xs text-muted-foreground">
                  Получать важные обновления на email
                </p>
              </div>
              <Switch 
                id="email-notif"
                checked={emailNotifications}
                onCheckedChange={setEmailNotifications}
              />
            </div>

            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label htmlFor="portfolio-alerts">Алерты по портфелю</Label>
                <p className="text-xs text-muted-foreground">
                  Уведомления о важных изменениях в портфеле
                </p>
              </div>
              <Switch 
                id="portfolio-alerts"
                checked={portfolioAlerts}
                onCheckedChange={setPortfolioAlerts}
              />
            </div>

            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label htmlFor="market-news">Новости рынка</Label>
                <p className="text-xs text-muted-foreground">
                  Получать дайджест новостей рынка
                </p>
              </div>
              <Switch 
                id="market-news"
                checked={marketNews}
                onCheckedChange={setMarketNews}
              />
            </div>

            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label htmlFor="weekly-report">Еженедельный отчет</Label>
                <p className="text-xs text-muted-foreground">
                  Получать сводку по портфелю раз в неделю
                </p>
              </div>
              <Switch 
                id="weekly-report"
                checked={weeklyReport}
                onCheckedChange={setWeeklyReport}
              />
            </div>
          </div>
        </div>

        <Separator />

        {/* Email */}
        <div className="space-y-4">
          <div className="flex items-center gap-2">
            <Mail className="h-5 w-5 text-primary" />
            <h3 className="font-semibold">Email адрес</h3>
          </div>
          
          <div className="space-y-3 pl-7">
            <div>
              <Label htmlFor="email">Текущий email</Label>
              <Input 
                id="email" 
                type="email" 
                defaultValue="user@example.com"
                className="mt-1.5"
              />
            </div>
            <Button variant="outline" size="sm">
              Изменить email
            </Button>
          </div>
        </div>

        <Separator />

        {/* Опасная зона */}
        <div className="space-y-4">
          <div>
            <h3 className="font-semibold text-red-600 flex items-center gap-2">
              <Trash2 className="h-5 w-5" />
              Опасная зона
            </h3>
          </div>
          
          <div className="pl-7 space-y-3">
            <p className="text-sm text-muted-foreground">
              Удаление аккаунта приведет к безвозвратной потере всех данных, 
              включая портфели, настройки и историю.
            </p>
            <Button 
              variant="destructive" 
              size="sm"
              onClick={onDeleteAccount}
            >
              <Trash2 className="h-4 w-4 mr-2" />
              Удалить аккаунт
            </Button>
          </div>
        </div>

        <Separator />

        {/* Сохранить изменения */}
        <div className="flex justify-end">
          <Button onClick={handleSave}>
            <Save className="h-4 w-4 mr-2" />
            Сохранить изменения
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}

