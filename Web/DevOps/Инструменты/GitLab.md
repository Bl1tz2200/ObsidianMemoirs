#Web #DevOps  #Инструменты 
#CI/CD

GitLab - аналог Jenkins от создателей Git, поэтому он максимально синергичен с ним. Его главное отличие от Jenkins: у него есть платные версии из-за принадлежности к Microsoft
GitLab работает как локальный GitHub, однако есть и бесплатная версия, которая отличается тем, что сам GitLab мы будем хостить на своем сервере, а не на облачном
# Указание EXTERNAL_URl
Когда мы будем указывать EXTERNAL_URL рекомендую указать протокол http, однако GitLab может сам создать сертификат для https
# Подключение

После скачивания необходимо будет либо зайти через доменное имя, либо через ip адрес на 8080 порт (Лучше всего хотя-бы в /etc/hosts указать домен GitLab для корректной работы GitLab, а в особенности GitLab Runners)

Для логина и пароля необходимо использовать информацию из консоли, предоставленную GitLab после установки. После входа рекомендую создать нового пользователя, от чьего лица мы будем работать

# Gitlab-ctl
Для управления GitLab из консоли необходимо использовать gitlab-ctl

Остановить сервис GitLab:
```
gitlab-ctl stop
```

Запустить сервис GitLab:
```
gitlab-ctl start
```


Другие команды, которые обычно используют для сервисов, так же сработают и здесь

# Создание проекта
Для создания самого обычного проекта, нужно выбрать blank project. После создания проекта GitLab предоставит инструкцию для начала работы:

Необходимо будет через git сделать push в наш gitlab, после этого в GitLab появятся наши файлы. Таким образом мы можем спокойно делать git clone наших файлов, использовать другие команды git так же разрешено

## Настройки доступа
Настройки доступа к репозиторию в GitLab настраиваются создателем, однако рекомендуется ставить доступ не по логину и паролю, а по ssh .pub ключу. Да, доступ по ssh такой же, как и в GitHub, поэтому если вы умеете делать это в GitHub, сделать так же в GitLab будет не трудно

# Gilab Pipelines
GitLab Pipelines прописываются в *.gitlab-ci.yml*
Из-за использования yml файлов, написание Pipelines похоже на Ansible и Kubernetes
## Начало работы
В самом начале мы прописываем все этапы, которые у нас будут, а затем расписываем каждый подробнее. Такие этапы выполняются на маленьких системах GitLab Runner, которые запускаются в разных системах в зависимости от выбранного Executor (исполнитель)

В качестве GitLab Runner Executor можно использовать SSH, Shell, VirtualBox, Parrallels, Docker, Kubernetes и тд
### Подключение Runner 
Для подключения Runner к вашему репозиторию необходимо открыть Settings -> CI/CD -> Runner, после чего там будет указана подробная инструкция по запуску Runner с необходимым для вас Executor

После этого Pipeline будут исполнять команды в *.gitlab-ci.yml* через этот Runner
Далее необходимо добавить файл *.gitlab-ci.yml* в ваш репозиторий, тогда GitLab сможет запускать выбранный Runner для этого репозитория

## Tags
Через tags можно связывать разные stages с разными Runner, если к репозиторию подключено сразу несколько Runner. Для этого необходимо указать одинаковые tags как в stage в .*gitlab-ci.yml*, так и в tags у Runner при его создании

## Пример *.gitlab-ci.yml* с Runner на Docker

```
stages:
  - build
  - test
  - deploy staging
  - deploy production

build:
  stage: build
  image: ubuntu:20.04
  tags:
    - runner1
  script:
    - echo "Build" >> index.html
  artifacts:
    paths:
      - build/

test_on_dev:
  stage: test
  image: ubuntu:20.04
  tags:
    - runner1
  script:
    - test -f build/index.html

deploy_on_staging:
  stage: deploy staging
  image: ubuntu:20.04
  tags:
    - runner1
  script:
    - cp build/index.html /var/www/index.html
    - test -f /var/www/index.html


deploy_on_production:
  stage: deploy production
  image: ubuntu:20.04
  tags:
    - runner1
  script:
    - cp build/index.html /var/www/html/index.html
    - test -f /var/www/html/index.html
```

Для связки различных этапов нужно указывать этап из stages в атрибуте stage внутри нашего этапа. Этапы будут идти по очереди, сверху вниз, так, как они указаны в блоке stages. Если на каком-то этапе произошла ошибка, другие этапы не запустятся. Атрибут artifacts в начале необходим для сохранения файлов

## Игнорирование ошибок
Для игнорирования ошибок необходимо указать атрибут allow_failure: true:

```
stages:
  - test

test_on_dev:
  stage: test
  image: ubuntu:20.04
  tags:
    - runner1
  script:
    - test -f build/index.html
  allow_failure: true
```

## Shell как Executor
Использование Shell в качестве Executor позволит запускать команды на компьютере, на котором хостится gitlab. Благодаря этому можно проводить тестирование на изолированной машине, а интеграцию уже на вашем сервере:

```
stages:
  - build
  - test
  - deploy

build:
  stage: build
  image: ubuntu:20.04
  tags:
    - runner1
  script:
    - echo "Build" >> index.html
  artifacts:
    paths:
      - build/

test_on_dev:
  stage: test
  image: ubuntu:20.04
  tags:
    - runner1
  script:
    - test -f build/index.html


deploy:
  stage: deploy
  tags:
    - shell_runner
  script:
    - cp build/index.html /var/www/html/index.html
    - test -f /var/www/html/index.html
    - sudo bash -c 'echo "Using Sudo" > /root/made_by_gitlab.txt'
```

## Sudo в pipelines
Для использования sudo внутри Pipeline необходимо в /etc/sudoers.tmp добавить для пользователя gitlab-runner права ALL=(ALL) NOPASSWD: ALL

Также, для использования команд с sudo, необходимо использовать `sudo bash -c 'Ваша команда'`

## Переменные
Для использования переменных необходимо указывать блок variables:

```
stages:
  - build
  - test
  - deploy

variables:
  TEST_DIR: /etc/testdir

build:
  stage: build
  image: ubuntu:20.04
  tags:
    - runner1
  script:
    - echo "Build" >> index.html
  artifacts:
    paths:
      - $TEST_DIR/

test_on_dev:
  stage: test
  image: ubuntu:20.04
  tags:
    - runner1
  script:
    - test -f $TEST_DIR/index.html


deploy:
  stage: deploy
  tags:
    - shell_runner
  script:
    - cp build/index.html /var/www/html/index.html
    - test -f /var/www/html/index.html
    - sudo bash -c 'echo "Using Sudo" > /'$TEST_DIR'/made_by_gitlab.txt'
  when: manual
```

Из текста пред последней строчки можно заметить, что переменные нельзя использовать внутри кавычек, поэтому их приходится закрывать, указывать переменную и снова открывать

Последняя строчка отвечает за запуск последнего этапа: его можно запустить только руками, только самому. Такой метод запуска существует на тот случай, есть надо ещё раз перепроверить всю систему перед интеграцией и после тестирования

### Глобальные перемнные
Для указания глобальных переменных необходимо перейти в Setting -> CI/CD -> Variables. Глобальные переменные будут видны в тех репозиториях, которые мы укажем при её создании. Такие глобальные переменные, например, удобны для AWS, так как их можно сделать скрытыми

# Docker Container Registry
Есть такое дополнение для GitLab - Docker Container Registry. Оно находится в пункте Pakages & Registries и его может не быть изначально. Необходим он для хранения специально настроенных образов Docker, которые нам нужно часто использовать

Как его установить лучше посмотреть в [обширной документации GitLab](https://docs.gitlab.com/ee/administration/packages/container_registry.html) 

# [Источник](https://www.youtube.com/playlist?list=PLqVeG_R3qMSzYe_s3-q7TZeawXxTyltGC)

В качестве источника представлен плейлист на ютубе

Также можно рассмотреть [сайт](https://stafwag.github.io/blog/blog/2023/11/15/getting_started_with_gitlab-ce/#:~:text=The%20installation%20is%20straightforward%2C%20you,initial%20root%20password%20to%20STDOUT) с помощью по установке GitLab