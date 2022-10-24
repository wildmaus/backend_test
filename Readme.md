# Backend test
Простой rest api сервер на Go и Psql.
### Commands
Для запуска сервера (скачает postgressql image и создаст image для go app):
```bash
docker compose up --build app
```
### Возникшие вопросы
- Не было понятно нужно ли обрабатывать запросы для изменения баланса в случае, если пользователь еще не был создан в системе. В данной реализации методы post/put не отличаются и позволяют создавать пользователей. То же самое касается получателя (но только получателя) при отправке денег другому пользователю.
- Не было понятно на основе каких дат строить отчет - дат резервирования средств или дат списания. В данной реализации при резервировании создается транзакция со статусом 2 (резерв) и она демонстрируется пользователю как информация о списаниях. А при подтверждении списания создается еще одна транзакция (деньги при этом не списываются, только резерв) со статусом 4 (списанно) и текущим временем, на основании которой строится отчет.
- Не было понятно могут ли повторяться заказы и услуги. В данной реалзиации подразумевается что в одном заказе может быть несколько услуг, но в сумме id услуги + id заказа являются уникальными для каждого резерва.
### Нюансы реализации
- Type в транзакциях: 0 - зачисление средств, 1 - перессылка между пользователями, 2 - резервирование, 3 - отмена резервирования, 4 - списание резерва в отчет
- Пользователь при запросе своих транзакций не будет видеть транзакции c type = 4
- Пользователь при запросе своих транзакций может отсортировать их или по объему или по дате как в убывающем, так и в возрастающем порядке.
- Пользователь получает по 5 своих транзакций за запрос
- Изменения в балансах реализованы с помощью сабзапросов, а не отдельных вызовов из соображений безопасности - субд отвечает за атомарность исполнения, так что таким образом либо весь запрос выполниться, либо не произойдет ничего. Считаю что это более безопасный подход (на случай неожиданного отключения сервера в середине исполнения запроса)
- Отчет генерируется в виде _service id, amount_
