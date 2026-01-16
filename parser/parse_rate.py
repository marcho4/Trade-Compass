from bs4 import BeautifulSoup
from datetime import datetime
from decimal import Decimal

with open('ставка.html', 'r', encoding='utf-8') as f:
    soup = BeautifulSoup(f.read(), 'html.parser')

rows = soup.find_all('tr')[1:]

raw_data = []
for row in rows:
    cells = row.find_all('td')
    if len(cells) == 2:
        date_str = cells[0].text.strip()
        rate_str = cells[1].text.strip().replace(',', '.')
        date = datetime.strptime(date_str, '%d.%m.%Y').date()
        rate = Decimal(rate_str)
        raw_data.append({'date': date, 'rate': rate})

raw_data.sort(key=lambda x: x['date'])

prev = None
values = []
for item in raw_data:
    if prev is None or item['rate'] != prev:
        values.append(f"('{item['date']}', {item['rate']})")
        prev = item['rate']

print('INSERT INTO CBRate(date, rate) VALUES')
print(',\n'.join(values) + ';')


