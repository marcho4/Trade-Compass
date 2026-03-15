import Link from "next/link"
import { ArrowRight } from "lucide-react"

export const FinalCTASection = () => {
  return (
    <section className="my-20">
      <div className="overflow-hidden rounded-3xl border border-primary/20 bg-gradient-to-br from-primary/10 via-primary/5 to-background p-12 text-center shadow-lg md:p-16">
        <h2 className="mb-6 text-4xl font-bold tracking-tight md:text-5xl">
          Хватит читать отчёты на 200 страниц
        </h2>

        <p className="mx-auto mb-8 max-w-2xl text-lg text-foreground md:text-xl">
          Trade Compass собирает данные, AI объясняет их смысл, ты принимаешь
          решения. Как должно быть.
        </p>

        <Link
          href="/auth/register"
          className="group inline-flex items-center justify-center rounded-full bg-primary px-10 py-5 text-lg font-bold text-primary-foreground shadow-md transition-colors hover:bg-primary/90 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
          aria-label="Анализ SBER бесплатно"
        >
          Проанализировать Сбербанк бесплатно
          <ArrowRight className="ml-2 h-5 w-5 transition-transform group-hover:translate-x-1" aria-hidden="true" />
        </Link>

        <p className="mt-6 text-sm text-foreground">
          10x быстрый анализ • Без карты
        </p>
      </div>
    </section>
  )
}
