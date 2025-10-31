#Web #DevOps #IaC #Инструменты 

Ansible это система автоматизации, которая позволяет удаленно управлять и настраивать сервера, скачивать ПО и развертывать приложения
# Представление кода 
Главное: весь код в Ansible можно представить в виде дерева:
├── production                # инвентарный файл для продакшн-серверов
├── stage                     # инвентарный файл для stage-окружения
│
├── group_vars/
│   ├── group1                # здесь назначаются переменные для
│   └── group2                # конкретных групп
├── host_vars/
│   ├── hostname1             # специфические переменные для хостов в
│   └── hostname2             # случае необходимости прописываются здесь
│
├── site.yml                  # основной сценарий
├── webservers.yml            # сценарий для веб-сервера
└─── dbservers.yml            # сценарий для сервера базы данных

# Управление системой
Для управления системой Ansible использует ssh протокол и ip адреса. Для их получения необходимо указывать `inventory` файл (его иногда называют `hosts`)

# Команды
```
ansible all -m copy -a "src=1 dest=/root mode=600" -b #копируем файл
```
-m - модуль (Все модули смотри [тут](https://docs.ansible.com/ansible/latest/index.html)) (Кстати, свои модули можно написать используя api Ansible и Python, вот [пример](https://habr.com/ru/articles/553908/))
-a - параметры
-b - become банальное sudo (я не помню как там пароль вводить с -b)
src - источник (на основном сервере ansible)
dest - путь куда копируем
mode - с каким доступом

```
ansible all -m shell -a "apt-get update" -b #выполняем команду (можно так же через модуль command) 
```

# Playbook
Ansible-playbook это набор сценариев, который выполняет Ansible на заданных устройствах. Для описания сценариев Ansible использует yml

Изначально Ansible-playbook будет запускать `inventory` файл файл из своей директории, которая указана у него по умолчанию, чтобы это изменить необходимо указывать через `-i` свой файл:
```
ansible-playbook -i ./hosts playbook.yml
```

Чтобы указать пароль к sudo нужно использовать `--ask-become-pass`, пароль Ansible спросит в терминале

Затем в `inventory`, указываем ip адреса (не забудь разблокировать пользователей, например, `root` изначально заблокирован в openssh, вот [решение](https://www.dmosk.ru/miniinstruktions.php?mini=ubuntu-ssh-root) как это исправить):

```
client01 ansible_host=192.168.1.202 ansible_user=test ansible_password=ptest       
```                                                                     
или
```
192.168.1.37 ansible_user=root ansible_ssh_private_key_file=/root/.ssh/id_rsa
```
так же можно разделять их на группы:

```
[test]
client01 ansible_host=192.168.1.202 ansible_user=root ansible_password=root

[group2]
192.168.1.37 ansible_user=root ansible_ssh_private_key_file=/root/.ssh/id_rsa


[all_groups:children] #Объединение нескольких групп в 1 (как детей)
test
group2     
```

## Групповые параметры
Групповые параметры, применяемые определенными группами, можно выносить. Для этого нужно создать папку `group_var` в которой в файлах с именами групп через ключ: значение указываем их, например указать пользователя root:
```
ansible_user: root
```
Они автоматически подставятся в соответствующие поля при запуске playbook

## Переменные
Чтобы в playbook указывать переменные, необходимо сначала их задать через

```
vars:
  название переменной: значение
```

А затем применить его как:

куда-то надо вставить переменную: "`{{название переменной}}`" (кстати, переменные можно вытаскивать и из inventory, например "`{{ansible_user}}`")

Так же через переменную "`{{item}}`" и `loop`: после можно делать циклы, например:

```
- name: Loops
  hosts: all
  become: yes

  tasks:
    - name: Create Folder
      file:
        path: /root/"{{item}}"
        state: directory
      loop:
        - dir1
        - dir2
```

Или даже можно использовать словари как тут:

```
- name: Items
  hosts: all
  become: yes

  tasks:
    - name: Create groups
      group:
          name: "{{item}}"
          state: present
      loop:
        - dev
        - test1

    - name: Create users
      user:
        name: "{{item.clientname}}"
        shell: /bin/bash
        groups: dev,test1
        append: yes
        home: "/home/{{item.homedir}}"
      with_items:
        - {clientname: client1, homedir: client1}
        - {clientname: client2, homedir: client2}
```

## Модули

### Set_fact
С помощью модуля `set_fact` так же можно создавать переменную:
```
- set_fact:
    message: "Ок"
```
### Register
C помощью подмодуля `register` можно записать вывод команды в переменную:
```
- shell: id client1
register: client_groups
```

### Debug
C помощью модуля `debug` можно логировать любую, необходимую для нас, инфу

## When
С помощью блоков when можно создавать условия выполнения модулей:

```
-tasks:
  - block
    - name: Создание чего-то там
      ...
  when: ansible_hostname == "ubuntu" #Для наглядного примера
```

## Jinja2
С помощью `jinja2 (j2)` и `template` можно связывать обычные файлы, позволяя использовать в них переменные ansible, главное чтобы он был с расширение j2 (переменные указываются через `{{переменная}}`)

```
template:
         src: ./file01.j2 (заранее созданный файл, в котором прописаны некоторые переменные, включая переменную, взятую из inventory, а именно ansible_hostname)
         dest: /home
         mode: 0600
```


# Ansible galaxy
Ansible galaxy - сайт, являющимся эдаким *Docker Hub*, откуда можно скачивать чужие роли (готовые сценарии) чтобы затем применять их у себя. Роли - упрощение работы со сложными многоструктурными разворотами
## Создание роли
Чтобы создать роли нужно:
```
ansible-galaxy init название_роли
```
Создастся папка с кучей папок внутри, с помощью которых можно связать все, что лежит в папке между собой.
## Запуск роли
Чтобы запустить роли нужно:

```
- name: Online Service
  hosts: all
  gather_facts: yes #получать настройки (переменные) хоста
  roles: 
          -Web (Например развертка веб сервиса) #Изначально прописанная и готовая роль с именем папки Web
          -Pay (Например развертка системы оплаты) #Изначально прописанная и готовая роль с именем папки Pay
```

# [Источник](https://www.youtube.com/watch?v=YYjCwLs-1hA)
В качестве источника представлен длинный видос на ютубе

Ещё можно глянуть [статью](https://habr.com/ru/companies/selectel/articles/196620/) на хабре, там больше полезной информации

