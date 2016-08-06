package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"github.com/cyberdelia/law/operator"
)

type walPush struct {
	segment *string
}

func (cmd *walPush) Name() string {
	return "wal-push"
}

func (cmd *walPush) DefineFlags(fs *flag.FlagSet) {
	cmd.segment = fs.String("segment", "", "Path to a WAL segment to upload")
}

func (cmd *walPush) Run() {
	if *cmd.segment == "" {
		log.Fatalln("law: wal segment required")
	}
	o, err := operator.NewOperator()
	if err != nil {
		log.Fatal(err)
	}
	if err = o.Archive(*cmd.segment); err != nil {
		log.Fatal(err)
	}
}

type walFetch struct {
	segment     *string
	destination *string
}

func (cmd *walFetch) Name() string {
	return "wal-fetch"
}

func (cmd *walFetch) DefineFlags(fs *flag.FlagSet) {
	cmd.segment = fs.String("segment", "", "Name of the WAL segment to download")
	cmd.destination = fs.String("destination", "", "Path of WAL segment locally")
}

func (cmd *walFetch) Run() {
	if *cmd.segment == "" {
		log.Fatalln("law: wal segment required")
	}
	if *cmd.destination == "" {
		log.Fatalln("law: wal destination required")
	}
	o, err := operator.NewOperator()
	if err != nil {
		log.Fatal(err)
	}
	if err = o.Unarchive(*cmd.segment, *cmd.destination); err != nil {
		log.Fatal(err)
	}
}

type backupPush struct {
	cluster *string
}

func (cmd *backupPush) Name() string {
	return "backup-push"
}

func (cmd *backupPush) DefineFlags(fs *flag.FlagSet) {
	cmd.cluster = fs.String("cluster", "", "Path of cluster directory")
}

func (cmd *backupPush) Run() {
	if *cmd.cluster == "" {
		log.Fatalln("law: cluster directory required")
	}
	o, err := operator.NewOperator()
	if err != nil {
		log.Fatal(err)
	}
	if err = o.Backup(*cmd.cluster); err != nil {
		log.Fatal(err)
	}
}

type backupFetch struct {
	cluster *string
	name    *string
	offset  *string
}

func (cmd *backupFetch) Name() string {
	return "backup-fetch"
}

func (cmd *backupFetch) DefineFlags(fs *flag.FlagSet) {
	cmd.cluster = fs.String("cluster", "", "Path of cluster directory")
	cmd.name = fs.String("name", "", "Name of backup")
	cmd.offset = fs.String("offset", "", "Offset of backup")
}

func (cmd *backupFetch) Run() {
	if *cmd.cluster == "" {
		log.Fatalln("law: cluster directory required")
	}
	if *cmd.name == "" {
		log.Fatalln("law: name of backup required")
	}
	o, err := operator.NewOperator()
	if err != nil {
		log.Fatal(err)
	}
	if err = o.Restore(*cmd.cluster, *cmd.name, *cmd.offset); err != nil {
		log.Fatal(err)
	}
}

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile = flag.String("memprofile", "", "write memory profile to this file")
)

func main() {
	log.SetFlags(0)
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	Parse(new(walPush), new(walFetch), new(backupPush), new(backupFetch))
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
	}
}
