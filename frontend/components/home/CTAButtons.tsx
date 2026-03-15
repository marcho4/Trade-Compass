import Link from "next/link"
import { ArrowRight } from "lucide-react"

export const CTAButtons = () => {
  return (
    <div className="flex animate-fade-in-up flex-col gap-4 [animation-delay:300ms]">
      <Link
        className="group inline-flex items-center justify-center rounded-full bg-primary px-10 py-5 text-lg font-bold text-primary-foreground shadow-md transition-colors hover:bg-primary/90 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
        href="/dashboard/screener"
        aria-label="Начать бесплатно"
      >
        Попробовать бесплатно
        <ArrowRight className="ml-2 h-5 w-5 transition-transform group-hover:translate-x-1" aria-hidden="true" />
      </Link>
      <p className="text-sm text-foreground">
        Без карты • Данные по всем крупным компаниям Мосбиржи
      </p>
    </div>
  )
}
