INSERT INTO "user" (username, email, password)
       VALUES ('test', 'test@gohst.com', 'password');

INSERT INTO command (commandstring)
       VALUES ('ls -l');

INSERT INTO context (hostname, username, shell, directory)
       VALUES ('laptop', 'me', '/bin/bash', '/home/me');

INSERT INTO "session" (contextid, "timestamp")
       VALUES (1, '1999-01-08 0:00:00');

INSERT INTO tag (name)
       VALUES ('test');

INSERT INTO tag (name)
       VALUES ('simple');

INSERT INTO tag (name)
       VALUES ('dir');

INSERT INTO invocation (userid, commandid, returnstatus, "timestamp", sessionid)
       VALUES(1, 1, 0, '1999-01-08 4:05:06', 1);

INSERT INTO invocationtag (tagid, invocationid)
       VALUES(2, 1);

INSERT INTO invocationtag (tagid, invocationid)
       VALUES(3, 1);