CREATE MATERIALIZED VIEW energy_usage_monthly
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('1 month', time) as month,
    (
        CAST(last(power_meter1, time) - first(power_meter1, time) AS FLOAT)
        + CAST(last(power_meter2, time) - first(power_meter2, time) AS FLOAT)
    ) / 1000 AS power_usage_kwh,
    CAST(last(gas_meter, time) - first(gas_meter, time) AS FLOAT) / 1000 AS gas_usage_m3
FROM energy_meter
GROUP BY month;


SELECT add_continuous_aggregate_policy('energy_usage_monthly',
    start_offset => INTERVAL '3 months',
    end_offset => INTERVAL '1 hour',
    schedule_interval => INTERVAL '1 day');
