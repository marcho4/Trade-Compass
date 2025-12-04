import { MetadataRoute } from 'next'

const BASE_URL = 'https://trade-compass.ru'

export default function sitemap(): MetadataRoute.Sitemap {
  const currentDate = new Date().toISOString()

  return [
    {
      url: BASE_URL,
      lastModified: currentDate,
      changeFrequency: 'weekly',
      priority: 1.0,
    },
    {
      url: `${BASE_URL}/auth`,
      lastModified: currentDate,
      changeFrequency: 'monthly',
      priority: 0.5,
    },
    {
      url: `${BASE_URL}/auth/register`,
      lastModified: currentDate,
      changeFrequency: 'monthly',
      priority: 0.5,
    },
  ]
}
