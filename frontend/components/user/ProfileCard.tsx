"use client"

import { User } from "@/types"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { UserCircle2, Mail, Calendar, Clock } from "lucide-react"

interface ProfileCardProps {
  user: User
  onEdit?: () => void
}

export const ProfileCard = ({ user, onEdit }: ProfileCardProps) => {
  const formatDate = (date: Date) => {
    return new Intl.DateTimeFormat("ru-RU", {
      day: "numeric",
      month: "long",
      year: "numeric",
    }).format(date)
  }

  const formatDateTime = (date: Date) => {
    return new Intl.DateTimeFormat("ru-RU", {
      day: "numeric",
      month: "long",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    }).format(date)
  }

  return (
    <Card>
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <CardTitle className="text-lg">Профиль</CardTitle>
          {onEdit && (
            <Button variant="outline" size="sm" onClick={onEdit}>
              Редактировать
            </Button>
          )}
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Аватар и имя */}
        <div className="flex items-center gap-3">
          <div className="h-14 w-14 rounded-full bg-gradient-to-br from-primary via-[hsl(var(--chart-1))] to-[hsl(var(--chart-3))] flex items-center justify-center shadow-lg">
            <UserCircle2 className="h-8 w-8 text-primary-foreground" />
          </div>
          <div>
            <h3 className="text-xl font-bold">{user.name}</h3>
            <p className="text-xs text-muted-foreground">ID: {user.id}</p>
          </div>
        </div>

        {/* Информация */}
        <div className="space-y-2">
          <div className="flex items-center gap-2 text-sm">
            <Mail className="h-3.5 w-3.5 text-muted-foreground" />
            <span className="text-muted-foreground">Email:</span>
            <span className="font-medium">{user.email}</span>
          </div>

          <div className="flex items-center gap-2 text-sm">
            <Calendar className="h-3.5 w-3.5 text-muted-foreground" />
            <span className="text-muted-foreground">Регистрация:</span>
            <span className="font-medium">{formatDate(user.createdAt)}</span>
          </div>

          <div className="flex items-center gap-2 text-sm">
            <Clock className="h-3.5 w-3.5 text-muted-foreground" />
            <span className="text-muted-foreground">Последний вход:</span>
            <span className="font-medium">{formatDateTime(user.lastLoginAt)}</span>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

