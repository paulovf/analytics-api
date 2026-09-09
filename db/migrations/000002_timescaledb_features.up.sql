SELECT create_hypertable('payment_events', 'occurred_at', migrate_data => true, if_not_exists => true);

CREATE INDEX IF NOT EXISTS idx_payment_events_status_time 
ON payment_events (status, occurred_at DESC);

CREATE MATERIALIZED VIEW IF NOT EXISTS payment_status_daily
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('1 day', occurred_at) AS day_bucket,
    status,
    COUNT(*) AS total_events
FROM payment_events
GROUP BY day_bucket, status
WITH NO DATA;

SELECT add_continuous_aggregate_policy('payment_status_daily',
    start_offset => INTERVAL '30 days',
    end_offset => INTERVAL '1 hour',
    schedule_interval => INTERVAL '1 hour',
    if_not_exists => true);
