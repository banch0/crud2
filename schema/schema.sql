CREATE table houses (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    code_id TEXT DEFAULT '',
    title TEXT NOT NULL,
    price INTEGER NOT NULL CHECK (price > 0),
    rooms TEXT DEFAULT '',
    description text DEFAULT '',
    city TEXT DEFAULT '',
    rayon TEXT DEFAULT '',
    infrastructure jsonb DEFAULT '{}',
    category TEXT DEFAULT '',
    addr TEXT DEFAULT '',
    owner_id uuid,
    images jsonb DEFAULT '{}',
    type TEXT DEFAULT 'active',
    created_at TIMESTAMP NOT NULL NOW(),
    PRIMARY KEY (id)
);

CREATE table agents_house (
    house_id uuid REFERENCES houses,
    agent_id uuid REFERENCES agents,
    PRIMARY KEY (house_id, agent_id)
);

CREATE table owners_houses (
    house_id uuid REFERENCES houses,
    owners_id uuid REFERENCES owners,
    PRIMARY KEY (house_id, owners_id)
);

CREATE table cities (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    city VARCHAR(40),
    PRIMARY KEY (id)
);

CREATE table rayons (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    rayon VARCHAR(50),
    city_id uuid REFERENCES cities ON DELETE CASCADE,
    PRIMARY KEY (id)
);

CREATE table infrastructure (
    id VARCHAR(50),
    title VARCHAR(50),
    PRIMARY KEY (id)
);

CREATE table orders (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    regdate timestamp DEFAULT now(),
    PRIMARY KEY (id)
);

CREATE table users (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    email VARCHAR(150) not null UNIQUE,
    username VARCHAR(50) NOT NULL,
    image text DEFAULT '',
    bio text DEFAULT '',
    password VARCHAR(64) NOT NULL,
    status VARCHAR(20) DEFAULT 'employer',
    regdate timestamp DEFAULT now(),
    PRIMARY KEY (id)
    UNIQUE(username, email)
);

CREATE TABLE IF NOT EXISTS owners (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR(64) not null,
    lastname VARCHAR(64) not null,
    email VARCHAR(150) not null UNIQUE,
    image text DEFAULT '',
    phone text DEFAULT '',
    craetedAt TIMESTAMP DEFAULT now(),
    password VARCHAR(64) NOT NULL,
    PRIMARY KEY (id)
);