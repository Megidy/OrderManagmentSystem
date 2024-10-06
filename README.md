Project Description:


this project was focused on microservice architecture with apache kafka broker and MySql database and building of project with docker and docker compose. It is a lightweight prototype of client communication and the procedure for ordering, picking up and viewing orders in such fast-food restaurants as, for example: McDonald's or KFC or something like that. There are 3 services here: ordering, kitchen and pickup. Each of them communicates with each other using the Kafka broker 
with different topics and key identifiers using http requests.

How to start a project :

1 first of all you need docker to be downloaded and opened.

2 write in terminal in main directory of project: 

docker compose up --build 

this will automaticly build the project with docker and connect zooker, kafka and mysql services.

to test manually :

1 POST make order:   http://localhost:8080/order/create

{
  "dishes": [
    {
      "name": "Pizza",
      "quantity": 2
    },
    {
      "name": "Pasta",
      "quantity": 1
    }
  ]
}


2 GET check orders: http://localhost:8080/orders

3 GET check my orders:  http://localhost:8080/myorders

4 DELETE take my order: http://localhost:8080/orders/take
