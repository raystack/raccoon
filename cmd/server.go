package cmd

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/raystack/raccoon/app"
	"github.com/raystack/raccoon/config"
	"github.com/raystack/raccoon/logger"
	"github.com/raystack/raccoon/metrics"
	"github.com/raystack/raccoon/middleware"
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
			middleware.Load()
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

type durationFlag struct {
	value *time.Duration
}

func (df durationFlag) String() string {
	if df.value == nil {
		return "0"
	}
	return fmt.Sprintf("%d", *df.value/time.Millisecond)
}

func (df durationFlag) Set(raw string) error {
	v, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing duration: %w", err)
	}
	*df.value = time.Millisecond * time.Duration(v)
	return nil
}

func (df durationFlag) Type() string {
	return "int"
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

func bindFlag(flag *pflag.FlagSet, ref any, meta reflect.StructField) {

	el := reflect.ValueOf(ref).Elem()
	kind := el.Kind()
	typ := el.Type()
	flagName := meta.Tag.Get("cmdx")
	desc := meta.Tag.Get("desc")

	switch {
	case typ.Name() == "Duration":
		v := ref.(*time.Duration)
		flag.Var(durationFlag{v}, flagName, desc)
	case typ.Name() == "AckType":
		v := ref.(*config.AckType)
		flag.Var(ackTypeFlag{v}, flagName, desc)
	case kind == reflect.String:
		v := ref.(*string)
		flag.StringVar(v, flagName, *v, desc)
	case kind == reflect.Int:
		v := ref.(*int)
		flag.IntVar(v, flagName, *v, desc)
	case kind == reflect.Int64:
		v := ref.(*int64)
		flag.Int64Var(v, flagName, *v, desc)
	case kind == reflect.Uint32:
		v := ref.(*uint32)
		flag.Uint32Var(v, flagName, *v, desc)
	case kind == reflect.Bool:
		v := ref.(*bool)
		flag.BoolVar(v, flagName, *v, desc)
	case kind == reflect.Slice && typ.Elem().String() == "string":
		v := ref.(*[]string)
		flag.StringSliceVar(v, flagName, *v, desc)
	default:
		msg := fmt.Sprintf("unsupport flag. kind = %s, type = %s", kind, typ)
		panic(msg)
	}

	// viper.BindPFlag(name, flag.Lookup(flagName))
}
