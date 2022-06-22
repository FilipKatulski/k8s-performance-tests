package main

import (
	"bufio"
	"flag"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"
	"reflect"
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
			fmt.Println("Implement all")
			plotTimeline(data, "Stateless")
			plotHistograms(data)
			break out
		case "histograms":
			fmt.Println("Implement histograms")
			plotHistograms(data)
		case "timeline":
			fmt.Println("Implement timeline")
			plotTimeline(data, "Stateless")
		default:
			fmt.Printf("Plot type '%s' is not implemented.\n", t)
		}
	}
}

type DataForPlotting struct {
	timeStamp  []float64
	numOfElems []float64
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
	createdMinimalVal, _ := CreatedDf.Elem(0, 0).Int()

	fmt.Println("createdDf", reflect.TypeOf(CreatedDf))
	fmt.Println("createdMinimalVal", reflect.TypeOf(createdMinimalVal))

	ScheduledDf := dataDf.Filter(
		dataframe.F{Colname: "Transition", Comparator: series.Eq, Comparando: "{create schedule 0s}"},
	)
	ScheduledDf = ScheduledDf.Select([]string{"ToUnix"})
	ScheduledDf = ScheduledDf.Arrange(dataframe.Sort("ToUnix"))
	scheduledMinimalVal, _ := ScheduledDf.Elem(0, 0).Int()
	fmt.Println(ScheduledDf)
	fmt.Println("MinimalvalScheduled: ", scheduledMinimalVal)

	RunDf := dataDf.Filter(
		dataframe.F{Colname: "Transition", Comparator: series.Eq, Comparando: "{schedule run 0s}"},
	)
	RunDf = RunDf.Select([]string{"ToUnix"})
	RunDf = RunDf.Arrange(dataframe.Sort("ToUnix"))
	runMinimalVal, _ := RunDf.Elem(0, 0).Int()

	WatchedDf := dataDf.Filter(
		dataframe.F{Colname: "Transition", Comparator: series.Eq, Comparando: "{run watch 0s}"},
	)
	WatchedDf = WatchedDf.Select([]string{"ToUnix"})
	WatchedDf = WatchedDf.Arrange(dataframe.Sort("ToUnix"))
	watchedMinimalVal, _ := WatchedDf.Elem(0, 0).Int()

	createdGroups := parseDf(&CreatedDf, createdMinimalVal, "FromUnix")
	schedulerGroups := parseDf(&ScheduledDf, scheduledMinimalVal, "ToUnix")
	runGroups := parseDf(&RunDf, runMinimalVal, "ToUnix")
	watchGroups := parseDf(&WatchedDf, watchedMinimalVal, "ToUnix")

	createdValues := createDataForPlotting(createdGroups)
	scheduledValues := createDataForPlotting(schedulerGroups)
	runValues := createDataForPlotting(runGroups)
	watchValues := createDataForPlotting(watchGroups)

	err := createTimelinePlot("timeline.png", createdValues, scheduledValues, runValues, watchValues)
	if err != nil {
		log.Fatalf("could not plot the data: %v", err)
	}
}

func parseDf(df *dataframe.DataFrame, minVal int, groupingCol string) map[string]dataframe.DataFrame {

	s := df.Rapply(func(s series.Series) series.Series {
		deposit, err := s.Elem(0).Int()
		if err != nil {
			return series.Ints("NAN")
		}
		return series.Ints(deposit - minVal)
	})

	mutatedDf := df.Mutate(s.Col("X0")).Rename("Z", "X0")

	fmt.Println(mutatedDf)
	groupedDf := mutatedDf.GroupBy(groupingCol)

	groups := groupedDf.GetGroups()
	fmt.Println("groups: ", groups)
	return groups
}

func createDataForPlotting(groups map[string]dataframe.DataFrame) DataForPlotting {
	var values DataForPlotting

	for _, elem := range groups {
		timeInteg := elem.Elem(0, 1).Float()
		values.timeStamp = append(values.timeStamp, timeInteg)
		values.numOfElems = append(values.numOfElems, float64(elem.Nrow()))
	}
	for i := 1; i < len(values.numOfElems); i++ {
		values.numOfElems[i] += values.numOfElems[i-1]
	}

	fmt.Println(values.timeStamp)
	fmt.Println(values.numOfElems)

	return values
}

func addPlotLine(lineName string, p *plot.Plot, dataPoints DataForPlotting) {

	pxys := make(plotter.XYs, len(dataPoints.numOfElems))

	for j, elem := range dataPoints.timeStamp {
		pxys[j].X = elem
	}
	for i, elem := range dataPoints.numOfElems {
		fmt.Println(elem)
		pxys[i].Y = elem
	}

	fmt.Println("PXYS", pxys)

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

	p := plot.New()
	p.Title.Text = "Timeline"
	p.X.Label.Text = "Time [s]"
	p.Y.Label.Text = "number of Pods"

	addPlotLine("Created", p, created)
	addPlotLine("Scheduled", p, scheduled)
	addPlotLine("Run", p, run)
	addPlotLine("Watch", p, watch)

	wt, err := p.WriterTo(512, 512, "png")
	if err != nil {
		return fmt.Errorf("could not create writer: %v", err)
	}

	_, err = wt.WriteTo(f)
	if err != nil {
		return fmt.Errorf("could not write plot to file: %v", err)
	}

	return nil
}

func plotHistograms(dat []TimelineData) {
	fmt.Println("TODO: Histograms")
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
