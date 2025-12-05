"use client"

import { useState, useEffect } from "react"
import Link from "next/link"

const COOKIE_CONSENT_KEY = "cookie_consent"

export function CookieBanner() {
  const [showBanner, setShowBanner] = useState(false)

  useEffect(() => {
    const consent = localStorage.getItem(COOKIE_CONSENT_KEY)
    if (!consent) {
      setShowBanner(true)
    }
  }, [])

  const acceptCookies = () => {
    localStorage.setItem(COOKIE_CONSENT_KEY, "accepted")
    setShowBanner(false)
  }

  const declineCookies = () => {
    localStorage.setItem(COOKIE_CONSENT_KEY, "declined")
    setShowBanner(false)
    // Optionally disable analytics here
    if (typeof window !== "undefined" && window.ym) {
      // Disable Yandex Metrika tracking
      window.ym(105649346, "notBounce")
    }
  }

  if (!showBanner) return null

  return (
    <div className="fixed bottom-0 left-0 right-0 z-50 p-4 bg-background border-t shadow-lg">
      <div className="container mx-auto max-w-4xl">
        <div className="flex flex-col sm:flex-row items-start sm:items-center gap-4">
          <div className="flex-1">
            <p className="text-sm text-foreground">
              Мы используем файлы cookies для улучшения работы сайта и анализа трафика.
              Продолжая использовать сайт, вы соглашаетесь с{" "}
              <Link href="/cookies" className="text-primary hover:underline">
                политикой использования cookies
              </Link>{" "}
              и{" "}
              <Link href="/privacy" className="text-primary hover:underline">
                политикой конфиденциальности
              </Link>
              .
            </p>
          </div>
          <div className="flex gap-2 shrink-0">
            <button
              onClick={declineCookies}
              className="px-4 py-2 text-sm border border-border rounded-md hover:bg-muted transition-colors"
            >
              Отклонить
            </button>
            <button
              onClick={acceptCookies}
              className="px-4 py-2 text-sm bg-primary text-primary-foreground rounded-md hover:bg-primary/90 transition-colors"
            >
              Принять
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}

// Extend Window interface for Yandex Metrika
declare global {
  interface Window {
    ym?: (id: number, action: string, ...args: unknown[]) => void
  }
}
