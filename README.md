# kube-apiproxy

If you ever had to set up a Kubernetes cluster you know that it is a fairly
simple task, apart from one thing - if API server goes away, you have to bring
it back and tell other components where it is. Some people use an ELB with a
persistent DNS entry which they update if kubernetes API service is moved
somewhere else, others use config management tools to update configuration and
restart Kubernetes components, etc.

A lot of people deploy Kubernetes on CoreOS using fleet which is a really great
way.

This proxy service talks to fleet daemon to find which node runs Kubernetes API
and proxies traffic through to the API server. When that node goes down, fleet
will then start the API service on another node which is automatically
discovered by `kube-apiproxy`.

This means that Kubernetes API is always accessible on `localhost:8081` by default.

## Building

You will need `gb` tool - http://getgb.io/.

```
git clone https://github.com/vaijab/kube-apiproxy.git
cd kube-apiproxy
gb build all
```

## Configuration

Configuration is done via command line arguments.

```bash
$ kube-apiproxy --help
Usage of bin/kube-apiproxy:
  -api-port="8080": kubernetes api port
  -fleet-endpoint="unix:///run/fleet.sock": fleet endpoint
  -proxy-listen="localhost:8081": proxy listen ip:port
  -unit-name="kube-apiserver.service": fleet unit name for kubernetes api server
```

## Contribution

Please feel free to send PRs and issues my way. I am very new to Go
programming, I am sure that the code can be improved a lot by keeping the core
principle the same.

