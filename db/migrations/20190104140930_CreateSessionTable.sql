
-- +goose Up
CREATE TABLE IF NOT EXISTS `session` (
  `id` varchar(50) NOT NULL,
  `updated_at` int(11) NOT NULL,
  `data` blob NOT NULL,
  PRIMARY KEY (`id`),
  KEY `updated_at` (`updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


-- +goose Down
DROP TABLE IF EXISTS `session`;

