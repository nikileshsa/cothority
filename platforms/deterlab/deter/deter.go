// deter is the deterlab process that should run on the boss node
//
// It spawns multiple timestampers and clients, while constructing
// the topology defined on cfg.json. It assumes that hosts.txt has
// the entire list of hosts to run timestampers on and that the final
// host is the designated logging server.
//
// The overall topology that is created is defined by cfg.json.
// The port layout for each node, however, is specified here.
// cfg.json will assign each node a port p. This is the port
// that each singing node is listening on. The timestamp server
// to which clients connect is listneing on port p+1. And the
// pprof server for each node is listening on port p+2. This
// means that in order to debug each client, you can forward
// the p+2 port of each node to your localhost.
//
// In the future the loggingserver will be connecting to the
// servers on the pprof port in order to gather extra data.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	dbg "github.com/ineiti/cothorities/helpers/debug_lvl"
	"github.com/ineiti/cothorities/helpers/cliutils"
	"github.com/ineiti/cothorities/helpers/config"
	"github.com/ineiti/cothorities/helpers/graphs"
)

var rootname string

var nmsgs string
var hpn string
var bf string
var debug int
var rate int
var failures int
var rFail int
var fFail int
var rounds int
var kill bool
var testConnect bool
var app string
var suite string
var pprofaddr string
var configFile string
var rootwait int

func init() {
	flag.StringVar(&nmsgs, "nmsgs", "100", "the number of messages per round")
	flag.StringVar(&hpn, "hpn", "", "number of hosts per node")
	flag.StringVar(&bf, "bf", "", "branching factor")
	flag.IntVar(&debug, "debug", 2, "set debug level")
	flag.IntVar(&rate, "rate", -1, "number of milliseconds between messages")
	flag.IntVar(&failures, "failures", 0, "percent showing per node probability of failure")
	flag.IntVar(&rFail, "rfail", 0, "number of consecutive rounds each root runs before it fails")
	flag.IntVar(&fFail, "ffail", 0, "number of consecutive rounds each follower runs before it fails")
	flag.IntVar(&rounds, "rounds", 100, "number of rounds to timestamp")
	flag.BoolVar(&kill, "kill", false, "kill everything (and don't start anything)")
	flag.BoolVar(&testConnect, "test_connect", false, "test connecting and disconnecting")
	flag.StringVar(&app, "app", "stamp", "app to run")
	flag.StringVar(&suite, "suite", "nist256", "abstract suite to use [nist256, nist512, ed25519]")
	flag.StringVar(&pprofaddr, "pprof", ":10000", "the address to run the pprof server at")
	flag.StringVar(&configFile, "config", "cfg.json", "the json configuration file")
	flag.IntVar(&rootwait, "rootwait", 30, "the amount of time the root should wait")
}

func main() {
	flag.Parse()
	dbg.DebugVisible = debug

	fmt.Println("running deter with nmsgs:", nmsgs, rate, rounds)

	virt, err := cliutils.ReadLines("remote/virt.txt")
	if err != nil {
		log.Fatal(err)
	}
	phys, err := cliutils.ReadLines("remote/phys.txt")
	if err != nil {
		log.Fatal(err)
	}
	vpmap := make(map[string]string)
	for i := range virt {
		vpmap[virt[i]] = phys[i]
	}
	// kill old processes
	var wg sync.WaitGroup
	for _, h := range phys {
		wg.Add(1)
		go func(h string) {
			defer wg.Done()
			cliutils.SshRun("", h, "sudo killall exec logserver timeclient scp ssh 2>/dev/null >/dev/null")
			time.Sleep(1 * time.Second)
			cliutils.SshRun("", h, "sudo killall exec 2>/dev/null >/dev/null")
		}(h)
	}
	wg.Wait()

	if kill {
		return
	}

	for _, h := range phys {
		wg.Add(1)
		go func(h string) {
			defer wg.Done()
			cliutils.Rsync("", h, "remote", "")
		}(h)
	}
	wg.Wait()

	nloggers := 3
	masterLogger := phys[0]
	slaveLogger1 := phys[1]
	slaveLogger2 := phys[2]
	loggers := []string{masterLogger, slaveLogger1, slaveLogger2}

	phys = phys[nloggers:]
	virt = virt[nloggers:]

	// Read in and parse the configuration file
	file, err := ioutil.ReadFile("remote/cfg.json")
	if err != nil {
		log.Fatal("deter.go: error reading configuration file: %v\n", err)
	}
	dbg.Lvl3("cfg file:", string(file))
	var cf config.ConfigFile
	err = json.Unmarshal(file, &cf)
	if err != nil {
		log.Fatal("unable to unmarshal config.ConfigFile:", err)
	}

	hostnames := cf.Hosts
	dbg.Lvl3("hostnames:", hostnames)

	depth := graphs.Depth(cf.Tree)
	var random_leaf string
	cf.Tree.TraverseTree(func(t *graphs.Tree) {
		if random_leaf != "" {
			return
		}
		if len(t.Children) == 0 {
			random_leaf = t.Name
		}
	})

	rootname = hostnames[0]

	dbg.Lvl3("depth of tree:", depth)

	// mapping from physical node name to the timestamp servers that are running there
	// essentially a reverse mapping of vpmap except ports are also used
	physToServer := make(map[string][]string)
	for _, virt := range hostnames {
		v, _, _ := net.SplitHostPort(virt)
		p := vpmap[v]
		ss := physToServer[p]
		ss = append(ss, virt)
		physToServer[p] = ss
	}

	// start up the logging server on the final host at port 10000
	fmt.Println("starting up logservers: ", loggers)
	// start up the master logger
	loggerports := make([]string, len(loggers))
	for i, logger := range loggers {
		loggerport := logger + ":10000"
		loggerports[i] = loggerport
		// redirect to the master logger
		master := masterLogger + ":10000"
		// if this is the master logger than don't set the master to anything
		if loggerport == masterLogger+":10000" {
			master = ""
		}

		go cliutils.SshRunStdout("", logger, "cd remote/logserver; sudo ./logserver -addr="+loggerport+
			" -hosts="+strconv.Itoa(len(hostnames))+
			" -depth="+strconv.Itoa(depth)+
			" -bf="+bf+
			" -hpn="+hpn+
			" -nmsgs="+nmsgs+
			" -rate="+strconv.Itoa(rate)+
			" -master="+master)
	}

	// wait a little bit for the logserver to start up
	time.Sleep(5 * time.Second)
	dbg.Lvl1("starting", len(physToServer), "time clients:", physToServer)
	// start up one timeclient per physical machine
	// it requests timestamps from all the servers on that machine

	i := 0
	for p, ss := range physToServer {
		if len(ss) == 0 {
			continue
		}
		servers := strings.Join(ss, ",")
		go func(i int, p string) {
			_, err := cliutils.SshRun("", p, "cd remote; sudo ./timeclient -nmsgs="+nmsgs+
				" -name=client@"+p+
				" -server="+servers+
				" -logger="+loggerports[i]+
				" -debug="+strconv.Itoa(debug)+
				" -rate="+strconv.Itoa(rate))
			if err != nil {
				dbg.Lvl3("Deter.go : timeclient error ", err)
			}
			dbg.Lvl3("Deter.go : Finished with timeclient", p)
		}(i, p)
		i = (i + 1) % len(loggerports)
	}

	for phys, virts := range physToServer {
		if len(virts) == 0 {
			continue
		}
		dbg.Lvl1("starting timestampers for", virts)
		cmd := GenExecCmd(rFail, fFail, failures, phys, virts, loggerports[i], random_leaf)
		i = (i + 1) % len(loggerports)
		wg.Add(1)
		time.Sleep(100 * time.Millisecond)
		go func(phys, cmd string) {
			//dbg.Lvl3("running on ", phys, cmd)
			defer wg.Done()
			dbg.Lvl3("deter.go Starting clients on physical machine ", phys, cmd)
			err := cliutils.SshRunStdout("", phys, cmd)
			if err != nil {
				log.Fatal("ERROR STARTING TIMESTAMPER:", err, phys, virts)
			}
			dbg.Lvl3("Finished with Timestamper", phys)
		}(phys, cmd)

	}

	// wait for the servers to finish before stopping
	wg.Wait()
	//time.Sleep(10 * time.Minute)
}

// Generate all commands on one single physicial machines to launch every "nodes"
func GenExecCmd(rFail, fFail, failures int, phys string, names []string, loggerport, random_leaf string) string {
	connect := false
	cmd := ""
	bg := " & "
	for i, name := range names {
		dbg.Lvl2("deter.go Generate cmd timestamper : names == %s, random_leaf == %s, testConnect = %t", name, random_leaf, testConnect)
		if name == random_leaf && testConnect {
			connect = true
		}
		amroot := " -amroot=false"
		if name == rootname {
			amroot = " -amroot=true"
		}

		if i == len(names)-1 {
			bg = ""
		}
		cmd += "(cd remote; sudo ./forkexec" +
		" -rfail=" + strconv.Itoa(rFail) +
		" -ffail=" + strconv.Itoa(fFail) +
		" -failures=" + strconv.Itoa(failures) +
		" -physaddr=" + phys +
		" -hostname=" + name +
		" -logger=" + loggerport +
		" -debug=" + strconv.Itoa(debug) +
		" -suite=" + suite +
		" -rounds=" + strconv.Itoa(rounds) +
		" -app=" + app +
		" -test_connect=" + strconv.FormatBool(connect) +
		" -config=" + configFile +
		" -pprof=" + pprofaddr +
		" -rootwait=" + strconv.Itoa(rootwait) +
		amroot + bg +
		" ); "
	}
	return cmd
}