#Web #Мониторинг #Инструменты 

Prometheus - это open source проект, написанный на языке Go, являющийся отличной бесплатной системой мониторинга и оповещения
Он хранит метрики в виде TSDB (Time Series Data Base), то есть информация о метриках сохраняется с отметкой времени
За счет открытости исходного кода Prometheus большое комьюнити с кучей дополнений
# Метрики
В Prometheus существует несколько типов метрик:
- Counter - метрика, которая может только увеличиваться, можно обнулить при перезапуске наблюдаемого сервиса
- Gauge - метрика, которая может увеличиваться и уменьшаться
- Histogram - распределение метрики по группам (если метрики в этом промежутке, то идут в одну группу, если нет, то в другую)
- Summary - показывает процентное/квантильное соотношение групп из Histogram

# Принцип работы
Prometheus работает следующим образом: он будет обращаться по их ip адресам и собирать данные из /metrics, затем будет записывать их себе в базу данных

Для запроса метрик используется PromQL (Prometheus Query Language). Этот язык используется в стандартном графическом интерфейсе Prometheus и нужен для поиска метрик

Обычно в качестве графического интерфейса используют Grafana, поскольку стандартный GUI Prometheus не достаточно гибок

# Exporter
Мониторинговый агент в Prometheus называется Exporter
Существует просто ОГРОМНОЕ количество готовых Exporters для разных систем
## Доступ
Exporter обычно доступен по определенному порту на ip адресе компьютере, на котором запущен Prometheus, а сами метрики хранятся в разделе /metrics
Пример: Node Exporter публикует информацию о cебе на `http://127.0.0.1:9100`, а метрики хранит на: `http://127.0.0.1:9100/metrics`

# Установка и запуск
Чтобы установить Prometheus Server или Exporter необходимо обратиться к их [официальному сайту](https://prometheus.io/)

## Установка Prometheus Server
На один компьютер устанавливается Prometheus Server (по умолчанию работает на порту 9090), на другие компьютеры, за которыми мы хотим следить, устанавливаются Exporters, после этого в конфигурационном файле Prometheus Server (называется этот файл prometheus.yml) мы прописывает какие метрики откуда собирать

Сам же Prometheus Server состоит из 3 файлов: бинарник prometheus (сам Prometheus Server), prometheus.yml (конфигурационный файл) и папка data (появляется после запуска, туда сохраняются все метрики)

На Prometheus Server можно установить Alert Manager, который будут отправлять алерты почти куда угодно

### Prometheus.yml
Стандартный *prometheus.yml* (после удаления лишних изначальных настроек для Alert Manager) 
выглядит следующим образом:

```
global:
  scrape_interval: 10s

scrape_configs:
  - job_name: "prometheus"

    static_configs:
      - targets: 
        - localhost:9090
        - 192.168.1.197:9100
```

global - глобальные настройки для prometheus (scrape_interval обозначает время, через которое prometheus обновляет метрики)
scrape_configs - настройки для Exporters, там мы указываем название и static_configs, в котором мы указываем цели: ip адреса, откуда мы будем получать метрики
По умолчанию prometheus получает свои метрики, в данном примере он получает метрики ещё и с *192.168.1.197:9100*

## Запуск Prometheus Server
Для запуска Prometheus Server нужно просто запустить бинарник prometheus после конфигурации
### Запуск Prometheus Server через сервис
Для облегчения запуска сервера его делают сервисом, для этого:

1. бинарник перемещают в `/usr/bin` для доступа отовсюду
2. *prometheus.yml* перемещают в `/etc/prometheus`, там же создают папку для данных, в моем случае `/etc/prometheus/data`

```
sudo su
mv prometheus /usr/bin
mkdir /etc/prometheus
mkdir /etc/prometheus/data
mv prometheus.yml /etc/prometheus/
```

3. Затем создают пользователя специально для prometheus и дают ему доступ ко всему prometheus

```
useradd -rs /bin/false prometheus
chown prometheus:prometheus /usr/bin/prometheus
chown -R prometheus:prometheus /etc/prometheus/
```

4. После создают сервис prometheus

```
vi /etc/systemd/system/prometheus.service
```

5. В этот файл вставляют самый обычный сервис (можно поменять папку сохранения данных в `--storage.tsdb.path`):

```
[Unit]
Description=Prometheus Server
After=network.target

[Service]
User=prometheus
Group=prometheus
Type=simple
Restart=on-failure
ExecStart=/usr/bin/prometheus \
  --config.file         /etc/prometheus/prometheus.yml \
  --storage.tsdb.path   /etc/prometheus/data

[Install]
WantedBy=multi-user.target
```

6. После этого перезагружают все систему сервисов и запускают prometheus.service:

```
systemctl daemon-reload
systemctl start prometheus
```

7. Для автоматического запуска при перезагрузке компьютера необходимо его активировать:

```
systemctl enable prometheus
```

Вот так и запускается prometheus профессионально
(Можно сделать и проще: [скачать файл](https://github.com/adv4000/prometheus) который установит все за нас, надо лишь указать папку для данных внутри файла .sh)
## Установка Exporter и запуск его как сервис
Установим Exporter на сервер Ubuntu: для этого просто скачаем его с их сайта. В примере будет установлен node_exporter который работает на 9100 порту

После установки у нас будет 1 бинарник, как и в случае установки Prometheus Server, однако закрепить в системе его стоит по профессиональному: опять в виде сервиса, для этого повторяем процедуру, которую мы делали в прошлом разделе про профессиональную установку Prometheus Server:

```
sudo su
useradd -rs /bin/false node_exporter
mv node_exporter /usr/bin/node_exporter
chown node_exporter:node_exporter /usr/bin/node_exporter
```

Затем создают сервис node_exporter:

```
vi /etc/systemd/system/node_exporter.service
```

И в нем опять указываем стандартный сервис:

```
[Unit]
Description=Node Exporter
After=network.target

[Service]
User=node_exporter
Group=node_exporter
Type=simple
Restart=on-failure
ExecStart=/usr/bin/node_exporter

[Install]
WantedBy=multi-user.target
```

Затем перезагружаем демона и запускаем сервис (активируя автозапуск при перезагрузке):
```
systemctl daemon-reload
systemctl enable node_exporter
systemctl start node_exporter
```

Теперь наш node_exporter будет постоянно показывать метрики с нашего сервера на порту 9100 в разделе /metrics
(Ах да, можно опять, с той же [ссылки](https://github.com/adv4000/prometheus), которую мы использовали при предыдущей установке, можно сделать все эти действия одним файлом, внутри которого необходимо будет лишь указать нужную нам папку)

А в самом *prometheus.yml* на нашем Prometheus Server укажем наш ubuntu_server в виде нового job:

```
global:
  scrape_interval: 10s

scrape_configs:
  - job_name: "prometheus"
    static_configs:
      - targets:
        - localhost:9090

  - job_name: "ubuntu_server"
    static_configs:
      - targets:
        - 192.168.1.197:9100
```

Теперь наш Prometheus Server раз в 10 секунд обновляет информацию, получаемую с node_exporter на *192.168.1.197:9100*

# Grafana
Перейдем к GUI в виде Grafana, для этого нам необходимо установить Grafana Server и связать его с Prometheus Server
Сама по себе Grafana - отличная утилита для веб визуализации различных данных
## Установка
Устанавливается она с официального сайта, главное скачивать grafana open source (OSS) (а не платную версию), работает по умолчанию на порту *3000*, сохраняет все данные в `/etc/grafana`, стандартные логин и пароль: admin admin

После входа в разделе Connections -> Data sources выбираем наш Prometheus
Далее указываем ссылку на наш Prometheus Server и другие различные настройки по необходимости

В конце рекомендуется в разделе Dashboard импортировать необходимые для нас страницы визуализации (хотя можно сделать и самому, но зачем заново изобретать велосипед). Для импортирования нужно скачать JSON файл Dashboard, а затем закинуть в разделе импорта

Теперь у нас будет просто прекраснейший инструмент визуализации нашей рабочей системы
### Автоматический запуск
Для автоматического запуска Grafana после перезагрузки компьютера нужно её активировать, но перед этим перезапустить демон:

```
systemctl enable grafana-server
systemctl daemon-reload
systemctl start grafana-server
```
# [Источник](https://www.youtube.com/playlist?list=PLg5SS_4L6LYu6qjwwwjt2WRsEudkzqB7J)
В качестве источника представлен видеокурс, также есть [ссылка на файлы для автоустановки prometheus и exporters в виде сервисов](https://github.com/adv4000/prometheus)