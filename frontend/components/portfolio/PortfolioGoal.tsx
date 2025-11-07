"use client"

import { useState } from "react"
import { Target, Edit2, Check, X } from "lucide-react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Progress } from "@/components/ui/progress"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"

interface PortfolioGoalProps {
  currentValue: number
  goalValue: number
  goalDescription?: string
  onUpdateGoal?: (goalValue: number, description: string) => void
}

export const PortfolioGoal = ({
  currentValue,
  goalValue,
  goalDescription = "Достичь целевой суммы инвестиций",
  onUpdateGoal,
}: PortfolioGoalProps) => {
  const [isEditing, setIsEditing] = useState(false)
  const [newGoalValue, setNewGoalValue] = useState(goalValue.toString())
  const [newDescription, setNewDescription] = useState(goalDescription)

  const progress = Math.min((currentValue / goalValue) * 100, 100)
  const remaining = Math.max(goalValue - currentValue, 0)

  const handleSave = () => {
    const parsedValue = parseFloat(newGoalValue)
    if (!isNaN(parsedValue) && parsedValue > 0) {
      onUpdateGoal?.(parsedValue, newDescription)
      setIsEditing(false)
    }
  }

  const handleCancel = () => {
    setNewGoalValue(goalValue.toString())
    setNewDescription(goalDescription)
    setIsEditing(false)
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex items-start justify-between">
          <div className="flex items-center gap-2">
            <Target className="h-5 w-5 text-primary" />
            <div>
              <CardTitle>Цель портфеля</CardTitle>
              <CardDescription className="mt-1">
                {isEditing ? "Редактирование цели" : goalDescription}
              </CardDescription>
            </div>
          </div>
          {!isEditing && (
            <Button
              variant="ghost"
              size="icon"
              onClick={() => setIsEditing(true)}
            >
              <Edit2 className="h-4 w-4" />
            </Button>
          )}
        </div>
      </CardHeader>
      <CardContent>
        {isEditing ? (
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="goal-value">Целевая сумма (₽)</Label>
              <Input
                id="goal-value"
                type="number"
                value={newGoalValue}
                onChange={(e) => setNewGoalValue(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="goal-description">Описание цели</Label>
              <Input
                id="goal-description"
                value={newDescription}
                onChange={(e) => setNewDescription(e.target.value)}
                placeholder="Например: Накопить на покупку квартиры"
              />
            </div>
            <div className="flex gap-2">
              <Button onClick={handleSave} className="flex-1">
                <Check className="h-4 w-4 mr-2" />
                Сохранить
              </Button>
              <Button onClick={handleCancel} variant="outline" className="flex-1">
                <X className="h-4 w-4 mr-2" />
                Отмена
              </Button>
            </div>
          </div>
        ) : (
          <div className="space-y-4">
            <div className="space-y-2">
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground">Прогресс</span>
                <span className="font-semibold">{progress.toFixed(1)}%</span>
              </div>
              <Progress value={progress} className="h-3" />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <p className="text-sm text-muted-foreground">Текущая стоимость</p>
                <p className="text-lg font-bold">
                  {new Intl.NumberFormat("ru-RU", {
                    style: "currency",
                    currency: "RUB",
                    maximumFractionDigits: 0,
                  }).format(currentValue)}
                </p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Целевая сумма</p>
                <p className="text-lg font-bold">
                  {new Intl.NumberFormat("ru-RU", {
                    style: "currency",
                    currency: "RUB",
                    maximumFractionDigits: 0,
                  }).format(goalValue)}
                </p>
              </div>
            </div>
            {remaining > 0 && (
              <div className="pt-2 border-t">
                <p className="text-sm text-muted-foreground">
                  Осталось до цели:{" "}
                  <span className="font-semibold text-foreground">
                    {new Intl.NumberFormat("ru-RU", {
                      style: "currency",
                      currency: "RUB",
                      maximumFractionDigits: 0,
                    }).format(remaining)}
                  </span>
                </p>
              </div>
            )}
            {progress >= 100 && (
              <div className="pt-2 border-t">
                <p className="text-sm font-semibold text-green-600">
                  🎉 Поздравляем! Цель достигнута!
                </p>
              </div>
            )}
          </div>
        )}
      </CardContent>
    </Card>
  )
}

