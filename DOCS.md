Day1 (17.09.2025)
1. Creating Architecture for project. (Done)
2. Creating Schema (Done)
3. Creating models 
4. Opening web server


1) Architecture of Project
zhake_site
|______config - Для чтение конфигурационного файла, в моем случае я бы хотел через .env работать чтобы db, api и   все такое здесь было
|______models - Думаю нужно чтобы работать с бд, то есть через модел я работаю с таблицами
|______repo - Здесь сама логика CRUD и операций с таблицами
|______handlers - обработчики api calls, то есть это связь с клиентом и то что ожидает клиент и все это передает сервису
|______router - это связь между запросом клиента и handler-ов
|______service - Это главная логика которое связывает repo и handlers
|______main.go - точка входа в программу
|______Dockerfile
|______README.md
|______DOCS.md - здесь буду каждое действие указывать, что сделал и для чего.
|______.env - для конфигурационных данных.

2) Database Schema

Tables:

users
___________________
user_id
user_name
____________________


friendships
______________
friends_id
user_id
friend_id
status  // pending/aplied/rejected
______________


media_items
______________________
media_id
type // film/anime/etc
name
year
author
_______________________

recommendations
_______________________________________
reccomendo_id
from_user_id // (from) who recommended
to_user_id // (to) who did you recommend
media_id // string(name of film/anime)
_________________________________________

3) models

user.go:
User {
    UserId int
    UserName string    
}

friendship.go:
Friendship {
    User User
    Fiend User
    Status string
}   

media_items.go:
MediaItems{
    Type string
    Name string
    Year int
    Author string
}

recommendations.go:
Recommendations{
    FromUser User
    ToUser User
    MediaItem MediaItems
}


Day2 (18.09.2025)

1. Реализовал роутер chi
2. Создал ендпоинты /users, /user/{userID}
3. Создал логику для добавление в друзья


Полностью изменил логику друзей. Теперь работает как в соц сетях. У каждого польхователя есть подписчики и подписки, если это взаимно то они считаются друзьями.


Day3 19.09.2025

Надо создать ендпоинт Get("/media/type/name) который вернет мне список фильмов по названию так как у фильмов название могут совпадать.
Надо создать Post("/recommend") в JSON {fromID, toID, media{}}
Список мне порекомендованных (/user/{userID}/recommendations) можно рядом указать кто порекомендовал
Список который я порекомендовал (/user/{userID}/recommended) можно рядом указать кому порекомендовал