# pigo

(Patcher and interceptor in GO) In escense `pigo` is a reverse *Proxy* cli with the ability to intercept and modify requests and responses defined on a configuration file.

## Use cases

Pigo becomes really useful when testing apps or services that are interconnected or dependent on external services. For example:

- Mock an entire backend with custom responses. This would allow testing services with limited connectivity or ahead of development.
- Intercept specific responses that change the behavior of the system under test. This allows to modify and test the behavior of the system under test withouht changing the behavior of it's dependencies.

## Install pigo

```bash
curl https://pigo.com
```

### Running pigo

Running pigo is as simple as execute the cli app passing the flags:

```bash
./pigo -p 8400 -f pigo.yaml
```

### Flags

| Flag | Name  | Description | Default |
|------|-------|-------------| ------------|
| p    | port  | Port number to run the proxy on | 8400
| f    | file  | Configuration file to use | pigo.yml

### Configuration file

This is the configuration file
