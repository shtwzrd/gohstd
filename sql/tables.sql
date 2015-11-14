CREATE TABLE "user"
(
  userid serial NOT NULL,
  username character varying(128) UNIQUE,
  email character varying(320),
  password text,
  CONSTRAINT user_pk PRIMARY KEY (userid)
);

CREATE TABLE tag
(
  tagid serial NOT NULL,
  name text,
  CONSTRAINT tag_pk PRIMARY KEY (tagid)
);

CREATE TABLE command
(
  commandid serial NOT NULL,
  commandstring text,
  CONSTRAINT command_pk PRIMARY KEY (commandid)
);

CREATE TABLE context
(
  contextid serial NOT NULL,
  hostname text,
  username text,
  shell text,
  directory text,
  CONSTRAINT context_pk PRIMARY KEY (contextid)
);

CREATE TABLE configuration
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

CREATE TABLE "session"
(
  sessionid serial NOT NULL,
  contextid integer,
  "timestamp" timestamp without time zone,
  CONSTRAINT session_pk PRIMARY KEY (sessionid),
  CONSTRAINT session_context_fk FOREIGN KEY (contextid)
      REFERENCES context (contextid) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION
);

create table servicelog
(
  requestid serial NOT NULL,
  userid integer,
  message text,
  CONSTRAINT servicelog_pk PRIMARY KEY (requestid),
  CONSTRAINT servicelog_user_fk FOREIGN KEY (userid)
      REFERENCES "user" (userid) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE TABLE notification
(
  notificationid serial NOT NULL,
  userid integer,
  message text,
  CONSTRAINT notification_pk PRIMARY KEY (notificationid),
  CONSTRAINT notification_user_fk FOREIGN KEY (userid)
      REFERENCES "user" (userid) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE TABLE invocation
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

CREATE TABLE invocationtag
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
