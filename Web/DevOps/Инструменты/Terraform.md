#Web #DevOps #Инструменты 
#Облако #IaC

Terraform позволяет собирать облачные виртуалки (aka облачный Ansible), управлять ими и так далее. Абсолютно бесплатен. Лучше всего использовать связку Terraform + Ansible: Terraform создает виртуалки у провайдера, Ansible производит редактирование внутри них

# Основные команды
Terraform имеет не много необходимых команд для запуска:

- `terraform init `- инициализация Terraform

 - `terraform plan` - дебаг перед запуском, показывает что изменит Terraform после запуска

- `terraform apply` - запуск того, что мы указали в файле .tf, то есть развертка

- `terraform destroy` - удаление виртуалок, указанных в файле

- `terraform console` - запускает консоль терраформа, которая позволяет проверять команды, тестировать их и так далее 

- `terraform show` - выводит файл .tfstate со всеми параметрами

- `terraform import` - сохраняет что-то, созданное вручную, в заданный resource (по поводу импорта подробнее в конце файла)

- `terraform state` - взаимодействие с .tfstate, который может находиться удаленно (подробнее в конце файла в пункте Refectoring)

- `terraform workspace` - взаимодействие с workspace, позволяющий проводить тесты (подробнее в конце файла в пункте Terraform Workspace)

- `terraform taint` - помечает ресурс для пересоздания, после чего ресурс пересоздается при apply (можно использовать аналог в виде `terraform apply -replace ресурс_для_пересоздания`)

# HCL
Для написания кода в Terraform используется язык Hashicorp Configuration Language (HCL), который напоминает yaml, то есть нам нужно указывать то, что мы хотим получить

# Исполняемые файлы
Terraform работает с файлами формата `.tf`, причем если таких файлов в папке много, то он просто объединит их в один во время запуска. Обычно разделяют все на 3 файла: *main.tf*, *variables.tf*, *outputs.tf*

Ещё стоит разделять все технологии по отдельным папкам и использовать их как модули, таким образом можно будет избежать неразберихи в файле большим количеством строчек. Так можно будет видеть: с какой папки пришла ошибка, где что-то пошло не так и так далее

# Resource и его параметры
Для создания виртуалок используется блок `resource`, в котором нужно указывать параметры для запуска

Перед запуском нужно прочитать документацию облачной компании, предоставляющей ресурсы, ведь она полезна тем, что показывает все возможные для использования параметры, покажу пример с AWS:

```
resource "aws_instance" "my_ubuntu" {
  ami           = "ami-03a71cec707bfc3d7"
  instance_type = "t3.micro"
}
```

Виртуалка у провайдера AWS с названием *my_ubuntu*, размером вируталки *t3.micro* и идентификатор образа машины Amazon для экземпляра EC2 *ami-03a71cec707bfc3d7* (тут уже не нужен внешний блок `terraform {}`)

# Подключение к провайдеру
Ниже представлен код для конфигурации провайдера, у которого мы будем создавать виртуалки нужно:

```
terraform {
	provider "aws" {
  		access_key = "ACCESS_KEY"
  		secret_key = "SECRET_KEY"
  		region     = "eu-central-1"
	}
}
```

Для работы с провайдером необходимо прочитать документацию по работе с ним на сайте Terraform.

## Множество провайдеров
Для использования множества провайдеров:

```
provider "aws" { # Провайдер по умолчанию
    region = "eu-central-1"
    assume_role { # Смена роли в данной сессии
       role_arn = "arn:aws:iam:1234567890:role/RemoteAdmin" # Указывем роль, на которую меняемся (аналог access_key, secret_key)
       session_name = "Terraform_Session" # Меняем сессию на указанную
    }
}

provider "aws" { # Провайдер Америки с псевдонимом USA
    region = "us-east-1"
    alias = "USA"
}

provider "aws" { # Провайдер Канады с псевдонимом Canada
    region = "ca-central-1"
    alias = "Canada"
}

resource "aws_instance" "servers" { # Создаем сервер с провейдером по умолчанию
  ami = "ami-03a71cec707bfc3d7"
  instance_type = "t3.micro"
}

resource "aws_instance" "servers" { 
  provider aws.USA # Создаем сервер с провейдером USA, указанным через .
  ami = "ami-03a71cec707bfc3d7"
  instance_type = "t3.micro"
}

resource "aws_instance" "servers" {
  provider aws.Canada # Создаем сервер с провейдером Canada, указанным через .
  ami = "ami-03a71cec707bfc3d7"
  instance_type = "t3.micro"
}
```

# Подстановка переменных окружения
Мы задаем параметры входа провайдера AWS через ключ доступа и паролю (Access/Secret keys) в регионе центральной европы, однако хранить секреты так открыто не стоит, поэтому стоит посмотреть документацию на секреты, а именно под какие переменные окружения терминала их можно подставить и заменяем:

```
terraform {
	provider "aws" {}
}
```
, а в терминале: 
```
export AWS_ACCESS_KEY_ID=ACCESS_KEY
export AWS_SECRET_ACCESS_KEY=SECRET_KEY
export AWS_DEFAULT_REGION=eu-central-1
```

И Terraform сам подставит переменные в провайдер, переменные будут защищены и смогут быть использованы только в терминале

# Использование переменных в файле
Для использования переменных в файлах можно использовать нужно использовать шаблоны, то есть вместо:

`user_data = file("./user_data.sh")`

Нужно переименовать user_data.sh в user_data.sh.tpl:

```
user_data = templatefile("user_data.sh.tpl", {
	f_name = "Denis",
	names = ["Oleg", "Petya", "Masha"]
})
```

Тогда:
*user_data.sh.tpl*
```
#!/bin/bash
yum update -y

%{ for x in names ~}
touch /${x};
echo ${f_name} > ${x};
%{ endfor ~}
```
Предыдущий файл создаст файлы с названиями Oleg, Petya, Masha с Denis внутри, то, что сделано через проценты - аналог цикла в Terraform. Таким образом мы смогли создать динамичные файлы

Однако для динамичных частей в самом `.tf` файле необходимо использовать функцию - предположим нам на сайте на вход необходимо открыть несколько портов:
```
dynamic "ingress" { # Задаем динамику, после запуска оно создаст новые блоки с заданным именем (ingress) и поставленными значениями
  for_each = ["80", "443", "8080"] # Значения для подстановки
 
  content { # модуль, который должен разместиться внутри нашего блока
    from_port   = ingress.value # Вытаскиваем значения с помощью названия динамической функции и атрибута value
    to_port     = ingress.value
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] 
  }
}
```
# Lifecycle
Для исключения лишних изменений существует блок `lifecycle` с множеством атрибутов:
```
lifecycle {
    prevent_destroy = true # Не даст удалить виртуалку
    ignore_changes = ["ami", "user_data"] # Игнорирует изменения в заданных атрибутах
    create_before_destroy = true # Сначало terraform создаст новую виртуалку, а затем уничтожит старую
}
```

Рекомендуется использовать:
```
resource "aws_eip" "my_static_eip" { # Elastick ip для сохранения статичного ip адресса даже при уничтожении виртуалки
     instance = aws_instance.web_server.id 
}
```

# Вывод в терминал
В Terraform есть вывод значений к нам в терминал:
```
output "web_server_id" {
  value = aws_instance.web_server.id # выводим id сервера
}

output "web_server_ip" {
  value = aws_eip.my_static_eip.public_ip # выводим ip адрес эластик ip
}
```

Вывод принято создавать в отдельном файле:

*outputs.tf*
```
output "web_server_id" {
  value = aws_instance.web_server.id # выводим id сервера
}

output "web_server_ip" {
  value = aws_eip.my_static_eip.public_ip # выводим ip адрес эластик ip
}
```

# Порядок создания
Для расставления приоритетов запуска (порядок создания) есть атрибут dependes_on:

*webserver.tf*
```
resource "aws_instance" "web_server" {
  ami                    = "ami-03a71cec707bfc3d7"
  instance_type          = "t3.micro"                                        
  user_data              = file("./user_data.sh")                           
  depends_on = [aws_instance.db] # Веб сервер создастся только после db, можно задавать несколько серверов при необходимости
}

resource "aws_instance" "db" {
  ami                    = "ami-03a71cec707bfc3d7" 
  instance_type          = "t3.micro"                                     
  user_data              = file("./user_data.sh")                           
}
```

# Модуль data
Существует в Terraform модуль `data`, который используется исключительно для вывода определенных параметров. У них нет атрибутов, а сам модуль импортирует информацию от провайдера:

```
data "aws_availability_zones" "working" {}

output "aws_zones"{
   value = data.aws_availability_zones.working.names # Выведет рабочие зоны для AWS, взятые у них самих
}
```

Можно и фильтровать:

```
data "aws_vpc" "blitz_vpc" {
    tags = {Name = "blitz"}
}

output "vpc_id"{
   value = "VPC id: ${data.aws_vpc.blitz_vpc.id}" # Выведет id всех VPC с тегом Name: blitz, для подстановки переменных в текст - ${}
}
```

Пример использования в работе:

```
resource "aws_instance" "web_server_on_ubuntu_22_04" {
  ami                    = data.aws_ami.ubuntu_22_04.id
  instance_type          = "t3.micro"
}

data "aws_ami" "ubuntu_22_04" {
  owners = ["099720109477"] # Кто создатель?
  most_recent = true # Найти самое последнее издание?
  filter { # отфильтруем по названию, то есть найдем ubuntu по имени
    name = "name"
    value = ["ubuntu/images/hvm-ssd/ubuntu-bionic-22.04-amd64-server-*"] # Найти последнюю версию ubuntu 22.04
  }
}
```


# Green/Blue deployment 
Green/Blue deployment - метод работы, когда для обновления на новую версию мы сначала создаем сервера с новой версией (Blue сервера), которые, после проверки на работоспособность, начинают потихоньку заменять Green сервера. После полной замены Green сервера себя уничтожают, а Blue сервера становятся Green серверами

Для использования G/B deployment в AWS через Terraform необходимо использовать:

```
resource "aws_default_vpc" "default" {}

resource "aws_security_group" "web_server_security_group" {
  name        = "web_server_security_group" 
  description = "Allow inbound traffic"
  vpc_id = aws_default_vpc.default.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress { 
    from_port   = 0 
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"] 
  }
}

resource "aws_launch_template" "web_server_template" {
    name = "web_server_template"
    image_id = data.aws_ami.ubuntu_22_04.id # Буду использовать модули из предыдущей программы чтобы не повторяться.
    instance_type = "t3.micro"
    vpc_security_group_ids = [aws_security_group.web_server_security_group.id]
    user_data = file("user_data.sh")
}

resource "aws_autoscaling_group" "web_server_a_g" { 
    name = "ASG-${aws_launch_template.web_server_template.latest_version}" 
    launch-configuration = aws_launch_configuration.web_server_launch_config.name
    min-size = 2
    max-size = 2
    min_elb_capacity = 2 
    vpc_zone_identifier = [aws_default_subnet.web_server_subnet1.id, aws_default_subnet.web_server_subnet2.id]
    health_check_type = "ELB" 
    load_balancers = [aws_elb.web_server_load_balancer]
    target_group_arns = [aws_lb_target_group.web_server_t_g.arn]

    launch_template {
        id = aws_launc_template.web_server_template.id
        version = aws_launc_template.web_server_template.latest_version
    }

    dynamic "tag" {
        for_each = {
            Name = "WebServerASG-${aws_launch_templdate.web_server_template.latest_version}"
            Owner = "Blitz"
        }
        content {
            key = tag.key
            value = tag.value
            propagate_at_launch = true
        }
    }

    lifecycle {
        create_before_destroy = true
    }
}

resource "aws_lb" "web_server_load_balancer" {
    name = "WebServer-HA-ALB"
    security_groups = [aws_security_group.web_server_security_group]
    load_balancer_type = "application"
    subnets = [aws_default_subnet.web_server_subnet1.id, aws_default_subnet.web_server_subnet2.id]
}

resource "aws_lb_target_group" "web_server_t_g" {
    name = "WebServer-T-G"
    vpc_id = 
    port = 80
    protocol = "HTTP"
    deregistration_delay = 10
}

resource "aws_lb_listener" "http" {
     load_balancer_arn = aws_lb.web_server_load_balancer.arn
     port = "80"
     protocol = "HTTP"
     default_action {
          type = "forward"
          target_group_arn = aws_lb_target_group.web_server_t_g.arn
     }
}

output "web_load_balancer_url" {
    value = aws_lb.web_server_load_balancer.dns_name
}
```

# Variables
Переменные через Variable нужны, чтобы люди не лазили в ваш main.tf:

*variables.tf*

```
variable "region" {} # Переменная без значения default - ввод из терминала (Terraform запросит её при запуске)

variable "region_1" {
    type = string # Добавляем тип переменной, можно не добавлять, тогда Terraform сам определит тип, но есть вы введете другой тип, ошибки не будет, поэтому лучше добавлять
    # Всего есть 5 type: number (число), string (строка), list (список), map (словарь), bool (true/false)
    default = "us-east-1" # Создали переменную со значением us-east-1
    description = "Us_region" # Описание для переменной
}

variable "default_tags" {
    type = map
    default = {
        Owner = "Blitz"
        Project = "Test_Terraform"
        Project_name = "Test_Server"
    }
    description = "Default_tags"
}

variable "instance_type"{
    type = string
    default = "t3.micro"
}
```

Для использования переменных нужно:

*main.tf*
```
terraform {
    provider "aws" {
        region = var.region_1
    }
}

resource "aws_instance" "test_server" {
  ami                    = "ami-03a71cec707bfc3d7"
  instance_type          = var.instance_type
  tags = merge(var.default_tags, {Name = "Server_test made by ${var.default_tags["Owner"]}"}) # Используем merge для соединения нескольких значений, например для объединения двух словарей
  # В нашем случае мы объединяем два словаря, причем в качества значения для Name второго словаря мы используем значение из первого словаря, обращаясь по ключу
}
```
## Изменение переменных
Для изменения переменных при запуске нужно указать тег -var:

```
terraform plan -var="region_1=eu-central-1" -var="instance_type=t2.micro"
terraform apply -var="region_1=eu-central-1" -var="instance_type=t2.micro"
```

Можно изменять переменные и по другому, уже известному нам способу (через глобальные переменные терминала):

```
export TF_VAR_region_1=us-west-2
```

*TF_VAR_ перед началом - обязательная часть для задачи переменных Terraform*
## Удаление переменных
Для удаления переменных:

```
unset TF_VAR_region_1
```

## Задание переменных через .tfvars
Есть ещё один способ задачи переменных, используя файл .tfvars, который возьмет terraform при запуске:

*terraform.tfvars*
```
region = "ca-cenral-1"
instance_type = "t2.micro"
```

Однако может возникнуть потребность в нескольких файлах:

*site1.tfvars*
```
region = "ca-cenral-1"
instance_type = "t2.micro"
```

*site2.tfvars*
```
region = "eu-central-1"
instance_type = "t2.small"
```

Для запуска:

```
terraform apply -var-file="site1.tfvars"
```

Для создания локальных переменных в файле нужно использовать locals:

*variables.tfvars*
```
env = "AWS"
name = "test_server_2332"
```

*main.tf*
```
data "aws_availability_zones" "available_zones" {}

locals {
    full_name = "${var.env}-${var.name}"
    a_z = join(";", data.aws_availability_zones.available_zones.names) # Соеденить все значения списка в 1 строку с разделением в виде ;
}

terraform {
    provider "aws" {
        region = var.region_1
    }
    default_tags {
        name = local.full_name
        zones = local.a_z
    }
}
```

# Provisioner
Provisioner используется для запуска команд на локальном хосте:

*(в качестве ресурса использую null_resource чтобы не выполнять никаких серверных взаимодействий)*

```
terraform {
    provider "aws" {
        region = var.region_1
    }
}

resource "null_resource" "command1" {
    provisioner "local-exec" {
        command = "echo Terraform Start: ${date} >> log.txt"
    }
}
```

Запустив такую команду Terraform создаст файл log.txt в текущей директории с Terraform Start: *сегодняшняя дата*
Таким образом можно запускать и ping, и другие команды.

Вот ещё примеры:

- Выведет `Hello world!` через питон:
```
resource "null_resource" "command2" {
    provisioner "local-exec" {
        interpreter = ["python", "-c"]
        command = "print('Hello world!')"
    }
}
```

- Сохранит имена из заданных в переменных в файл *names.txt*:

```
resource "null_resource" "command3" {
    provisioner "local-exec" {
        command = "echo $NAME1 $NAME2 >> names.txt"
    }
    environment = {
        NAME1 = "Vasya"
        NAME2 = "Petya"
    }
}
```


- Сохранит дату запуска сервера:

```
resource "aws_instance" "command4" {
  ami = "ami-03a71cec707bfc3d7"
  instance_type = "t3.micro"
    provisioner "local-exec" {
        command = "echo Server started in ${date} >> log.txt"
    }
}
```


# Создание паролей
Создание паролей происходит через специальные ресурсы:

*main.tf*
```
provider "aws" {
    region = "ca-central-1"
}

resource "random_string" "rds_password" { # Создаем случайный пароль
    length = 12 # Длина 12 символов
    special = true # Использовать спецсимволы?
    override_special = "!&#" # Какие спец символы можно использовать
    keepers = { # Сам по себе пароль будет постоянно один и тот же, для изменения нужно поменять какой-либо параметр. Для этого в ресурсе рандома есть атрибут keepers, где при смене переменной меняется пароль
        keeper1 = "Change_password"
    }
}

resource "aws_ssm_parameter" "rds_password" { # Сохраним пароль в AWS
    name = "/prod/mysql" 
    descriprion = "Password for mysql db"
    type = "SecureString" # Указываем тип: защищенная строка
    value =  random_string.rds_password.result # Указываем наш пароль
}

data "aws_ssm_parameter" "my_rds_password" { # Получаем пароль по имени
    name = "/prod/mysql"
    depends_on = [random_string.rds_password] # обязательно получаем пароль только после его создания
}

output "password" {
    value = data.aws_ssm_parameter.my_rds_password.value
}
```

# Lookups, validation и conditions
Lookups, validation и conditions в Terraform заменяют программные команды

## Lookups
Lookup работает как x in list:

```
variable "size" {
    default = {
        "prod" = "t3.medium"
        "dev" = "t3.micro"
        "staging" = "t3.small"
    }
}

instance_type = lookup(var.size, "dev") # подставим значение dev, то есть t3.micro
```

Можно совмещать c переменными:

```
instance_type = lookup(var.size, var.env)
```

## Validation
Validation позволяет выводить ошибки:
```
validation {
    condition = substr(var.region, 0, 3) == "eu-" # Если первые 3 символа var.region равны eu-, тогда все в норме
    error_message = "Must be EU region" # Иначе выведет ошибку
}
```

## Conditions
Conditions - сокращенный if:
```
variable "env" {
    default = "prod"
}

instance_type = (var.env == "prod" ? "t2.small":"t3.small") # есть переменная env является prod, то мы создадит маленький t2, иначе маленький t3 сервер
count  = (var.env == "prod" ? 1:0) # Если prod, то создать сервер, если нет, то не создавать, то есть создать 0 серверов
```

То есть: X = (CONDITION ? VALUE_IF_TRUE:VALUE_IF_FALSE)

# Count и For in
Сосредоточимся по подробнее на атрибутах `count` и `For in`. 

`Count` не просто число, а подсчет от нуля до числа, за счет этого можно производить перебор списков. Для вывода числа count как пересчет, нужно использовать его атрибут `.index`, однако его можно заменить [*] для перебора всех значений в списке

`For in` вообще не отличается от питона, как и проверка условия через `if`

Пример:

```
terraform {
provider "aws" {
    region = "ca-central-1"
}
}

variable "aws_users" {
     description = "List of AWS users"
     default = ["Vasya", "Petya", "Kolya", "Lena"]
}

resource "aws_iam_user" "users" { # Создаем пользователей
    count = length(var.aws_users) # создадим столько пользователей, сколько длина списка их имен 
    name = element(var.aws_users, count.index) # element берет елемент массива по его номеру, а count.index выводит индекс count, то есть его значение на определенный момент
}

resource "aws_instance" "servers" { # Создаст 3 сервера с нумерацией сервера в тегах 0, 1, 2
  ami = "ami-03a71cec707bfc3d7"
  instance_type = "t3.micro"
  count = 3
  tags = {Name = "Server number: ${count.index}"}
}

output "created_users" { # Выведет всю инфу по пользователям, созданными нами
     value = aws_iam_user.users
}

output "user_ids" {
    value = aws_iam_users.users[*].id # Выведем все id пользователей. Для указания всех пользователей по очереди используем * (аналог count.index)
}

output "user_custom" {
    value = [ # Перебираем список используя for x in как в питоне, так и работает For in в terraform
      for x in aws_iam_users.users:
      "Username: ${x.name} ARN: ${x.arn}"
    ]
}

output "user_custom_map" {
    value = {# Так же можно делать и со словарями
      for i in aws_iam_users.users:
      i.unique_id => i.id # Создаст словать с ключами i.unique_id, значениями i.id
    }
}

output "user_custom_map" {
    value = {# Так же можно делать и со словарями
      for i in aws_iam_users.users:
      i.unique_id => i.id # Создаст словать с ключами i.unique_id, значениями i.id
    }
}

output "user_custom_length" { # Выведет имена пользователей, если их имена состоят из 4 символов
    value = [
      for user in aws_iam_users.users:
      user.name 
      if length(user.name) == 4
    ]
}
```

# Remote state
Terraform remote state - очень важная функция в Terraform. Основной её смысл: хранить `.tfstate` в удаленном хранилище, это необходимо по понятным причинам: чтобы разные разработчики своими развертками не мешали друг другу; так же это делается для безопасности: там, как ни как, хранятся все данные, в том числе и пароли

Вот как это делается:

```
terraform {
provider "aws" {
    region = "ca-sentral-1"
}

backend "s3" { # Сохраняем .tfstate в AWS S3 bucket
    bucket = "your_aws_bucket" # Наш bucket
    key = "dev/network/terraform.tfstate" # Куда в bucket сохранится .tfstate
    region = "us-east-1" # Регион нашего bucket
}

}

variable "vpc_cidr" {
    default = "10.0.0.0/16"
}

resource "aws_vpc" "main" { # Создаем самую простую сеть в AWS для примера
    cidr_block = var.vpc_cidr
    tags = {
        Name = "My VPC"
    }
}

resource "aws_internet_gateway" "main" {
    vpc_id = aws_vpc.main.id
}

output "vpc_id" { # Используем output для общения между разными Terraform файлами (то есть аналог работы в команде)
    value = aws_vpc.main.id
}

output "vpc_cidr" { # Используем output для общения между разными Terraform файлами (то есть аналог работы в команде)
    value = aws_vpc.main.cidr_block
}
```

*Через output можно работать в команде, для этого нужно использовать общий `.tfstate`*

Для получения данных из `output` нужно использовать `data`:

```
data "terraform_remote_state" "network" { # Предположим что мы другой разраб в одной компании, мы работаем с другим разрабом, который сделал для нас сеть выше
    backend = "s3" # Указываем имя backend
    config = { # Указываем конфигурацию backend, такую же, как и на сохранении данных
        bucket = "your_aws_bucket"
        key = "dev/network/terraform.tfstate"
        region = "us-east-1"
    }
}

output "all_what_we_see" {
     value = data.terraform_remote_state.network
}

output "get_vpc_id" {
     value = data.terraform_remote_state.network.outputs.vpc_id # Выводим vpc_id, переданный в output другим разработчиком
}

output "get_vpc_cidr" {
     value = data.terraform_remote_state.network.outputs.vpc_cidr # Выводим vpc_cidr, переданный в output другим разработчиком
}

resource "aws_security_group" "new_security_group" { 
  name        = "new_security_group"
  vpc_id = data.terraform_remote_state.network.outputs.vpc_id                

  ingress {   
    to_port     = 80 
    protocol    = "tcp"
    cidr_blocks = [data.terraform_remote_state.network.outputs.vpc_cidr] 
  }


  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"] 
  }
}

backend "s3" { # Сохраняем .tfstate в AWS S3 bucket
    bucket = "your_aws_bucket" # Наш bucket
    key = "dev/security_group/terraform.tfstate" # Сохраняем уже в другую папку, чтобы не сломать .tfstate другого разработчика
    region = "us-east-1" # Регион нашего bucket
}

output "security_group_id" {
    value = aws_security_group.new_security_group.id # Сохраняем id security group в output для других разработчиков 
}
```

Благодаря связям Terraform_remote_state можно связывать работы разных разработчиков, запросто работать в команде, держать данные в сохранности и многое другое, самое главное: иметь bucket в облаке провайдера. Так же лучше сохранять в облаке файлы `.terraform.lock.hcl`, они являются хэшами провайдеров.

# Модули
Модули похожи на функции: они принимают значения и возвращают что-то. Для создания модуля требуется создать отдельную папку, в которую мы закидываем наши .`tf` файлы, формирующие проект, затем необходимо удалить провайдера. Таким образом мы создали модуль, который можно использовать

Модули в Terraform - основа всей работы. Так мы можем выполнять сложный код, развертка поочередно которого заняла бы много времени.

Модули можно хранить где угодно, хоть на гитхабе, главное правильно указывать путь, откуда мы его импортируем:
```
source = "<ссылка на страницу в гитхабе>//<название папки с модулем>"
```

## Интеграция модулей
Для использования модуля необходимо создать новый `.tf` файл вне директории с модулем, указав внутри:

```
provider "aws" { # Указываем провайдера
   region = "eu-north-1"
}

module "vpc-dev" { # хотим использовать существующий модуль, который есть у нас на компе
   source = "/terraform/module" # указываем путь до модуля (можно использовать хоть ссылки гитхаба)
   for_each = { # Пример использования цикла для добавления нескольких модулей с разными переменными
      Name = Oleg
      Second_name = Petya
   }
   part = each.key # Ключ
   name = each.value # Значение
}

module "vpc-prod" { # импортируем ещё один модуль
   source = "/terraform/module" # указываем путь до модуля
   vpc_cidr = "10.100.0.0/16" # заменяем переменную vpc_cidr
   env = "development" # заменяем переменную env
   public_subnet_cidrs = ["10.100.1.0/24", 10.100.2.0/24"] # заменяем переменную public_subnet_cidrs
}

output "public_subnet_ids_from_module" {
   value = module.vpc-prod.public_subnet_ids # Выводим output, прописанный в output в нашем модуле
}
```

## Разные провайдеры в модуле
Для использования разных провайдеров в модуле необходимо при создании модуля указать атрибут `providers`:

```
module "vpc-prod" {
   source = "/terraform/module" 
   providers = { # прокидываем провайдеров в модуль по принципу: провайдер_внутри_модуля = наш_провайдер_в_файле_где_используем_модуль
      aws = aws
      aws.USA = aws.USA
      aws.Canada = aws.Canada
   }
}
```

А в самом модуле нужно указать блок:

```
terraform {
   required_providers {
      aws = {
         source = "hashicorp/aws"
         configuration_aliases = [aws.USA, aws.Canada]
      }
   }
}
```

# Хранение глобальных переменных
Для хранения глобальных переменных можно использовать уже известную нам технологию terraform_remote_state.
Для этого необходимо указать наши глобальные переменные в `output`, а затем сохранить их в `bucket`:

```
terraform {
   provider "aws" {
      region = "ca-central-1"
   }

   backend "s3" { 
       bucket = "your_aws_bucket" # Наш bucket
       key = "dev/global_vars/terraform.tfstate" # Сохраняем в папку с глобальными переменными
       region = "us-east-1"
   }
}

output "global_name" {
   value = "Blitz"
}

output "global_tags" {
   value = {
      Project_version = "1"
      Owner = "Blitz"
   }
}
```

Рекомендую использовать глобальные перменные, переводя их в локальные для того, чтобы не тратить лишнее время на прописывание пути до переменной:

```
data "terraform_remote_state" "global" {
   backend = "s3"
   config {
      bucket = "your_aws_bucket"
      key = "dev/global_vars/terraform.tfstate" 
      region = "us-east-1"
   }
}

locals {
   my_name = data.terraform_remote_state.global.outputs.global_name
   tags = data.terraform_remote_state.global.outputs.global_tags
}
```

# Встроенный Refactoring
В Terraform есть Refactoring кода, однако он является очень опасной функцией, которую необходимо несколько раз перепроверить

Refactoring необходим для создания более лаконичного и понятного кода из уже существующего, например: нам необходимо из старой большой папки выделить часть на prod и часть на dev, которые не связаны друг с другом, но в лежат в одной папке и созданы с одним `.tfstate`

Завязан рефакторинг на командах terraform state:

1. Смотрим что у нас есть:
```
terraform state list 
```

2. Смотрим state файл изнутри:
```
terraform state pull
```

3. Удаляем из всего state файла resource с названием *aws_eip.prod-ip1* и сохраняем его в локальный terraform.tfstate с названием *prod-ip1*:
```
terraform state mv -state-out="terraform.tfstate" aws_eip.prod-ip1 prod-ip1
```

4. Сохраним `.tfstate` на наш AWS S3 bucket (главное помнить, что это делается из другого main файл в другой папке с прописанным backend):

```
terraform init
```

5. Обновляем наш созданный `terraform.tfstate` с удаленным после редактирования (если мы редактировали его):

```
terraform state push terraform.tfstate
```

# Terraform Workspaces
Terraform Workspaces создан для тестирования нашей системы, которую мы хотим поднять, точнее он позволяет тестировать в другом рабочем месте, которое изолированно от нашего. Таким образом мы можем запросто тестировать наши файлы, не задев уже существующую систему

Все не `main` Workspace будут добавлять приписку с названием Workspace в сохраняемых файлах в `bucket`

Все команды:

- `terraform workspace show` - покажет в каком мы Workspace

- `terraform workspace new name` - создать новый Workspace с названием name

- `terraform workspace list` - показывает список всех Workspace

- `terraform workspace select name` - выбрать конкретный Workspace name

- `terraform workspace delete name` - удалить конкретный Workspace name
# Terraform Cloude
Terraform Cloud необходим для запуска команд на удаленном компьютере, то есть запускать свои `apply` на удаленном облачном компе. Это очень полезно для БОЛЬШИХ команд облачных инженеров, так как позволяет хранить `.tfstate` сразу внутри облака, не сохраняя куда-либо (по факту оно уже лежит на облаке)

Ещё для групп полезно тем, что можно отправлять `plan` руководителю, который будет решать: делать `apply` или нет, так же можно добавлять комментарии к чужим запускам

Можно прописать различные правила для поднятия инфраструктуры (эта функция платная), а ещё она позволяет показать стоимость поднимаемой инфраструктуры

Terraform Cloude видит все изменения, даже те, что были созданы в связанном с проектом github репозиторие

*Terraform Cloud Enterprise - платная версия Terraform Cloud с возможностью использовать Terraform Cloud в вашей локальной сети*
# Import To
Помимо команды импорт в Terraform так же есть и следующий блок кода:
```
import {
   id = "id_of_security_group" # импортируем по id в AWS
   to = aws_security_group.web_server # Куда импортируем
}
```
Далее запускаем с созданием их конфига:
```
terraform init
terraform plan -generate-config-out=generated-sg.tf
terraform apply
```

# Пример веб сервера
Для запуска веб сервера создадим виртуалку на AWS и обновим её, также пропишем политику доступа:

```
terraform {
  provider "aws" {}
}

resource "aws_instance" "web_server" {
  ami                    = "ami-03a71cec707bfc3d7" # Amazon Linux AMI
  instance_type          = "t3.micro" # Micro size
  vpc_security_group_ids = [aws.security_group.web_server_security_group.id] # Указываем id нашего дальнейшего security group
  user_data              = <<EOF
#!/bin/bash
yum update -y
EOF 
# После создания будет исполняться код между EOF, таким образом мы сразу обновим систему на всех машинах
}

resource "aws_security_group" "web_server_security_group" { # Создаем политику доступа, скопировано с http://man.hubwiz.com/docset/Terraform.docset/Contents/Resources/Documents/docs/providers/aws/r/security_group.html
  name        = "web_server_security_group"                 # Указываем имя
  description = "Allow inbound traffic"                     # Указываем описание

  ingress {          # На входящий поток
    from_port   = 80 # Доступно с чужого порта 80
    to_port     = 80 # На наш порт 80
    protocol    = "tcp"         # Протокол tcp
    cidr_blocks = ["0.0.0.0/0"] # Доступ со всех ip адрессов
  }

  egress { # На исходящий поток (на отправку данных с виртуалки)
    from_port   = 0 # Собственно мы ничего не отправляем, так что ничего не указываем
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"] 
  }
}
```

Однако shell скрипты лучше выносить в отдельные файлы, тогда код будет выглядеть так:

*user_data.sh*
```
#!/bin/bash
yum update -y
```
*web_site.tf*
```
terraform {
  provider "aws" {}
}


resource "aws_instance" "web_server" {
  ami                    = "ami-03a71cec707bfc3d7"                           # Amazon Linux AMI
  instance_type          = "t3.micro"                                        # Micro size
  vpc_security_group_ids = [aws.security_group.web_server_security_group.id] # Указываем id нашего дальнейшего security group (да, код в Terraform работает асинхронно, то есть создает все параллельно)
  user_data              = file("./user_data.sh")                            # Подставляем код для оболочки из user_data.sh
}

resource "aws_security_group" "web_server_security_group" { # Создаем политику доступа, скопировано с http://man.hubwiz.com/docset/Terraform.docset/Contents/Resources/Documents/docs/providers/aws/r/security_group.html
  name        = "web_server_security_group"                 # Указываем имя
  description = "Allow inbound traffic"                     # Указываем описание

  ingress {          # На входящий поток, чтобы открыть ещё один порт, нужно просто скопировать ingress и вставить ещё раз
    from_port   = 80 # Доступно с чужого порта 80
    to_port     = 80 # На наш порт 80
    protocol    = "tcp"         # Протокол tcp
    cidr_blocks = ["0.0.0.0/0"] # Доступ со всех ip адрессов
  }

  ingress {          # Открываем и 8080 порт
    from_port   = 8080 
    to_port     = 8080 
    protocol    = "tcp"         
    cidr_blocks = ["0.0.0.0/0"] 
  }

  egress { # На исходящий поток (на отправку данных с виртуалки)
    from_port   = 0 # Собственно мы ничего не отправляем, так что ничего не указываем
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"] 
  }
}
```

# Пример Web сервера с Zero DownTime с провайдером AWS

*main.tf*
```
terraform{
   provider "aws" {
      default_tags { # Применяется ко всем созданным блокам
         tags = {
            version = v1
         }
      }
   }
}

data "aws_ami" "ubuntu_22_04" {
  most_recent = true
  filter {
    name = "name"
    value = ["ubuntu/images/hvm-ssd/ubuntu-bionic-22.04-amd64-server-*"] # Найти последнюю версию ubuntu 22.04 
  }
}

data "aws_availability_zones" "available" {}

resource "aws_security_group" "web_server_security_group" {
  name        = "web_server_security_group" 
  description = "Allow inbound traffic"

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress { 
    from_port   = 0 
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"] 
  }
}

resource "aws_launch_configuration" "web_server_launch_config" {
    name-prefix = "WebServer-Highly-Available-LC-" # Название не будет постоянным, то есть AWS сам добавит случайные числа
    image-id = data.aws_ami.ubuntu_22_04.id
    instance_type = "t3.micro"
    security_groups = [aws_security_group.web_server_security_group]
    user_data = file("user_data.sh") # Что угодно можно написать в shell файле при необходимости
    
    lifecycle {
        create_before_destroy = true
    }
}

resource "aws_autoscaling_group" "web_server_a_g" { # Автораспределение нагрузки между серверами балансировщика (то есть load_balancer их создает, а autoscaling распределяет нагрузку)
    name = "ASG-${aws_launch_configuration.web_server_launch_config.name}" # Название будет зависеть от названия сервера
    launch-configuration = aws_launch_configuration.web_server_launch_config.name
    min-size = 2 # Количество серверов мин
    max-size = 2 # Количество серверов макс
    min_elb_capacity = 2 # Сколько раз сервера должны пройти health check чтобы подствердить, что сервер жив, цел, работает
    vpc_zone_identifier = [aws_default_subnet.web_server_subnet1.id, aws_default_subnet.web_server_subnet2.id]
    health_check_type = "ELB" # Проверка состояния методом ping (ELB)
    load_balancers = [aws_elb.web_server_load_balancer]

    dynamic "tag" { # Сделаем динамические теги
        for_each = {
            Name = "WebServer"
            Owner = "Blitz"
        }
        content {
            key = tag.key
            value = tag.value
            propagate_at_launch = true
        }
    }

    lifecycle {
        create_before_destroy = true
    }
}

resource "aws_elb" "web_server_load_balancer" { # Балансировщик нагрузки, серверная часть
    name = "WebServer-HA-ELB"
    availability_zones = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1]]
    security_groups = [aws_security_group.web_server_security_group]
    listener { # Что будет слушать load_balancer
        lb_port = 80
        lb_protocol = "http"
        instance_port = 80
        instance_protoc = "http"
    }
    health_check { # Отправляем запросы на проверку сервера: жив ли он?
        healthy_theshold = 2
        unhealthy_theshold = 2
        timeout = 2
        target = "HTTP:80/"
        interval = 10
    }

    tags = {
        name = "web_server_load_balancer"
    }
}

resource "aws_default_subnet" "web_server_subnet1" { # В каких сетях создавать новые виртуалки балансировщика нагрузки
    availability_zone = data.aws_availability_zones.available.names[0] 
}

resource "aws_default_subnet" "web_server_subnet2" {
    availability_zone = data.aws_availability_zones.available.names[1]
}

output "web_load_balancer_url" { # Ссылка на наш сервер
    value = aws_elb.web_server_load_balancer.dns_name
}
```

# Пример развертки реальной сети
Примерно так должна выглядеть работа опытного Terraform специалиста:

*outputs.tf*
```
output "vpc_id" {
   value = aws_vpc.main.id
}

output "public_subnet_id" {
   value = aws_subnet.public_subnet.id
}

output "private_subnet_id" {
   value = aws_subnet.private_subnet.id
}

output "aws_eip_id" {
   value = aws_eip.nat_eip.id
}

output "aws_nat_gateway_id" {
   value = aws_nat_gateway.nat.id
}

output "aws_route_table_private_id" {
   value = aws_route_table.private.id
}

output "aws_route_table_public_id" {
   value = aws_route_table.public.id
}

output "public_internet_gateway_id" {
   value = aws_route.public_internet_gateway.id
}

output "private_internet_gateway_id" {
   value = aws_route.private_internet_gateway.id
}

output "route_table_association_public_id" {
   value = aws_route_table_association.public.id
}

output "route_table_association_private_id" {
   value = aws_route_table_association.private.id
}

output "instance_my_new_web_server_id" {
    value = aws_instance.my_new_web_server.id
}

output "security_group_web_server_sec_group_id" {
    value = aws_security_group.web_server_sec_group.id
}
```

*terraform.tf*
```
terraform {
    provider "aws" {
        region = "${var.aws_region}"
    }

    backend "s3" {
    bucket = "mybucket"
    key    = "path/to/my/key"
    region = "us-east-1"
    }
}

locals {
  availability_zones = ["${var.aws_region}a", "${var.aws_region}b"]
}

resource "aws_vpc" "main" {
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = true
  enable_dns_support   = true
}

resource "aws_subnet" "public_subnet" {
  vpc_id                  = aws_vpc.main.id
  count                   = length(var.public_subnets_cidr)
  cidr_block              = element(var.public_subnets_cidr, count.index)
  availability_zone       = element(local.availability_zones, count.index)
  map_public_ip_on_launch = true

}

resource "aws_subnet" "private_subnet" {
  vpc_id                  = aws_vpc.main.id
  count                   = length(var.private_subnets_cidr)
  cidr_block              = element(var.private_subnets_cidr, count.index)
  availability_zone       = element(local.availability_zones, count.index)
  map_public_ip_on_launch = false

}

#Internet gateway
resource "aws_internet_gateway" "ig" {
  vpc_id = aws_vpc.main.id
}

# Elastic-IP (eip) for NAT
resource "aws_eip" "nat_eip" {
  vpc        = true
  depends_on = [aws_internet_gateway.ig]
}

# NAT Gateway
resource "aws_nat_gateway" "nat" {
  allocation_id = aws_eip.nat_eip.id
  subnet_id     = element(aws_subnet.public_subnet.*.id, 0)
}

# Routing tables to route traffic for Private Subnet
resource "aws_route_table" "private" {
  vpc_id = aws_vpc.main.id
}

# Routing tables to route traffic for Public Subnet
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id
}

resource "aws_route" "public_internet_gateway" {
  route_table_id         = aws_route_table.public.id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.ig.id
}

# Route for NAT Gateway
resource "aws_route" "private_internet_gateway" {
  route_table_id         = aws_route_table.private.id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_nat_gateway.nat.id
}

# Route table associations for both Public & Private Subnets
resource "aws_route_table_association" "public" {
  count          = length(var.public_subnets_cidr)
  subnet_id      = element(aws_subnet.public_subnet.*.id, count.index)
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "private" {
  count          = length(var.private_subnets_cidr)
  subnet_id      = element(aws_subnet.private_subnet.*.id, count.index)
  route_table_id = aws_route_table.private.id
}

resource "aws_instance" "my_new_web_server" {
    ami = "ami-identificator"
    instance_type = "t3.micro"
    security_group_ids = [aws_security_group.web_server_sec_group.id]
}

resource "aws_security_group" "web_server_sec_group" {
   name        = "Web_Server_sec_group"
   vpc_id      = "${aws_vpc.main.id}"
   subnet_id   = aws_subnet.private_subnet.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = var.private_subnets_cidr
  }

  egress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    cidr_blocks     = ["0.0.0.0/0"]
  }
}
```

*variables.tf*
```
variable "aws_region" {
    default = "us-central-1"
}

variable "vpc_cidr" {
    default = "10.0.0.0/16"
}

variable "public_subnets_cidr" {
   default = ["10.0.1.0/24", "10.0.2.0/24"]
}

variable "private_subnets_cidr" {
   default = ["10.0.3.0/24", "10.0.4.0/24"]
}
```

В этой работе разворачивается две сети: private и public
Публичная сеть доступна из интернета, приватная же доступна только через NAT, связанный с публичной сетью, сайт же доступен только из приватной сети, то есть находится в сети, попасть в которую можно только через NAT, усиливая защиту (да, в instance не указан user_data, потому что  её необходимо будет заполнить по заданию, так же там не заданы ami)

# [Источник](https://www.youtube.com/playlist?list=PLg5SS_4L6LYujWDTYb-Zbofdl44Jxb2l8)
В качестве источника представлен плейлист на ютубе
Помимо этого источника также можно ознакомиться с [основами Terraform](https://habr.com/ru/companies/otus/articles/696694/) на хабре
Или ознакомиться с [полезными ссылками](https://github.com/adv4000/terraform-lessons/tree/master/Lesson-25) в репозитории на гитхабе