package rrdtool

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ziutek/rrd"
)

const step = 1

func FileSearcher(dirname string) (files []string, err error) {

	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		log.Fatalf("FileSearcher error: %s", err)
	}

	err = filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if filepath.Ext(path) == ".rrd" {
				files = append(files, path)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		err = errors.New("No files in directory")
		return nil, err
	}

	return files, nil
}

func RRDFetch(dbfile string) rrd.FetchResult {
	inf, err := rrd.Info(dbfile)
	if err != nil {
		log.Fatal(err)
	}

	end := time.Unix(int64(inf["last_update"].(uint)), 0)
	start := end.Add(-20 * step * time.Second)

	fetchRes, err := rrd.Fetch(dbfile, "AVERAGE", start, end, step*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	return fetchRes
}

func RRDProcess(dbfiles []string, gpreset *GraphPreset, config Config, output string) {

	var fetchResults []rrd.FetchResult
	for idx, dbfile := range dbfiles {
		fetchRes := RRDFetch(dbfile)
		fetchResults = append(fetchResults, fetchRes)
		defer fetchRes.FreeValues()
		if !config.MultipleSources {
			BuildGraph(fetchRes, idx, gpreset, config, output)
		}
	}
	if config.MultipleSources {
		BuildGraphMulty(fetchResults, gpreset, config, output)

	}
}

func BuildGraph(fetchres rrd.FetchResult, fileidx int, gpreset *GraphPreset, config Config, output string) {
	base := filepath.Base(fetchres.Filename)
	ext := filepath.Ext(fetchres.Filename)
	name := base[0 : len(base)-len(ext)]
	dir := filepath.Base(filepath.Dir(fetchres.Filename))

	for _, dsname := range fetchres.DsNames {
		varname := strings.ToUpper(dsname)
		gpreset.VarsNames = append(gpreset.VarsNames, varname)
	}
	//Graph
	grapher := rrd.NewGrapher()
	grapher.SetVLabel(config.VLabel)
	grapher.SetSize(gpreset.Width, gpreset.Height)
	grapher.SetRigid()
	grapher.SetSlopeMode()
	grapher.SetAltAutoscaleMax()
	grapher.Comment(fmt.Sprintf(gpreset.HeadFormat, "Average", "Maximum", "Minimum"))

	if len(gpreset.VarsNames) == 0 {
		err := errors.New("No elements in VarsNames!")
		log.Fatalf("%s", err)
	}

	for idx, dsname := range fetchres.DsNames {
		grapher.Def(gpreset.VarsNames[2*fileidx+idx], fetchres.Filename, dsname, "AVERAGE")
		grapher.AddOptions(gpreset.Options...)
		varsnames := fmt.Sprintf(gpreset.LegendFormat, gpreset.VarsNames[2*fileidx+idx])
		if gpreset.IsArea {
			grapher.Area(gpreset.VarsNames[2*fileidx+idx], GraphPallete[idx], varsnames)
		}
		grapher.Line(2, gpreset.VarsNames[2*fileidx+idx], GraphPallete[idx], varsnames)
		for _, parg := range gpreset.PrintArg {
			grapher.GPrint(gpreset.VarsNames[2*fileidx+idx], parg)
		}
	}
	var filename string
	filename = name

	if len(config.BadPrefixes) != 0 {
		filename = RemovePrefix(name, config.BadPrefixes)
	}

	now := time.Now()
	timePeriods := GetTimePeriods(now)

	for key, period := range timePeriods {
		title := dir + "-" + filename + "-" + key
		grapher.SetTitle(title)
		_, buf, err := grapher.Graph(period, now)
		if err != nil {
			log.Fatal(err)
		}
		graphfile := output + title + ".png"
		err = ioutil.WriteFile(graphfile, buf, 0666)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func BuildGraphMulty(fetchresults []rrd.FetchResult, gpreset *GraphPreset, config Config, output string) {
	//Graph
	grapher := rrd.NewGrapher()
	grapher.SetVLabel(config.VLabel)
	grapher.SetSize(gpreset.Width, gpreset.Height)
	grapher.SetRigid()
	grapher.SetSlopeMode()
	grapher.SetAltAutoscaleMax()
	grapher.Comment(fmt.Sprintf(gpreset.HeadFormat, "Average", "Maximum", "Minimum"))

	var cdef string
	var dir string
	for fetchidx, fetchres := range fetchresults {
		for _, dsname := range fetchres.DsNames {
			var varname string
			base := filepath.Base(fetchres.Filename)
			ext := filepath.Ext(fetchres.Filename)
			name := base[0 : len(base)-len(ext)]
			dir = filepath.Base(filepath.Dir(fetchres.Filename))
			varname = strings.ToUpper(name)
			if len(config.BadPrefixes) != 0 {
				trimmedname := RemovePrefix(name, config.BadPrefixes)
				varname = strings.ToUpper(trimmedname)
			}
			gpreset.VarsNames = append(gpreset.VarsNames, varname)

			grapher.Def(gpreset.VarsNames[fetchidx], fetchres.Filename, dsname, "AVERAGE")
			grapher.AddOptions(gpreset.Options...)
			varsnames := fmt.Sprintf(gpreset.LegendFormat, gpreset.VarsNames[fetchidx])
			if gpreset.IsArea {
				grapher.Area(gpreset.VarsNames[fetchidx], GraphPallete[fetchidx], varsnames)
			}
			grapher.Line(2, gpreset.VarsNames[fetchidx], GraphPallete[fetchidx], varsnames)

			for _, parg := range gpreset.PrintArg {
				grapher.GPrint(gpreset.VarsNames[fetchidx], parg)
			}

			if fetchidx == 0 {
				cdef = gpreset.VarsNames[0]
				break
			} else {
				cdef = cdef + "," + gpreset.VarsNames[fetchidx] + ",+"
			}
		}
	}

	grapher.CDef(dir, cdef)

	now := time.Now()
	timePeriods := GetTimePeriods(now)
	for key, period := range timePeriods {
		grapher.SetTitle(dir + "-" + key)
		_, buf, err := grapher.Graph(period, now)
		if err != nil {
			log.Fatal(err)
		}
		graphfile := output + dir + "-" + key + ".png"
		err = ioutil.WriteFile(graphfile, buf, 0666)
		if err != nil {
			log.Fatal(err)
		}
	}
}
