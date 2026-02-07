"use client"

import { TopNavbar } from "@/components/layout/TopNavbar"
import { AIChatPanel } from "@/components/layout/AIChatPanel"
import { usePathname } from "next/navigation"
import { portfolioPrompts, screenerPrompts, companyAnalysisPrompts, defaultPrompts } from "@/lib/ai-prompts"
import { MessageSquare, X } from "lucide-react"
import { useState } from "react"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"

export default function Layout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname()
  const isLandingPage = pathname === "/"
  const [isChatOpen, setIsChatOpen] = useState(false)

  // Определяем какие промпты показывать в зависимости от страницы
  const getPromptExamples = () => {
    if (pathname.includes("/portfolio")) {
      return portfolioPrompts
    }
    // Страница анализа компании (формат: /dashboard/TICKER)
    if (pathname.match(/\/dashboard\/[A-Z]{4,}/)) {
      return companyAnalysisPrompts
    }
    if (pathname.includes("/screener")) {
      return screenerPrompts
    }
    if (pathname.includes("/dashboard")) {
      return defaultPrompts
    }
    return defaultPrompts
  }

  if (isLandingPage) {
    return <main>{children}</main>
  }

  return (
    <div className="flex flex-col min-h-screen bg-background">
      <TopNavbar />

      <div className="flex flex-1 pt-16">
        <main className="flex-1 px-4 md:px-8 py-8 overflow-y-auto">
          {children}
        </main>
      
        {/* Desktop - фиксированная панель справа */}
        <aside className="hidden lg:block w-[420px] py-8 pr-8 sticky top-16 h-[calc(100vh-4rem)]">
          <AIChatPanel promptExamples={getPromptExamples()} />
        </aside>
      </div>

      {/* Mobile/Tablet - плавающая кнопка и выезжающая панель */}
      <div className="lg:hidden">
        {/* Кнопка открытия чата */}
        <Button
          onClick={() => setIsChatOpen(true)}
          className={cn(
            "fixed bottom-6 right-6 z-50 h-14 w-14 rounded-full shadow-2xl",
            "hover:scale-110 transition-all duration-300",
            isChatOpen && "hidden"
          )}
          aria-label="Открыть AI чат"
        >
          <MessageSquare className="h-6 w-6" />
        </Button>

        {/* Оверлей */}
        {isChatOpen && (
          <div
            className="fixed inset-0 bg-black/50 z-40 backdrop-blur-sm"
            onClick={() => setIsChatOpen(false)}
            style={{
              animation: "fadeIn 0.3s ease-out"
            }}
          />
        )}

        {/* Панель чата */}
        <div
          className={cn(
            "fixed inset-x-0 bottom-0 z-50 transition-transform duration-300 ease-out",
            "h-[85vh] md:h-[80vh]",
            "md:inset-x-auto md:right-4 md:bottom-4 md:w-[420px] md:rounded-t-3xl",
            isChatOpen ? "translate-y-0" : "translate-y-full"
          )}
          style={{
            maxHeight: "calc(100vh - 2rem)"
          }}
        >
          <div className="relative h-full">
            {/* Кнопка закрытия - показывается только когда чат открыт */}
            {isChatOpen && (
              <Button
                onClick={() => setIsChatOpen(false)}
                variant="ghost"
                size="icon"
                className="absolute -top-12 right-4 z-10 h-10 w-10 rounded-full bg-background/90 backdrop-blur-sm shadow-lg"
                aria-label="Закрыть чат"
              >
                <X className="h-5 w-5" />
              </Button>
            )}

            <div className="h-full overflow-hidden rounded-t-3xl md:rounded-3xl">
              <AIChatPanel promptExamples={getPromptExamples()} />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}