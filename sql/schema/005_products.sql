-- +goose Up
CREATE TABLE products (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    supplier TEXT NOT NULL,
    category TEXT NOT NULL,
    price TEXT NOT NULL,
    image_url TEXT NOT NULL,
    description TEXT NOT NULL,
    product_location TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE products;