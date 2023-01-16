# RIT Assignments

## Задача 1
Запуск и тестирование проще всего делать через `Makefile`:

- `make run` запустит HTTP-сервер на `http://localhost:8080` и консоль для команд `send` и `sendc`. В браузере можно открыть и проверить отправку сообщений клиенту или бродкаст хабу.
- `make test` запустит тесты.

После запуска сервера можно вводить команды в консоль. Айдишники хабов и клиентов будут видны в консоли, их можно использовать для отправки сообщений. Например:

```sh
> send 3e593c30-9593-471d-9ac1-0cca0deab35b  # broadcast to hub
> sendc 8f9f8ab5-1ccc-458a-91e1-494ba2b4010e # send to client
```



https://github.com/gorilla/websocket/tree/master/examples/chat

## Задача 2
Решение для PostgreSQL. Сначала создадим таблицу и заполним её данными из задания:

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

Если нужно нагляднее показать дни, а не сотрудников, то можно аггрегировать по-другому:

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

Если табличный вывод из SQL нужен именно такой, как в задании, можно динамически генерировать таблицу с именами колонок как дни месяца и для кадой проверять наличие рабочих часов. Но у меня это так сделать не получилось 🤷🏻‍♂️