SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for pastebin
-- ----------------------------
DROP TABLE IF EXISTS `pasteBin`;
CREATE TABLE `pasteBin` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `poster` varchar(30) DEFAULT NULL,
  `syntax` varchar(30) DEFAULT NULL,
  `content` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;
