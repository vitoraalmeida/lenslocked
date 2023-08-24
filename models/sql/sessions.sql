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

-- Faz em uma query o que seria necessário em duas:
-- Primeiro buscar qual é o user id vinculado a uma sessão
-- depois buscar os dados daquele usuário baseado no id
SELECT
    users.id,
    users.email,
    users.password_hash
FROM
    sessions
    JOIN users ON users.id = sessions.user_id
WHERE
    sessions.token_hash = 1;

-- podemos usar sql indexes para buscar mais rapidamente linhas que são frequentemente necessárias
-- Se vamos fazer uma busca de usuário pelo email muito frequentemente
-- é mais rápido caso esses emails estejam ordenados de forma a facilitar a 
-- busca. Em colunas que possuem primary/foreign keys e unique o postgres
-- já cria automaticamente indexes, então caso colocarmos na tabela que
-- o email é unico, será mais performatico quando fizermos a busca de um usuário
-- pelo email
-- para criar indexes
CREATE TABLE shirts (
    id SERIAL PRIMARY KEY,
    color TEXT,
    size INT,
);

CREATE INDEX shirts_color_size_idx ON shirts(color, size);

-- caso percebamos que muito frequentemente as pessoas comprem blusas de 
-- certas cores e tamanhos mais frequentemente, podemos otimizar essa busca
SELECT
    *
FROM
    shirts
WHERE
    color = 'black'
    AND size > 4;

-- para tentar inserir dados e caso já exista, atualizar o existente
INSERT INTO
    sessions (user_id, token_hash)
VALUES
    ($ 1, $ 2) ON CONFLICT (user_id) DO
UPDATE
SET
    token_hash = $ 2 RETURNING id;