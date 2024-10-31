DROP TABLE public.alerts CASCADE;

DROP TABLE public.tickets CASCADE;

DROP INDEX IF EXISTS idx_alert_ticket_alert_id;
