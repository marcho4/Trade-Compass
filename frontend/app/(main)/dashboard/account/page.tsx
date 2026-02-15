"use client"

// import { useState } from "react"
// import {
//   ProfileCard,
//   SubscriptionCard,
//   UsageLimitsCard,
//   AccountSettings,
// } from "@/components/user"
// import { mockUser, mockSubscription, mockUsageLimits } from "@/lib/mock-data"

export default function AccountPage() {
  // TODO: Загружать данные пользователя, подписки и лимитов из API
  // const [user] = useState(mockUser)
  // const [subscription] = useState(mockSubscription)
  // const [limits] = useState(mockUsageLimits)

  return (
    <div className="space-y-6">
      {/* Заголовок */}
      <div>
        <h1 className="text-3xl font-bold">Мой аккаунт</h1>
        <p className="text-muted-foreground mt-1 text-sm">
          Управляйте профилем, подпиской и настройками
        </p>
      </div>

      {/* TODO: Вернуть компоненты когда будет API */}
      <div className="rounded-lg border bg-card p-12 text-center">
        <p className="text-muted-foreground">
          Раздел аккаунта будет доступен после подключения API
        </p>
      </div>
    </div>
  )
}

