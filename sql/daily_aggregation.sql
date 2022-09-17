CREATE MATERIALIZED VIEW energy_usage_daily
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('1 day', time) as day,
    last(power_meter1, time) - first(power_meter1, time) as power_usage1,
    last(power_meter2, time) - first(power_meter2, time) as power_usage2,
    (last(gas_meter, time) - first(gas_meter, time)) as gas_usage
FROM energy_meter
GROUP BY day;

SELECT add_continuous_aggregate_policy('energy_usage_daily',
    start_offset => INTERVAL '1 week',
    end_offset => INTERVAL '1 hour',
    schedule_interval => INTERVAL '1 hour');
