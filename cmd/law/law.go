package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"github.com/cyberdelia/law/operator"
)

func init() {
	log.SetPrefix("law: ")
}

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
		log.Fatalln("wal segment required")
	}
	log.Printf("uploading wal segment %s", *cmd.segment)
	o, err := operator.NewOperator(*storage)
	if err != nil {
		log.Fatal(err)
	}
	if err = o.Archive(*cmd.segment); err != nil {
		log.Fatal(err)
	}
	log.Printf("uploaded wal segment %s", *cmd.segment)
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
		log.Fatalln("wal segment required")
	}
	if *cmd.destination == "" {
		log.Fatalln("wal destination required")
	}
	log.Printf("downloading wal segment %s", *cmd.segment)
	o, err := operator.NewOperator(*storage)
	if err != nil {
		log.Fatal(err)
	}
	if err = o.Unarchive(*cmd.segment, *cmd.destination); err != nil {
		log.Fatal(err)
	}
	log.Printf("downloaded wal segment %s", *cmd.segment)
}

type backupPush struct {
	cluster *string
	rate    *int
}

func (cmd *backupPush) Name() string {
	return "backup-push"
}

func (cmd *backupPush) DefineFlags(fs *flag.FlagSet) {
	cmd.cluster = fs.String("cluster", "", "Path of cluster directory")
	cmd.rate = fs.Int("rate-limit", 0, "Rate-limit i/o")
}

func (cmd *backupPush) Run() {
	if *cmd.cluster == "" {
		log.Fatalln("cluster directory required")
	}
	log.Printf("backuping %s", *cmd.cluster)
	o, err := operator.NewOperator(*storage)
	if err != nil {
		log.Fatal(err)
	}
	if err = o.Backup(*cmd.cluster, *cmd.rate); err != nil {
		log.Fatal(err)
	}
	log.Printf("backuped %s", *cmd.cluster)
}

type backupFetch struct {
	cluster *string
	name    *string
}

func (cmd *backupFetch) Name() string {
	return "backup-fetch"
}

func (cmd *backupFetch) DefineFlags(fs *flag.FlagSet) {
	cmd.cluster = fs.String("cluster", "", "Path of cluster directory")
	cmd.name = fs.String("name", "", "Name of backup")
}

func (cmd *backupFetch) Run() {
	if *cmd.cluster == "" {
		log.Fatalln("law: cluster directory required")
	}
	if *cmd.name == "" {
		log.Fatalln("law: name of backup required")
	}
	log.Printf("restoring backup %s to %s", *cmd.name, *cmd.cluster)
	o, err := operator.NewOperator(*storage)
	if err != nil {
		log.Fatal(err)
	}
	if err = o.Restore(*cmd.cluster, *cmd.name); err != nil {
		log.Fatal(err)
	}
	log.Printf("restored backup %s to %s", *cmd.name, *cmd.cluster)
}

var (
	cpuprofile = flag.String("cpuprofile", "", "CPU profile filepath")
	memprofile = flag.String("memprofile", "", "Memory profile filepath")
	storage    = flag.String("storage", os.Getenv("STORAGE_URL"), "Storage Source Name")
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

	if *storage == "" {
		log.Fatalln("storage source name required")
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
