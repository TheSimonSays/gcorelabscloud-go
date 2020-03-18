package loadbalancers

import (
	"fmt"
	"gcloud/gcorecloud-go"
	"gcloud/gcorecloud-go/gcore/loadbalancer/v1/loadbalancers"
	"gcloud/gcorecloud-go/gcore/task/v1/tasks"
	"gcloud/gcorecloud-go/gcoreclient/flags"
	"gcloud/gcorecloud-go/gcoreclient/utils"

	"github.com/urfave/cli/v2"
)

var loadBalancerIDText = "loadbalancer_id is mandatory argument"

var loadBalancerListSubCommand = cli.Command{
	Name:     "list",
	Usage:    "Loadbalancers list",
	Category: "loadbalancer",
	Action: func(c *cli.Context) error {
		client, err := utils.BuildClient(c, "loadbalancers", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		pages, err := loadbalancers.List(client).AllPages()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		results, err := loadbalancers.ExtractLoadBalancers(pages)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var loadBalancerGetSubCommand = cli.Command{
	Name:      "show",
	Usage:     "Show loadbalancer",
	ArgsUsage: "<loadbalancer_id>",
	Category:  "loadbalancer",
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstArg(c, loadBalancerIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := utils.BuildClient(c, "loadbalancers", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		result, err := loadbalancers.Get(client, clusterID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var loadBalancerDeleteSubCommand = cli.Command{
	Name:      "delete",
	Usage:     "Show loadbalancer",
	ArgsUsage: "<loadbalancer_id>",
	Category:  "loadbalancer",
	Action: func(c *cli.Context) error {
		loadBalancerID, err := flags.GetFirstArg(c, loadBalancerIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := utils.BuildClient(c, "loadbalancers", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := loadbalancers.Get(client, loadBalancerID).ExtractTasks()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if results == nil {
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, client, results, false, func(task tasks.TaskID) (interface{}, error) {
			_, err := loadbalancers.Get(client, loadBalancerID).Extract()
			if err == nil {
				return nil, fmt.Errorf("cannot delete loadbalancer with ID: %s", loadBalancerID)
			}
			switch err.(type) {
			case gcorecloud.ErrDefault404:
				return nil, nil
			default:
				return nil, err
			}
		})
	},
}

var loadBalancerUpdateSubCommand = cli.Command{
	Name:      "update",
	Usage:     "Update loadbalancer",
	ArgsUsage: "<loadbalancer_id>",
	Category:  "loadbalancer",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Loadbalancer name",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstArg(c, loadBalancerIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := utils.BuildClient(c, "loadbalancers", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := loadbalancers.UpdateOpts{Name: c.String("name")}

		result, err := loadbalancers.Update(client, clusterID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if result == nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var LoadBalancerCommands = cli.Command{
	Name:  "loadbalancers",
	Usage: "GCloud loadbalancers API",
	Subcommands: []*cli.Command{
		&loadBalancerListSubCommand,
		&loadBalancerGetSubCommand,
		&loadBalancerUpdateSubCommand,
		&loadBalancerDeleteSubCommand,
	},
}