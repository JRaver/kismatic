# Add-ons

Once the core Kubernetes cluster components are up and running, KET installs
additional components that are required to produce a fully operational cluster.
These extra components are collectively referred to as add-ons, and they include
the networking plugin, cluster DNS, heapster monitoring, and others.

Add-on installation and configuration can be controlled via the plan file under 
the `add-ons` section. KET supports the following add-ons:

- [CNI](#cni)
- [DNS](#dns)
- [Heapster](#heapster)
- [Dashboard](#dashboard)
- [Package Manager](#package-manager)
- [Rescheduler](#rescheduler)

## CNI
The Container Networking Interface (CNI) enables the use of different
networking solutions with a Kubernetes cluster. KET supports various 
CNI providers, enabling you to choose the most appropriate for your environment.

Plan file options:

| Field | Description | 
|-------|-------------|
| `add_ons.cni.disable` | Set to true to disable the installation of CNI | 
| `add_ons.cni.provider` | Choose the CNI provider. Options: `calico`, `weave`, `contiv`, `custom` |
| `add_ons.cni.options.calico.mode` | The Calico networking mode. Options: `bridged`, `routed` |
| `add_ons.cni.options.calico.log_level` | Set the log level for Calico components. Options: `warning`, `info`, `debug` |
| `add_ons.cni.options.calico.workload_mtu` | Set Calico's CNI config `mtu` |
| `add_ons.cni.options.calico.felix_input_mtu` | Set Calico's `FELIX_IPINIPMTU` env |

### Disabled CNI
When CNI is disabled, KET will skip the installation of the CNI binaries and CNI plugin.
Furthermore, KET will also skip the cluster smoke test, as it will fail without a working
pod network.

### Custom CNI provider
If the providers supported by KET do not fit your needs, you may bring your own
CNI-compliant provider by setting the CNI provider to `custom` in the plan file.

In this configuration, KET will skip the cluster smoke test, as it will fail without
a working pod network.

## DNS
DNS provides service discovery to pods running on the cluster, and is a required component for a functional cluster. 

There are different options available on the market, though KET deploys either [KubeDNS](https://github.com/kubernetes/dns) or [CoreDNS](https://github.com/coredns/coredns) as the DNS service on the cluster. If you chose to deploy an alternative DNS solution, you can disable the installation and validation of KubeDNS by setting the below flag.

**NOTE** switching providers between upgrades is not supported at this point.

Plan file options:

| Field | Description |
|-------|-------------|
| `add_ons.dns.disable` | Set to true to disable the installation of KubeDNS in the cluster |
| `add_ons.dns.provider` | Options: `kubedns`, `coredns` |

## Heapster
[Heapster](https://github.com/kubernetes/heapster) is a monitoring solution that enables container monitoring throughout
the cluster. When heapster is running on the cluster, `kubectl` and the Kubernetes 
dashboard surface resource utilization metrics for all pods.

**Important:** If you wish to persist the gathered metrics, you must set the `add_ons.heapster.options.influxdb.pvc_name` option.

Plan file options:

| Field | Description |
|---------------|-------------|
| `add_ons.heapster.disable` | Set to true if heapster should not be deployed during installation |
| `add_ons.heapster.options.heapster.replicas`  | Number of replicas for the heapster deployment |
| `add_ons.heapster.options.heapster.serviceType` | Set the service type for the Heapster service |
| `add_ons.heapster.options.heapster.sink` | The location where Heapster will store it's data |
| `add_ons.heapster.options.influxdb.pvc_name` | Name of a persistent volume claim that will be used by the influxdb database for persistence. This PVC must be manually created after installation. |


## Dashboard
The [Kubernetes dashboard](https://github.com/kubernetes/dashboard) is a web-based UI for managing Kubernetes clusters.

Plan file options:

| Field | Description | 
|-------|-------------|
| `add_ons.dashboard.disable` | Set to true to skip the deployment of the Dashboard |


## Package Manager
[Helm](https://github.com/kubernetes/helm) is the official package manager for Kubernetes. KET includes the `helm` client-side binary in the distribution package. KET also installs the server-side agent, Tiller, on the cluster during installation. 

Plan file options:

| Field | Description |
|---------------|-------------|
| `add_ons.package_manager.disable` | Set to true if the package manager should not be deployed during installation |
| `add_ons.package_manager.provider` | The package manager that should be deployed. Options: `helm` |
| `add_ons.package_manager.options.helm.namespace` | Configure the kubernetes `namespace` to deploy `tiller` to |

## Rescheduler
The [Critical Add-Ons Rescheduler](https://kubernetes.io/docs/tasks/administer-cluster/guaranteed-scheduling-critical-addon-pods/) is a pod that runs in the cluster and ensures that ciritcal add-ons are always scheduled.

Plan file options:

| Field | Description | 
|-------|-------------|
| `add_ons.rescheduler.disable` | Set to true to skip the deployment of the Rescheduler |