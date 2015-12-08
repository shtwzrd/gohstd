-- name: create-user-table
CREATE TABLE IF NOT EXISTS "user"
(
  userid serial NOT NULL,
  username character varying(128) UNIQUE,
  email character varying(320),
  password text,
  CONSTRAINT user_pk PRIMARY KEY (userid)
);

-- name: create-tag-table
CREATE TABLE IF NOT EXISTS tag
(
  tagid serial NOT NULL,
  name text,
  CONSTRAINT tag_pk PRIMARY KEY (tagid)
);

-- name: create-command-table
CREATE TABLE IF NOT EXISTS command
(
  commandid serial NOT NULL,
  commandstring text,
  CONSTRAINT command_pk PRIMARY KEY (commandid)
);

-- name: create-context-table
CREATE TABLE IF NOT EXISTS context
(
  contextid serial NOT NULL,
  hostname text,
  username text,
  shell text,
  directory text,
  CONSTRAINT context_pk PRIMARY KEY (contextid)
);

-- name: create-configuration-table
CREATE TABLE IF NOT EXISTS configuration
(
  configurationid serial NOT NULL,
  userid integer,
  key text,
  value text,
  CONSTRAINT configuration_pk PRIMARY KEY (configurationid),
  CONSTRAINT configuration_user_fk FOREIGN KEY (userid)
      REFERENCES "user" (userid) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- name: create-session-table
CREATE TABLE IF NOT EXISTS "session"
(
  sessionid serial NOT NULL,
  contextid integer,
  "timestamp" timestamp without time zone,
  CONSTRAINT session_pk PRIMARY KEY (sessionid),
  CONSTRAINT session_context_fk FOREIGN KEY (contextid)
      REFERENCES context (contextid) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- name: create-servicelog-table
create table IF NOT EXISTS servicelog
(
  requestid serial NOT NULL,
  userid integer,
  message text,
  CONSTRAINT servicelog_pk PRIMARY KEY (requestid),
  CONSTRAINT servicelog_user_fk FOREIGN KEY (userid)
      REFERENCES "user" (userid) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- name: create-notification-table
CREATE TABLE IF NOT EXISTS notification
(
  notificationid serial NOT NULL,
  userid integer,
  message text,
  CONSTRAINT notification_pk PRIMARY KEY (notificationid),
  CONSTRAINT notification_user_fk FOREIGN KEY (userid)
      REFERENCES "user" (userid) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- name: create-invocation-table
CREATE TABLE IF NOT EXISTS invocation
(
  invocationid serial NOT NULL,
  userid integer,
  commandid integer,
  returnstatus smallint,
  "timestamp" timestamp without time zone,
  sessionid integer,
  CONSTRAINT invocation_pk PRIMARY KEY (invocationid),
  CONSTRAINT invocation_command_fk FOREIGN KEY (commandid)
      REFERENCES command (commandid) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT invocation_session_fk FOREIGN KEY (sessionid)
      REFERENCES session (sessionid) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT invocation_user_fk FOREIGN KEY (userid)
      REFERENCES "user" (userid) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- name: create-invocationtag-table
CREATE TABLE IF NOT EXISTS invocationtag
(
  refid serial NOT NULL,
  tagid integer,
  invocationid integer,
  CONSTRAINT invocationtag_pk PRIMARY KEY (refid),
  CONSTRAINT invocationtag_invocation_fk FOREIGN KEY (invocationid)
      REFERENCES invocation (invocationid) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT invocationtag_tag_fk FOREIGN KEY (tagid)
      REFERENCES tag (tagid) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- name: create-commandhistory-view
CREATE OR REPLACE VIEW commandhistory AS
SELECT U.username AS user,I.invocationid,I.sessionid,I.returnstatus,I.timestamp,
CT.hostname,CT.username,CT.shell,CT.directory,CM.commandstring,
ARRAY(SELECT TA.name
FROM(invocationtag TI LEFT OUTER JOIN tag T
ON (TI.tagid = T.tagid)) TA
WHERE TA.invocationid = I.invocationid) AS tags
FROM invocation I INNER JOIN "session" S ON (I.sessionid = S.sessionid)
INNER JOIN context CT ON (S.contextid = CT.contextid)
INNER JOIN command CM ON (I.commandid = CM.commandid)
INNER JOIN "user" U ON (I.userid = U.userid)
ORDER BY I.invocationid DESC;

