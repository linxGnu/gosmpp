/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table receive_sms
# ------------------------------------------------------------

DROP TABLE IF EXISTS `receive_sms`;

CREATE TABLE `receive_sms` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `isdn` bigint(20) unsigned NOT NULL,
  `content` varchar(450) NOT NULL,
  `is_processed` tinyint(1) NOT NULL,
  `process_expired` datetime NOT NULL DEFAULT current_timestamp(),
  `request_check_errcode` int(11) DEFAULT NULL,
  `request_check_message` varchar(200) DEFAULT NULL,
  `request_charge_errcode` int(11) DEFAULT NULL,
  `request_charge_message` varchar(200) DEFAULT NULL,
  `request_pay_errcode` int(11) DEFAULT NULL,
  `request_pay_message` varchar(200) DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `indx_process_status` (`process_expired`, `is_processed`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table send_sms
# ------------------------------------------------------------

DROP TABLE IF EXISTS `send_sms`;

CREATE TABLE `send_sms` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `alias` varchar(20) NOT NULL,
  `isdn` bigint(20) unsigned NOT NULL,
  `content` varchar(450) NOT NULL,
  `submit_status` smallint(6) NOT NULL,
  `submit_expired` datetime NOT NULL DEFAULT current_timestamp(),
  `submit_status_code` int(11) NOT NULL DEFAULT 0,
  `smsc_message_id` bigint(20) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `indx_submit_status` (`submit_expired`, `submit_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;




/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
