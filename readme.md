# RIT Assignments
## Задача 1
⌛️ Потрачено: ~8 часов.

Наглядное демо на полторы минуты: https://www.loom.com/share/3ece7a950f684c678fdae3d54d417362

Запуск и тестирование проще всего делать через `Makefile`:

- `make run` запустит HTTP-сервер на `http://localhost:8080` и консоль для команд `send` и `sendc`. В браузере можно открыть и проверить отправку сообщений клиенту или бродкаст хабу. По умолчанию размер хаба равен 3. Размер можно поменять флагом `-s`.
- `make test` запустит тесты.

После запуска сервера можно вводить команды в консоль. Айдишники хабов и клиентов будут логироваться в консоли. Например:

```sh
client 8703bcc6-2dbf-4fc6-b36c-7114e05c57ad was added to hub 5803267d-c693-40c6-824c-fde259b888cc
> send 5803267d-c693-40c6-824c-fde259b888cc  # broadcast to hub
> sendc 8703bcc6-2dbf-4fc6-b36c-7114e05c57ad # send to client
```

## Задача 2
⌛️ Потрачено: ~3 часа.

Решение для PostgreSQL.

Сначала создадим таблицу и заполним её данными из задания:

```sql
CREATE TABLE IF NOT EXISTS timework(
  id SERIAL PRIMARY KEY,
  employee_id INTEGER NOT NULL,
  start_period TIMESTAMPTZ,
  end_period TIMESTAMPTZ
);

INSERT INTO timework (id, employee_id, start_period, end_period)
VALUES  (1, 1, '2021-08-11 12:00:00.000', '2021-08-11 18:00:00.000'),
        (2, 1, '2021-08-12 12:00:00.000', '2021-08-12 18:00:00.000'),
        (3, 1, '2021-08-16 12:00:00.000', '2021-08-16 18:00:00.000'),
        (4, 2, '2021-08-11 12:00:00.000', '2021-08-11 18:00:00.000'),
        (5, 2, '2021-08-14 12:00:00.000', '2021-08-14 18:00:00.000');
```

Для каждого сотрудника можно получить кортеж с днями, когда он работал:

```sql
SELECT
  employee_id as employee,
  array_agg(dd::date) as working_days
FROM
  generate_series('2021-08-01', '2021-08-31', '1 day'::interval) as dd
JOIN timework ON
  timework.start_period >= dd
  AND timework.end_period < dd + interval '1' day
GROUP BY employee_id
ORDER BY employee_id;
```

Получим таблицу вида:

```
employee | working_days
1        | {2021-08-11,2021-08-12,2021-08-16}
2        | {2021-08-11,2021-08-14}
```

Если нужно нагляднее показать дни, а не сотрудников, то можно агрегировать по-другому:

```sql
SELECT
  date_trunc('day', dd)::date AS day,
  array_agg(employee_id)
    FILTER (WHERE employee_id IS NOT NULL) AS employees
FROM
  generate_series ('2021-08-01', '2021-08-31', '1 day'::interval) AS dd
LEFT JOIN timework ON timework.start_period >= dd
  AND timework.end_period < dd + interval '1' day
GROUP BY dd
ORDER BY day;
```

Получим таблицу, где для кадого дня будет подставлен кортеж с сотрудниками, которые в этот день работали:

```
day        | employees
2021-08-01 |
...
2021-08-10 |
2021-08-11 | {1,2}
2021-08-12 | {1}
2021-08-13 |
2021-08-14 | {2}
2021-08-15 |
2021-08-16 | {1}
2021-08-17 |
...
2021-08-31  
```

Если табличный вывод из SQL нужен именно такой, как в задании, можно динамически генерировать таблицу с именами колонок как дни месяца и для каждой проверять наличие рабочих часов. Но у меня так сделать не получилось за адекватное время 🤷🏻‍♂️