package plugin

type PerfPlugin func(metric string)

var PerfPluginer *PerfPlugin

func DoPerfPlugin(metric string) {
	if PerfPluginer == nil {
		return
	}
	(PerfPlugin)(*PerfPluginer)(metric)
	return
}
