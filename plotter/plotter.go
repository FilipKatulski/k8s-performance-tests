package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"gonum.org/v1/plot"
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

func removeDuplicateInt(intSlice []int) []int {
	allKeys := make(map[int]bool)
	list := []int{}
	for _, item := range intSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func plotTimeline(dat []TimelineData, PodStateFilterSelector string) {
	dataDf := dataframe.LoadStructs(dat)
	dataDf = dataDf.Arrange(dataframe.Sort("FromUnix"))
	dataDf = dataDf.Filter(
		dataframe.F{Colname: "PodStateFilter", Comparator: series.Eq, Comparando: PodStateFilterSelector},
	)
	fmt.Println(dataDf)

	selectedDf := dataDf.Filter(
		dataframe.F{Colname: "Transition", Comparator: series.Eq, Comparando: "{create schedule 0s}"},
	)

	selectedDf = selectedDf.Select([]string{"FromUnix"})

	minimalVal := selectedDf.Elem(0, 0)
	fmt.Println(minimalVal, reflect.TypeOf(minimalVal))

	s := selectedDf.Rapply(func(s series.Series) series.Series {
		deposit, err := s.Elem(0).Int()
		if err != nil {
			return series.Ints("NAN")
		}
		withdrawal, err := minimalVal.Int()
		if err != nil {
			return series.Ints("NAN")
		}
		return series.Ints(deposit - withdrawal)
	})

	df := selectedDf.Mutate(s.Col("X0")).Rename("time_diff", "X0")

	fmt.Println(df)

	groupedDf := df.GroupBy("FromUnix")

	groups := groupedDf.GetGroups()
	fmt.Println("GROUPS:")
	fmt.Println(groups)

	var values [2][]int

	for _, elem := range groups {
		elemInteg, _ := elem.Elem(0, 0).Int()
		fmt.Println(elem.Elem(0, 0), reflect.TypeOf(elem.Elem(0, 0)))
		fmt.Println(elemInteg, reflect.TypeOf(elemInteg))
		fmt.Println(elem.Nrow())
		timeInteg, _ := elem.Elem(0, 1).Int()
		values[0] = append(values[0], timeInteg)
		values[1] = append(values[1], elem.Nrow())
	}

	fmt.Println(values[0])
	fmt.Println(values[1])
	fmt.Println(values)

	err := plotData("created.png", values)
	if err != nil {
		log.Fatalf("could not plot the data: %v", err)
	}
}

func plotHistograms(dat []TimelineData) {
	fmt.Println("TODO: Histograms")
}

func plotData(path string, xy [2][]int) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create %s.png file: %v", path, err)
	}
	defer f.Close()

	p := plot.New()

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
