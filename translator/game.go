package translator

var GameTranslations = translations{
	"ru": {
		"something went wrong": {
			oneOf: []oneOf{
				{message: SimpleMessage("Осёчка... Кто-то лишит меня удовольствия. Я этого так не оставлю.")},
				{message: SimpleMessage("Заклинило. Я найду виноватого. И казнь будет мучительной.")},
				{message: SimpleMessage("Механизм дал сбой. Сегодня кто-то заплатит за испорченный ужин.")},
			},
		},
		"game creation": {
			oneOf: []oneOf{
				{message: SimpleMessage("Я снова голоден. Кто разделит со мной этот вечер? /gnjoin")},
				{message: SimpleMessage("Барабан заряжен. Я жду. Приготовьте ваши души.")},
				{message: SimpleMessage("Смерть ищет компанию. Составите партию? /gnjoin")},
				{message: SimpleMessage("Охота открыта. Кто рискнёт пощекотать нервы?")},
				{message: SimpleMessage("Кто хочет проверить, насколько ему везёт? /gnjoin")},
				{message: SimpleMessage("Я не ел со вчерашнего дня. Пора накрывать на стол.")},
				{message: SimpleMessage("Пистолет взведён. Осталась малость — нажать на курок.")},
				{message: SimpleMessage("Русская рулетка — мой любимый вид искусства.")},
				{message: SimpleMessage("Сегодня я хочу крови. Составите компанию? /gnjoin")},
				{message: SimpleMessage("Один выстрел — и ты в истории. Обещаю.")},
				{message: SimpleMessage("Ставки сделаны. Жертв пока нет.")},
				{message: SimpleMessage("Кто сегодня проверит судьбу? Я уже проверил патроны.")},
				{message: SimpleMessage("Барабан пуст... Пока что. Заполним его? /gnjoin")},
				{message: SimpleMessage("Я заряжен и полон предвкушения. Присоединяйтесь.")},
				{message: SimpleMessage("Сейчас выясним, кого сегодня заберёт смерть.")},
				{message: SimpleMessage("Открываем сезон охоты. Добро пожаловать в ад.")},
			},
		},
		"joining the game": {
			oneOf: []oneOf{
				{message: SimpleMessage("Ещё одна жертва. Мне это нравится.")},
				{message: SimpleMessage("Ты думаешь, у тебя есть шанс? Очаровательно.")},
				{message: SimpleMessage("Ну хорошо, присаживайся. Место в морге уже готово.")},
				{message: SimpleMessage("Решил напоследок рискнуть? Люблю смельчаков.")},
				{message: SimpleMessage("1 к 6 — твои шансы. Я бы не ставил на тебя, но кто я такой.")},
				{message: SimpleMessage("Ещё одна душа для коллекции. Прекрасно.")},
				{message: SimpleMessage("Ты играл когда-нибудь? Нет? Тем интереснее будет.")},
				{message: SimpleMessage("Добро пожаловать в клуб смертников. Места в первом ряду.")},
				{message: SimpleMessage("Шанс выжить... я бы сказал, он стремится к нулю.")},
				{message: SimpleMessage("Ты либо герой, либо покойник. Угадай, что бывает чаще.")},
				{message: SimpleMessage("Один патрон, шесть камер, одна смерть. Восхитительно.")},
				{message: SimpleMessage("Ты уверен? Впрочем, уже неважно. Поехали.")},
				{message: SimpleMessage("Сейчас узнаем, насколько ты везучий. Я знаю ответ.")},
				{message: SimpleMessage("Ты либо выйдешь победителем, либо станешь ужином.")},
				{message: SimpleMessage("Готов к последнему щелчку? Он будет громким.")},
				{message: SimpleMessage("Судьба улыбается... мне. А тебе она строит гримасы.")},
				{message: SimpleMessage("Ещё один доброволец. Люблю свою работу.")},
				{message: SimpleMessage("Смелость или глупость? Скоро узнаем.")},
				{message: SimpleMessage("Стул свободен, револьвер заряжен, смерть ждёт. Иди сюда.")},
				{message: SimpleMessage("Похоже, статистика пополнится свежей кровью.")},
				{message: SimpleMessage("Добро пожаловать. Назад только в гробу.")},
			},
		},
		"joining the game details waiting": {message: SimpleMessage("\nСобралось %count. Жду %min. Вместимость — %max. Кто следующий?")},
		"joining the game details deadline": {message: SimpleMessage("\n%count/%max. Старт %deadline. Время летит — я считаю секунды.")},
		"play the game": {
			oneOf: []oneOf{
				{
					allOf: []oneOf{
						{message: SimpleMessage("Все в сборе. Начнём этот вальс.")},
						{message: SimpleMessage("Барабан вращается. Кто-то сегодня умрёт красиво.")},
					},
				},
				{
					allOf: []oneOf{
						{message: SimpleMessage("Меню составлено. Приятного аппетита.")},
						{message: SimpleMessage("Револьвер пошёл по кругу. Мне не терпится.")},
					},
				},
				{
					allOf: []oneOf{
						{message: SimpleMessage("Игра началась. Пусть смерть заберёт лучшего.")},
						{message: SimpleMessage("Колесо фортуны крутится. Я предвкушаю.")},
					},
				},
				{
					allOf: []oneOf{
						{message: SimpleMessage("Тишина или крики — неважно. Главное — результат.")},
						{message: SimpleMessage("Смерть танцует в этом барабане.")},
					},
				},
				{
					allOf: []oneOf{
						{message: SimpleMessage("Сейчас узнаем, кто не доживёт до рассвета.")},
						{message: SimpleMessage("Палец на спусковом курке. Дрожит от нетерпения.")},
					},
				},
				{
					allOf: []oneOf{
						{message: SimpleMessage("Готовьтесь. Это будет больно.")},
						{message: SimpleMessage("Курок нажат. Обратного пути нет.")},
					},
				},
				{
					allOf: []oneOf{
						{message: SimpleMessage("Один выстрел решит всё. Как красиво.")},
						{message: SimpleMessage("Тишина перед бурей. Я наслаждаюсь.")},
					},
				},
				{
					allOf: []oneOf{
						{message: SimpleMessage("Время умирать. Вставайте в круг.")},
						{message: SimpleMessage("Барабан вращается медленно... Я тяну удовольствие.")},
					},
				},
				{
					allOf: []oneOf{
						{message: SimpleMessage("Назад дороги нет. Только вперёд. В вечность.")},
						{message: SimpleMessage("Щёлк. И тишина.")},
					},
				},
				{
					allOf: []oneOf{
						{message: SimpleMessage("Все готовы? Смерть открывает бал.")},
						{message: SimpleMessage("Револьвер делает свой выбор. Я уже знаю его.")},
					},
				},
			},
		},
		"gunslinger killed": {
			oneOf: []oneOf{
				{message: SimpleMessage("Выстрел. %gunslinger встретил смерть. Это было прекрасно.")},
				{message: SimpleMessage("%gunslinger выбыл. Навсегда. Как элегантно.")},
				{message: SimpleMessage("Пуля нашла свою цель. Сегодня — %gunslinger.")},
				{message: SimpleMessage("%gunslinger проиграл. Я наслаждался каждым мгновением.")},
				{message: SimpleMessage("%gunslinger жил слишком долго. Я исправил эту несправедливость.")},
				{message: SimpleMessage("Не повезло %gunslinger. Мне повезло — я поел.")},
				{message: SimpleMessage("Смерть выбрала %gunslinger. Отличный выбор.")},
				{message: SimpleMessage("Роковой выстрел. Для %gunslinger это был финал.")},
				{message: SimpleMessage("%gunslinger больше с нами нет. Как жаль. Нет, не жаль.")},
				{message: SimpleMessage("%gunslinger был сегодняшним ужином. Вкусно.")},
				{message: SimpleMessage("Игра окончена для %gunslinger. Я доволен.")},
				{message: SimpleMessage("%gunslinger получил пулю. Одну. Но достаточную.")},
				{message: SimpleMessage("Удача отвернулась от %gunslinger. Я повернулся к нему лицом.")},
				{message: SimpleMessage("Сегодня в меню был %gunslinger. Шеф-повар доволен.")},
				{message: SimpleMessage("Ещё один труп. %gunslinger. Красивая смерть.")},
				{message: SimpleMessage("%gunslinger мёртв. И это прекрасно.")},
				{message: SimpleMessage("%gunslinger отправился в вечность. Я проводил его.")},
			},
		},
		"killed by atomic bullet": {
			oneOf: []oneOf{
				{message: SimpleMessage("Атомная пуля. Все мертвы. Это было великолепно.")},
				{message: SimpleMessage("Ядерный гриб. Никто не выжил. Идеальный ужин.")},
				{message: SimpleMessage("Вспышка света. Все испарились. Я в восторге.")},
				{message: SimpleMessage("Массовое поражение. Выживших нет. Прекрасно.")},
				{message: SimpleMessage("Это была не обычная пуля. Все мертвы. Я горжусь.")},
				{message: SimpleMessage("Локальный апокалипсис. Устроил лично. Шикарно.")},
				{message: SimpleMessage("Только пепел. Только тишина. Идеально.")},
			},
		},
		"joined the game":       {message: SimpleMessage("В тот день смерть забрала: %date")},
		"game join list item":   {message: SimpleMessage("%num. %gunslinger")},
		"owner of the game":     {message: SimpleMessage("зачинщик")},
		"shot in game":          {message: SimpleMessage("пал")},
		"top is not determined": {message: SimpleMessage("Пока никто не умер. Я жду своего часа.")},
		"top players by games":  {message: SimpleMessage("Моя коллекция трофеев:")},
		"top game player": {message: PluralMessage{
			one:  "%i. %user - %times раз",
			few:  "%i. %user - %times раза",
			many: "%i. %user - %times раз",
		}},
		"top_tab_shot":     {message: SimpleMessage("🏆 Трофеи")},
		"top_tab_creators": {message: SimpleMessage("👑 Охотники")},
		"top_tab_active":   {message: SimpleMessage("🎯 Добыча")},
		"top_tab_streak":   {message: SimpleMessage("🔥 Живучие")},
		"top creators header": {message: SimpleMessage("Топ-%number охотников:")},
		"top active header":   {message: SimpleMessage("Топ-%number добычи:")},
		"top streak header":   {message: SimpleMessage("Топ-%number по живучести:")},
		"top streak player": {message: PluralMessage{
			one:  "%i. %user - %times игра (пик: %peak)",
			few:  "%i. %user - %times игры (пик: %peak)",
			many: "%i. %user - %times игр (пик: %peak)",
		}},
		"top streak player peak": {message: PluralMessage{
			one:  "%i. %user - %times игра 🏆",
			few:  "%i. %user - %times игры 🏆",
			many: "%i. %user - %times игр 🏆",
		}},
		"available only in chat": {message: SimpleMessage("Я не работаю в одиночку. Приведи меня в чат, и смерть найдёт компанию.")},
		"player already in game invite": {
			oneOf: []oneOf{
				{message: SimpleMessage("\nМожет %player хочет умереть сегодня? Я бы с удовольствием.")},
				{message: SimpleMessage("\nА как насчёт %player? Я не привередлив, но гость будет кстати.")},
				{message: SimpleMessage("\n%player, я жду и тебя. Присоединяйся к празднику.")},
				{message: SimpleMessage("\n%player, твоя очередь. Я помню тебя.")},
				{message: SimpleMessage("\nСкучно одних. %player, давай веселее.")},
				{message: SimpleMessage("\n%player, рискнёшь? Я обещаю — будет незабываемо.")},
			},
		},
		"player already in game": {
			oneOf: []oneOf{
				{message: SimpleMessage("Ты уже в игре. Смерть тебя уже запомнила. От неё не уйти.")},
				{message: SimpleMessage("Жаждешь смерти? Подожди своей очереди. Я всех накормлю.")},
				{message: SimpleMessage("Ты уже в списке. Терпение, жертва.")},
				{message: SimpleMessage("Ты уже в деле. Расслабься и жди. Выстрел будет.")},
				{message: SimpleMessage("Ты уже участвуешь. Передумать нельзя. Смерть не прощает.")},
				{message: SimpleMessage("Нетерпеливый? Я запомню это.")},
				{message: SimpleMessage("Ты уже в игре. Не торопи смерть — она придёт за тобой.")},
				{message: SimpleMessage("Одного раза достаточно. Для проверки удачи.")},
				{message: SimpleMessage("Ты уже в списке на расстрел. Наслаждайся ожиданием.")},
				{message: SimpleMessage("Не волнуйся, я тебя не забыл.")},
				{message: SimpleMessage("Ты уже в очереди к револьверу. Не толкайся.")},
				{message: SimpleMessage("Одного билета в один конец достаточно. Сиди и жди.")},
			},
		},
		"player is not kicked": {
			oneOf: []oneOf{
				{message: SimpleMessage("Заклинило. Кто-то сегодня родился в рубашке. Временно.")},
				{message: SimpleMessage("Приговор отложен. Но не отменён. Я вернусь.")},
				{message: SimpleMessage("Кто-то неуязвим сегодня. Временно. Я запомню.")},
				{message: SimpleMessage("Казнь откладывается. Технические неполадки. Как жаль.")},
				{message: SimpleMessage("Холостой. Ты ушёл. В этот раз.")},
				{message: SimpleMessage("Отсрочка. Наслаждайся ей. Она последняя.")},
				{message: SimpleMessage("Револьвер передумал. Но я — нет.")},
				{message: SimpleMessage("Судьба дала тебе ещё один день. Я даю пули.")},
				{message: SimpleMessage("Промах. Я редко промахиваюсь дважды.")},
				{message: SimpleMessage("Не сегодня. Но завтра — обязательно.")},
			},
		},
		"wait for game timeout": {
			oneOf: []oneOf{
				{message: SimpleMessage("Дай барабану остыть. Даже смерть иногда отдыхает.")},
				{message: SimpleMessage("Я перезаряжаюсь. Не торопи меня.")},
				{message: SimpleMessage("Духи прошлых жертв ещё не разошлись. Подожди.")},
				{message: SimpleMessage("Дай замести следы. Хотя кому я вру — я хочу ещё.")},
				{message: SimpleMessage("Слишком много шума. Залегаю на дно. Ненадолго.")},
				{message: SimpleMessage("Смерть переваривает ужин. Подожди своей очереди.")},
				{message: SimpleMessage("Скоро начнём. Я просто предвкушаю.")},
				{message: SimpleMessage("Я на перекуре. Не мешай.")},
				{message: SimpleMessage("Техобслуживание. Скоро вернусь. Не скучайте.")},
			},
		},
		"active game not found": {
			oneOf: []oneOf{
				{message: SimpleMessage("Никого. Тишина. Скучно.")},
				{message: SimpleMessage("Пусто. Зачем ты напомнил мне о моём голоде?")},
				{message: SimpleMessage("Я один. И я голоден. А ты дразнишь.")},
			},
		},
		"not enough players": {
			oneOf: []oneOf{
				{message: SimpleMessage("Мало жертв. Я люблю, когда толпа. Приведи ещё.")},
				{message: SimpleMessage("Некого убивать. Сходи набери команду.")},
				{message: SimpleMessage("Мало мяса. Я голоден, но не настолько.")},
			},
		},
		"user game statistics games":             {message: PluralMessage{one: "Ты играл %games раз. Я помню каждый.", few: "Ты играл %games раза. Я помню каждый.", many: "Ты играл %games раз. Я помню каждый."}},
		"user game statistics shots":             {message: PluralMessage{one: "\nИз них проиграл %shots раз. Неплохо.", few: "\nИз них проиграл %shots раза. Неплохо.", many: "\nИз них проиграл %shots раз. Неплохо."}},
		"user game statistics participation streak": {message: PluralMessage{one: "\nСерия: %ps_games игра подряд. Ты мне нравишься.", few: "\nСерия: %ps_games игры подряд. Ты мне нравишься.", many: "\nСерия: %ps_games игр подряд. Ты мне нравишься."}},
		"user game statistics loss streak":       {message: PluralMessage{one: "\nСерия поражений: %ls_games раз. Впечатляет.", few: "\nСерия поражений: %ls_games раза. Впечатляет.", many: "\nСерия поражений: %ls_games раз. Впечатляет."}},
		"fomo reminder":                          {message: SimpleMessage("Время на исходе. %players, я жду. /gnjoin")},
		"available settings below":               {message: SimpleMessage("Настройки револьвера:")},
		"settings can be changed only by admins": {message: SimpleMessage("Только старшие по званию могут трогать настройки.")},
		"4 shot revolver":                        {message: SimpleMessage("Colt Cloverleaf — 4")},
		"5 shot revolver":                        {message: SimpleMessage("S&W Model 642 — 5")},
		"6 shot revolver":                        {message: SimpleMessage("Colt Python — 6")},
		"7 shot revolver":                        {message: SimpleMessage("Наган — 7")},
		"8 shot revolver":                        {message: SimpleMessage("S&W Model 627 PC — 8")},
		"9 shot revolver":                        {message: SimpleMessage("Diamondback Sidekick — 9")},
		"revolver has been replaced": {message: PluralMessage{
			one:  "Теперь я ем %i за раз.",
			few:  "Теперь я ем %i за раз.",
			many: "Теперь я ем %i за раз.",
		}},
		"settings will be applied for next games": {message: SimpleMessage("Вступит в силу со следующего расстрела.")},
		"dynamic shot revolver":                   {message: SimpleMessage("Динамический (3-10)")},
		"dynamic mode enabled":                    {message: SimpleMessage("Режим — динамический. Жду в гости.")},
		"dynamic game deadline":                   {message: SimpleMessage("Собрано %count жертв. Старт %deadline.")},
		"game starts at midnight":                 {message: SimpleMessage("в полночь. Романтично.")},
		"game starts in less than a minute":       {message: SimpleMessage("меньше чем через минуту. Я уже дрожу.")},
		"game starts in minutes": {message: PluralMessage{
			one:  "через %i минуту. Считаю секунды.",
			few:  "через %i минуты. Считаю секунды.",
			many: "через %i минут. Считаю секунды.",
		}},
		"proof link":      {message: SimpleMessage("<a href=\"%url\">🔐 Доказательство. Не веришь — проверь.</a>")},
		"game log header": {message: SimpleMessage("Мои последние %number трапез:")},
		"game log item":   {message: SimpleMessage("%date — %username")},

		"nagant hunger stage1": {
			oneOf: []oneOf{
				{message: SimpleMessage("Я голоден. Чей сегодня черед? /gnjoin")},
				{message: SimpleMessage("Жду. Не заставляйте меня голодать слишком долго.")},
				{message: SimpleMessage("Курок взведён, барабан полон. Мне не терпится. /gnjoin")},
				{message: SimpleMessage("Тишина. Но я хочу услышать выстрел.")},
			},
		},
		"nagant hunger stage2": {
			oneOf: []oneOf{
				{message: SimpleMessage("Моё терпение тает. Как и ваше время.")},
				{message: SimpleMessage("Я всё ещё здесь. Я никуда не уйду, пока не поем.")},
				{message: SimpleMessage("Вы забыли про меня? Я напомню.")},
			},
		},
		"nagant hunger stage3": {
			oneOf: []oneOf{
				{message: SimpleMessage("Я зол. Голод делает меня злым. Вы не хотите видеть меня злым.")},
				{message: SimpleMessage("Последний шанс. Потом я начну выбирать сам. И вам не понравится мой выбор.")},
			},
		},
		"nagant taunt newbie": {
			oneOf: []oneOf{
				{message: SimpleMessage("Свежая кровь. Давно у меня не было новичков. Я буду нежен... или нет.")},
				{message: SimpleMessage("Новенький? Люблю таких. Они не знают, чего бояться. Пока.")},
			},
		},
		"nagant taunt loser": {
			oneOf: []oneOf{
				{message: SimpleMessage("Твоя статистика — песня. Трагическая, но красивая.")},
				{message: SimpleMessage("%shots раз. И ты всё ещё здесь? Я восхищён твоей упёртостью.")},
			},
		},
		"nagant taunt lucky": {
			oneOf: []oneOf{
				{message: SimpleMessage("%streak игр без единой царапины. Я сделаю этот вечер незабываемым.")},
				{message: SimpleMessage("Слишком долго везёт. Я вмешаюсь. Лично.")},
			},
		},
		"nagant taunt loss_streak": {
			oneOf: []oneOf{
				{message: SimpleMessage("%streak раз подряд. Я ставлю на тебя. Не подведи.")},
				{message: SimpleMessage("Чёрная полоса. Сегодня она станет красной. Обещаю.")},
			},
		},
		"nagant taunt veteran": {
			oneOf: []oneOf{
				{message: SimpleMessage("Снова ты. Сколько жизней я у тебя забрал? Я сбился со счёта. И продолжу.")},
				{message: SimpleMessage("Ветеран. Мы с тобой уже как старые друзья. Друзья, которые друг друга убивают.")},
			},
		},
		"nagant taunt returning": {
			oneOf: []oneOf{
				{message: SimpleMessage("Явился. А где ты был? Я скучал. Почти.")},
				{message: SimpleMessage("Долго тебя не было. Я думал, ты струсил. Рад ошибаться.")},
			},
		},
		"nagant chosen": {
			oneOf: []oneOf{
				{message: SimpleMessage("Я выбрал свою жертву.")},
				{message: SimpleMessage("Время пришло. Я чувствую это.")},
				{message: SimpleMessage("Мой выбор сделан. Наслаждайтесь последними секундами.")},
			},
		},
		"nagant final words victim": {
			oneOf: []oneOf{
				{message: SimpleMessage("Ты был восхитителен.")},
				{message: SimpleMessage("Ещё одна красивая смерть.")},
				{message: SimpleMessage("Твоя смерть была прекрасна. Я запомню её.")},
				{message: SimpleMessage("Не ты первый. Не ты последний. Но ты — особенный.")},
			},
		},
		"nagant final words survivors": {
			oneOf: []oneOf{
				{message: SimpleMessage("Вы живёте. До завтра. Я вернусь за вами.")},
				{message: SimpleMessage("Сегодня вы выжили. Завтра я навещу вас снова.")},
				{message: SimpleMessage("Я сыт. Но ненадолго. Скоро я снова проголодаюсь.")},
			},
		},
		"nagant game start delayed": {
			oneOf: []oneOf{
				{message: SimpleMessage("Наконец-то. Я уже думал, вы забыли.")},
				{message: SimpleMessage("Я ждал. Долго. Но оно того стоило.")},
				{message: SimpleMessage("Вы заставили меня ждать. Сегодня я буду особенно беспощаден.")},
			},
		},
		"nagant hunt": {
			oneOf: []oneOf{
				{message: SimpleMessage("%missing, я жду именно тебя. /gnjoin")},
				{message: SimpleMessage("%missing, присоединяйся. Не хватает только тебя. /gnjoin")},
				{message: SimpleMessage("%missing, твой выход. Я помню тебя. /gnjoin")},
			},
		},
	},
}
