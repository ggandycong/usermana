CREATE TABLE `tbl_user_info`(
    `id` bigInt(20) NOT NULL AUTO_INCREMENT,
    `user_name` varchar(255) NOT NULL DEFAULT '',
    `nick_name` varchar(255) NOT NULL DEFAULT '',
    `pic_name` varchar(255) DEFAULT '',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`user_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `tbl_login_info`(
    `id` bigInt(20) NOT NULL AUTO_INCREMENT,
    `user_name` varchar(255) NOT NULL DEFAULT '',
    `password` varchar(255) NOT NULL DEFAULT '',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`user_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;