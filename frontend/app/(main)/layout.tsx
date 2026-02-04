"use client"

import { FloatingSidebar } from "@/components/layout/FloatingSidebar"
import { AIChatPanel } from "@/components/layout/AIChatPanel"
import { usePathname } from "next/navigation"
import { portfolioPrompts, screenerPrompts, companyAnalysisPrompts, defaultPrompts } from "@/lib/ai-prompts"

export default function Layout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname()
  const isLandingPage = pathname === "/"

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
    <div className="flex min-h-screen bg-background">
      <FloatingSidebar />

      <main className="flex-1 px-8 py-8 ml-20 overflow-y-auto">
        {children}
      </main>
      
      <aside className="w-[420px] py-8 pr-8 sticky top-0 h-screen">
        <AIChatPanel promptExamples={getPromptExamples()} />
      </aside>
    </div>
  );
}