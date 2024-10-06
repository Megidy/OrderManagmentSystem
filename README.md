Project Description:


this project was focused on microservice architecture with apache kafka broker and MySql database and building of project with docker and docker compose. It is a lightweight prototype of client communication and the procedure for ordering, picking up and viewing orders in such fast-food restaurants as, for example: McDonald's or KFC or something like that. There are 3 services here: ordering, kitchen and pickup. Each of them communicates with each other using the Kafka broker 
with different topics and key identifiers using http requests.

How to start a project :

1 first of all you need docker to be downloaded and opened.

2 write in terminal in main directory of project: 

docker compose up --build 

this will automaticly build the project with docker and connect zooker, kafka and mysql services.

To create database and tables you have execute mysql service :

docker exec -it ordermanagmentsystem-mysql-1 mysql -uroot -ppassword

Create db and tables :

create database ordermanagmentsystem_db;use ordermanagmentsystem_db;create table orders(order_id INT PRIMARY KEY,customer_id INT ,created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,status ENUM('pending', 'in_progress', 'completed', 'taken') NOT NULL);create table dishes(order_id INT  ,name VARCHAR(255),quantity INT);
