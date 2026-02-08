import Link from "next/link";

export const Header = (showLogin: boolean = true) => {
  return (
    <header className="fixed top-0 left-0 right-0 z-50 w-full border-b border-border/40 bg-background/80 backdrop-blur-xl">
      <nav
        aria-label="Главная навигация"
        className="mx-auto flex max-w-7xl items-center justify-between px-6 py-4"
      >
        <Link
          href="/"
          className="text-xl font-bold tracking-tight text-foreground transition-colors hover:text-primary"
        >
          TradeCompass
        </Link>

        {showLogin && (
          <Link
            href="/auth"
            className="text-lg font-bold tracking-tight text-gray-900 transition-all hover:text-primary focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
          > 
            Войти
          </Link>
        )}
        
      </nav>
    </header>
  );
};

