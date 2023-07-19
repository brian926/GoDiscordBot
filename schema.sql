CREATE DATABASE discordBot;

USE discordBot;

--

-- Table structure for table `discord_messages`

--

DROP TABLE IF EXISTS `discord_messages`;

CREATE TABLE
    `discord_messages` (
        `id` int NOT NULL AUTO_INCREMENT,
        `payload` json NOT NULL,
        `user_id` bigint NOT NULL,
        PRIMARY KEY (`id`)
    ) ENGINE = InnoDB AUTO_INCREMENT = 2 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;