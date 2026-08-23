package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/chronos-mesh/chronos-mesh/internal/clockmodel"
	"github.com/chronos-mesh/chronos-mesh/internal/configuration"
	"github.com/chronos-mesh/chronos-mesh/internal/control"
	"github.com/chronos-mesh/chronos-mesh/internal/leapsecond"
	"github.com/chronos-mesh/chronos-mesh/internal/ntp"
	"github.com/chronos-mesh/chronos-mesh/internal/platform"
	"github.com/chronos-mesh/chronos-mesh/internal/ratelimit"
	"github.com/chronos-mesh/chronos-mesh/internal/selection"
	"github.com/chronos-mesh/chronos-mesh/internal/source"
	"os"
	"sync"
	"time"
)

func main() {
	path := flag.String("config", "configs/config.yaml", "config path")
	flag.Parse()
	cfg, e := configuration.Load(*path)
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
	if e = configuration.Validate(cfg); e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
	log := platform.Logger()
	ctx, cancel := platform.SignalContext(context.Background())
	defer cancel()
	model := clockmodel.New(platform.RealClock{}, cfg.Holdover)
	mgr := source.NewManager(source.UDPProvider{Timeout: 500 * time.Millisecond})
	sc := make([]struct{ Name, Address string }, len(cfg.Sources))
	for i, s := range cfg.Sources {
		sc[i] = struct{ Name, Address string }{s.Name, s.Address}
	}
	go mgr.Run(ctx, sc)
	go func() {
		t := time.NewTicker(time.Second)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				samples := mgr.List()
				d := selection.Select(samples)
				model.Update(d, samples)
			}
		}
	}()
	leap := &leapsecond.Table{}
	leap.Load(1, time.Now().Add(24*time.Hour), []leapsecond.Entry{})
	udp := ntp.NewServer(cfg.NTPAddr, model, log)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		if e := udp.Start(ctx); e != nil {
			log.Error("ntp stopped", "error", e)
		}
	}()
	ctl := &control.Server{Addr: cfg.HTTPAddr, Model: model, Sources: mgr, Leap: leap, Limit: ratelimit.New(cfg.RatePerSecond), Started: time.Now()}
	go func() {
		defer wg.Done()
		if e := control.Run(ctx, ctl); e != nil && ctx.Err() == nil {
			log.Error("control stopped", "error", e)
		}
	}()
	log.Info("chronos-mesh started", "http", cfg.HTTPAddr, "ntp", cfg.NTPAddr)
	<-ctx.Done()
	udp.Close()
	wg.Wait()
}
