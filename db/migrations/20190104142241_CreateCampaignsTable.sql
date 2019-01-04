
-- +goose Up
CREATE TABLE IF NOT EXISTS `campaigns` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `UserId` int(11) NOT NULL,
  `AccountId` int(11) NOT NULL,
  `VideoId` int(11) NOT NULL,
  `Title` varchar(250) NOT NULL,
  `TotalTitles` INT(11) NOT NULL DEFAULT 0,
  `CompleteTitles` INT(11) NOT NULL DEFAULT 0,
  `TitlesGenerating` bool NOT NULL DEFAULT FALSE,
  `TitlesGenerated` bool NOT NULL DEFAULT FALSE,
  `IpAddress` varchar(32) NOT NULL,
  `DateStart_at` TIMESTAMP NULL,
  `DateComplete_at` TIMESTAMP NULL,
  `Created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `Updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`Id`),
  KEY (`UserId`)
  ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


-- +goose Down
DROP TABLE IF EXISTS `campaigns`;

