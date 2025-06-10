# Modules for Prometheus Modbus exporter

Custom configuration modules for @RichiH's [Modbus Prometheus exporter
(`RichiH/modbus_exporter`)](https://github.com/RichiH/modbus_exporter) for
personal use.

The modules are available in the [`module` directory](./module/).

## Notes for manual Modbus testing

```shell
docker run -it --rm --network host docker.io/library/debian:stable
```

```shell
apt-get update && \
apt-get install -yy python3-pymodbus python3-prompt-toolkit
```

```shell
pymodbus.console --verbose tcp --host host --port 502

client.read_holding_registers slave=1 count=1 address=1234
```

<!-- vim: set sw=2 sts=2 et : -->
