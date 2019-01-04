
-- +goose Up
CREATE TABLE IF NOT EXISTS `videos` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `UserId` int(11) NOT NULL,
  `AccountId` int(11) NOT NULL,
  `Title` varchar(250) NOT NULL,
  `Description` text NOT NULL,
  `Tags` varchar(500) NULL,
  `Category` varchar(32) NOT NULL,
  `Language` varchar(32) NOT NULL DEFAULT 'EN',
  `File` varchar(255) NOT NULL,
  `TmpFile` varchar(255) NULL,
  `Playlist` varchar(100) NULL,
  `Title_gen` bool NOT NULL DEFAULT FALSE,
  `Created_at` int(11) NOT NULL,
  `Updated_at` int(11) NULL,
  `Deleted` bool NOT NULL DEFAULT FALSE,
  `Pending` bool NOT NULL DEFAULT FALSE,
  `IpAddress` varchar(250) NULL,
  `Status` bool NOT NULL DEFAULT TRUE,
  PRIMARY KEY (`Id`),
  KEY `videos_UserId` (`UserId`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


-- +goose Down
DROP TABLE IF EXISTS `videos`;

