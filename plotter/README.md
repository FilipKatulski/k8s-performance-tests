# ClusterLoader's Plotting module

This module was created to plot data created with "cl2 density" tests. Creates PNG files of the Timeline, histograms of each phase and a Piechart with cumulative time of main phases.

Created by Filip Katulski <filip.katulski@cern.ch> 

## Running Plotter 

To run the Plotter type with default settings:

```
go run plotter.go -filepath path/to/timeline.csv -outputpath path/to/folder 
```
Additional parameters:
```
  -additional string
    	Specify additional parameteres for plotting, separate with ',' 
  -plots string
    	Specify types of plots, separate with ','  (default "all")
  -podstate string
    	Specify the state of Pods.  (default "Stateless")
```

