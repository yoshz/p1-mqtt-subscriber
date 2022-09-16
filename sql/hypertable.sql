CREATE TABLE energy_meter (
  time TIMESTAMPTZ NOT NULL,
  location VARCHAR(100) NOT NULL,
  power_meter1 INT NOT NULL,
  power_meter2 INT NOT NULL,
  gas_meter INT NOT NULL
);
ALTER TABLE energy_meter OWNER TO "metrics";
SELECT create_hypertable('energy_meter','time');
