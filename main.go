package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/ayn2op/discordo/cmd"

	"net/http"
	_ "net/http/pprof"
	"runtime/pprof"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}

	// Start CPU profiling
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	m, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer m.Close()

	if err := cmd.Run(); err != nil {
		slog.Error("failed to run command", "err", err)
	}

	// Write the heap profile after the program runs
	if err := pprof.WriteHeapProfile(m); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}

	log.Println("Memory profile saved to mem.prof")
}
