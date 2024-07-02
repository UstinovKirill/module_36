DROP TABLE IF EXISTS news;

CREATE TABLE news (
                      id SERIAL PRIMARY KEY,
                      title TEXT, -- заголовок
                      content TEXT NOT NULL UNIQUE, -- текст
                      published BIGINT DEFAULT 0, -- время публикации
                      link TEXT -- ссылка );