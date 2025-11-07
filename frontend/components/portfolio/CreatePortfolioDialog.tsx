"use client"

import { useState } from "react"
import { Plus } from "lucide-react"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"

interface CreatePortfolioDialogProps {
  onCreatePortfolio?: (name: string, initialAmount: number) => void
}

export const CreatePortfolioDialog = ({
  onCreatePortfolio,
}: CreatePortfolioDialogProps) => {
  const [open, setOpen] = useState(false)
  const [name, setName] = useState("")
  const [initialAmount, setInitialAmount] = useState("")

  const handleCreate = () => {
    if (name && initialAmount) {
      onCreatePortfolio?.(name, parseFloat(initialAmount))
      setName("")
      setInitialAmount("")
      setOpen(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button size="lg" className="gap-2">
          <Plus className="h-5 w-5" />
          Создать портфель
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Создание нового портфеля</DialogTitle>
          <DialogDescription>
            Введите название портфеля и начальную сумму инвестиций
          </DialogDescription>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="grid gap-2">
            <Label htmlFor="name">Название портфеля</Label>
            <Input
              id="name"
              placeholder="Например: Дивидендный портфель"
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="amount">Начальная сумма (₽)</Label>
            <Input
              id="amount"
              type="number"
              placeholder="100000"
              value={initialAmount}
              onChange={(e) => setInitialAmount(e.target.value)}
            />
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => setOpen(false)}>
            Отмена
          </Button>
          <Button onClick={handleCreate} disabled={!name || !initialAmount}>
            Создать
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

