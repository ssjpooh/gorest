

CREATE TABLE oauth_client_details (
  user_idx                CHAR(36)      NOT NULL,
  client_id               CHAR(36)      NOT NULL PRIMARY KEY,
  client_secret           CHAR(36)      NOT NULL,
  scope                   VARCHAR(255)  DEFAULT NULL,
  authorized_grant_types  VARCHAR(255)  DEFAULT NULL,
  web_server_redirect_uri VARCHAR(255)  DEFAULT NULL,
	cdate					          VARCHAR(14)	  NULL,
	mdate					          VARCHAR(14)   NULL,
  INDEX (user_idx),
  INDEX (client_id),
  INDEX (cdate),
	INDEX (mdate)
);


CREATE TABLE oauth_client_tokens (
  client_id               CHAR(36)      NOT NULL PRIMARY KEY,
  expires_at              INT(11)       NOT NULL,
  token                   VARCHAR(1000) NOT NULL,
  refresh_token           VARCHAR(1000) NOT NULL,
  server_address          VARCHAR(50)   DEFAULT NULL,
  cdate					          VARCHAR(14)		NULL,
	mdate					          VARCHAR(14)		NULL,
  INDEX (client_id),
  INDEX (cdate),
	INDEX (mdate)
);