import Link from "next/link"

export function Footer() {
  const currentYear = new Date().getFullYear()

  return (
    <footer className="border-t bg-background">
      <div className="container mx-auto px-4 py-8">
        <div className="flex flex-col md:flex-row justify-between items-center gap-4">
          <div className="text-sm text-muted-foreground">
            © {currentYear} Trade Compass. Все права защищены.
          </div>

          <nav className="flex flex-wrap justify-center gap-4 md:gap-6">
            <Link
              href="/terms"
              className="text-sm text-muted-foreground hover:text-foreground transition-colors"
            >
              Пользовательское соглашение
            </Link>
            <Link
              href="/privacy"
              className="text-sm text-muted-foreground hover:text-foreground transition-colors"
            >
              Политика конфиденциальности
            </Link>
            <Link
              href="/cookies"
              className="text-sm text-muted-foreground hover:text-foreground transition-colors"
            >
              Cookies
            </Link>
          </nav>
        </div>

        <div className="mt-6 pt-6 border-t">
          <p className="text-xs text-muted-foreground text-center">
            Информация на сайте носит исключительно информационный характер и не является
            индивидуальной инвестиционной рекомендацией. Инвестиции в ценные бумаги связаны
            с риском потери капитала.
          </p>
        </div>
      </div>
    </footer>
  )
}