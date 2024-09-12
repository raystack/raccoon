package cmd

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/raystack/raccoon/app"
	"github.com/raystack/raccoon/config"
	"github.com/raystack/raccoon/pkg/logger"
	"github.com/raystack/raccoon/pkg/metrics"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func serverCommand() *cobra.Command {
	var configFile = "config.yaml"
	command := &cobra.Command{
		Use:   "server",
		Short: "Start raccoon server",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := config.Load(configFile)
			if err != nil {
				return err
			}
			metrics.Setup()
			defer metrics.Close()
			logger.SetLevel(config.Log.Level)
			return app.Run()
		},
	}
	command.Flags().SortFlags = false
	command.Flags().StringVarP(&configFile, "config", "c", configFile, "path to config file")
	for _, cfg := range config.Walk() {
		bindFlag(command.Flags(), cfg.Ref, cfg.Meta)
	}
	return command
}

func bindFlag(flag *pflag.FlagSet, ref any, meta reflect.StructField) {

	flagName := meta.Tag.Get("cmdx")
	desc := meta.Tag.Get("desc")

	switch v := ref.(type) {
	case *config.AckType:
		flag.Var(ackTypeFlag{v}, flagName, desc)
	case *string:
		flag.StringVar(v, flagName, *v, desc)
	case *int:
		flag.IntVar(v, flagName, *v, desc)
	case *int64:
		flag.Int64Var(v, flagName, *v, desc)
	case *uint32:
		flag.Uint32Var(v, flagName, *v, desc)
	case *bool:
		flag.BoolVar(v, flagName, *v, desc)
	case *[]string:
		flag.StringSliceVar(v, flagName, *v, desc)
	default:
		msg := fmt.Sprintf("unsupport flag of type %T", ref)
		panic(msg)
	}
}

type ackTypeFlag struct {
	value *config.AckType
}

func (af ackTypeFlag) String() string {
	if af.value == nil {
		return "0"
	}
	return fmt.Sprintf("%d", *af.value)
}

func (af ackTypeFlag) Set(raw string) error {
	v, err := strconv.ParseInt(raw, 10, 0)
	if err != nil {
		return fmt.Errorf("error parsing bool: %w", err)
	}
	*af.value = config.AckType(v)
	return nil
}

func (af ackTypeFlag) Type() string {
	return "int"
}
