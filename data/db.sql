DROP TABLE IF EXISTS `operator`;
CREATE TABLE `operator`
(
    `id`         INTEGER PRIMARY KEY,
    `username`   TEXT     NOT NULL DEFAULT '',
    `password`   TEXT     NOT NULL,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS `openai`;
CREATE TABLE `openai`
(
    `id`                 INTEGER PRIMARY KEY,
    `desc`               TEXT     NOT NULL DEFAULT '',
    `session_token`      TEXT     NOT NULL DEFAULT '',
    `session_token_exp`  TEXT NOT NULL DEFAULT '',
    `created_at`         datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`
(
    `id`         INTEGER PRIMARY KEY,
    `name`       TEXT        NOT NULL DEFAULT '',
    `token`      TEXT        NOT NULL DEFAULT '',
    `expire_at`  datetime,
    `oid`        INTEGER              DEFAULT '0',
    `created_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS `conversation`;
CREATE TABLE `conversation`
(
    `id`              INTEGER PRIMARY KEY,
    `cid`             TEXT     NOT NULL DEFAULT '',
    `title`           TEXT     NOT NULL,
    `uid`             INTEGER           DEFAULT '0',
    `oid`             INTEGER           DEFAULT '0',
    `msg_num`         INTEGER           DEFAULT '0',
    `created_at`      datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`      datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS `message_log`;
CREATE TABLE `message_log`
(
    `id`         INTEGER PRIMARY KEY,
    `cid`        TEXT     NOT NULL DEFAULT '',
    `msgid`      TEXT     NOT NULL DEFAULT '',
    `uid`        INTEGER           DEFAULT '0',
    `oid`        INTEGER           DEFAULT '0',
    `model`      TEXT     NOT NULL DEFAULT '',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- test data
INSERT INTO `operator` (`username`, `password`)
VALUES ('nextshare', 'cns@0001');
INSERT INTO `openai` (`id`, `desc`, `session_token`, `session_token_exp`)
VALUES 
(1, '车1', '请替换为你的session_token', '');

INSERT INTO `user` (`name`, `token`, `oid`, `expire_at`)
VALUES 
('user1', 'cns0001', 1, '2024-10-10 00:00:00'),
('user2', 'cns0002', 1, '2024-10-10 00:00:00');
