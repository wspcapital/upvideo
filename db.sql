CREATE TABLE IF NOT EXISTS `user` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `Email` varchar(50) NOT NULL,
  `PasswordHash` varchar(32) NOT NULL,
  `APIKey` varchar(36) NOT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;


CREATE TABLE IF NOT EXISTS `session` (
  `id` varchar(50) NOT NULL,
  `updated_at` int(11) NOT NULL,
  `data` blob NOT NULL,
  PRIMARY KEY (`id`),
  KEY `updated_at` (`updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE IF NOT EXISTS `videos` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `UserId` int(11) NOT NULL, 
  `Title` varchar(250) NOT NULL, 
  `Description` text NOT NULL, 
  `Tags` varchar(500) NULL, 
  `Category` varchar(32) NOT NULL, 
  `Language` varchar(32) NOT NULL DEFAULT 'EN', 
  `File` varchar(100) NOT NULL, 
  `Playlist` varchar(100) NULL, 
  `Title_gen` bool NOT NULL DEFAULT FALSE, 
  `Created_at` int(11) NOT NULL, 
  `Updated_at` int(11) NULL, 
  `Deleted` bool NOT NULL DEFAULT FALSE, 
  `Pending` bool NOT NULL DEFAULT FALSE, 
  `IpAddress` varchar(250) NULL, 
  `Status` bool NOT NULL DEFAULT TRUE,
  PRIMARY KEY (`Id`),
  KEY `UserId` (`UserId`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


CREATE TABLE IF NOT EXISTS `accounts` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `UserId` int(11) NOT NULL,
  `Username` varchar(250) NOT NULL,
  `Password` varchar(250) NOT NULL,
  `Channelname` varchar(250) NOT NULL,
  `Channelurl` varchar(250) NOT NULL,
  `Clientsecrets` varchar(250) NOT NULL,
  `Requesttoken` varchar(250) NOT NULL,
  `AuthUrl` text NULL,
  `OnetimeCode` varchar(255) NULL,
  `Note` text NULL,
  `Created_at` int(11) NOT NULL, 
  `Updated_at` int(11) NULL, 
  `Deleted` bool NOT NULL DEFAULT FALSE,
  PRIMARY KEY (`Id`),
  KEY `UserId` (`UserId`)
  ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


CREATE TABLE IF NOT EXISTS `titles` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `UserId` int(11) NOT NULL,
  `VideoId` int(11) NOT NULL,
  `Title` varchar(250) NOT NULL,
  `Tags` varchar(500) NOT NULL,
  `File` varchar(250) NOT NULL,
  `Posted` bool NOT NULL DEFAULT FALSE,
  `Converted` bool NOT NULL DEFAULT FALSE,
  `Pending` bool NOT NULL DEFAULT FALSE,
  `IpAddress` text NULL,
  `Created_at` int(11) NOT NULL, 
  `Updated_at` int(11) NULL, 
  PRIMARY KEY (`Id`),
  KEY `UserId` (`UserId`)
  ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

