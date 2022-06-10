package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	//"io/ioutil"
	//"os"
	//"path"
	//"time"
	//"gonum.org/v1/plot"
	//"gonum.org/v1/plot/plotter"
	//"gonum.org/v1/plot/vg"
)

var (
	datafile     string
	path_to_file string
)

func displayHelp() {

}

func displayHeader() {

}

type xy struct{ x, y string }

type timeline_data struct {
	Name              string
	Transition        string
	Namespace_PodName string
	NodeName          string
	pod_state_filter  string
	// TODO check if types are OK:
	diff      int
	from_unix int
	to_unix   int
}

func parseDataFile(path string) ([]timeline_data, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var textlines []string

	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		textlines = append(textlines, s.Text())
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("Could not scan: %v", err)
	}

	for _, eachline := range textlines {
		fmt.Println(eachline)
	}

	return nil, nil
}

func plotTypeSelection() {

}

func parseCommandLineArgs() {

}

func initFlags() {

}

func main() {
	fmt.Println("Plotter")
	path_to_file = "../clusterloader2/timeline.csv"
	xy, err := parseDataFile(path_to_file)
	if err != nil {
		log.Fatalf("Could not read file %s: %v", path_to_file, err)
	}
	fmt.Println(xy)

}
