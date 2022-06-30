# ClusterLoader's Plotting module

This module was created to plot data created with "cl2 density" tests. Creates PNG files of the Timeline, histograms of each phase and a Piechart with cumulative time of main phases.

Created by Filip Katulski <filip.katulski@cern.ch> 

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
pod_startup, {create watch 5s}, test-XXXXXX-1/latency-deployment-0-5f6455489f-8f4dr, node01, State, 6224, 1647419225, 1647419231
```