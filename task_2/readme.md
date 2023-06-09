Работат с сетью и файловой системой.

Задача: обойти сайты данные на вход, получить от них коды ответы 2/3/4/5, и записать результат в файл, построчно:
https://google.com 200
https://someOtherSite.com 500
и тд

Сайты в вход выглядят так:
`sites := string[]{"https://www.avito.ru/", "https://www.ozon.ru/", "https://vk.com/", "https://yandex.ru/", "https://www.google.com/"}`,
вы можете записать больше сайтов и запрашивать тысячи сайтов.

Что следует реализовать:

- [x] Получение результата (как минимум)
- [x] Конкурентные запросы - обходить кучу сайтов последовательно очень долго
- [x] Лимит на паралельные запросы - каждая горутина, это занятый файловый дескриптер, необходимо ограничить их
  максимальное
  количество, например 32/64 параллельных сетевых запроса
- [x] Таймаут - если сайт долго не отвечает, нужно прервать запрос. Установите таймаут в 2 секунды или любой другой
- [x] Обработка ошибок - при ошибке можно игнорировать результат, не оставаливая всю программу
- [x] Ретрай - при провале запроса имеет смысл попробовать снова еще несколько раз
- [x] Постепенная запись в файл - имея слишком много данных на вход мы не можем сначала обойти все сайты, а потом за раз
  записать результат в файл. Есть риск, что наш процесс потребует слишком много оперативной памяти и будет убит OOM.
  Следует писать результаты в файл бачами, например по 128/256/512 строк за раз
