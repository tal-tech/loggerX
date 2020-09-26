// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log4go

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type xmlProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

type xmlFilter struct {
	Enabled  string        `xml:"enabled,attr"`
	Tag      string        `xml:"tag"`
	Level    string        `xml:"level"`
	Type     string        `xml:"type"`
	Property []xmlProperty `xml:"property"`
}

type xmlLoggerConfig struct {
	Filter []xmlFilter `xml:"filter"`
}

type LogConfig struct {
	//存储路径
	LogPath string
	//日志级别
	Level string
	//日志标签 多日志时使用
	Tag string
	//日志格式
	Format string
	//最大行数切割
	RotateLines string
	//最大容量切割
	RotateSize string
	//按日期切割
	RotateHourly bool
	//是否启用切割
	Rotate bool
	//保留时间，day
	Retention string
	//Console输出
	Console bool
}

func (log *LogConfig) SetConfigMap(conf map[string]string) {
	for k, v := range conf {
		switch k {
		case "logPath":
			log.LogPath = v
		case "level":
			log.Level = v
		case "tag":
			log.Tag = v
		case "format":
			log.Format = v
		case "rotateLines":
			log.RotateLines = v
		case "rotateSize":
			log.RotateSize = v
		case "rotateHourly":
			log.RotateHourly = (v == "true")
		case "rotate":
			log.Rotate = (v != "false")
		case "retention":
			log.Retention = v
		case "console":
			log.Console = (v == "true")
		}
	}
}

// Use LogConfig
func (log Logger) LoadLogConfig(config *LogConfig) {
	log.Close()

	var lvl Level
	switch config.Level {
	case "FINEST":
		lvl = FINEST
	case "FINE":
		lvl = FINE
	case "DEBUG":
		lvl = DEBUG
	case "TRACE":
		lvl = TRACE
	case "INFO":
		lvl = INFO
	case "WARNING":
		lvl = WARNING
	case "ERROR":
		lvl = ERROR
	case "CRITICAL":
		lvl = CRITICAL
	default:
		fmt.Fprintf(os.Stderr, "LoadLogConfig: Error: Required child <%s> for filter has unknown value %s\n", "level", config.Level)
		os.Exit(1)
	}
	file := strings.Trim(config.LogPath, " \r\n")
	format := strings.Trim(config.Format, " \r\n")
	maxlines := strToNumSuffix(strings.Trim(config.RotateLines, " \r\n"), 1000)
	maxsize := strToNumSuffix(strings.Trim(config.RotateSize, " \r\n"), 1024)
	hourly := config.RotateHourly
	rotate := config.Rotate
	retention, _ := strconv.Atoi(config.Retention)

	flw := NewFileLogTraceWriter(file, rotate)
	flw.SetFormat(format)
	flw.SetRotateLines(maxlines)
	flw.SetRotateSize(maxsize)
	flw.SetRotateHourly(hourly)
	flw.SetRetention(retention)
	log[config.Tag] = &Filter{lvl, flw}
	if config.Console {
		filt, _ := xmlToConsoleLogWriter("defaultConsole", []xmlProperty{}, true)
		log["Console"] = &Filter{lvl, filt}
	}
	return
}

// Load XML configuration; see examples/example.xml for documentation
func (log Logger) LoadConfiguration(filename string) {
	log.Close()

	// Open the configuration file
	fd, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not open %q for reading: %s\n", filename, err)
		os.Exit(1)
	}

	contents, err := ioutil.ReadAll(fd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not read %q: %s\n", filename, err)
		os.Exit(1)
	}

	xc := new(xmlLoggerConfig)
	if err := xml.Unmarshal(contents, xc); err != nil {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not parse XML configuration in %q: %s\n", filename, err)
		os.Exit(1)
	}

	for _, xmlfilt := range xc.Filter {
		var filt LogWriter
		var lvl Level
		bad, good, enabled := false, true, false

		// Check required children
		if len(xmlfilt.Enabled) == 0 {
			fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Required attribute %s for filter missing in %s\n", "enabled", filename)
			bad = true
		} else {
			enabled = xmlfilt.Enabled != "false"
		}
		if len(xmlfilt.Tag) == 0 {
			fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Required child <%s> for filter missing in %s\n", "tag", filename)
			bad = true
		}
		if len(xmlfilt.Type) == 0 {
			fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Required child <%s> for filter missing in %s\n", "type", filename)
			bad = true
		}
		if len(xmlfilt.Level) == 0 {
			fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Required child <%s> for filter missing in %s\n", "level", filename)
			bad = true
		}

		switch xmlfilt.Level {
		case "FINEST":
			lvl = FINEST
		case "FINE":
			lvl = FINE
		case "DEBUG":
			lvl = DEBUG
		case "TRACE":
			lvl = TRACE
		case "INFO":
			lvl = INFO
		case "WARNING":
			lvl = WARNING
		case "ERROR":
			lvl = ERROR
		case "CRITICAL":
			lvl = CRITICAL
		default:
			fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Required child <%s> for filter has unknown value in %s: %s\n", "level", filename, xmlfilt.Level)
			bad = true
		}

		// Just so all of the required attributes are errored at the same time if missing
		if bad {
			os.Exit(1)
		}

		switch xmlfilt.Type {
		case "console":
			filt, good = xmlToConsoleLogWriter(filename, xmlfilt.Property, enabled)
		case "file":
			filt, good = xmlToFileLogWriter(filename, xmlfilt.Property, enabled)
		case "file2":
			filt, good = xmlToFileLogWriter2(filename, xmlfilt.Property, enabled)
		case "file_trace":
			filt, good = xmlToFileLogWriterTrace(filename, xmlfilt.Property, enabled)
		case "xml":
			filt, good = xmlToXMLLogWriter(filename, xmlfilt.Property, enabled)
		case "socket":
			filt, good = xmlToSocketLogWriter(filename, xmlfilt.Property, enabled)
		default:
			fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not load XML configuration in %s: unknown filter type \"%s\"\n", filename, xmlfilt.Type)
			os.Exit(1)
		}

		// Just so all of the required params are errored at the same time if wrong
		if !good {
			os.Exit(1)
		}

		// If we're disabled (syntax and correctness checks only), don't add to logger
		if !enabled {
			continue
		}

		log[xmlfilt.Tag] = &Filter{lvl, filt}
	}
}

func xmlToConsoleLogWriter(filename string, props []xmlProperty, enabled bool) (*ConsoleLogWriter, bool) {

	format := "[%D %T] [%L] (%S) %M"

	// Parse properties
	for _, prop := range props {
		switch prop.Name {
		case "format":
			format = strings.Trim(prop.Value, " \r\n")
		default:
			fmt.Fprintf(os.Stderr, "LoadConfiguration: Warning: Unknown property \"%s\" for console filter in %s\n", prop.Name, filename)
		}
	}

	// If it's disabled, we're just checking syntax
	if !enabled {
		return nil, true
	}

	clw := NewConsoleLogWriter()
	clw.SetFormat(format)

	return clw, true
}

// Parse a number with K/M/G suffixes based on thousands (1000) or 2^10 (1024)
func strToNumSuffix(str string, mult int) int {
	num := 1
	if len(str) > 1 {
		switch str[len(str)-1] {
		case 'G', 'g':
			num *= mult
			fallthrough
		case 'M', 'm':
			num *= mult
			fallthrough
		case 'K', 'k':
			num *= mult
			str = str[0 : len(str)-1]
		}
	}
	parsed, _ := strconv.Atoi(str)
	return parsed * num
}
func xmlToFileLogWriter(filename string, props []xmlProperty, enabled bool) (*FileLogWriter, bool) {
	file := ""
	format := "[%D %T] [%L] (%S) %M"
	maxlines := 0
	maxsize := 0
	daily := false
	rotate := false

	// Parse properties
	for _, prop := range props {
		switch prop.Name {
		case "filename":
			file = strings.Trim(prop.Value, " \r\n")
		case "format":
			format = strings.Trim(prop.Value, " \r\n")
		case "maxlines":
			maxlines = strToNumSuffix(strings.Trim(prop.Value, " \r\n"), 1000)
		case "maxsize":
			maxsize = strToNumSuffix(strings.Trim(prop.Value, " \r\n"), 1024)
		case "daily":
			daily = strings.Trim(prop.Value, " \r\n") != "false"
		case "rotate":
			rotate = strings.Trim(prop.Value, " \r\n") != "false"
		default:
			fmt.Fprintf(os.Stderr, "LoadConfiguration: Warning: Unknown property \"%s\" for file filter in %s\n", prop.Name, filename)
		}
	}

	// Check properties
	if len(file) == 0 {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Required property \"%s\" for file filter missing in %s\n", "filename", filename)
		return nil, false
	}

	// If it's disabled, we're just checking syntax
	if !enabled {
		return nil, true
	}

	flw := NewFileLogWriter(file, rotate)
	flw.SetFormat(format)
	flw.SetRotateLines(maxlines)
	flw.SetRotateSize(maxsize)
	flw.SetRotateDaily(daily)
	return flw, true
}
func xmlToFileLogWriter2(filename string, props []xmlProperty, enabled bool) (*FileLog2Writer, bool) {
	file := ""
	format := "[%D %T] [%L] (%S) %M"
	maxlines := 0
	maxsize := 0
	daily := false
	rotate := false

	// Parse properties
	for _, prop := range props {
		switch prop.Name {
		case "filename":
			file = strings.Trim(prop.Value, " \r\n")
		case "format":
			format = strings.Trim(prop.Value, " \r\n")
		case "maxlines":
			maxlines = strToNumSuffix(strings.Trim(prop.Value, " \r\n"), 1000)
		case "maxsize":
			maxsize = strToNumSuffix(strings.Trim(prop.Value, " \r\n"), 1024)
		case "daily":
			daily = strings.Trim(prop.Value, " \r\n") != "false"
		case "rotate":
			rotate = strings.Trim(prop.Value, " \r\n") != "false"
		default:
			fmt.Fprintf(os.Stderr, "LoadConfiguration: Warning: Unknown property \"%s\" for file filter in %s\n", prop.Name, filename)
		}
	}

	// Check properties
	if len(file) == 0 {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Required property \"%s\" for file filter missing in %s\n", "filename", filename)
		return nil, false
	}

	// If it's disabled, we're just checking syntax
	if !enabled {
		return nil, true
	}

	flw := NewFileLog2Writer(file, rotate)
	flw.SetFormat(format)
	flw.SetRotateLines(maxlines)
	flw.SetRotateSize(maxsize)
	flw.SetRotateDaily(daily)
	return flw, true
}

func xmlToFileLogWriterTrace(filename string, props []xmlProperty, enabled bool) (*FileLogTraceWriter, bool) {
	file := ""
	format := "[%D %T] [%L] (%S) %M"
	maxlines := 0
	maxsize := 0
	hourly := false
	rotate := false
	retention := 0

	// Parse properties
	for _, prop := range props {
		switch prop.Name {
		case "filename":
			file = strings.Trim(prop.Value, " \r\n")
		case "format":
			format = strings.Trim(prop.Value, " \r\n")
		case "maxlines":
			maxlines = strToNumSuffix(strings.Trim(prop.Value, " \r\n"), 1000)
		case "maxsize":
			maxsize = strToNumSuffix(strings.Trim(prop.Value, " \r\n"), 1024)
		case "hourly":
			hourly = strings.Trim(prop.Value, " \r\n") != "false"
		case "rotate":
			rotate = strings.Trim(prop.Value, " \r\n") != "false"
		case "retention":
			retention, _ = strconv.Atoi(strings.Trim(prop.Value, " \r\n"))
		default:
			fmt.Fprintf(os.Stderr, "LoadConfiguration: Warning: Unknown property \"%s\" for file filter in %s\n", prop.Name, filename)
		}
	}

	// Check properties
	if len(file) == 0 {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Required property \"%s\" for file filter missing in %s\n", "filename", filename)
		return nil, false
	}

	// If it's disabled, we're just checking syntax
	if !enabled {
		return nil, true
	}

	file = tryUpdateFilename(file)
	flw := NewFileLogTraceWriter(file, rotate)
	flw.SetFormat(format)
	flw.SetRotateLines(maxlines)
	flw.SetRotateSize(maxsize)
	flw.SetRotateHourly(hourly)
	flw.SetRetention(retention)
	return flw, true
}

func tryUpdateFilename(filename string) string {
	exePath, err := os.Executable()
	if err != nil {
		return filename
	}

	//  /root/.../cwisciapi.xesv5.com/bin/...
	if len(exePath) <= len("/bin/kafkaWorker") {
		return filename
	}

	binIndex := strings.LastIndex(exePath, "/bin")
	if binIndex == -1 {
		return filename
	}

	exePath = exePath[0:binIndex]
	index := strings.LastIndex(exePath, "/")
	if index == -1 {
		return filename
	}

	suffix := "-" + exePath[index+1:]

	dotIndex := strings.LastIndex(filename, ".log")
	if dotIndex == -1 {
		return filename + suffix
	}

	if dotIndex > 1 && (filename[dotIndex-2:dotIndex] == "_%") {
		return filename[0:dotIndex-2] + suffix + ".log"
	}
	return filename
}

func xmlToXMLLogWriter(filename string, props []xmlProperty, enabled bool) (*FileLogWriter, bool) {
	file := ""
	maxrecords := 0
	maxsize := 0
	daily := false
	rotate := false

	// Parse properties
	for _, prop := range props {
		switch prop.Name {
		case "filename":
			file = strings.Trim(prop.Value, " \r\n")
		case "maxrecords":
			maxrecords = strToNumSuffix(strings.Trim(prop.Value, " \r\n"), 1000)
		case "maxsize":
			maxsize = strToNumSuffix(strings.Trim(prop.Value, " \r\n"), 1024)
		case "daily":
			daily = strings.Trim(prop.Value, " \r\n") != "false"
		case "rotate":
			rotate = strings.Trim(prop.Value, " \r\n") != "false"
		default:
			fmt.Fprintf(os.Stderr, "LoadConfiguration: Warning: Unknown property \"%s\" for xml filter in %s\n", prop.Name, filename)
		}
	}

	// Check properties
	if len(file) == 0 {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Required property \"%s\" for xml filter missing in %s\n", "filename", filename)
		return nil, false
	}

	// If it's disabled, we're just checking syntax
	if !enabled {
		return nil, true
	}

	xlw := NewXMLLogWriter(file, rotate)
	xlw.SetRotateLines(maxrecords)
	xlw.SetRotateSize(maxsize)
	xlw.SetRotateDaily(daily)
	return xlw, true
}

func xmlToSocketLogWriter(filename string, props []xmlProperty, enabled bool) (SocketLogWriter, bool) {
	endpoint := ""
	protocol := "udp"

	// Parse properties
	for _, prop := range props {
		switch prop.Name {
		case "endpoint":
			endpoint = strings.Trim(prop.Value, " \r\n")
		case "protocol":
			protocol = strings.Trim(prop.Value, " \r\n")
		default:
			fmt.Fprintf(os.Stderr, "LoadConfiguration: Warning: Unknown property \"%s\" for file filter in %s\n", prop.Name, filename)
		}
	}

	// Check properties
	if len(endpoint) == 0 {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Required property \"%s\" for file filter missing in %s\n", "endpoint", filename)
		return nil, false
	}

	// If it's disabled, we're just checking syntax
	if !enabled {
		return nil, true
	}

	return NewSocketLogWriter(protocol, endpoint), true
}
