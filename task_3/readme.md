Дан код, делающий асинхронные вызовы функции: https://go.dev/play/p/DSku2XonSrT

- Нужно реализовать rate limiter, который ограничит максимально количество запросов: 10 в секунду
- Пусть лимит будет настраиваемым: 1, 10, 50 ... в секунду
- Выберите алгоритм для реализации. Например: Token Bucket

Если хотите больше челленжа, то реализуйте лимитер на ином алгоритме:
- Token Bucket
- Leaky Bucket
- Fixed Window
- Sliding Window

Несколько типов алгоритмов: https://betterprogramming.pub/4-rate-limit-algorithms-every-developer-should-know-7472cb482f48