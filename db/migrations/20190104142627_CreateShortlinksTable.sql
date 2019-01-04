
-- +goose Up
CREATE TABLE IF NOT EXISTS `shortlinks` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `UserId` int(11) NOT NULL,
  `UniqId` varchar(25) NOT NULL,
  `Url` text NOT NULL,
  `Counter` bigint(20) NOT NULL DEFAULT 0,
  `Disabled` bool NOT NULL DEFAULT FALSE,
  `Created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `Updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`Id`),
  UNIQUE (`UniqId`)
  ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


-- +goose Down
DROP TABLE IF EXISTS `shortlinks`;

