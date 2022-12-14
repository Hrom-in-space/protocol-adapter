# Информация для доработки сервиса


## Правила оформления Merge request
Название должно соответствовать маске {TASK-id}_{короткое_описание}
Размер МР должен быть минимально возможным 
Один МР - одна причина для его существования(одно изменение)

## Спецификация
* Protocol-adapter использует спецификации из репозитория <link>
* gRPC Сервисы расположены в <link>
* Спецификация привязана к проекту git-сабмодулем `proto`
* protobuf генерируется командой `make proto` в `internal/generated/...`
   * в `internal/generated/server/` - пакеты с protobuf файлами для сервисов
   * в `internal/generated/common/` - пакеты с protobuf файлами для вспомагательных сообщений

## запуск тестов и линтеров
`make lint` - линтер
`make test` - тесты

## F.A.Q.
### Куда нужно вносить изменения чтобы чтобы добавить/доработать сервис
* [спецификации](#Спецификация) - внести изменения в репозиторий со спецификациями, обновить сабмодуль и сгенерировать новые protobuf
* настройки(package config) - добавить/изменить настройки сервиса
* Gateway(package gateway) - добавить/изменить методы обращения к удаленным сервисам
* Service(package services) - добавить/изменить RPC которые используют Gateway
* ./main.go() - инициализировать и зарегистрировать сервис, добавить создание нового клиента если необходимо.
