package main

import (
	"bufio"
	"flag"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

var (
	path_to_file string
	filepath     string
	plots        string
	additional   string
)

type TimelineData struct {
	Name              string
	Transition        string
	Namespace_PodName string
	NodeName          string
	PodStateFilter    string
	Difference        int
	FromUnix          int
	ToUnix            int
}

func displayHeader() {
	//TODO Make it more fancy
	fmt.Println("\n\n  PLOTTER\n\n")
}

func parseDataFile(path string) ([]TimelineData, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var textlines []string
	var timelineData []TimelineData

	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		textlines = append(textlines, s.Text())
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("could not scan: %v", err)
	}

	for _, eachline := range textlines[1:] {
		split := strings.Split(eachline, ", ")

		diff_int, err := strconv.Atoi(split[5])
		if err != nil {
			return nil, fmt.Errorf("could not parse 'diff' string value to integer in csv file: %v", err)
		}
		from_unix_int, err := strconv.Atoi(split[6])
		if err != nil {
			return nil, fmt.Errorf("could not parse 'from_unix' string value to integer in csv file: %v", err)
		}
		to_unix_int, err := strconv.Atoi(split[7])
		if err != nil {
			return nil, fmt.Errorf("could not parse 'to_unix' string value to integer in csv file: %v", err)
		}

		line := TimelineData{
			Name:              split[0],
			Transition:        split[1],
			Namespace_PodName: split[2],
			NodeName:          split[3],
			PodStateFilter:    split[4],
			Difference:        diff_int,
			FromUnix:          from_unix_int,
			ToUnix:            to_unix_int,
		}

		timelineData = append(timelineData, line)
	}

	return timelineData, nil
}

func plotTypeSelection(plotlist string, data []TimelineData) {
	plots := strings.Split(plotlist, ",")
out:
	for _, t := range plots {
		switch t {
		case "all":
			plotTimeline(data, "Stateless")
			plotHistograms(data, "Stateless")
			break out
		case "histograms":
			fmt.Println("Implement histograms")
			plotHistograms(data, "Stateless")
		case "timeline":
			fmt.Println("Implement timeline")
			plotTimeline(data, "Stateless")
		default:
			fmt.Printf("Plot type '%s' is not implemented.\n", t)
		}
	}
}

type DataForPlotting []DataPoint

type DataPoint struct {
	timeStamp  float64
	numOfElems float64
}

func plotTimeline(dat []TimelineData, PodStateFilterSelector string) {
	dataDf := dataframe.LoadStructs(dat)
	dataDf = dataDf.Filter(
		dataframe.F{Colname: "PodStateFilter", Comparator: series.Eq, Comparando: PodStateFilterSelector},
	)

	CreatedDf := dataDf.Filter(
		dataframe.F{Colname: "Transition", Comparator: series.Eq, Comparando: "{create schedule 0s}"},
	)

	CreatedDf = CreatedDf.Select([]string{"FromUnix"})
	CreatedDf = CreatedDf.Arrange(dataframe.Sort("FromUnix"))
	minimalVal, _ := CreatedDf.Elem(0, 0).Int()

	ScheduledDf := dataDf.Filter(
		dataframe.F{Colname: "Transition", Comparator: series.Eq, Comparando: "{create schedule 0s}"},
	)
	ScheduledDf = ScheduledDf.Select([]string{"ToUnix"})
	ScheduledDf = ScheduledDf.Arrange(dataframe.Sort("ToUnix"))

	RunDf := dataDf.Filter(
		dataframe.F{Colname: "Transition", Comparator: series.Eq, Comparando: "{schedule run 0s}"},
	)
	RunDf = RunDf.Select([]string{"ToUnix"})
	RunDf = RunDf.Arrange(dataframe.Sort("ToUnix"))

	WatchedDf := dataDf.Filter(
		dataframe.F{Colname: "Transition", Comparator: series.Eq, Comparando: "{run watch 0s}"},
	)
	WatchedDf = WatchedDf.Select([]string{"ToUnix"})
	WatchedDf = WatchedDf.Arrange(dataframe.Sort("ToUnix"))

	createdGroups := parseTimelineDf(&CreatedDf, minimalVal, "FromUnix")
	schedulerGroups := parseTimelineDf(&ScheduledDf, minimalVal, "ToUnix")
	runGroups := parseTimelineDf(&RunDf, minimalVal, "ToUnix")
	watchGroups := parseTimelineDf(&WatchedDf, minimalVal, "ToUnix")

	createdValues := createDataForTimelinePlotting(createdGroups)
	scheduledValues := createDataForTimelinePlotting(schedulerGroups)
	runValues := createDataForTimelinePlotting(runGroups)
	watchValues := createDataForTimelinePlotting(watchGroups)

	err := createTimelinePlot("timeline.png", createdValues, scheduledValues, runValues, watchValues)
	if err != nil {
		log.Fatalf("could not plot the data: %v", err)
	}
}

func parseTimelineDf(df *dataframe.DataFrame, minVal int, groupingCol string) map[string]dataframe.DataFrame {

	s := df.Rapply(func(s series.Series) series.Series {
		deposit, err := s.Elem(0).Int()
		if err != nil {
			return series.Ints("NAN")
		}
		return series.Ints(deposit - minVal)
	})

	mutatedDf := df.Mutate(s.Col("X0")).Rename("Z", "X0")

	groupedDf := mutatedDf.GroupBy(groupingCol)

	groups := groupedDf.GetGroups()
	return groups
}

func createDataForTimelinePlotting(groups map[string]dataframe.DataFrame) DataForPlotting {
	var values DataForPlotting

	for _, elem := range groups {
		var dp DataPoint

		timeInteg := elem.Elem(0, 1).Float()
		dp.timeStamp = timeInteg
		dp.numOfElems = float64(elem.Nrow())

		values = append(values, dp)
	}

	sort.SliceStable(values, func(i, j int) bool {
		return values[i].timeStamp < values[j].timeStamp
	})

	for i := 1; i < len(values); i++ {
		values[i].numOfElems += values[i-1].numOfElems
	}

	return values
}

func addNewTimeLine(lineName string, p *plot.Plot, dataPoints DataForPlotting) {

	pxys := make(plotter.XYs, len(dataPoints))

	for j, elem := range dataPoints {
		pxys[j].X = elem.timeStamp
	}
	for i, elem := range dataPoints {
		pxys[i].Y = elem.numOfElems
	}

	fmt.Println("XYs of ", lineName, ": ", pxys)

	line, err := plotter.NewLine(pxys)
	if err != nil {
		log.Fatalf("could not add new line %s: %v", lineName, err)
	}

	line.LineStyle.Color = color.RGBA{
		R: uint8(rand.Intn(255)),
		G: uint8(rand.Intn(255)),
		B: uint8(rand.Intn(255)),
		A: 255}
	p.Add(line)
	p.Legend.Add(lineName, line)
}

func createTimelinePlot(path string, created DataForPlotting, scheduled DataForPlotting, run DataForPlotting, watch DataForPlotting) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create %s.png file: %v", path, err)
	}
	defer f.Close()

	fmt.Println("Creating Timeline plot: ")

	p := plot.New()
	p.Title.Text = "Timeline"
	p.X.Label.Text = "Time [s]"
	p.Y.Label.Text = "Number of Pods"

	addNewTimeLine("Created", p, created)
	addNewTimeLine("Scheduled", p, scheduled)
	addNewTimeLine("Run", p, run)
	addNewTimeLine("Watch", p, watch)

	wt, err := p.WriterTo(1024, 512, "png")
	if err != nil {
		return fmt.Errorf("could not create writer: %v", err)
	}

	_, err = wt.WriteTo(f)
	if err != nil {
		return fmt.Errorf("could not write plot to file: %v", err)
	}

	return nil
}

func plotHistograms(dat []TimelineData, PodStateFilterSelector string) {
	fmt.Println("TODO: Implement Histograms. ")
	dataDf := dataframe.LoadStructs(dat)
	dataDf = dataDf.Filter(
		dataframe.F{Colname: "PodStateFilter", Comparator: series.Eq, Comparando: PodStateFilterSelector},
	)

	//Transition from Create to Schedule
	createToScheduleDf := dataDf.Filter(
		dataframe.F{Colname: "Transition", Comparator: series.Eq, Comparando: "{create schedule 0s}"})
	createToScheduleDf = createToScheduleDf.Select([]string{"Difference"})
	createToScheduleDf = createToScheduleDf.Arrange(dataframe.Sort("Difference"))
	createToScheduleValues := createDataForHistogramPlotting(parseHistogramDf(&createToScheduleDf))
	err := createHistogramPlot("createtoschedule-hist.png", "{create schedule 0s}", createToScheduleValues)
	if err != nil {
		log.Fatalf("could not plot the data: %v", err)
	}

	//Transition from Schedule to Run
	scheduleToRunDf := dataDf.Filter(
		dataframe.F{Colname: "Transition", Comparator: series.Eq, Comparando: "{schedule run 0s}"},
	)
	scheduleToRunDf = scheduleToRunDf.Select([]string{"Difference"})
	scheduleToRunDf = scheduleToRunDf.Arrange(dataframe.Sort("Difference"))
	scheduleToRunValues := createDataForHistogramPlotting(parseHistogramDf(&scheduleToRunDf))
	err = createHistogramPlot("scheduletorun-hist.png", "{schedule run 0s}", scheduleToRunValues)
	if err != nil {
		log.Fatalf("could not plot the data: %v", err)
	}

	//Transition from Run to Watch
	runToWatchDf := dataDf.Filter(
		dataframe.F{Colname: "Transition", Comparator: series.Eq, Comparando: "{run watch 0s}"},
	)
	runToWatchDf = runToWatchDf.Select([]string{"Difference"})
	runToWatchDf = runToWatchDf.Arrange(dataframe.Sort("Difference"))
	runToWatchValues := createDataForHistogramPlotting(parseHistogramDf(&runToWatchDf))
	err = createHistogramPlot("runtowatch-hist.png", "{run watch 0s}", runToWatchValues)
	if err != nil {
		log.Fatalf("could not plot the data: %v", err)
	}

	//Transition from Create to Watch
	createToWatchDf := dataDf.Filter(
		dataframe.F{Colname: "Transition", Comparator: series.Eq, Comparando: "{create watch 5s}"},
	)
	createToWatchDf = createToWatchDf.Select([]string{"Difference"})
	createToWatchDf = createToWatchDf.Arrange(dataframe.Sort("Difference"))
	createToWatchValues := createDataForHistogramPlotting(parseHistogramDf(&createToWatchDf))
	err = createHistogramPlot("createtowatch-hist.png", "{create to watch 5s}", createToWatchValues)
	if err != nil {
		log.Fatalf("could not plot the data: %v", err)
	}
}

func parseHistogramDf(df *dataframe.DataFrame) map[string]dataframe.DataFrame {

	groupedDf := df.GroupBy("Difference")

	groups := groupedDf.GetGroups()
	return groups
}

func createDataForHistogramPlotting(groups map[string]dataframe.DataFrame) DataForPlotting {
	var values DataForPlotting

	for _, elem := range groups {
		var dp DataPoint

		timeInteg := elem.Elem(0, 0).Float()
		dp.timeStamp = timeInteg
		dp.numOfElems = float64(elem.Nrow())

		values = append(values, dp)
	}

	sort.SliceStable(values, func(i, j int) bool {
		return values[i].timeStamp < values[j].timeStamp
	})

	return values
}

func createHistogramPlot(path string, histogramName string, data DataForPlotting) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create %s.png file: %v", path, err)
	}
	defer f.Close()

	p := plot.New()
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Number of Pods"

	addHistogram(histogramName, p, data)

	wt, err := p.WriterTo(1024, 512, "png")
	if err != nil {
		return fmt.Errorf("could not create writer: %v", err)
	}

	_, err = wt.WriteTo(f)
	if err != nil {
		return fmt.Errorf("could not write plot to file: %v", err)
	}

	return nil
}

func addHistogram(histogramName string, p *plot.Plot, dataPoints DataForPlotting) {
	pxys := make(plotter.XYs, len(dataPoints))

	binsNumber := 200

	for j, elem := range dataPoints {
		pxys[j].X = elem.timeStamp
	}
	for i, elem := range dataPoints {
		pxys[i].Y = elem.numOfElems
	}

	hist, err := plotter.NewHistogram(pxys, binsNumber)
	if err != nil {
		log.Fatalf("could not add new line %s: %v", histogramName, err)
	}
	p.Add(hist)
	p.Title.Text = "\"" + histogramName + "\" histogram"
}

func initFlags() {
	flag.StringVar(&filepath, "filepath", "", "Specify path to the timeline file ")
	flag.StringVar(&plots, "plots", "all", "Specify types of plots, separate with ',' ")
	flag.StringVar(&additional, "additional", "", "Specify additional parameteres for plotting, separate with ',' ")
	flag.Parse()
}

func main() {
	displayHeader()
	initFlags()

	data, err := parseDataFile(filepath)
	if err != nil {
		log.Fatalf("Could not read file %s: %v", path_to_file, err)
	}

	plotTypeSelection(plots, data)

}
