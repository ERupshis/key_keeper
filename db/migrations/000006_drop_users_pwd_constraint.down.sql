ALTER TABLE users
    ADD CONSTRAINT users_password_key UNIQUE (password);