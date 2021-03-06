package convert

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

// DistributedCliConfig _
func DistributedCliConfig() *cli.Command {
	return &cli.Command{
		Name:    "distributed-convert",
		Aliases: []string{"dconv"},
		Subcommands: []*cli.Command{
			{
				Name:  "add",
				Flags: convertParamsFlags(),
				Action: func(c *cli.Context) error {
					app := &DistributedConvertApp{}
					err := app.Init()

					if err != nil {
						return errors.Wrap(err, "Initializing app")
					}

					// if err != nil {
					// 	return errors.Wrap(err, "Initializing app")
					// }

					// err = app.StartContracter(c)

					// if err != nil {
					// 	return errors.Wrap(err, "Starting contracter")
					// }

					// <-app.Wait()

					// return nil

					err = app.AddTask(c)

					if err != nil {
						return errors.Wrap(err, "Adding task to queue")
					}

					app.cancel()
					<-app.Wait()
					return nil
				},
			},
			{
				Name: "work",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name: "worker",
					},
				},
				Action: func(c *cli.Context) error {
					app := &DistributedConvertApp{}

					err := app.Init()

					if err != nil {
						return errors.Wrap(err, "Initializing app")
					}

					err = app.StartWorker()

					if err != nil {
						return errors.Wrap(err, "Starting worker")
					}

					err = app.StartContracter()

					if err != nil {
						return errors.Wrap(err, "Starting contracter")
					}

					<-app.Wait()

					return nil
				},
			},
			{
				Name:    "list-orders",
				Aliases: []string{"lo"},
				Flags:   []cli.Flag{},
				Action: func(c *cli.Context) error {
					app := &DistributedConvertApp{}

					err := app.Init()

					if err != nil {
						return errors.Wrap(err, "Initializing app")
					}

					str, err := app.ListOrders()

					fmt.Println(str)

					if err != nil {
						return err
					}

					app.cancel()
					<-app.Wait()
					return nil
				},
			},
			{
				Name:    "show-order",
				Aliases: []string{"so"},
				Flags:   []cli.Flag{},
				Action: func(c *cli.Context) error {
					app := &DistributedConvertApp{}

					err := app.Init()

					if err != nil {
						return errors.Wrap(err, "Initializing app")
					}

					str, err := app.ShowOrder(c.Args().Get(0))

					fmt.Println(str)

					if err != nil {
						return err
					}

					app.cancel()
					<-app.Wait()
					return nil
				},
			},
			{
				Name:    "list-segments",
				Aliases: []string{"ls"},
				Flags:   []cli.Flag{},
				Action: func(c *cli.Context) error {
					app := &DistributedConvertApp{}

					err := app.Init()

					if err != nil {
						return errors.Wrap(err, "Initializing app")
					}

					str, err := app.ListSegments(c.Args().Get(0))

					fmt.Println(str)

					if err != nil {
						return err
					}

					app.cancel()
					<-app.Wait()
					return nil
				},
			},
		},
	}
}
