
-- +goose Up
CREATE TABLE IF NOT EXISTS `jobs` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `UserId` int(11) NOT NULL,
  `RelatedId` int(11) NOT NULL,
  `Type` varchar(250) NOT NULL,
  `ProcessId` int(11) NOT NULL,
  `Progress` int(5) NOT NULL DEFAULT 0,
  `Created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `Updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`Id`),
  KEY `jobs_UserId` (`UserId`),
  UNIQUE (`RelatedId`, `Type`)
  ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


-- +goose Down
DROP TABLE IF EXISTS `jobs`;

