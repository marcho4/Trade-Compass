"use client"

import { useState } from "react"
import {
  ProfileCard,
  SubscriptionCard,
  UsageLimitsCard,
  AccountSettings,
} from "@/components/user"
import { mockUser, mockSubscription, mockUsageLimits } from "@/lib/mock-data"

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
    <div className="space-y-6">
      {/* Заголовок */}
      <div>
        <h1 className="text-3xl font-bold">Мой аккаунт</h1>
        <p className="text-muted-foreground mt-1 text-sm">
          Управляйте профилем, подпиской и настройками
        </p>
      </div>

      {/* Профиль - полная ширина */}
      <div>
        <ProfileCard user={user} onEdit={handleEditProfile} />
      </div>

      {/* Подписка и лимиты - две колонки на больших экранах */}
      <div className="grid gap-6 lg:grid-cols-2">
        <SubscriptionCard
          subscription={subscription}
          onUpgrade={handleUpgradeSubscription}
          onCancel={handleCancelSubscription}
        />
        <UsageLimitsCard limits={limits} />
      </div>

      {/* Настройки - полная ширина */}
      <div>
        <AccountSettings
          onSave={handleSaveSettings}
          onDeleteAccount={handleDeleteAccount}
        />
      </div>
    </div>
  )
}

