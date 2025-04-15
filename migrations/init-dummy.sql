-- Create users table if not exists
CREATE TABLE IF NOT EXISTS users (
    user_id serial primary key,
    username varchar(50) unique not null,
    email varchar(100) unique not null,
	password varchar(255) not null,
	is_admin boolean default true,
    created_at timestamp default current_timestamp,
	updated_at timestamp default current_timestamp
);

-- Clean up existing users (optional: only for dev/test environment)
TRUNCATE TABLE users RESTART IDENTITY CASCADE;

-- Insert dummy data
INSERT INTO users (username, email, password, is_admin)
VALUES 
  ('admin', 'admin@example.com', '$2a$10$lZJxAlpuKr0uK1vGcP6/MedxbLfsLZAZMOYbJBmr6OD8l7n2Mtydi', true), -- password: admin
  ('user1', 'user1@example.com', '$2a$10$Y5AyyW0.m8eB5UL618CQI.Fqg9kizv0lGFR4jNJgPnrnBplL4WnhG', false); -- password: password

