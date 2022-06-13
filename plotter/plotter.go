package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	//"io/ioutil"
	//"path"
	//"time"
	//"gonum.org/v1/plot"
	//"gonum.org/v1/plot/plotter"
	//"gonum.org/v1/plot/vg"
)

var (
	datafile     string
	path_to_file string
	filepath     string
	plots        string
)

type TimelineData struct {
	Name              string
	Transition        string
	Namespace_PodName string
	NodeName          string
	pod_state_filter  string
	diff              int
	from_unix         int
	to_unix           int
}

func displayHelp() {

}

func displayHeader() {
	//TODO Make it more fancy
	fmt.Println("Plotter")
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
		return nil, fmt.Errorf("Could not scan: %v", err)
	}

	for _, eachline := range textlines[1:] {
		split := strings.Split(eachline, ", ")

		diff_int, err := strconv.Atoi(split[5])
		if err != nil {
			return nil, fmt.Errorf("Could not parse 'diff' string value to integer in csv file: %v", err)
		}
		from_unix_int, err := strconv.Atoi(split[6])
		if err != nil {
			return nil, fmt.Errorf("Could not parse 'from_unix' string value to integer in csv file: %v", err)
		}
		to_unix_int, err := strconv.Atoi(split[7])
		if err != nil {
			return nil, fmt.Errorf("Could not parse 'to_unix' string value to integer in csv file: %v", err)
		}

		line := TimelineData{
			Name:              split[0],
			Transition:        split[1],
			Namespace_PodName: split[2],
			NodeName:          split[3],
			pod_state_filter:  split[4],
			diff:              diff_int,
			from_unix:         from_unix_int,
			to_unix:           to_unix_int,
		}

		timelineData = append(timelineData, line)
	}

	return timelineData, nil
}

func plotTypeSelection(plotlist string) []string {
	plots := strings.Split(plotlist, ",")
	for _, t := range plots {
		switch t {
		case "all":
			fmt.Println("Implement all")
		case "histograms":
			fmt.Println("Implement histograms")
		case "timeline":
			fmt.Println("Implement timeline")
		default:
			fmt.Printf("Plot type '%s' is not implemented.\n", t)
		}
	}

	return plots
}

func plotTimeline([]TimelineData) {

}

func plotHistograms([]TimelineData) {

}

func initFlags() {
	flag.StringVar(&filepath, "filepath", "", "Specify path to the timeline file. ")
	flag.StringVar(&plots, "plots", "", "Specify types of plots, separate by ',' ")
	flag.Parse()
}

func main() {
	displayHeader()
	initFlags()
	plts := plotTypeSelection(plots)
	data, err := parseDataFile(filepath)
	if err != nil {
		log.Fatalf("Could not read file %s: %v", path_to_file, err)
	}

	fmt.Println(data)
	fmt.Println(plts)

}
