-- everytime insert in invocation check how many invocation use that command id
-- if >20 then create a record in notification table

DROP TRIGGER IF EXISTS command_notification ON invocation;
CREATE TRIGGER command_notification
AFTER INSERT ON invocation
FOR EACH ROW EXECUTE PROCEDURE command_notification_insert();
