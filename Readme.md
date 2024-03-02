# Оценка сложности паролей для рускоязычных пользователей

Практика на основе условия: https://github.com/cs-itmo/practice-2024/blob/master/tracks/passwords.md

## Принцип работы

Проверка сложности пароля состоит из 4 этапов.
1. Проверка регулярными выражениями: используется классическая проверка на "минимум одна цифра, прописная и заглавная буква, символ" и длину. Можно настраивать кастомные регулярные выражения и менять требуемую длину при помощи переменных окружения.
2. Проверка по словарю утекших паролей: производится простой поиск по словарю утекших паролей. Источник можно настраивать.
3. Расчет энтропии: рассчитывается энтропия пароля на основе встроенного словаря и словаря из предыдущего пункта. На основании энтропии считается время, за которое пароль может быть взломан. Время, являющееся минимально допустимым также настраивается.
4. Проверка транслита: пароль разбивается на логические части и проверяется по четырем условиям: наличие в словаре английского языка, наличие в словаре русского языка при переводе "по клавиатуре" (ака ghbdtn -> привет), наличие в словаре русского языка при транслите (ака privet -> привет), наличие в словаре русского (английского языка) после символьных подстановок (ака c@t -> cat).

Если все четыре пункта пройдены, то пароль считается надежным.

## Использование

Собрать: `go build cmd/main.go cmd/config.go`
Запустить: `./main --verbose privet@MIR37` 

## Простор для улучшения

1. Оптимизация: сейчас на время выполнения плавает от 1 до 16 секунд в зависимости от длины пароля. Если используется API, то время растет. Также есть смысл обрезать множество тестируемых подстрок и оптимизировать spell check.
2. Есть небольшая неупорядоченность с флагами и переменными окружения.
3. При маленьких границах находятся слова, которые в общем то не должны находиться, например "bbochk" читается как "мочь", что явно не то, что хотел пользователь.

## Замечания

Чем лучше словари (утечек, частот, слов) используются, тем лучше результаты :)