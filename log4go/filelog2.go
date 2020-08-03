// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log4go

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// This log writer sends output to a file
type FileLog2Writer struct {
	rec       chan *LogRecord
	out       chan *LogBuffer
	rot       chan bool
	logBuffer *LogBuffer

	// The opened file
	filename string
	file     *os.File

	// The logging format
	format string
	pieces [][]byte

	// File header/trailer
	header, trailer string

	// Rotate at linecount
	maxlines          int
	maxlines_curlines int

	// Rotate at size
	maxsize         int
	maxsize_cursize int

	// Rotate daily
	daily          bool
	daily_opendate int

	// Keep old logfiles (.001, .002, etc)
	rotate    bool
	maxbackup int
}

var g_levelMapper = map[Level]string{
	FINEST:   "finest",
	FINE:     "fine",
	DEBUG:    "debug",
	TRACE:    "trace",
	INFO:     "info",
	WARNING:  "warning",
	ERROR:    "error",
	CRITICAL: "critical",
}

// This is the FileLog2Writer's output method
func (w *FileLog2Writer) LogWrite(rec *LogRecord) {
	w.rec <- rec
}

func (w *FileLog2Writer) Close() {
	close(w.rec)
	if w.logBuffer != nil {
		w.out <- w.logBuffer
	}
	time.Sleep(10 * time.Millisecond) //wait for write to file
	close(w.out)
	w.file.Close()
	w.file.Sync()
}

// NewFileLog2Writer creates a new LogWriter which writes to the given file and
// has rotation enabled if rotate is true.
//
// If rotate is true, any time a new log file is opened, the old one is renamed
// with a .### extension to preserve it.  The various Set* methods can be used
// to configure log rotation based on lines, size, and daily.
//
// The standard log-line format is:
//   [%D %T] [%L] (%S) %M
func NewFileLog2Writer(fname string, rotate bool) *FileLog2Writer {
	w := &FileLog2Writer{
		rec:       make(chan *LogRecord, LogBufferLength),
		out:       make(chan *LogBuffer, LogBufferLength),
		rot:       make(chan bool),
		filename:  fname,
		format:    "[%G] [%L] (%S) %M",
		pieces:    bytes.Split([]byte("[%G] [%L] (%S) %M"), []byte{'%'}),
		rotate:    rotate,
		maxbackup: 999,
	}

	// open the file for the first time
	if err := w.intRotate(); err != nil {
		fmt.Fprintf(os.Stderr, "FileLog2Writer(%q): %s\n", w.filename, err)
		return nil
	}

	go func() {
		defer func() {
			if w.file != nil {
				fmt.Fprint(w.file, FormatLogRecord(w.trailer, &LogRecord{Created: time.Now()}))
			}
		}()

		for {
			select {
			case rec, ok := <-w.rec:
				if !ok {
					return
				}
				now := time.Now()
				if (w.maxlines > 0 && w.maxlines_curlines >= w.maxlines) ||
					(w.maxsize > 0 && w.maxsize_cursize >= w.maxsize) ||
					(w.daily && now.Day() != w.daily_opendate) {
					if w.logBuffer != nil {
						_, err := w.logBuffer.buf.WriteTo(w.file)
						if err != nil {
							fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.filename, err)
						}
						bufferPool.Put(w.logBuffer)
						w.logBuffer = nil
					}
					if err := w.intRotate(); err != nil {
						fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.filename, err)
						return
					}
				}

				if w.logBuffer == nil {
					w.logBuffer = bufferPool.Get().(*LogBuffer)
					w.logBuffer.init()
				}

				w.logBuffer.Encode(w.pieces, rec)
				if len(w.rec) == 0 || w.logBuffer.Flush() {
					w.out <- w.logBuffer
					w.logBuffer = nil
				}

				// Update the counts
				w.maxlines_curlines++
			}
		}
	}()

	go func() {
		for {
			select {
			case <-w.rot:
				if err := w.intRotate(); err != nil {
					fmt.Fprintf(os.Stderr, "FileLog2Writer(%q): %s\n", w.filename, err)
					return
				}
			case lb, ok := <-w.out:
				if !ok {
					return
				}
				n, err := lb.buf.WriteTo(w.file)
				if err != nil {
					fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.filename, err)
				}
				bufferPool.Put(lb)
				// Update the counts
				w.maxsize_cursize += int(n)
			}
		}
	}()
	return w
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		lb := new(LogBuffer)
		lb.buf = bytes.NewBuffer(make([]byte, 0, 4096))
		return lb
	},
}

type LogBuffer struct {
	buf      *bytes.Buffer
	inittime time.Time
}

func (lb *LogBuffer) init() {
	lb.inittime = time.Now()
}

func (lb *LogBuffer) Flush() bool {
	return time.Now().Sub(lb.inittime) > time.Millisecond
}

func (lb *LogBuffer) Encode(pieces [][]byte, rec *LogRecord) error {
	if rec == nil {
		return nil
	}
	if len(pieces) == 0 {
		return nil
	}

	secs := rec.Created.UnixNano() / 1e9

	cache := getFormatCache()
	if cache.LastUpdateSeconds != secs {
		month, day, year := rec.Created.Month(), rec.Created.Day(), rec.Created.Year()
		hour, minute, second := rec.Created.Hour(), rec.Created.Minute(), rec.Created.Second()
		zone, _ := rec.Created.Zone()
		updated := &formatCacheType{
			LastUpdateSeconds: secs,
			shortTime:         fmt.Sprintf("%02d:%02d", hour, minute),
			shortDate:         fmt.Sprintf("%02d/%02d/%02d", day, month, year%100),
			longTime:          fmt.Sprintf("%02d:%02d:%02d %s", hour, minute, second, zone),
			longDate:          fmt.Sprintf("%04d/%02d/%02d", year, month, day),
			isoTime:           fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d+08:00", year, month, day, hour, minute, second),
			unixStr:           fmt.Sprintf("%10d", rec.Created.Unix()),
		}

		cache = updated
		setFormatCache(updated)
	}

	rec.Message = strings.Replace(rec.Message, "\"", "'", -1)
	// first write the open parenthese
	lb.buf.WriteString(`{"devid":"-",`)

	// Iterate over the pieces, replacing known formats
	for i, piece := range pieces {
		if i > 0 && len(piece) > 0 {
			switch piece[0] {
			case 'T':
				lb.buf.WriteString(cache.longTime)
			case 't':
				lb.buf.WriteString(cache.shortTime)
			case 'D':
				lb.buf.WriteString(cache.longDate)
			case 'd':
				lb.buf.WriteString(cache.shortDate)
			case 'L':
				lb.buf.WriteString(`"loglevel":"`)
				lb.buf.WriteString(g_levelMapper[rec.Level])
				lb.buf.WriteString(`",`)
			case 'S':
				if ret := strings.Split(rec.Source, "\t"); len(ret) == 2 {
					lb.buf.WriteString(`"position":"`)
					lb.buf.WriteString(ret[0])
					lb.buf.WriteString(`",`)

					lb.buf.WriteString(`"host":"`)
					lb.buf.WriteString(ret[1])
					lb.buf.WriteString(`",`)
				}
			case 's':
				slice := strings.Split(rec.Source, "/")
				lb.buf.WriteString(slice[len(slice)-1])
			case 'M':
				tagValue := ""
				costValue := "0"

				costIndex := strings.LastIndex(rec.Message, "COST:")
				infoEndIndex := len(rec.Message)
				if costIndex != -1 {
					infoEndIndex = costIndex
				}

				tabIndex := strings.Index(rec.Message, "\t")
				if tabIndex != -1 {
					tagValue = rec.Message[0:tabIndex]
					if costIndex != -1 {
						costValue = rec.Message[costIndex+5:]
					}
				}

				lb.buf.WriteString(`"timestamp":`)
				lb.buf.WriteString(cache.unixStr)
				lb.buf.WriteString(`,`)
				lb.buf.WriteString(`"date":"`)
				lb.buf.WriteString(cache.longDate)
				lb.buf.WriteString(" ")
				lb.buf.WriteString(cache.longTime)
				lb.buf.WriteString(`",`)
				lb.buf.WriteString(`"tag":"`)
				lb.buf.WriteString(tagValue)
				lb.buf.WriteString(`",`)

				if tabIndex != -1 {
					lb.buf.WriteString(`"info":"`)
					lb.buf.WriteString(rec.Message[tabIndex+1 : infoEndIndex])
					lb.buf.WriteString(`",`)
				} else {
					lb.buf.WriteString(`"info":"",`)
				}
				// no more fields
				lb.buf.WriteString(`"cost":`)
				lb.buf.WriteString(costValue)

			case 'G':
				lb.buf.WriteString(`"logdate":"`)
				lb.buf.WriteString(cache.isoTime)
				lb.buf.WriteString(`",`)
			}
			if len(piece) > 1 {
				lb.buf.Write(piece[1:])
			}
		} else if len(piece) > 0 {
			lb.buf.Write(piece)
		}
	}

	//  write the close parenthese
	lb.buf.WriteString("}")
	lb.buf.WriteByte('\n')
	return nil
}

// Request that the logs rotate
func (w *FileLog2Writer) Rotate() {
	w.rot <- true
}

// If this is called in a threaded context, it MUST be synchronized
func (w *FileLog2Writer) intRotate() error {
	// Close any log file that may be open
	if w.file != nil {
		fmt.Fprint(w.file, FormatLogRecord(w.trailer, &LogRecord{Created: time.Now()}))
		w.file.Close()
	}

	// If we are keeping log files, move it to the next available number
	if w.rotate {
		_, err := os.Lstat(w.filename)
		if err == nil { // file exists
			// Find the next available number
			num := 1
			fname := ""
			if w.daily && time.Now().Day() != w.daily_opendate {
				yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

				for ; err == nil && num <= w.maxbackup; num++ {
					fname = w.filename + fmt.Sprintf(".%s.%03d", yesterday, num)
					_, err = os.Lstat(fname)
				}
				// return error if the last file checked still existed
				if err == nil {
					return fmt.Errorf("Rotate: Cannot find free log number to rename %s\n", w.filename)
				}
			} else {
				num = w.maxbackup - 1
				for ; num >= 1; num-- {
					fname = w.filename + fmt.Sprintf(".%d", num)
					nfname := w.filename + fmt.Sprintf(".%d", num+1)
					_, err = os.Lstat(fname)
					if err == nil {
						os.Rename(fname, nfname)
					}
				}
			}

			w.file.Close()
			// Rename the file to its newfound home
			err = os.Rename(w.filename, fname)
			if err != nil {
				return fmt.Errorf("Rotate: %s\n", err)
			}
		}
	}

	// Open the log file
	fd, err := os.OpenFile(w.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		return err
	}
	w.file = fd

	now := time.Now()
	fmt.Fprint(w.file, FormatLogRecord(w.header, &LogRecord{Created: now}))

	// Set the daily open date to the current date
	w.daily_opendate = now.Day()

	// initialize rotation values
	w.maxlines_curlines = 0
	w.maxsize_cursize = 0

	return nil
}

// Set the logging format (chainable).  Must be called before the first log
// message is written.
func (w *FileLog2Writer) SetFormat(format string) *FileLog2Writer {
	w.format = format
	w.pieces = bytes.Split([]byte(format), []byte{'%'})
	return w
}

// Set the logfile header and footer (chainable).  Must be called before the first log
// message is written.  These are formatted similar to the FormatLogRecord (e.g.
// you can use %D and %T in your header/footer for date and time).
func (w *FileLog2Writer) SetHeadFoot(head, foot string) *FileLog2Writer {
	w.header, w.trailer = head, foot
	if w.maxlines_curlines == 0 {
		fmt.Fprint(w.file, FormatLogRecord(w.header, &LogRecord{Created: time.Now()}))
	}
	return w
}

// Set rotate at linecount (chainable). Must be called before the first log
// message is written.
func (w *FileLog2Writer) SetRotateLines(maxlines int) *FileLog2Writer {
	//fmt.Fprintf(os.Stderr, "FileLog2Writer.SetRotateLines: %v\n", maxlines)
	w.maxlines = maxlines
	return w
}

// Set rotate at size (chainable). Must be called before the first log message
// is written.
func (w *FileLog2Writer) SetRotateSize(maxsize int) *FileLog2Writer {
	//fmt.Fprintf(os.Stderr, "FileLog2Writer.SetRotateSize: %v\n", maxsize)
	w.maxsize = maxsize
	return w
}

// Set rotate daily (chainable). Must be called before the first log message is
// written.
func (w *FileLog2Writer) SetRotateDaily(daily bool) *FileLog2Writer {
	//fmt.Fprintf(os.Stderr, "FileLog2Writer.SetRotateDaily: %v\n", daily)
	w.daily = daily
	return w
}

// Set max backup files. Must be called before the first log message
// is written.
func (w *FileLog2Writer) SetRotateMaxBackup(maxbackup int) *FileLog2Writer {
	w.maxbackup = maxbackup
	return w
}

// SetRotate changes whether or not the old logs are kept. (chainable) Must be
// called before the first log message is written.  If rotate is false, the
// files are overwritten; otherwise, they are rotated to another file before the
// new log is opened.
func (w *FileLog2Writer) SetRotate(rotate bool) *FileLog2Writer {
	//fmt.Fprintf(os.Stderr, "FileLog2Writer.SetRotate: %v\n", rotate)
	w.rotate = rotate
	return w
}
