package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/watchtower/internal/cmd/base"
	"github.com/hashicorp/watchtower/internal/cmd/commands/controller"
	"github.com/hashicorp/watchtower/internal/cmd/commands/dev"
	"github.com/hashicorp/watchtower/internal/cmd/commands/worker"
	"github.com/mitchellh/cli"
)

// Commands is the mapping of all the available commands.
var Commands map[string]cli.CommandFactory

func initCommands(ui, serverCmdUi cli.Ui, runOpts *RunOptions) {
	/*
		getBaseCommand := func() *base.Command {
			return &base.Command{
				UI:      ui,
				Address: runOpts.Address,
			}
		}
	*/

	Commands = map[string]cli.CommandFactory{
		"controller": func() (cli.Command, error) {
			return &controller.Command{
				Command: &base.Command{
					UI:      serverCmdUi,
					Address: runOpts.Address,
				},
				ShutdownCh: MakeShutdownCh(),
				SighupCh:   MakeSighupCh(),
				SigUSR2Ch:  MakeSigUSR2Ch(),
			}, nil
		},
		"worker": func() (cli.Command, error) {
			return &worker.Command{
				Command: &base.Command{
					UI:      serverCmdUi,
					Address: runOpts.Address,
				},
				ShutdownCh: MakeShutdownCh(),
				SighupCh:   MakeSighupCh(),
				SigUSR2Ch:  MakeSigUSR2Ch(),
			}, nil
		},
		"dev": func() (cli.Command, error) {
			return &dev.Command{
				Command: &base.Command{
					UI:      serverCmdUi,
					Address: runOpts.Address,
				},
				ShutdownCh: MakeShutdownCh(),
				SighupCh:   MakeSighupCh(),
				SigUSR2Ch:  MakeSigUSR2Ch(),
			}, nil
		},
	}
}

// MakeShutdownCh returns a channel that can be used for shutdown
// notifications for commands. This channel will send a message for every
// SIGINT or SIGTERM received.
func MakeShutdownCh() chan struct{} {
	resultCh := make(chan struct{})

	shutdownCh := make(chan os.Signal, 4)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-shutdownCh
		close(resultCh)
	}()
	return resultCh
}

// MakeSighupCh returns a channel that can be used for SIGHUP
// reloading. This channel will send a message for every
// SIGHUP received.
func MakeSighupCh() chan struct{} {
	resultCh := make(chan struct{})

	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, syscall.SIGHUP)
	go func() {
		for {
			<-signalCh
			resultCh <- struct{}{}
		}
	}()
	return resultCh
}