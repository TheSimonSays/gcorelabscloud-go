package main

import (
	"fmt"
	"gcloud/gcorecloud-go/gcoreclient/flags"
	"gcloud/gcorecloud-go/gcoreclient/flavors"
	"gcloud/gcorecloud-go/gcoreclient/heat/resources"
	"gcloud/gcorecloud-go/gcoreclient/keypairs"
	"gcloud/gcorecloud-go/gcoreclient/magnum/clusters"
	"gcloud/gcorecloud-go/gcoreclient/magnum/nodegroups"
	"gcloud/gcorecloud-go/gcoreclient/magnum/templates"
	"gcloud/gcorecloud-go/gcoreclient/networks"
	"gcloud/gcorecloud-go/gcoreclient/subnets"
	"gcloud/gcorecloud-go/gcoreclient/tasks"
	"gcloud/gcorecloud-go/gcoreclient/volumes"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var commands = []*cli.Command{
	&networks.NetworkCommands,
	&tasks.TaskCommands,
	&keypairs.KeypairCommands,
	&volumes.VolumeCommands,
	&subnets.SubnetCommands,
	&flavors.FlavorCommands,
	{
		Name:  "heat",
		Usage: "Gcloud Heat API",
		Subcommands: []*cli.Command{
			&resources.ResourceCommands,
		},
	},
	{
		Name:  "magnum",
		Usage: "Gcloud Magnum API",
		Subcommands: []*cli.Command{
			&clusters.ClusterCommands,
			&templates.ClusterTemplatesCommands,
			&nodegroups.ClusterNodeGroupCommands,
		},
	},
}

func main() {

	flags.AddOutputFlags(commands)

	app := cli.NewApp()
	app.Version = "0.0.1"
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		{
			Name:  "password",
			Usage: fmt.Sprintf("GCloud API client\n%s", flags.PasswordClientHelpText),
			Flags: flags.PasswordClientFlags,
			Before: func(c *cli.Context) error {
				return c.Set("client-type", "password")
			},
			Subcommands: commands,
		},
		{
			Name:        "token",
			Usage:       fmt.Sprintf("GCloud API client\n%s", flags.TokenClientHelpText),
			Flags:       flags.TokenClientFlags,
			Subcommands: commands,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		logrus.Errorf("Cannot initialize application: %+v", err)
		os.Exit(1)
	}
}