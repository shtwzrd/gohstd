package main

const SqlCommandNotificationTrigger = `
CREATE TRIGGER command_notification
AFTER INSERT ON invocation
FOR EACH ROW EXECUTE PROCEDURE command_notification_insert();
`

const SqlInsertNotificationFunction = `
CREATE OR REPLACE FUNCTION command_notification_insert()
  RETURNS trigger AS
$BODY$
	DECLARE
		command_count INTEGER;
		notification_min INTEGER := 20;
		command_string TEXT;
		notification_message TEXT;
	BEGIN
		SELECT COUNT(*) INTO command_count
		FROM invocation
		WHERE commandid= NEW.commandid;

		IF command_count = notification_min THEN

			SELECT commandstring INTO command_string
			FROM command
			WHERE commandid = NEW.commandid;

			notification_message = 'You have issued the command ' || command_string::text || ' verbatim over 20 times. Consider making an alias.';

			INSERT INTO public.notification
			(notificationid, userid, message)
			VALUES(nextval('notification_notificationid_seq'::regclass), NEW.userid, notification_message);


		END IF;
		RETURN NEW;
	END;
$BODY$
  LANGUAGE plpgsql VOLATILE
  COST 100;
`
