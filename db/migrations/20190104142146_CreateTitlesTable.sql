
-- +goose Up
CREATE TABLE IF NOT EXISTS `titles` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `UserId` int(11) NOT NULL,
  `CampaignId` int(11) NOT NULL,
  `Title` varchar(250) NOT NULL,
  `Tags` varchar(500) NOT NULL,
  `File` varchar(255) NOT NULL,
  `TmpFile` varchar(255) NOT NULL,
  `YoutubeId` varchar(255) NOT NULL,
  `YoutubeUrl` varchar(255) NOT NULL,
  `Posted` bool NOT NULL DEFAULT FALSE,
  `Converted` bool NOT NULL DEFAULT FALSE,
  `Pending` bool NOT NULL DEFAULT FALSE,
  `FrameRate` int(2) NOT NULL DEFAULT 25,
  `Resolution` int(5) NOT NULL DEFAULT 1280,
  `IpAddress` varchar(32) NOT NULL,
  `Created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `Updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `Published_at` TIMESTAMP NULL,
  PRIMARY KEY (`Id`),
  KEY `titles_UserId` (`UserId`),
  UNIQUE (`Title`),
  UNIQUE (`CampaignId`, `FrameRate`, `Resolution`)
  ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


-- +goose Down
DROP TABLE IF EXISTS `titles`;

