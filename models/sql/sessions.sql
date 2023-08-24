CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE,
    token_hash TEXT UNIQUE NOT NULL
);

INSERT INTO
    sessions (user_id, token_hash)
VALUES
    ($ 1, $ 2) RETURNING id;

-- estabelecer relacionamento com outra tabela
-- faz com que ids que colocarmos nesse tabela de fato precisam existir na tabela de usuário e se deletarmos um usuário, não podemos ter uma sessão para aquele user id pois uma sessão que usa aquele user_id impede que o usuário seja deletado, então para deletar o usuário precisamos deletar antes a sessão
CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE REFERENCES users (id),
    token_hash TEXT UNIQUE NOT NULL
);

-- ou 
CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE,
    token_hash TEXT UNIQUE NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users (id)
);

-- para adicionar uma relação a uma tabela já existente
ALTER TABLE
    sessions
ADD
    CONSTRAINT sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id);

-- Para deletar sessões quando o usuário relacionado for deletado
CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE REFERENCES users (id) ON DELETE CASCADE,
    token_hash TEXT UNIQUE NOT NULL
);

SELECT
    users.id,
    users.email,
    users.password_hash
FROM
    sessions
    JOIN users ON users.id = sessions.user_id
WHERE
    sessions.token_hash = $ 1