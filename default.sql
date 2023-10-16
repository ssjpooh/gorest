CREATE TABLE `user_tbl` (
  `user_id` varchar(100) DEFAULT NULL,
  `owner_idx` varchar(100) DEFAULT NULL,
  `user_passwd` varchar(100) DEFAULT NULL,
  `kor_user_name` varchar(100) DEFAULT NULL,
  `eng_user_name` varchar(100) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


CREATE TABLE conf_tbl 
(
	owner_idx VARCHAR(32) NOT NULL, 
	idx VARCHAR(32) NOT NULL, 
	room_no INT(10) AUTO_INCREMENT PRIMARY KEY, 
	title VARCHAR(100) NOT NULL,
	start_date VARCHAR(16) NOT NULL, 
	end_date VARCHAR(16) NOT NULL,
	is_deleted TINYINT(1) NOT NULL

);

CREATE TABLE `oauth_client_details` (
  `owner_idx` varchar(100) NOT NULL,
  `client_id` varchar(255) NOT NULL,
  `client_secret` varchar(255) NOT NULL,
  `scope` varchar(255) DEFAULT NULL,
  `authorized_grant_types` varchar(255) DEFAULT NULL,
  `web_server_redirect_uri` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`owner_idx`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


CREATE TABLE `OAUTH_CLIENT_TOKENS` (
  `client_id` varchar(36) NOT NULL,
  `expires_at` int(11) NOT NULL,
  `token` varchar(1000) NOT NULL,
  `refresh_token` varchar(1000) NOT NULL,
  `server_address` varchar(50) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;