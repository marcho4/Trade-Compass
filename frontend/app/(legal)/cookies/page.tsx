import { Metadata } from "next"

export const metadata: Metadata = {
  title: "Политика использования cookies - Trade Compass",
  description: "Информация об использовании файлов cookies на сайте Trade Compass",
}

export default function CookiesPage() {
  return (
    <article className="prose prose-gray dark:prose-invert max-w-none">
      <h1>Политика использования файлов cookies</h1>
      <p className="text-muted-foreground">Последнее обновление: 5 декабря 2024 г.</p>

      <section>
        <h2>1. Что такое cookies</h2>
        <p>
          Файлы cookies (куки) — это небольшие текстовые файлы, которые сохраняются
          на вашем устройстве при посещении веб-сайтов. Cookies позволяют сайту
          «запоминать» ваши действия и настройки в течение определённого времени.
        </p>
      </section>

      <section>
        <h2>2. Как мы используем cookies</h2>
        <p>Сервис Trade Compass использует следующие типы cookies:</p>

        <h3>2.1. Необходимые cookies</h3>
        <p>
          Эти cookies обеспечивают базовую функциональность сайта и не могут быть
          отключены. Они включают:
        </p>
        <ul>
          <li>Cookies для поддержания сессии авторизации</li>
          <li>Cookies для запоминания согласия на использование cookies</li>
          <li>Cookies для обеспечения безопасности</li>
        </ul>

        <h3>2.2. Аналитические cookies</h3>
        <p>
          Мы используем аналитические cookies для понимания того, как посетители
          взаимодействуют с нашим сайтом. Это помогает нам улучшать качество Сервиса.
        </p>

        <h4>Яндекс.Метрика</h4>
        <p>
          Мы используем сервис веб-аналитики Яндекс.Метрика, предоставляемый
          ООО «ЯНДЕКС» (119021, Россия, г. Москва, ул. Льва Толстого, д. 16).
        </p>
        <p>Яндекс.Метрика собирает следующую информацию:</p>
        <ul>
          <li>Источник перехода на сайт</li>
          <li>Просмотренные страницы и время на сайте</li>
          <li>Технические данные (браузер, устройство, разрешение экрана)</li>
          <li>Географическое положение (на уровне города)</li>
          <li>Действия на сайте (клики, прокрутка)</li>
        </ul>
        <p>
          Подробнее о политике конфиденциальности Яндекса:{" "}
          <a
            href="https://yandex.ru/legal/confidential/"
            target="_blank"
            rel="noopener noreferrer"
            className="text-primary hover:underline"
          >
            https://yandex.ru/legal/confidential/
          </a>
        </p>
      </section>

      <section>
        <h2>3. Перечень используемых cookies</h2>
        <div className="overflow-x-auto">
          <table>
            <thead>
              <tr>
                <th>Название</th>
                <th>Назначение</th>
                <th>Срок хранения</th>
                <th>Тип</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>cookie_consent</td>
                <td>Хранение согласия на использование cookies</td>
                <td>1 год</td>
                <td>Необходимый</td>
              </tr>
              <tr>
                <td>session_token</td>
                <td>Поддержание сессии авторизации</td>
                <td>До закрытия браузера или 30 дней</td>
                <td>Необходимый</td>
              </tr>
              <tr>
                <td>_ym_uid</td>
                <td>Идентификатор пользователя Яндекс.Метрики</td>
                <td>1 год</td>
                <td>Аналитический</td>
              </tr>
              <tr>
                <td>_ym_d</td>
                <td>Дата первого визита (Яндекс.Метрика)</td>
                <td>1 год</td>
                <td>Аналитический</td>
              </tr>
              <tr>
                <td>_ym_isad</td>
                <td>Определение блокировщика рекламы</td>
                <td>2 дня</td>
                <td>Аналитический</td>
              </tr>
              <tr>
                <td>_ym_visorc</td>
                <td>Данные для Вебвизора (Яндекс.Метрика)</td>
                <td>30 минут</td>
                <td>Аналитический</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <section>
        <h2>4. Управление cookies</h2>
        <p>
          Вы можете управлять файлами cookies через настройки вашего браузера.
          Большинство браузеров позволяют:
        </p>
        <ul>
          <li>Просматривать установленные cookies</li>
          <li>Удалять cookies (все или выборочно)</li>
          <li>Блокировать cookies от определённых сайтов</li>
          <li>Блокировать все cookies</li>
          <li>Удалять все cookies при закрытии браузера</li>
        </ul>
        <p>
          <strong>Важно:</strong> Отключение cookies может повлиять на функциональность
          Сервиса. Некоторые функции могут стать недоступны.
        </p>

        <h3>Инструкции по управлению cookies в популярных браузерах:</h3>
        <ul>
          <li>
            <a
              href="https://support.google.com/chrome/answer/95647"
              target="_blank"
              rel="noopener noreferrer"
              className="text-primary hover:underline"
            >
              Google Chrome
            </a>
          </li>
          <li>
            <a
              href="https://support.mozilla.org/ru/kb/udalenie-kukov-i-dannyh-sajtov-v-firefox"
              target="_blank"
              rel="noopener noreferrer"
              className="text-primary hover:underline"
            >
              Mozilla Firefox
            </a>
          </li>
          <li>
            <a
              href="https://support.apple.com/ru-ru/guide/safari/sfri11471/mac"
              target="_blank"
              rel="noopener noreferrer"
              className="text-primary hover:underline"
            >
              Safari
            </a>
          </li>
          <li>
            <a
              href="https://support.microsoft.com/ru-ru/microsoft-edge/удаление-файлов-cookie-в-microsoft-edge-63947406-40ac-c3b8-57b9-2a946a29ae09"
              target="_blank"
              rel="noopener noreferrer"
              className="text-primary hover:underline"
            >
              Microsoft Edge
            </a>
          </li>
        </ul>
      </section>

      <section>
        <h2>5. Отказ от аналитических cookies</h2>
        <p>
          Вы можете отказаться от сбора данных Яндекс.Метрикой, установив
          специальное расширение для браузера:{" "}
          <a
            href="https://yandex.ru/support/metrica/general/opt-out.html"
            target="_blank"
            rel="noopener noreferrer"
            className="text-primary hover:underline"
          >
            Блокировка Яндекс.Метрики
          </a>
        </p>
      </section>

      <section>
        <h2>6. Изменения в политике</h2>
        <p>
          Мы можем обновлять данную политику использования cookies. Актуальная
          версия всегда доступна на этой странице. Рекомендуем периодически
          проверять её на предмет изменений.
        </p>
      </section>

      <section>
        <h2>7. Контактная информация</h2>
        <p>
          Если у вас есть вопросы о нашей политике использования cookies,
          свяжитесь с нами:
        </p>
        <ul>
          <li>Email: support@tradecompass.ru</li>
        </ul>
      </section>

      <section>
        <h2>8. Связанные документы</h2>
        <ul>
          <li>
            <a href="/privacy" className="text-primary hover:underline">
              Политика конфиденциальности
            </a>
          </li>
          <li>
            <a href="/terms" className="text-primary hover:underline">
              Пользовательское соглашение
            </a>
          </li>
        </ul>
      </section>
    </article>
  )
}
