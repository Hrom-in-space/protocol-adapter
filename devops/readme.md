## Старт сервиса через docker-compose
docker-compose -f ./devops/local/docker-compose.yml up
## Линтеры
[golangci-lint](https://github.com/golangci/golangci-lint) - линтер\
https://github.com/golangci/golangci-lint/releases/tag/v1.46.2 - текущая версия \
make link - команда запуска(требует наличие исполняемого файла golangci-lint в path)
## Unit-тесты
make test
## Тест работоспособности
Проверка работоспособности осуществляется утилитой `grpc_health_probe`\
[docs](https://github.com/grpc-ecosystem/grpc-health-probe)\
Example: `/app/grpc_health_probe -addr=localhost:50051`
