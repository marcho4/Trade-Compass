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
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle>Профиль</CardTitle>
          {onEdit && (
            <Button variant="outline" size="sm" onClick={onEdit}>
              Редактировать
            </Button>
          )}
        </div>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* Аватар и имя */}
        <div className="flex items-center gap-4">
          <div className="h-20 w-20 rounded-full bg-gradient-to-br from-primary via-[hsl(var(--chart-1))] to-[hsl(var(--chart-3))] flex items-center justify-center shadow-lg">
            <UserCircle2 className="h-12 w-12 text-primary-foreground" />
          </div>
          <div>
            <h3 className="text-2xl font-bold">{user.name}</h3>
            <p className="text-sm text-muted-foreground">ID: {user.id}</p>
          </div>
        </div>

        {/* Информация */}
        <div className="space-y-4">
          <div className="flex items-center gap-3 text-sm">
            <Mail className="h-4 w-4 text-muted-foreground" />
            <span className="text-muted-foreground">Email:</span>
            <span className="font-medium">{user.email}</span>
          </div>

          <div className="flex items-center gap-3 text-sm">
            <Calendar className="h-4 w-4 text-muted-foreground" />
            <span className="text-muted-foreground">Дата регистрации:</span>
            <span className="font-medium">{formatDate(user.createdAt)}</span>
          </div>

          <div className="flex items-center gap-3 text-sm">
            <Clock className="h-4 w-4 text-muted-foreground" />
            <span className="text-muted-foreground">Последний вход:</span>
            <span className="font-medium">{formatDateTime(user.lastLoginAt)}</span>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

