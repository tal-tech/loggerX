package plugin

type PerfPlugin func(metric string)

var PerfPluginer *PerfPlugin

//perf plugin
//count the error log in falcon
func DoPerfPlugin(metric string) {
	if PerfPluginer == nil {
		return
	}
	(PerfPlugin)(*PerfPluginer)(metric)
	return
}
