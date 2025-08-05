#Программирование #Разметка
#Web

CSS не обычный язык программирования, как минимум потому, что он не является языком.
Основная его функция: добавление дополнительных элементов оформления на странице
# Принцип интеграции
Интеграция происходит через выделение объекта в виде тега, указанного в основном файле html и добавления его свойств, например:

*site.html*
```
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <link rel="stylesheet" href="style.css">
        <title>MyFirstSite</title>
    </head>
    <body>
        <p>Alert</p>
    </body>
</html>
```

*style.css*
```
/* Делаем объект p красным (да, это комментарий в css файле) */
p {color: #ff0000;}
```
Код выше делает текст Alert красным (#ff0000 - красный)

# Свойства в CSS
У css есть много различных свойств (помимо цвета):

*style.css*
```
p {
    color: #ff0000; /* красный */
    font-size: 24px; /* размер 24 пикселя */
    font-family: 'Times New Roman', Times, serif; /* шрифт */
} 
```

# Взаимодействие с классами html
Обращение в CSS к классам происходит через ., такое обращение нужно для группировки множества объектов:

```
#page { /* Весь div, обращаемся к нему через имя с решеткой в начале, так нужно обращаться ко всем через id */
    width: 960px;
    margin: 0 auto;
}

.strong-red { /* класс строгий красный, обращаемся через . к нему */
    font-weight: bolder;
    color: red;
}

.light-blue { /* класс голубой, также обращаемся через . к нему */
    font-weight: lighter;
    color: blue;
}
```
## Определение более важного класса
В CSS действует принцип: кто последний, тот и прав, то есть все, что было написано ниже будет важнее чем то, что было выше, однако наследование будет происходить всегда, например:

```
.light-blue {
    font-weight: lighter;
    color: blue;
}

.light-blue {
    font-size: 24px; /* размер теста */
}
```

Вроде два класса, нижний не преопределяет ничего из верхнего, поэтому конечный результат будет как:

```
.light-blue {
    font-weight: lighter;
    color: blue;
    font-size: 24px; /* размер теста */
}
```

Но если мы переопределим элемент:

```
.light-blue {
    font-weight: lighter; /* стиль теста (тонкий, жирный и тд) */
    color: blue;
}

.light-blue {
    color: red;
    font-size: 24px; /* размер теста */
}
```

То результат будет:
```
.light-blue {
    font-weight: lighter; /* стиль теста (тонкий, жирный и тд) */
    color: red; /* Тот, что был написан последним, применится */
    font-size: 24px;/* размер теста */
}
```

Тоже самое происходит и в style контейнере в html файлах:
```
<style>
	.light-blue {
		color: red; 
	}
</style>
<link rel="stylesheet" href="style.css">
```

Нижняя ссылка на css файл будет главнее, поэтому цвет поменяется на тот, который был прописан в файле, однако если он не изменялся, то останется таким же, как и тот, что прописан в style

# Background
Для редактирования заднего фона нужно использовать body:

```
body{
    background-color: #333333; /* цвет */
    background-image: url(CSS_background_images/Tasks.png); /* картинка для заднего фона */
    background-position: 80% 50%; /* местонахождение картинки на фоне (первое - X, второе - Y, указан в процентах, можно и в пикселях) */
    background-repeat: no-repeat; /* будет ли повторяться картинка заднего фона */
    background-attachment: scroll; /* картинка на заднем фоне будет прокручиваться со всей страницей, можно задать fixed, тогда она будет статичной */
}
```
Так же можно всё накидать просто в один background:
```
body {
	background: #333333 url(CSS_background_images/Tasks.png) no-repeat 80% 50% scroll
}
```

# Редактирование текста
Для редактирования текста в CSS существует множество свойств:

```
.text {
    margin-top: 100px; /* отступ сверху */
    margin-left: 5px; /* отступ слева */
    width: 600px; /* ширина теста */
    text-transform: uppercase; /* сделать все буквы в тексте ЗАГЛАВНЫМИ */
    text-decoration: none; /* добавить ли декорацию теста (перечеркивание, подчеркивание, надчеркивание) */
    word-spacing: 10px; /* расстояние между словами в тесте */
    line-height: 30px; /* расстояние между строчками теста */
    font-weight: lighter; /* стиль теста (тонкий, жирный и тд) */
    letter-spacing: 5px; /* расстояние между буквами */
    text-indent: 30px; /* красная строчка */
    text-align: right; /* выравнивание теста (можно по центру, слева, справа, по ширине) */
    font-family:'Lucida Sans', sans-serif; /* шрифт, sans-serif - шрифт, который будет, если не загрузится другие, установлен во всех браузерах*/
    color: rgb(55, 68, 255); /* цвет теста */
    font-size: 24px; /* размер теста */
    font-style: oblique; /* стиль шрифта */
}
```

# Кастомизация списка
В списке можно настроить немало параметров, даже добавить собственное изображение для иконок списка:

## Пример 1
*site.html*
```
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <link rel="stylesheet" href="style.css">
        <title>MyFirstSite</title>
    </head>
    <body>
        <p>
            <ul id="styles">
                <li>Пункт 1</li>
                <li>Пункт 2</li>
                <li>Пункт 3</li>
                <li>Пункт 4</li>
                <li>Пункт 5</li>
            </ul>
        </p>
    </body>
</html>
```

*style.css*
```
#styles li {
    color: blueviolet; 
    font-family: "Ubuntu", sans-serif;
    list-style-type: upper-alpha; /* Стиль иконок списка (в этом случае - заглавные буквы алфавита) */
    line-height: 30px;
    padding: 5px; /* Ширина центральной части, где пишется текст (не путать с шириной текста), по другому поле с текстом */
    list-style-position: inside; /* Местонахождение символов списка (внутри поля с текстом, вне поля с тестом) */
    border: 1px solid red; /* Толщина и цвет обводки поля с текстом (solid делает обводку видимой) */
    
    list-style-image: url(CSS_background_images/Ionic-Ionicons-Shield.16.png); /* Используем свою иконку списка (заменяет list-style-type) */
}

#styles li:hover {
    background: bisque; /* обозначаем цвет заливки при наводке на поле списка */
}
```

и опять можно вынести все в один стиль списка:

*style.css*
```
#styles li {
    color: blueviolet;
    font-family: "Ubuntu", sans-serif;
    list-style: upper-alpha url(CSS_background_images/Ionic-Ionicons-Shield.16.png) inside;
    line-height: 30px;
    padding: 5px;
    border: 1px solid red;

}

#styles li:hover {
    background: bisque;
}
```

## Пример 2
Ещё пример работы со списком и ссылками:

*site.html*
```
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <link rel="stylesheet" href="style.css">
        <title>MyFirstSite</title>
    </head>
    <body>
        <div id="navigation">
            <ul id="styles">
                <li><a href="https://meyerweb.com/eric/tools/css/reset/">Пункт 1</a></li>
                <li>
		    <a href="https://webcademy.ru/blog/739/">Пункт 2</a>
		    <ul class="dropdown">
                        <li><a href="#">Ссылка 1</a></li>
                        <li><a href="#">Ссылка 2</a></li>
                        <li><a href="#">Ссылка 3</a></li>
                        <li><a href="#">Ссылка 4</a></li>
                    </ul>
		</li>
                <li><a href="https://habr.com/ru/companies/otus/articles/580442/">Пункт 3</a></li>
                <li><a href="https://snipp.ru/html-css/css-reset">Пункт 4</a></li>
                <li><a href="https://www.youtube.com/">Пункт 5</a></li>
            </ul>
        </div>
    </body>
</html>
```

*style.css*
```
#styles li {
    color: rgb(255, 255, 255);
    font-family: "Ubuntu", sans-serif;
    list-style-type: upper-alpha;
    line-height: 30px;
    padding: 10px;
    list-style-position: inside;
    
    list-style-image: url(CSS_background_images/Ionic-Ionicons-Shield.16.png);
}

#styles li:hover {
    background: rgb(128, 128, 128);
}

#navigation {
    height: 45px;
}

#navigation ul li {
    float: left; /* Выравниваем список на лево (делаем горизонтальным) */
    position: relative; /* Способ позиционирования (в этом случае устанавливается относительно его исходного места) */
}

#navigation ul li a {
    color: white;
    padding: 5px;
    font-family: "Ubuntu", sans-serif;
    text-decoration: underline;
    line-height: 30px;
}

#navigation ul li a:hover {
    background: rgb(117, 117, 117);
}

#navigation .dropdown {
    position: absolute; /* Способ позиционирования (в этом случае устанавливается абсолютный, то есть так, как оно должно идти по стандартной html разметке игнорируя остальные элементы) */
    top: 100%;
    width: 150px;
    background-color: #8a8a8a;
    display: none; /* Изначально задаем .dropdown невидимым */
}

#navigation .dropdown li {
    width: 100%;
}

#navigation li:hover .dropdown {
    display: block; /* Когда срабатывает li:hover дисплей .dropdown меняется на block, то есть делает его видимым */
}
```

# Fixed
Способ позиционирования fixed заставляет объект неперемещаемым при прокрутке сайта:

```
#fixed { /* Возьмем простой квадрат с id fixed */
    font-size: 50px;
    background: rgba(255, 255, 255, 0.5);
    top: 50%;
    left: 50%; /* Хотим сделать квадрат поверх всего посередине, однако мы указываем позицию верхнего лувого угла квадрата, а не его центра */
    margin-left: -25px; /* Сдвигаем наверх на половину всей ширины квадрата */
    margin-top: -25px; /* Сдвигаем наверх на половину всей высоты квадрата */
	/* Теперь квадрат отцентрализован */
    position: fixed; /* Делает квадрат фиксированным, поверх всей страницы, не будет перемещаться со скролл баром*/
}
```

# Псевдо классы
Псевдо классы отвечают за взаимодействие с пользователем, вот основные из них:

```
.pclass { /* контейнер в html с названием класса pclass */
    color: red;
}

.pclass:hover { /* Проверка наводки мышки на объект */
    color:green;
}

.pclass:visited { /* Проверка перехода по ссылке объекта */
    color: blue;
}

.pclass.focus { /* Проверка взаимодействия с ссылкой (перетаскивание, наводка через Tab) */
    color: yellow;
}

#article p:first-child { /* Выбрать первый элемент из предложенных */
    color: aquamarine;
}

#article p:last-child { /* Выбрать последний элемент из предложенных */ 
    color: black;
}
```

# Псевдо элементы
Псевдо элементы отвечают за оформление маленьких деталей:

```
#pelem::before { 
    content: "<br>";
    display: block;
} /* Добавление текста к контейнеру pelem, before значит добавить до контейнера */

#pelem::after { /* Добавление текста к контейнеру pelem, after значит добавить после контейнера */
    content: "<br>";
    display: block;
}

#pelem::first-letter { /* Изменить стиль первой буквы в контейнере текста */
    color: rebeccapurple;
    font-size: large;
}

#fl::first-line { /* Изменить стиль первой строчки в контейнере текста */
    font-size: 5px;
    color: #3cff00;
}

::selection { /* Изменить стиль выделения текста (без объекта применяется ко всему) */
    color: #3cff00;
    background-color: black;
}
```

# Использование спецсимволов
Для оформления могу использоваться спецсимволы +, >, ~:
*Подробнее [здесь](https://techbrij.com/css-selector-adjacent-child-sibling)*

```
#article {
    width: 760px;
    padding-right: 20px;
    float: left;
}

#article > p { /* Оформляем только то, что находится конктретно в контейнере article*/
    background: rgb(0, 255, 42);
    font-family: "Ubuntu", sans-serif;
}

#article .combinator + p { /* Оформляем только один контейнер, находящийся после обозначенного .combinator */
    background-color: blue;
}

#article .combinator ~ p { /* Оформляем все контейнера, находящиеся после обозначенного .combinator */
    background-color: rgb(255, 0, 200);
}
```

# Комбинации
Комбинации позволяют более точно настраивать стили для определенных элементов, например:

*site.html*
```
<input type="text">
<input type="button" value="Find">
<p><a href="#" target="_blank"> AAAAAAAAAAAAAA</a></p>
<p><a href="#" target="_blank" title="ссылка внеш"> AAAAAAAAAAAAAA</a></p>
<p><a href="#" class="social-vk"> AAAAAAAAAAAAAA</a></p>
<p><a href="#" title="social vk"> AAAAAAAAAAAAAA</a></p>
<p><a href="#" title="buttonicon"> AAAAAAAAAAAAAA</a></p>
<p><a href="#" title="button icon"> AAAAAAAAAAAAAA</a></p>
```

*site.css*
```
a[target = "_blank"] { /* Если target равен _blank, тогда применить свойства */
    color: rgb(111, 0, 255);
}

a[title~="внеш"] { /* Если в title присутствует слово "внеш", отделенное пробелом (вообще отделение пробелом все равно что отделение запятой, то есть разные слова) */
    background-color: red;
}

a[class|="social"] { /* Если в названии класса в начале стоит social */
    background-color: rgb(255, 217, 0);
}

a[title^="social"] { /* Если первое в title слово social */
    background-color: rgb(55, 0, 255);
}

a[title$="icon"] { /* Если последнее в title слово icon */
    background-color: rgb(226, 95, 95);
}

a[title*="icon"] { /* Если в title присутсвтует хотя-бы одно упоминание icon */
    background-color: rgb(212, 255, 163);
}

input[type="text"]{ /* Так же можно изменять значение инпута и других объектов */
    background-color: black;
}
```

# @media
Правило `@media something {/* Тут ваши стили */}` используется для определения характеристик устройства, с которого открыт сайт, или размеров страницы. Это позволяет настраивать страницу, если она, например, будет ужата до определенной ширины. Можно указывать, например, стили для принтера используя `@media print` или через проверку минимальной ширины страницы `@media (min-width:400px)`

Точно так же можно делать и в html файле, для применения разных css файлов при разных максимальных/минимальных размерах страницы сайта.
```
<link rel="stylesheet" href="style.css" media="(max-width:980px) and (min-width: 600px)">
```
# [Источник](https://www.youtube.com/watch?v=zCxwcJsglJQ)
В качестве источника представлен очень длинный видос на ютубе, полностью объясняющий весь CSS