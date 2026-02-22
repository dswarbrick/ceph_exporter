package osd

import "github.com/alecthomas/kingpin/v2"

var (
	asokDir = kingpin.Flag(
		"collector.osd.socket-dir", "Directory containing OSD admin sockets.",
	).Default("/run/ceph").String()
)
