package main

const SqlCreateCommandHistoryView = `
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
`
