CREATE VIEW CommandHistory AS
SELECT I.sessionid,I.returnstatus,I.timestamp,CT.hostname,CT.username,CT.shell,CT.directory,CM.commandstring,IT.refid
FROM invocation I, "session" S, context CT, command CM, invocationtag IT;