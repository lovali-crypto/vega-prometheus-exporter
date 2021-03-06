# vega-prometheus-exporter

The aim of the project is to create a simple yet configurable Prometheus exporter to simplify the creation of Grafana dashboards.

## Configuration file

Please customize the .env file. Available options are:

- VEGA_ENDPOINT: the Vega endpoint

## Output

Here is an output example:

```bash
# HELP vega_sync_cytching_up Is the node catching uo?
# TYPE vega_sync_cytching_up gauge
vega_sync_cytching_up 0
# HELP vega_up Was the last vega query successful.
# TYPE vega_up gauge
vega_up 1
# HELP vega_validator_signing Flag indicating if a validator is signing or not (per validator).
# TYPE vega_validator_signing gauge
vega_validator_signing{validator="B-Harvest"} 1
vega_validator_signing{validator="Chorus One"} 0
vega_validator_signing{validator="Commodum"} 0
vega_validator_signing{validator="Figment"} 1
vega_validator_signing{validator="Greenfield One"} 1
vega_validator_signing{validator="Lovali"} 1
vega_validator_signing{validator="Nala DAO"} 1
vega_validator_signing{validator="Nodes Guru"} 1
vega_validator_signing{validator="P2P.ORG Validator"} 1
vega_validator_signing{validator="Rockaway Blockchain Fund"} 1
vega_validator_signing{validator="Ryabina"} 1
vega_validator_signing{validator="Staking Facilities"} 1
vega_validator_signing{validator="XPRV"} 0
```