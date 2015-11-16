package main

const SqlCreateCommandHistoryView = `
CREATE OR REPLACE VIEW CommandHistory AS
SELECT DISTINCT I.invocationid,I.sessionid,I.returnstatus,I.timestamp,
                CT.hostname,CT.username,CT.shell,CT.directory,CM.commandstring,
                ARRAY(SELECT TA.name
                      FROM(invocationtag TI LEFT OUTER JOIN tag T
                           ON (TI.tagid = T.tagid)) TA
                      WHERE TA.invocationid = I.invocationid) AS tags
FROM invocation I, "session" S, context CT, command CM;
`
