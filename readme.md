# pass-keeper - приложение для хранения доступов с работой из консоли

Позволяет хранить различные типы доступов к сервисам. Имеет функционал автопроверки валидности.

### Описание защиты данных:

Основной упор защиты делается на мастер пароль, который вводит пользователь при начале работы с приложением.

Для шифрования в приложении используется не сам мастер пароль, а его хеш (мастер ключ). Это позволяет иметь 128 битный ключ при любой длине мастер пароля.
Для исключения возможности кражи файла с мастер ключем и использования его в дальнейшем, мастер ключ хранится в зашифрованном виде. Для шифрования мастер ключа используется уникальный идентификатор компьютера пользователя и закрытый ключ приложения. В случае простой кражи базы данных и файла с мастер ключем, у злоумышленника не появляется доступа к паролям в базе данных. Так как не известен идентификатор компьютера, с помощью которого зашифрован мастер ключ.

Для хеширования мастер пароля используется алгоритм md5. Для шифрования логина и пароля используется алгоритм AES с 256 битным ключем (128 бит - мастер ключ и 128 бит закрытый ключ приложения)

Чтобы получить логины и пароли необходимо иметь:
1. Базу данных
2. Файл с мастер ключем, который был создан для текущего ПК / либо файл с мастер ключем, который был создан для другого ПК и уникальный идентификатор ПК, для которого был создан ключ.
3. Приложение с зашитым ключем / либо ключ, с которым собиралось приложение и внутренний алгоритм работы приложения.

### Roadmap
- [x] Создание структуры базы данных
- [x] Создание типа SSH
- [x] CRUD операции над доступами
- [x] Putty интеграция (добавление доступов из putty, запуск putty по имеющимся доступам)
- [x] Функционал валидации доступов
- [ ] Создание типа FTP
- [ ] FileZilla интеграция
- [ ] PhpStorm интеграция
- [ ] Notepad++ интеграция
- [ ] Создание типа Битрикс