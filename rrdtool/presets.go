package rrdtool

type GraphPreset struct {
	DataType     string
	Width        uint
	Height       uint
	VarsNames    []string
	HeadFormat   string
	LegendFormat string
	PrintArg     []string
	Options      []string
	IsArea       bool
}

var DiskPreset = GraphPreset{
	DataType: "DERIVE",
	Width:    600,
	Height:   300,
	IsArea:   false,
	Options: []string{
		"--font", "TITLE:12:",
		"--font", "AXIS:8:",
		"--font", "LEGEND:10:Courier New",
		"--font", "UNIT:8:",
	},
	PrintArg: []string{
		"AVERAGE: %5.2lf %s\\t",
		"MAX: %5.2lf %s\\t",
		"MIN: %5.2lf %s\\n",
	},
	HeadFormat:   "% 27s%19s%19s",
	LegendFormat: "%-15s",
}

var NetworkPreset = GraphPreset{
	DataType: "DERIVE",
	Width:    600,
	Height:   300,
	IsArea:   false,
	Options: []string{
		"--font", "TITLE:12:",
		"--font", "AXIS:8:",
		"--font", "LEGEND:10:Courier New",
		"--font", "UNIT:8:",
	},
	PrintArg: []string{
		"AVERAGE: %6.2lf %s\\t",
		"MAX: %6.2lf %s\\t",
		"MIN: %6.2lf %s\\n",
	},
	HeadFormat:   "% 27s%19s%19s",
	LegendFormat: "%-15s",
}

var SystemPreset = GraphPreset{
	DataType: "DERIVE",
	Width:    600,
	Height:   300,
	IsArea:   false,
	Options: []string{
		"--font", "TITLE:12:",
		"--font", "AXIS:8:",
		"--font", "LEGEND:10:Courier New",
		"--font", "UNIT:8:",
	},
	PrintArg: []string{
		"AVERAGE: %6.2lf %s\\t",
		"MAX: %6.2lf %s\\t",
		"MIN: %6.2lf %s\\n",
	},
	HeadFormat:   "% 37s%18s%18s",
	LegendFormat: "%-23s",
}
