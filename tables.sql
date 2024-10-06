CREATE TABLE orders(
    order_id INT PRIMARY KEY ,
    customer_id INT ,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status ENUM('pending', 'in_progress', 'completed', 'cancelled') NOT NULL
)

CREATE TABLE dishes(
    order_id INT  ,
    name VARCHAR(255),
    quantity INT
)