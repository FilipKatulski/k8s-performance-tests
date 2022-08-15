# ClusterLoader

This clusterloader module was modified by Filip Katulski <filip.katulski@cern.ch>, TDAQ ATLAS project, CERN.

It was suited to run cluster density tests and to retreive timeline data. The data analysis adn plotting are described further down this file.  

## Running ClusterLoader

To run ClusterLoader type:

- Set the KUBECONFIG and GOPATH env:
```
export KUBECONFIG=~/example/filepath/config-tbed
export GOPATH=~/example/path/to/pkg
```
- Then follow one of these commands:
```
go run cmd/clusterloader.go --kubeconfig=kubeConfig.yaml --testconfig=config.yaml --provider=local
```

```
./run-e2e.sh --testconfig=config.yaml --provider=local
```


```
./run.sh
```

Or run it manually according to the structure provided in run.sh and run-e2e.sh files.

### Quick start

The simplest way to get acquainted with ClusterLoader is using [kind](https://kind.sigs.k8s.io/). Provision a cluster and ensure you can SSH to localhost. Eg. running the load test:

```
go run cmd/clusterloader.go --testconfig=testing/load/config.yaml --nodes=1 --provider=kind --kubeconfig=<path-to-kubeconfig> --masterip=127.0.0.1 --mastername=kind-control-plane --master-internal-ip=127.0.0.1
```

### Flags

#### Required

These flags are required for any test to be run.
 - kubeconfig - path to the kubeconfig file.
 - testconfig - path to the test config file. This flag can be used multiple times
if more than one test should be run.
 - provider - Cluster provider, options are: gce, gke, kind, kubemark, aws, local, vsphere, skeleton

#### Optional

 - nodes - number of nodes in the cluster.
If not provided, test will assign the number of schedulable cluster nodes.
 - report-dir - path to directory, where summaries files should be stored.
If not specified, summaries are printed to standard log.
 - mastername - Name of the master node
 - masterip - DNS Name / IP of the master node
 - testoverrides - path to file with overrides.
 - kubelet-port - Port of the kubelet to use (*default: 10250*)

## Tests

### Test definition

Test definition is an instantiation of this [api] (in json or yaml).
The motivation and description of the API can be found in [design doc].
Definitions of test as well as definitions of individual objects support templating.
Templates for test definition come with one predefined value - ```{{.Nodes}}```,
which represents the number of schedulable nodes in the cluster. \
Example of a test definition can be found here: [load test].

### Modules

ClusterLoader2 supports modularization of the test configs via the [Module API](https://github.com/kubernetes/perf-tests/blob/1bbb8bd493e5ce6370b0e18f3deaf821f3f28fd0/clusterloader2/api/types.go#L77).
With the Module API, you can divide a single test config file into multiple
module files. A module can be parameterized and used multiple times by
the test or other module. This provides a convenient way to avoid copy-pasting
and maintaining super-long, unreadable test configs<sup id="a1">[1](#f1)</sup>.

[TODO(mm4tt)]: <> (Point to the load config based on modules here once we migrate it)

### Object template

Object template is similar to standard kubernetes object definition
with the only difference being templating mechanism.
Parameters can be passed from the test definition to the object template
using the ```templateFillMap``` map.
Two always available parameters are ```{{.Name}}``` and ```{{.Index}}```
which specifies object name and object replica index respectively. \
Example of a template can be found here: [load deployment template].

### Overrides

Overrides allow to inject new variables values to the template. \
Many tests define input parameters. Input parameter is a variable
that potentially will be provided by the test framework. Cause input parameters are optional,
each reference has to be opaqued with ```DefaultParam``` function that will
handle case if given variable doesn't exist. \
Example of overrides can be found here: [overrides]

#### Passing environment variables

Instead of using overrides in file, it is possible to depend on environment
variables. Only variables that start with `CL2_` prefix will be parsed and
available in script.

Environment variables can be used with `DefaultParam` function to provide sane
default values.

##### Setting variables in shell
```shell
export CL2_ACCESS_TOKENS_QPS=5
```

##### Usage from test definition
```yaml
{{$qpsPerToken := DefaultParam .CL2_ACCESS_TOKENS_QPS 0.1}}
```

## Measurement

Currently available measurements are:
- **APIAvailabilityMeasurement** \
This measurement collects information about the availability of cluster's control plane. \
There are two slightly different ways this is measured:
  - cluster-level availability, where we periodically issue an API call to `/readyz`,
  - host-level availability, where we periodically poll each of the control plane's host `/readyz` endpoint.
    - this requires the [exec service](https://github.com/kubernetes/perf-tests/tree/master/clusterloader2/pkg/execservice) to be enabled.
- **APIResponsivenessPrometheusSimple** \
This measurement creates percentiles of latency and number for server api calls based on the data collected by the prometheus server. 
Api calls are divided by resource, subresource, verb and scope. \
This measurement verifies if [API call latencies SLO] is satisfied.
If prometheus server is not available, the measurement will be skipped.
- **APIResponsivenessPrometheus** \
This measurement creates summary for latency and number for server api calls
based on the data collected by the prometheus server.
Api calls are divided by resource, subresource, verb and scope. \
This measurement verifies if [API call latencies SLO] is satisfied.
If prometheus server is not available, the measurement will be skipped.
- **CPUProfile** \
This measurement gathers the cpu usage profile provided by pprof for a given component.
- **EtcdMetrics** \
This measurement gathers a set of etcd metrics and its database size.
- **MemoryProfile** \
This measurement gathers the memory profile provided by pprof for a given component.
- **MetricsForE2E** \
The measurement gathers metrics from kube-apiserver, controller manager,
scheduler and optionally all kubelets.
- **PodStartupLatency** \
This measurement verifies if [pod startup SLO] is satisfied.
- **ResourceUsageSummary** \
This measurement collects the resource usage per component. During gather execution,
the collected data will be converted into summary presenting 90th, 99th and 100th usage percentile
for each observed component. \
Optionally resource constraints file can be provided to the measurement.
Resource constraints file specifies cpu and/or memory constraint for a given component.
If any of the constraint is violated, an error will be returned, causing test to fail.
- **SchedulingMetrics** \
This measurement gathers a set of scheduler metrics.
- **SchedulingThroughput** \
This measurement gathers scheduling throughput.
- **Timer** \
Timer allows for measuring latencies of certain parts of the test
(single timer allows for independent measurements of different actions).
- **WaitForControlledPodsRunning** \
This measurement works as a barrier that waits until specified controlling objects
(ReplicationController, ReplicaSet, Deployment, DaemonSet and Job) have all pods running.
Controlling objects can be specified by label selector, field selector and namespace.
In case of timeout test continues to run, with error (causing marking test as failed) being logged.
- **WaitForRunningPods** \
This is a barrier that waits until required number of pods are running.
Pods can be specified by label selector, field selector and namespace.
In case of timeout test continues to run, with error (causing marking test as failed) being logged.
- **Sleep** \
This is a barrier that waits until requested amount of the time passes.

## Prometheus metrics

There are two ways of scraping metrics from pods within cluster:
- **ServiceMonitor** \
Allows to scrape metrics from all pods in service. Here you can find example [Service monitor]
- **PodMonitor** \
Allows to scrape metrics from all pods with specific label. Here you can find example [Pod monitor]

## Vendor

Vendor is created using [Go modules].

## Data Analysis and Plotter module

Plotter go module can be used to plot Timeline data gathered by clusterloader density tests (i.e. cl2-density-config.yaml). Plots are saved as PNG files and consist of timelines, histograms and a piechart.

A Jupyter Notebook was created for interactive visualisations created by Python and Plotly library. For instructions on how to create plots please follow the instructions inside the folder. 

Standalone version of the plotter module can be found [here](https://gitlab.cern.ch/fkatulsk/cl2-plotter). 

This module was created by Filip Katulski <filip.katulski@cern.ch>, TDAQ ATLAS project, CERN 


---

<sup><b id="f1">1.</b> As an example and anti-pattern see the 900 line [load test config.yaml](https://github.com/kubernetes/perf-tests/blob/92cc27ff529ae3702c87e8f154ea62f3f2d8e837/clusterloader2/testing/load/config.yaml) we ended up maintaining at some point. [↩](#a1)</sup>


[api]: https://github.com/kubernetes/perf-tests/blob/master/clusterloader2/api/types.go
[API call latencies SLO]: https://github.com/kubernetes/community/blob/master/sig-scalability/slos/api_call_latency.md
[exec service]: https://github.com/kubernetes/perf-tests/tree/master/clusterloader2/pkg/execservice
[design doc]: https://github.com/kubernetes/perf-tests/blob/master/clusterloader2/docs/design.md
[Go modules]: https://blog.golang.org/using-go-modules
[load deployment template]: https://github.com/kubernetes/perf-tests/blob/master/clusterloader2/testing/load/deployment.yaml
[load test]: https://github.com/kubernetes/perf-tests/blob/master/clusterloader2/testing/load/config.yaml
[overrides]: https://github.com/kubernetes/perf-tests/blob/master/clusterloader2/testing/density/scheduler/pod-affinity/overrides.yaml
[pod startup SLO]: https://github.com/kubernetes/community/blob/master/sig-scalability/slos/pod_startup_latency.md
[Service monitor]: https://github.com/kubernetes/perf-tests/blob/master/clusterloader2/pkg/prometheus/manifests/default/prometheus-serviceMonitorKubeProxy.yaml
[Pod monitor]: https://github.com/kubernetes/perf-tests/blob/master/clusterloader2/pkg/prometheus/manifests/default/prometheus-podMonitorNodeLocalDNS.yaml
