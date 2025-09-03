package main

import ( // Импортируем все необходимые библиотеки
	"bufio"   // Библиотека для сканера
	"fmt"     // Выводит строки + иногда считывает строку для ожидания
	"os"      // В этом случае работает в тандеме с bufio
	"strings" // Форматирует строки
	"time"    // Засекает секундомер
)

type Salat struct { // Структура салата цезаря со всеми необходимыми атрибутами
	sliced_eggs  bool
	eggs         bool
	tomato       bool
	crackers     bool
	shrimps      bool
	salat_leaves bool
	parmesan     bool
	water        bool
}

type Borsch struct { // Структура борща со всеми необходимыми атрибутами
	cabbage     bool
	potato      bool
	carrot      bool
	beet        bool
	garlic      bool
	onion       bool
	meat        bool
	boiled_meat bool
	water       bool
}

type Lasagna struct { // Структура лазаньи со всеми необходимыми атрибутами
	heated          bool
	basement        bool
	second_basement bool
	third_basement  bool
	minced_meat     bool
	tomatoes        bool
	meat            bool
	butter          bool
	ricotta         bool
	mozzarella      bool
	first_layer     bool
	second_layer    bool
	third_layer     bool
}

func counter_to_start() { // Выводит обратный отсчет
	fmt.Println("До старта:\n3")
	time.Sleep(1 * time.Second)
	fmt.Println("2")
	time.Sleep(1 * time.Second)
	fmt.Println("1")
	time.Sleep(1 * time.Second)
	fmt.Println("НАЧАЛИ ГОТОВКУ!!!!")
}

func await(ch chan bool, sec time.Duration) { // Асинхронная функция ожидания для иммитации работы духовки/варящегося мяса
	time.Sleep(sec * time.Second)
	ch <- true // По завершению возвращаем true в канал
}

func main() { // Основная функция
	scanner := bufio.NewScanner(os.Stdin) // Инициализируем сканер

	fmt.Println("Выбери блюдо:\n\tСалат цезарь\n\tБорщ\n\tЛазанья") // Выводим выбор
	scanner.Scan()                                                  // Считываем ответ от пользователя
	dish_to_cook := scanner.Text()                                  // Записываем ответ в переменную

	switch strings.ToLower(dish_to_cook) { // Теперь найдем то, что выбрал пользователь
	case "салат цезарь": // Если это был салат цезарь
		ch := make(chan bool)                                                    // Открываем канал связи
		cooking := Salat{false, false, false, false, false, false, false, false} // Инициализируем салат цезарь

		fmt.Println("Вы собираетесь приготовить салат цезарь, поэтому собрали все ингредиенты:\n\tяйца\n\tпомидорки\n\tайсберг\n\tсухарики\n\tсоус\n\tкреветки\n\tпармезан")
		fmt.Println("Так же у вас есть кран с водой, салатница, плита и кастрюля")
		fmt.Println("(Не забывайте про команды:\n\tНаполнить кастрюлю\n\tПоставить яйца вариться\n\tНарезать креветки\n\tНарезать помидорки\n\tНарезать яйца\n\tТщательно промыть листья салата айсберг\n\tНатереть пармезан\n\tРаспаковать сухарики\n\tСобрать салат)")
		fmt.Println("****Напишите любой символ для старта****")
		// Вывели инструктаж
		scanner.Scan() // После ввода начинается обратный отсчет со стартом

		counter_to_start()                        // Обратный отсчет
		for now := range time.Tick(time.Second) { // Запускаем секундомер
			fmt.Println("Введите действие:")
			scanner.Scan()         // Получаем действие игрока
			line := scanner.Text() // Записываем действие в переменную

			switch strings.ToLower(line) { // Находим нужное действие

			case "поставить яйца вариться":
				if cooking.water == true { // Проверяем, наполнена ли кастрюля
					if cooking.eggs != true { // Проверяем, ставили ли мы вариться яйца
						cooking.eggs = true // Доказываем, что яйца на варку мы поставили
						go await(ch, 30)    // Имитируем варку яиц
					} else {
						fmt.Println("Вы уже поставили яйца вариться")
					}
				} else {
					fmt.Println("Вы не наполнили кастрюлю")
				}

			case "собрать салат":

				switch { // Проверяем, что все подготовлено к сборке

				case cooking.crackers != true:
					fmt.Println("Вы не распаковали сухарики")

				case cooking.salat_leaves != true:
					fmt.Println("Вы не промыли салат")

				case cooking.shrimps != true:
					fmt.Println("Вы не нарезали креветки")

				case cooking.sliced_eggs != true:
					fmt.Println("Вы не нарезали яйца")

				case cooking.parmesan != true:
					fmt.Println("Вы не натерли пармезан")

				case cooking.tomato != true:
					fmt.Println("Вы не нарезали помидорки")

				default: // Если все готово, начинаем собирать салат

					fmt.Println("Вы начали собирать салат")
					time.Sleep(1 * time.Second) // Имитируем сборку салата
					fmt.Println(`Салат "Цезарь" готов!`)
					fmt.Println("У вас ушло", now.String()[41:47], "секунд") // Выводим время приготовления
					fmt.Println("****Напишите любой символ для выхода****")
					scanner.Scan() // После ответа пользователя программа закрывается
					return         // Завершаем работу программы
				}

			case "наполнить кастрюлю":
				if cooking.water != true { // Проверяем, совершали ли мы это действие
					time.Sleep(2 * time.Second) // Иммитируем наполнение кастрюли
					cooking.water = true
				} else {
					fmt.Println("Вы уже наполнили кастрюлю")
				}

			case "нарезать помидорки":
				if cooking.tomato != true { // Проверяем, совершали ли мы это действие
					time.Sleep(2 * time.Second) // Иммитируем нарезание помидорок
					cooking.tomato = true
				} else {
					fmt.Println("Вы уже нарезали помидорки")
				}

			case "распаковать сухарики":
				if cooking.crackers != true { // Проверяем, совершали ли мы это действие
					time.Sleep(2 * time.Second) // Иммитируем распаковку сухариков
					cooking.crackers = true
				} else {
					fmt.Println("Вы уже распаковали сухарики")
				}

			case "нарезать креветки":
				if cooking.shrimps != true { // Проверяем, совершали ли мы это действие
					time.Sleep(2 * time.Second) // Иммитируем нарезание креветок
					cooking.shrimps = true
				} else {
					fmt.Println("Вы уже нарезали креветки")
				}

			case "натереть пармезан":
				if cooking.parmesan != true { // Проверяем, совершали ли мы это действие
					time.Sleep(2 * time.Second) // Иммитируем натирание пармезана
					cooking.parmesan = true
				} else {
					fmt.Println("Вы уже натерли пармезан")
				}

			case "тщательно промыть листья салата айсберг":
				if cooking.salat_leaves != true { // Проверяем, совершали ли мы это действие
					time.Sleep(5 * time.Second) // Иммитируем тщательное промывание листьев салата айсберг
					cooking.salat_leaves = true
				} else {
					fmt.Println("Вы уже промыли листья салата айсберг")
				}

			case "нарезать яйца":
				if cooking.sliced_eggs != true { // Проверяем, совершали ли мы это действие
					if cooking.eggs == true { // Проверяем, ставил ли игрок яйца на варку
						if <-ch == true {
							time.Sleep(2 * time.Second) // Иммитируем нарезание яиц
							cooking.sliced_eggs = true
						}
					} else {
						fmt.Println("Сначала сварите яйца")
					}

				} else {
					fmt.Println("Вы уже нарезали яйца")
				}

			default: // Отрабатываем несуществующее действие
				fmt.Println("Вы не можете этого сделать")
			}
		}

	case "борщ": // Если это был борщ
		ch := make(chan bool)                                                            // Открываем канал связи
		cooking := Borsch{false, false, false, false, false, false, false, false, false} // Инициализируем борщ

		fmt.Println("Вы собираетесь приготовить борщ, поэтому собрали все ингредиенты:\n\tкапусту\n\tкартошку\n\tморковь\n\tсвеклу\n\tчеснок\n\tлук\n\tмясо")
		fmt.Println("Так же у вас есть кран с водой, плита и кастрюля")
		fmt.Println("(Не забывайте про команды:\n\tНаполнить кастрюлю\n\tПоставить мясо вариться\n\tНарезать *название овоща/мясо*\n\tДобавить все нарезанные овощи)")
		fmt.Println("****Напишите любой символ для старта****")
		// Вывели инструктаж
		scanner.Scan() // После ввода начинается обратный отсчет со стартом

		counter_to_start()                             // Обратный отсчет
		for now := range time.Tick(time.Millisecond) { // Запускаем секундомер
			fmt.Println("Введите действие:")
			scanner.Scan()         // Получаем действие игрока
			line := scanner.Text() // Записываем действие в переменную

			switch strings.ToLower(line) { // Находим нужное действие

			case "поставить мясо вариться":
				if cooking.water == true { // Проверяем, наполна ли кастрюля
					if cooking.meat == true { // Проверяем, нарезано ли мясо
						if cooking.boiled_meat != true { // Проверяем, ставили ли мы мясо вариться
							cooking.boiled_meat = true // Доказываем, что мясо на варку мы поставили
							go await(ch, 30)           // Имитием варку мяса
						} else {
							fmt.Println("Вы уже поставили мясо вариться")
						}
					} else {
						fmt.Println("Вы не нарезали мясо")
					}
				} else {
					fmt.Println("Вы не наполнили кастрюлю")
				}

			case "добавить все нарезанные овощи":
				if cooking.boiled_meat == true { // Проверяем, ставил ли игрок мясо на варку
					if <-ch == true { // Когда мясо сварилось
						switch { // Проверяем, что все овощи нарезаны
						case cooking.cabbage != true:
							fmt.Println("Вы не нарезали капусту")

						case cooking.potato != true:
							fmt.Println("Вы не нарезали картошку")

						case cooking.carrot != true:
							fmt.Println("Вы не нарезали морковь")

						case cooking.beet != true:
							fmt.Println("Вы не нарезали свеклу")

						case cooking.garlic != true:
							fmt.Println("Вы не нарезали чеснок")

						case cooking.onion != true:
							fmt.Println("Вы не нарезали лук")

						default: // Если все овощи нарезаны, то начинаем варить сам суп

							fmt.Println("Суп начал готовиться")
							time.Sleep(5 * time.Second) // Имитируем варку супа
							fmt.Println("Суп готов!")
							fmt.Println("У вас ушло", now.String()[41:47], "секунд") // Выводим время приготовления
							fmt.Println("****Напишите любой символ для выхода****")
							scanner.Scan() // После ответа пользователя программа закрывается
							return         // Завершаем работу программы
						}
					}
				} else {
					fmt.Println("Сначала сварите бульон")
				}

			case "наполнить кастрюлю":
				if cooking.water != true { // Проверяем, совершали ли мы это действие
					time.Sleep(2 * time.Second) // Иммитируем наполнение кастрюли
					cooking.water = true
				} else {
					fmt.Println("Вы уже наполнили кастрюлю")
				}

			case "нарезать капусту":
				if cooking.cabbage != true { // Проверяем, совершали ли мы это действие
					time.Sleep(2 * time.Second) // Иммитируем нарезание капусты
					cooking.cabbage = true
				} else {
					fmt.Println("Вы уже нарезали капусту")
				}

			case "нарезать картошку":
				if cooking.potato != true { // Проверяем, совершали ли мы это действие
					time.Sleep(2 * time.Second) // Иммитируем нарезание картошки
					cooking.potato = true
				} else {
					fmt.Println("Вы уже нарезали картошку")
				}

			case "нарезать морковь":
				if cooking.carrot != true { // Проверяем, совершали ли мы это действие
					time.Sleep(2 * time.Second) // Иммитируем нарезание моркови
					cooking.carrot = true
				} else {
					fmt.Println("Вы уже нарезали морковь")
				}

			case "нарезать свеклу":
				if cooking.beet != true { // Проверяем, совершали ли мы это действие
					time.Sleep(2 * time.Second) // Иммитируем нарезание свеклы
					cooking.beet = true
				} else {
					fmt.Println("Вы уже нарезали свеклу")
				}

			case "нарезать чеснок":
				if cooking.garlic != true { // Проверяем, совершали ли мы это действие
					time.Sleep(2 * time.Second) // Иммитируем нарезание чеснока
					cooking.garlic = true
				} else {
					fmt.Println("Вы уже нарезали чеснок")
				}

			case "нарезать лук":
				if cooking.onion != true { // Проверяем, совершали ли мы это действие
					time.Sleep(2 * time.Second) // Иммитируем нарезание лука
					cooking.onion = true
				} else {
					fmt.Println("Вы уже нарезали лук")
				}

			case "нарезать мясо":
				if cooking.meat != true { // Проверяем, совершали ли мы это действие
					time.Sleep(2 * time.Second) // Иммитируем нарезание мясо
					cooking.meat = true
				} else {
					fmt.Println("Вы уже нарезали мясо")
				}

			default: // Отрабатываем несуществующее действие
				fmt.Println("Вы не можете этого сделать")
			}
		}

	case "лазанья": // Если это была лазанья
		ch := make(chan bool)                                                                                         // Открываем канал связи
		cooking := Lasagna{false, false, false, false, false, false, false, false, false, false, false, false, false} // Инициализируем лазанью

		fmt.Println("Вы собираетесь приготовить лазанью, поэтому собрали все ингредиенты:\n\tфарш\n\tпротертые томаты\n\tсливочное масло\n\tлисты для лазаньи\n\tрикотту\n\tмоцареллу")
		fmt.Println("Так же у вас есть сковородка, противень, плита и духовка")
		fmt.Println("(Не забывайте про команды:\n\tВключить духовку на разогрев\n\tПоставить фарш обжариваться\n\tДобавить протертые томаты к фаршу\n\tСмазать сливочным маслом противень\n\tДобавить лист на слой\n\tДобавить фарш на слой\n\tДобавить рикотту на слой\n\tДобавить моцареллу на слой\n\tПоставить лазанью запекаться)")
		fmt.Println("****Напишите любой символ для старта****")
		// Вывели инструктаж
		scanner.Scan() // После ввода начинается обратный отсчет со стартом

		counter_to_start()                             // Обратный отсчет
		for now := range time.Tick(time.Millisecond) { // Запускаем секундомер
			fmt.Println("Введите действие:")
			scanner.Scan()         // Получаем действие игрока
			line := scanner.Text() // Записываем действие в переменную

			switch strings.ToLower(line) { // Находим нужное действие

			case "включить духовку на разогрев":
				if cooking.heated != true { // Проверяем, включили ли мы духовку
					cooking.heated = true // Доказываем, что яйца на варку мы поставили
					go await(ch, 40)      // Иммитируем разогрев духовки
				} else {
					fmt.Println("Вы уже включили духовку на разогрев")
				}

			case "поставить фарш обжариваться":
				if cooking.minced_meat != true { // Проверяем, уже обжаривали ли мы фарш
					time.Sleep(2 * time.Second) // Иммитируем обжарку фарша
					cooking.minced_meat = true
				} else {
					fmt.Println("Вы уже обжарили фарш")
				}

			case "добавить протертые томаты к фаршу":

				if cooking.minced_meat == true { // Проверяем, обжарен ли фарш
					if cooking.tomatoes != true { // Проверяем, уже добавляли ли мы томаты
						time.Sleep(2 * time.Second) // Иммитируем добавление томатов к фаршу
						cooking.tomatoes = true
					} else {
						fmt.Println("Вы уже добавили томаты")
					}

				} else {
					fmt.Println("Вы ещё не обжарили фарш")
				}

			case "смазать сливочным маслом противень":
				time.Sleep(1 * time.Second) // Иммитируем смажку противеня
				cooking.butter = true

			case "добавить лист на слой":
				if cooking.butter == true { // Проверяем, смазан ли противень

					if cooking.basement == true { // В случае, если мы положили листы для первого слоя

						switch { // Проверяем наполненость слоя перед тем, как вновь положить листы

						case cooking.meat == false:
							fmt.Println("Сначала добавьте фарш")

						case cooking.ricotta == false:
							fmt.Println("Сначала добавьте рикотту")

						case cooking.mozzarella == false:
							fmt.Println("Сначала добавьте моцареллу")

						default: // Если все же слой полон, то
							cooking.first_layer = true // Указываем завершённость первого слоя
							cooking.meat = false       // На втором слое отсутствует мясо, поэтому ставим в атрибуте false
							cooking.ricotta = false    // На втором слое отсутствует рикотта, поэтому ставим в атрибуте false
							cooking.mozzarella = false // На втором слое отсутствует моцарелла, поэтому ставим в атрибуте false
							cooking.basement = false   // Указываем, что больше не надо обрабатывать проверку первого слоя
						}

					} else {

						if cooking.first_layer == false && cooking.basement == false { // Проверяем: положили ли мы листы когда ещё не закончили первый слой
							time.Sleep(1 * time.Second) // Иммитируем положение листа
							cooking.basement = true
						}
					}

					if cooking.second_basement == true { // В случае, если мы положили листы для второго слоя

						switch { // Проверяем наполненость слоя перед тем, как вновь положить листы

						case cooking.meat == false:
							fmt.Println("Сначала добавьте фарш")

						case cooking.ricotta == false:
							fmt.Println("Сначала добавьте рикотту")

						case cooking.mozzarella == false:
							fmt.Println("Сначала добавьте моцареллу")

						default: // Если все же слой полон, то
							cooking.second_layer = true     // Указываем завершённость второго слоя
							cooking.meat = false            // На третьем слое отсутствует мясо, поэтому ставим в атрибуте false
							cooking.ricotta = false         // На третьем слое отсутствует рикотта, поэтому ставим в атрибуте false
							cooking.mozzarella = false      // На третьем слое отсутствует моцарелла, поэтому ставим в атрибуте false
							cooking.second_basement = false // Указываем, что больше не надо обрабатывать проверку второго слоя
						}

					} else {

						if cooking.second_layer == false && cooking.second_basement == false && cooking.first_layer == true { // Проверяем: положили ли мы листы когда ещё не закончили второй слой, так же проверяем завершенность предыдушего слоя
							time.Sleep(1 * time.Second) // Иммитируем положение листа
							cooking.second_basement = true
						}
					}

					if cooking.third_basement == true { // В случае, если мы положили листы для третьего слоя

						switch { // Проверяем наполненость слоя перед тем, как вновь положить листы

						case cooking.meat == false:
							fmt.Println("Сначала добавьте фарш")

						case cooking.ricotta == false:
							fmt.Println("Сначала добавьте рикотту")

						case cooking.mozzarella == false:
							fmt.Println("Сначала добавьте моцареллу")

						default: // Если все же слой полон, то
							cooking.third_layer = true     // Указываем завершённость третьего слоя
							cooking.third_basement = false // Указываем, что больше не надо обрабатывать проверку второго слоя
						}

					} else {

						if cooking.third_layer == false && cooking.third_basement == false && cooking.second_layer == true { // Проверяем: положили ли мы листы когда ещё не закончили третий слой, так же проверяем завершенность предыдушего слоя
							time.Sleep(1 * time.Second) // Иммитируем положение листа
							cooking.third_basement = true
						} else if cooking.third_layer == true {
							fmt.Println("Вы закончили лазанью, вам больше не требуется добавлять листы")

						}
					}

				} else { // Отработка того, что противень не смазан
					fmt.Println("Вы ещё не смазали противнь сливочным маслом")
				}

			case "добавить фарш на слой":
				if cooking.third_layer != true { // Проверяем, завершена ли лазанья
					if cooking.tomatoes == true { // Проверяем, добавили ли мы томаты к фаршу
						if cooking.basement == true || cooking.second_basement == true || cooking.third_basement == true { // Проверяем, есть ли активный слой, куда мы будем добавлять мясо
							switch { // Проверяем порядок, есть ли то, что идет до и нет того, что идет после

							case cooking.mozzarella == true:
								fmt.Println("После моцареллы идет лист")

							case cooking.ricotta == true:
								fmt.Println("После рикотты идет моцарелла")

							case cooking.meat == true:
								fmt.Println("Вы уже добавили фарш, теперь добавьте рикотту")

							default:
								time.Sleep(2 * time.Second) // Иммитируем добавление мяса
								cooking.meat = true
							}

						} else { // Если активного слоя, значит мы ещё не начали делать лазанью
							fmt.Println("Вы ещё не добавили лист")
						}
					} else {
						fmt.Println("Вы не добавили томаты к фаршу")
					}
				} else { // Отрабатываем завершение лазаньи
					fmt.Println("Вы закончили лазанью, вам больше не требуется добавлять фарш")
				}

			case "добавить рикотту на слой":
				if cooking.third_layer != true { // Проверяем, завершена ли лазанья
					if cooking.basement == true || cooking.second_basement == true || cooking.third_basement == true { // Проверяем, есть ли активный слой, куда мы будем добавлять рикотту
						switch { // Проверяем порядок, есть ли то, что идет до и нет того, что идет после

						case cooking.mozzarella == true:
							fmt.Println("После моцареллы идет лист")

						case cooking.ricotta == true:
							fmt.Println("Вы уже добавили рикотту, теперь добавьте моцареллу")

						case cooking.meat == false:
							fmt.Println("Вы ещё не добавили фарш")

						default:
							time.Sleep(2 * time.Second) // Иммитируем добавление рикотты
							cooking.ricotta = true
						}

					} else { // Если активного слоя, значит мы ещё не начали делать лазанью
						fmt.Println("Вы ещё не добавили лист")
					}
				} else { // Отрабатываем завершение лазаньи
					fmt.Println("Вы закончили лазанью, вам больше не требуется добавлять фарш")
				}

			case "добавить моцареллу на слой":
				if cooking.third_layer != true { // Проверяем, завершена ли лазанья
					if cooking.basement == true || cooking.second_basement == true || cooking.third_basement == true { // Проверяем, есть ли активный слой, куда мы будем добавлять моцареллу
						switch { // Проверяем порядок, есть ли то, что идет до и нет того, что идет после

						case cooking.mozzarella == true:
							fmt.Println("Вы уже добавили моцареллу, теперь добавьте лист")

						case cooking.ricotta == false:
							fmt.Println("Вы ещё не добавили рикотту")

						case cooking.meat == false:
							fmt.Println("Вы ещё не добавили фарш")

						default:
							time.Sleep(2 * time.Second) // Иммитируем добавление моцареллы
							cooking.mozzarella = true
						}

					} else { // Если активного слоя, значит мы ещё не начали делать лазанью
						fmt.Println("Вы ещё не добавили лист")
					}
				} else { // Отрабатываем завершение лазаньи
					fmt.Println("Вы закончили лазанью, вам больше не требуется добавлять фарш")
				}

			case "поставить лазанью запекаться":
				if cooking.heated == true { // Проверяем, поставил ли игрок духовку на разогрев
					if <-ch == true { // Когда духовка нагрелась
						switch { // Отрабатываем незавершенность лазаньи по слоям
						case cooking.first_layer != true:
							fmt.Println("Вы не закончили первый слой лазаньи")

						case cooking.second_layer != true:
							fmt.Println("Вы не закончили второй слой лазаньи")

						case cooking.third_layer != true:
							fmt.Println("Вы не закончили третий слой лазаньи")

						default: // Если лазанья завершена - запекаем её

							fmt.Println("Лазанья начала запекаться")
							time.Sleep(5 * time.Second) // Иммитируем запекание лазаньи
							fmt.Println("Лазанья готова!")
							fmt.Println("У вас ушло", now.String()[41:47], "секунд")
							fmt.Println("****Напишите любой символ для выхода****")
							scanner.Scan() // После ответа пользователя программа закрывается
							return
						}
					}
				} else {
					fmt.Println("Сначала нагрете духовку")
				}

			default: // Отрабатываем несуществующее действие
				fmt.Println("Вы не можете этого сделать")
			}
		}

	default: // Отрабатываем несуществующее блюдо
		fmt.Println("Вы не можете приготовить", dish_to_cook, "\nПопробуйте заново")
		fmt.Println("****Напишите любой символ для выхода****")
		scanner.Scan() // После ответа пользователя программа закрывается
	}
}
