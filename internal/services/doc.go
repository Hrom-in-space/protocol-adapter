/*
Package services - содержит имплементации RPC сервисов
Публичные только RPC

Каждый сервис должен распологаться в своем файле
Каждый сервис должен имплементировать grpcserver.Servicer для его регистрации в сервере

Методы имплементирующие RPC должны называться также как и в спецификации
Можно вызывать только один gateway своего сервиса
*/
package services