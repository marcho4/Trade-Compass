"use client"

import { useState } from "react"
import {
  ProfileCard,
  SubscriptionCard,
  UsageLimitsCard,
  AccountSettings,
} from "@/components/user"
import { mockUser, mockSubscription, mockUsageLimits } from "@/lib/mock-data"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Button } from "@/components/ui/button"
import { User, CreditCard, BarChart3, Settings } from "lucide-react"

export default function AccountPage() {
  const [user] = useState(mockUser)
  const [subscription] = useState(mockSubscription)
  const [limits] = useState(mockUsageLimits)

  const handleEditProfile = () => {
    console.log("Edit profile")
    // TODO: Открыть модальное окно редактирования профиля
  }

  const handleUpgradeSubscription = () => {
    console.log("Upgrade subscription")
    // TODO: Перенаправить на страницу выбора плана
  }

  const handleCancelSubscription = () => {
    console.log("Cancel subscription")
    // TODO: Показать подтверждение отмены подписки
  }

  const handleSaveSettings = (settings: {
    emailNotifications: boolean
    portfolioAlerts: boolean
    marketNews: boolean
    weeklyReport: boolean
  }) => {
    console.log("Save settings:", settings)
    // TODO: Сохранить настройки через API
  }

  const handleDeleteAccount = () => {
    console.log("Delete account")
    // TODO: Показать подтверждение удаления аккаунта
  }

  return (
    <div className="space-y-8">
      {/* Заголовок */}
      <div>
        <h1 className="text-4xl font-bold">Мой аккаунт</h1>
        <p className="text-muted-foreground mt-2">
          Управляйте профилем, подпиской и настройками
        </p>
      </div>

      {/* Tabs */}
      <Tabs defaultValue="profile" className="w-full">
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="profile" className="flex items-center gap-2">
            <User className="h-4 w-4" />
            <span className="hidden sm:inline">Профиль</span>
          </TabsTrigger>
          <TabsTrigger value="subscription" className="flex items-center gap-2">
            <CreditCard className="h-4 w-4" />
            <span className="hidden sm:inline">Подписка</span>
          </TabsTrigger>
          <TabsTrigger value="usage" className="flex items-center gap-2">
            <BarChart3 className="h-4 w-4" />
            <span className="hidden sm:inline">Лимиты</span>
          </TabsTrigger>
          <TabsTrigger value="settings" className="flex items-center gap-2">
            <Settings className="h-4 w-4" />
            <span className="hidden sm:inline">Настройки</span>
          </TabsTrigger>
        </TabsList>

        {/* Профиль */}
        <TabsContent value="profile" className="space-y-6 mt-6">
          <ProfileCard user={user} onEdit={handleEditProfile} />
        </TabsContent>

        {/* Подписка */}
        <TabsContent value="subscription" className="space-y-6 mt-6">
          <SubscriptionCard
            subscription={subscription}
            onUpgrade={handleUpgradeSubscription}
            onCancel={handleCancelSubscription}
          />

          {/* Дополнительная информация о планах */}
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {/* Можно добавить карточки с планами для сравнения */}
          </div>
        </TabsContent>

        {/* Лимиты использования */}
        <TabsContent value="usage" className="space-y-6 mt-6">
          <UsageLimitsCard limits={limits} />

          {/* Дополнительная информация */}
          <div className="rounded-lg border bg-card p-6">
            <h3 className="font-semibold mb-3">О лимитах</h3>
            <div className="space-y-2 text-sm text-muted-foreground">
              <p>
                <strong className="text-foreground">AI запросы</strong> — количество вопросов,
                которые вы можете задать AI-ассистенту в месяц.
              </p>
              <p>
                <strong className="text-foreground">Анализы компаний</strong> — количество
                подробных анализов компаний, доступных для просмотра.
              </p>
              <p>
                <strong className="text-foreground">Портфели</strong> — максимальное количество
                портфелей, которые вы можете создать.
              </p>
              <p>
                <strong className="text-foreground">Алерты</strong> — количество уведомлений
                о важных событиях в ваших портфелях.
              </p>
              <p className="pt-2 border-t mt-4">
                Нужно больше? 
                <Button
                  variant="link"
                  className="px-2"
                  onClick={handleUpgradeSubscription}
                >
                  Улучшите план подписки
                </Button>
              </p>
            </div>
          </div>
        </TabsContent>

        {/* Настройки */}
        <TabsContent value="settings" className="space-y-6 mt-6">
          <AccountSettings
            onSave={handleSaveSettings}
            onDeleteAccount={handleDeleteAccount}
          />
        </TabsContent>
      </Tabs>
    </div>
  )
}

