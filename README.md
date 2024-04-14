Инструкция по запуску проекта:<br>
1. Иметь установленный язык Go <br>
2. Установить и запустить БД PostgreSQL <br>
3. Создать базу данных в PostgreSQL <br>
4. Поменять данные в файле .env <br>
   4.1. Переменная PORT. Пример заполнения - PORT=8080 <br>
   4.2. Переменная DB_URL. Из чего состоит - postgres://Login_from_db:password_from_db@db_address:db_port/db_name. Если подключение по http, то нужно в конце, после db_name дописать ?sslmode=disable. Пример заполнения - DB_URL=postgres://lim:123@localhost:5432/car_catalog?sslmode=disable <br>
   4.3. Переменная API_URL. Полный путь к API для получения информации о машине по ее регистрационному номеру. Пример заполнения - API_URL=http://localhost:8000/info <br>
   4.4. Переменная SERVER_URL. Адрес на котором будет запускаться этот проект (без порта). Нужен для запуска swagger ui. Пример заполнения - SERVER_URL=http://localhost <br>
5. Открыть корневую папку проекта в командной строке <br>
6. Ввести в командную строку go build <br>
7. Ввести в командную строку ./<folder_name>.exe. <folder_name> заменить на название папки с проектом <br>
8. После этого запустятся миграции в БД (в БД, которая была указана в .env файле, появится таблица cars), а потом запустится сам сервер. Также, если перейти на сайт <SERVER_URL_from_dot_env>:<PORT_from_dot_env>/swagger/index.html можно будет попасть на страницу swagger ui. Если перейти на сайт <SERVER_URL_from_dot_env>:<PORT_from_dot_env>/swagger/doc.json, то можно увидеть swagger документацию в формате JSON (или можно открыть в самом проекте папку docs и там будет файл swagger.json и swagger.yaml)<br>



Instructions for launching the project:<br>
1. Have the Go language installed <br>
2. Install and run the PostgreSQL database <br>
3. Create a database in PostgreSQL <br>
4. Change the data in the file.env <br>
4.1. The PORT variable. Example of filling - PORT=8080 <br>
4.2. The DB_URL variable. What does it consist of - postgres://Login_from_db:password_from_db@db_address:db_port/db_name . If the connection is over http, then you need to add at the end, after db_name ?sslmode=disable. Example of filling - DB_URL=postgres://lim:123@localhost:5432/car_catalog?sslmode=disable <br>
4.3. The API_URL variable. The full path to the API for getting information about the machine by its registration number. Example of filling - API_URL=http://localhost:8000/info <br>
4.4. The SERVER_URL variable. The address where this project will be launched (without the port). It is needed to launch the swagger ui. Example of filling - SERVER_URL=http://localhost <br>
5. Open the root folder of the project in the command line <br>
6. Enter go build into the command line <br>
7. Enter into the command line./<folder_name>.exe. <folder_name> replace with the name of the project folder <br>
8. After that, migrations to the database will start (the cars table will appear in the database specified in the .env file), and then the server itself will start. Also, if you go to the site <SERVER_URL_from_dot_env>:<PORT_from_dot_env>/swagger/index.html you will be able to get to the swagger ui page. If you go to the site <SERVER_URL_from_dot_env>:<PORT_from_dot_env>/swagger/doc.json, you can see the swagger documentation in JSON format (or you can open the docs folder in the project itself and there will be a swagger file.json and swagger.yaml)<br>
