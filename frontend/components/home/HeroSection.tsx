import Image from "next/image"
import GridBackground from "@/components/ui/grid-background"
import { CTAButtons } from "./CTAButtons"

export const HeroSection = () => {
  return (
    <GridBackground className="min-h-screen overflow-hidden">
      <div className="relative z-20 flex h-full items-center justify-center px-4 py-12 sm:px-6 md:px-12 md:py-16">
        <div className="flex w-full max-w-7xl flex-col items-center gap-8 md:flex-row md:items-center md:gap-12 lg:gap-16">
          <div className="flex flex-1 flex-col items-start gap-6 text-left md:gap-8">
            <h1 className="animate-fade-in-up text-gray-900 bg-clip-text text-4xl font-bold leading-tight [animation-delay:100ms] sm:text-5xl md:text-5xl lg:text-6xl">
              Укажем направление для ваших инвестиций
            </h1>

            <p className="max-w-xl animate-fade-in-up text-lg leading-relaxed text-muted-foreground [animation-delay:200ms] md:text-xl">
              Анализируйте рынок в 10 раз быстрее. Без финансового образования.
            </p>

            <CTAButtons />
          </div>

          <div className="flex flex-1 items-center justify-center animate-fade-in-up [animation-delay:300ms]">
            <div className="relative overflow-hidden rounded-[2.5rem] border-2 border-gray-300 shadow-xl">
              <Image
                src="/landing-image.png"
                alt="Trade Compass - AI анализ акций российского фондового рынка"
                width={550}
                height={400}
                className="h-auto w-full max-w-[600px] object-cover"
                priority
              />
            </div>
          </div>
        </div>
      </div>
    </GridBackground>
  )
}

