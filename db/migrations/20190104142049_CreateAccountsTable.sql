
-- +goose Up
CREATE TABLE IF NOT EXISTS `accounts` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `UserId` int(11) NOT NULL,
  `Username` varchar(250) NULL,
  `Password` varchar(250) NULL,
  `Channelname` varchar(250) NOT NULL,
  `Channelurl` varchar(250) NOT NULL,
  `ClientId` varchar(250) NOT NULL,
  `Clientsecrets` varchar(250) NOT NULL,
  `Clientsecretsrow` text NULL,
  `Requesttoken` varchar(250) NOT NULL,
  `Requesttokenrow` text NULL,
  `AuthUrl` text NULL,
  `OnetimeCode` varchar(255) NULL,
  `Note` text NULL,
  `OperationId` varchar(32) NOT NUll,
  `Created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `Updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `Deleted` bool NOT NULL DEFAULT FALSE,
  PRIMARY KEY (`Id`),
  UNIQUE (`UserId`, `ClientId`),
  UNIQUE (`OperationId`),
  KEY `accounts_UserId` (`UserId`)
  ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


-- +goose Down
DROP TABLE IF EXISTS `accounts`;

