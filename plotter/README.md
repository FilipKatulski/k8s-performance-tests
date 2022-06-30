# ClusterLoader's plotting module

This module was created to plot data created with "cl2 density" tests. Creates PNG files of the Timeline, histograms of each phase and a Piechart with cumulative time of main phases.

Created by Filip Katulski <filip.katulski@cern.ch>, TDAQ ATLAS project, CERN 

## Running Plotter 

To run the Plotter type with default settings:

```
go run plotter.go -filepath path/to/timeline.csv -outputpath path/to/folder 
```

Default settings create timeline plot, histograms and pie chart, based on the data labeled with "Stateless" Pods.

However, specific plot types or podstate can also be specified with respective flags. 

To see "help" message run: 
```
go run plotter.go -h 
```
Output:
```
Usage of plotter:
  -filepath string
    	Specify the path to the timeline CSV file. 
  -outputpath string
    	Specify the path for the output PNG files.  (default ".")
  -plots string
    	Specify types of plots, separate with ','  (default "all")
  -podstate string
    	Specify the state of Pods.  (default "Stateless")
```


## Format of CSV files

Input files have to follow a specific format in order to be processed successfully: 

```
Name, Transition, Namespace/PodName, NodeName, pod_state_filter, diff, from_unix, to_unix
pod_startup, {create watch 5s}, test-XXXXXX-1/latency-deployment-0-XXXXXXXXXX-XXXXX, node01, State, 6224, 1647419225, 1647419231
```

As each column means:

Name - name of the Pod's phase

Tranisition - describes Pod's phase as a transition between two states

Namespace/PodName - The name of the Namespace where the testing Pod was deployed

NodeName - Pod was deployed on this node

pod_state_filter - State of the Pod. Possible options: Stateless, Statefull, MatchAll

diff - Time that the Pod spent in the transition between two states

from_unix - Starting point of the Pod's phase

to_unix - Ending point of the Pod's phase 


## Plot examples

![Timeline plot](/plotter/example-plots/timeline.png "Timeline plot")

![Create to Schedule histogram](/plotter/example-plots/createtoschedule-hist.png "Create to Schedule histogram")

![Pie Chart](/plotter/example-plots/piechart.png "Pie Chart")
