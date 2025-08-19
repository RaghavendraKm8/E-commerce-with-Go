-- Initializes three databases and creates tables for each service.
-- Runs automatically on first start of the postgres container.

CREATE DATABASE ecom_users;
CREATE DATABASE ecom_products;
CREATE DATABASE ecom_orders;

\connect ecom_users
CREATE TABLE IF NOT EXISTS users(
  id BIGSERIAL PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  age INT NOT NULL DEFAULT 0
);

\connect ecom_products
CREATE TABLE IF NOT EXISTS products(
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  price BIGINT NOT NULL
);

\connect ecom_orders
CREATE TABLE IF NOT EXISTS orders(
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  product_id BIGINT NOT NULL,
  qty INT NOT NULL,
  total BIGINT NOT NULL
);