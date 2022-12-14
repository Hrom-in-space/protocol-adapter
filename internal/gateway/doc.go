/*
Package gateway - отвечает за сетевые запросы в другие системы

На каждое выполнение метода должно выполняться не больше одного сетевого запроса
На каждый сервис должен приходиться один Gateway

Из чего состоит gateway:
* gateway каждого сервиса расположен в отдельном файле и названом по имени сервиса в который он обращается
* У gateway есть конструктор который вызывается в сервисах
* Публичные методы делающие запросы в сервис
* Для десериализации результатов определены структуры ответа (как успешного, так и содержащего ошибки).
*/
package gateway
