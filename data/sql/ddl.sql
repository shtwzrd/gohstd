-- name: create-user-table
CREATE TABLE IF NOT EXISTS "user"
(
  username character varying(128),
  email character varying(320),
  password character varying(60),
  isadmin boolean,
  created timestamp without time zone,
  image OID,
  CONSTRAINT user_pk PRIMARY KEY (username)
);

-- name: create-post-table
CREATE TABLE IF NOT EXISTS post
(
  postid serial NOT NULL,
  username character varying(128),
  title character varying(255) NOT NULL,
  message text,
  "timestamp" timestamp without time zone,
  CONSTRAINT post_pk PRIMARY KEY (postid)
  CONSTRAINT username_fk FOREIGN KEY (username)
           REFERENCES "user" (username) MATCH SIMPLE
           ON UPDATE NO ACTION ON DELETE NO ACTION
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

-- name: create-invocation-table
CREATE TABLE IF NOT EXISTS invocation
(
  invocationid serial NOT NULL,
  username character varying(128),
  commandid integer,
  exitcode smallint,
  hostname text,
  "user" text,
  shell text,
  directory text,
  "timestamp" timestamp without time zone,
  CONSTRAINT invocation_pk PRIMARY KEY (invocationid),
  CONSTRAINT invocation_command_fk FOREIGN KEY (commandid)
      REFERENCES command (commandid) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT invocation_user_fk FOREIGN KEY (username)
      REFERENCES "user" (username) MATCH SIMPLE
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
  SELECT I.username,I.invocationid,I.timestamp,I.exitcode,
    I.hostname,I.user,I.shell,I.directory,CM.commandstring,
    ARRAY(SELECT TA.name
      FROM(invocationtag TI LEFT OUTER JOIN tag T
        ON (TI.tagid = T.tagid)) TA
  WHERE TA.invocationid = I.invocationid) AS tags
  FROM invocation I INNER JOIN command CM ON (I.commandid = CM.commandid)
  ORDER BY I.timestamp DESC;

-- name: create-timestamp-index
DO $$
BEGIN
IF NOT EXISTS (
  SELECT to_regclass('public.timestampindex')
  ) THEN
  CREATE INDEX timestampindex ON public.invocation ("timestamp");
END IF;
END$$;
