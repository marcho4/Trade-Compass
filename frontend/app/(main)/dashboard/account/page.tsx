"use client"

import { useAuth } from "@/contexts/AuthContext"

export default function AccountPage() {
  const { user, isLoading } = useAuth()

  if (isLoading) {
    return (
      <div className="flex items-center justify-center p-12">
        <p className="text-muted-foreground">Загрузка...</p>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">
          Привет, {user?.name || "Пользователь"}!
        </h1>
        <p className="text-muted-foreground mt-1 text-sm">
          Управляйте профилем, подпиской и настройками
        </p>
      </div>
    </div>
  )
}

