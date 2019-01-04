
-- +goose Up
CREATE TABLE IF NOT EXISTS `user` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `Email` varchar(50) NOT NULL,
  `PasswordHash` varchar(32) NOT NULL,
  `APIKey` varchar(36) NOT NULL,
  `AccountId` int(11) NOT NULL,
  `ForgotPasswordToken` varchar(36) NULL,
  `ForgotPasswordTokenExpiredAt` TIMESTAMP NULL,
  `Created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`Id`),
  UNIQUE (`Email`),
  UNIQUE (`APIKey`),
  KEY (`Email`, `ForgotPasswordToken`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;


-- +goose Down
DROP TABLE IF EXISTS `user`;

