CREATE TABLE users (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    hashed_password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE groups (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE documents (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_uuid UUID NOT NULL REFERENCES groups(uuid) ON DELETE CASCADE ON UPDATE CASCADE,
    name VARCHAR(255) NOT NULL,
    content TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE user_groups (
    user_uuid UUID NOT NULL REFERENCES users(uuid) ON DELETE CASCADE ON UPDATE CASCADE,
    group_uuid UUID NOT NULL REFERENCES groups(uuid) ON DELETE CASCADE ON UPDATE CASCADE,
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (user_uuid, group_uuid)
);

CREATE INDEX idx_documents_group_uuid ON documents(group_uuid);
CREATE INDEX idx_user_groups_user_uuid ON user_groups(user_uuid);
CREATE INDEX idx_user_groups_group_uuid ON user_groups(group_uuid);
CREATE INDEX idx_users_login ON users(login);
CREATE INDEX idx_users_email ON users(email);
