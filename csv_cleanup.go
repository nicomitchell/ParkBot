package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	r, err := os.Open("VEHICLES_NOPERMIT_CLEANED.csv")
	if err == nil {
		// readFromCSVFile(bufio.NewReader(r))
		var objs []*licenseObj
		objs, err = readFromCSVFile(bufio.NewReader(r))
		if err == nil {
			fmt.Println(len(objs))
			for _, v := range objs {
				fmt.Println(v.state + " " + v.plateNumber + "\t" + strconv.Itoa(v.vehicleUID))
				for _, p := range v.permits {
					fmt.Println("\t\t\t" + p)
				}
			}
		} else {
			fmt.Println(err.Error())
		}
		r.Close()
	}
}

type licenseObj struct {
	vehicleUID  int
	state       string
	plateNumber string
	permits     []string
}

func readFromCSVFile(r *bufio.Reader) ([]*licenseObj, error) {
	var out []*licenseObj
	line, err := r.ReadBytes('\n')
	for {
		if err == nil {
			values := getValuesFromLine(string(line))
			var vuid int
			vuid, err = strconv.Atoi(values[0])
			if err == nil {
				obj := &licenseObj{
					vehicleUID:  vuid,
					state:       values[1][:2],
					plateNumber: values[1][3:],
					permits:     []string{},
				}
				for {
					line, err = r.ReadBytes('\n')
					if err == nil {
						values = getValuesFromLine(string(line))
						if len(values) > 1 {
							break
						} else {
							obj.permits = append(obj.permits, values[0])
						}
					} else {
						break
					}
				}
				if err == nil {
					out = append(out, obj)
					continue
				}
			}
		}
		if err != io.EOF {
			return nil, err
		}
		break
	}
	return out, nil
}

func getValuesFromLine(line string) []string {
	line = strings.Trim(line, ",\n")
	for strings.Contains(line, ",,") {
		line = strings.Replace(line, ",,", ",", -1)
	}
	split := strings.Split(line, ",")
	if len(split) == 1 {
		return split
	}
	return split[:2]
}
