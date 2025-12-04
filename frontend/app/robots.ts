import { MetadataRoute } from 'next'

const BASE_URL = 'https://trade-compass.ru'

export default function robots(): MetadataRoute.Robots {
  return {
    rules: [
      {
        userAgent: '*',
        allow: '/',
        disallow: [
          '/dashboard/',
          '/welcome/',
          '/auth/forgot-password/',
          '/api/',
        ],
      },
      {
        userAgent: 'Googlebot',
        allow: '/',
        disallow: [
          '/dashboard/',
          '/welcome/',
          '/auth/forgot-password/',
          '/api/',
        ],
      },
      {
        userAgent: 'Yandexbot',
        allow: '/',
        disallow: [
          '/dashboard/',
          '/welcome/',
          '/auth/forgot-password/',
          '/api/',
        ],
      },
    ],
    sitemap: `${BASE_URL}/sitemap.xml`,
    host: BASE_URL,
  }
}
