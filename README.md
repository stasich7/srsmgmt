# SRS management

Задача сервиса **SRS management** - управление через API видеопотоками, обрабатываемыми медиасервером SRS.

**Simple Realtime Server ((SRS)** позволяет производить стриминг RTMP, SRT и отдавать его в виде HLS потока с предварительным транскодированием с мультибитрейтом. Однако управлять потоками у **SRS** нет возможности. Также нет возможности ограничить доступ стримингу видео на базе авторизации. 

**SRS management** добавляет такую функциональность:
- создание потока с заданным паролем публикации для защиты от несанкционированного стриминга;
- получение информации о потоке , его состоянии , дате начала и окончания стриминга;
- остановка потока, с разрывом входящего потока;
- формирование HLS плейлистов с мультибитрейтом;
- формирование live HLS плейлиста, содержащего только 6 последних чанков видео, без возможности перемотки;
- формаирование DVR HLS плейлиста, позволяющего осуществлять перемотку;
- удаление потока с удалением всех видеозаписей данного потока.

Данные функции реализованы в API по HTTP REST и gRPC протоколам с авторизацией по техническому токену, передаваемому в Authorization-header.

Сервис **SRS management** взаимодействует с файлами, генерируемыми **SRS** (ts, m3u8), поэтому запущен на том же сервере и имеет права доступа к файлами **SRS**.

Для хранения информации о потоках сервис использует БД Postgres.

Базовая конфигурация **SRS management** сервиса:

```
DATABASE_URI=host=<host> port=<port> user=<dbuser> password=<dbpassword> dbname=<dbuser> sslmode=disable
HTTP_ADDR=127.0.0.1:8887
GRPC_ADDR=127.0.0.1:9087
SRS_ADDR=<SRS media server API>
HLS_ADDR=<Nginx media location>
LIVE_TS_PATH=<SRS playlists and TS location>
RTMP_ADDR=<RTMP publish location>
SRT_ADDR=<SRT publish location>
APIKEY=<secret API key>
SRS_CONF_PATH=<SRS config location on disk>
CACHE_TTL=3
PPROF_ENABLED=true
DEBUG=true
```

Базовая конфигурация **SRS** 
[srs_base.tpl](config/srs_base.tpl)
