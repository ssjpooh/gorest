CREATE TABLE USER_tbl 
(
	idx VARCHAR(32) NOT NULL, 
	user_id VARCHAR(100) NOT NULL, 
	user_passwd VARCHAR(100) NOT NULL, 
	email VARCHAR(100) NULL,
	kor_name VARCHAR(100) NULL,
	eng_name VARCHAR(100) NULL,
	dept_code VARCHAR(100) NULL
)

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

CREATE TABLE oauth_client_details 
(
    owner_idx VARCHAR(32) NOT NULL, 
    client_id VARCHAR(32) NOT NULL, 
    client_secret VARCHAR(32) NOT NULL, 

)

CREATE TABLE oauth_tokens 
(

    client_id VARCHAR(32) NOT NULL, 
	expires_at INT NOT NULL,
	token VARCHAR(1000) NOT NULL,
	refresh_tokne VARCHAR(1000) NOT NULL,
	server_address VARCHAR(50)
)