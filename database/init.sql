USE `user_management`;

CREATE TABLE `users` (
  `id` serial PRIMARY KEY,
  `name` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL UNIQUE,
  `password` text NOT NULL
);