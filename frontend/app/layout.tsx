import type { Metadata } from "next";
import { Geist, Geist_Mono, Archivo, Instrument_Serif } from "next/font/google";
import { Providers } from "@/components/providers/Providers";
import Script from "next/script";
import "./globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

const archivo = Archivo({
  variable: "--font-archivo",
  subsets: ["latin", "latin-ext"],
  weight: ["400", "600", "700", "800"],
});

const instrumentSerif = Instrument_Serif({
  variable: "--font-instrument-serif",
  subsets: ["latin", "latin-ext"],
  weight: ["400"],
  style: ["normal", "italic"],
});

export const metadata: Metadata = {
  title: "Trade Compass - Платформа для анализа российских акций",
  description: "Собирайте данные с MOEX, анализируйте финансовые отчеты, формируйте анализ на базе AI и держите портфель в балансе.",
  applicationName: "Trade Compass",
  openGraph: {
    title: "Trade Compass - Платформа для анализа российских акций",
    description: "AI помощник разберет отчетность за 2 минуты и покажет, что видят институционалы. Без финансового образования.",
    siteName: "Trade Compass",
    locale: "ru_RU",
    type: "website",
  },
  twitter: {
    card: "summary_large_image",
    title: "Trade Compass - Платформа для анализа российских акций",
    description: "AI помощник разберет отчетность за 2 минуты и покажет, что видят институционалы.",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ru">
      <body className={`${geistSans.variable} ${geistMono.variable} ${archivo.variable} ${instrumentSerif.variable} antialiased`}>
        <Script
          id="yandex-metrika"
          strategy="afterInteractive"
          dangerouslySetInnerHTML={{
            __html: `
              (function(m,e,t,r,i,k,a){
                m[i]=m[i]||function(){(m[i].a=m[i].a||[]).push(arguments)};
                m[i].l=1*new Date();
                for (var j = 0; j < document.scripts.length; j++) {if (document.scripts[j].src === r) { return; }}
                k=e.createElement(t),a=e.getElementsByTagName(t)[0],k.async=1,k.src=r,a.parentNode.insertBefore(k,a)
              })(window, document,'script','https://mc.yandex.ru/metrika/tag.js?id=105649346', 'ym');
              ym(105649346, 'init', {ssr:true, clickmap:true, ecommerce:"dataLayer", accurateTrackBounce:true, trackLinks:true});
            `,
          }}
        />
        <noscript>
          <div>
            <img 
              src="https://mc.yandex.ru/watch/105649346" 
              style={{position: "absolute", left: "-9999px"}} 
              alt="" 
            />
          </div>
        </noscript>

        <Providers>
          <a
            href="#main-content"
            className="fixed left-4 top-4 z-[100] -translate-y-20 rounded-lg bg-primary px-4 py-2 text-sm font-semibold text-primary-foreground transition-transform focus:translate-y-0"
          >
            Перейти к содержимому
          </a>
          <div id="main-content" className="relative min-h-screen bg-background text-foreground">
            {children}
          </div>
        </Providers>
      </body>
    </html>
  );
}
