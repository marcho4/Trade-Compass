"use client"

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group"
import { Label } from "@/components/ui/label"

export type RiskLevel = "conservative" | "moderate" | "aggressive"

interface RiskSelectorProps {
  value: RiskLevel
  onChange: (value: RiskLevel) => void
}

const riskOptions = [
  {
    value: "conservative" as RiskLevel,
    label: "Консервативный",
    description: "Низкий риск, стабильный доход. Больший акцент на дивидендные акции голубых фишек.",
    color: "text-blue-600",
  },
  {
    value: "moderate" as RiskLevel,
    label: "Умеренный",
    description: "Средний риск, сбалансированный портфель. Комбинация стабильных и растущих компаний.",
    color: "text-yellow-600",
  },
  {
    value: "aggressive" as RiskLevel,
    label: "Агрессивный",
    description: "Высокий риск, высокая потенциальная доходность. Фокус на растущих компаниях.",
    color: "text-red-600",
  },
]

export const RiskSelector = ({ value, onChange }: RiskSelectorProps) => {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Уровень риска портфеля</CardTitle>
        <CardDescription>
          Выберите стратегию управления рисками для вашего портфеля
        </CardDescription>
      </CardHeader>
      <CardContent>
        <RadioGroup value={value} onValueChange={(val) => onChange(val as RiskLevel)}>
          <div className="space-y-3">
            {riskOptions.map((option) => (
              <div
                key={option.value}
                className="flex items-start space-x-3 space-y-0 rounded-md border p-4 transition-colors hover:bg-accent"
              >
                <RadioGroupItem value={option.value} id={option.value} />
                <div className="space-y-1 leading-none flex-1">
                  <Label
                    htmlFor={option.value}
                    className={`text-base font-semibold cursor-pointer ${option.color}`}
                  >
                    {option.label}
                  </Label>
                  <p className="text-sm text-muted-foreground">
                    {option.description}
                  </p>
                </div>
              </div>
            ))}
          </div>
        </RadioGroup>
      </CardContent>
    </Card>
  )
}

