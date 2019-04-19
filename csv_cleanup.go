package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	r, err := os.Open("VEHICLES_NOPERMIT_CLEANED.csv")
	if err == nil {
		readFromCSVFile(bufio.NewReader(r))
		//var objs []licenseObj
		//objs, err = readFromCSVFile(bufio.NewReader(r))
		// for _, v := range objs {
		// 	fmt.Println(v.state + " " + v.plateNumber + "\t" + string(v.vehicleUID))
		// 	fmt.Println("\t" + strings.Join(v.permits, " "))
		// }
	}
}

type licenseObj struct {
	vehicleUID  int
	state       string
	plateNumber string
	permits     []string
}

func readFromCSVFile(r *bufio.Reader) { //([]licenseObj, error) {
	for {
		line, err := r.ReadBytes('\n')
		if err == nil {
			values := getValuesFromLine(string(line))
			fmt.Println(len(values), "\t", strings.Join(values, ", "))
		} else {
			return
		}
	}
}

func getValuesFromLine(line string) []string {
	line = strings.Trim(line, ",\n")
	for strings.Contains(line, ",,") {
		line = strings.Replace(line, ",,", ",", -1)
	}
	fmt.Println(line)
	split := strings.Split(line, ",")
	if len(split) == 1 {
		return split
	}
	return split[:2]

}
