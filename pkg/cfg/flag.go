package cfg

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"reflect"
)

type Flag[T any] struct {
	Name       string
	Short      string
	Usage      string
	Default    T
	Cfg        string
	Persistent bool
	Required   bool
}

func NewFlag[T any](name, cfg, usage string) *Flag[T] {
	return &Flag[T]{Name: name, Cfg: cfg}
}

func (f *Flag[T]) SetShort(v string) *Flag[T] {
	f.Short = v
	return f
}

func (f *Flag[T]) SetDefault(v T) *Flag[T] {
	f.Default = v
	return f
}

func (f *Flag[T]) SetPersistent(v bool) *Flag[T] {
	f.Persistent = v
	return f
}

func (f *Flag[T]) SetRequired(v bool) *Flag[T] {
	f.Required = v
	return f
}

func (f *Flag[T]) Bind(cmd *cobra.Command) {
	var flagSet *pflag.FlagSet
	if f.Persistent {
		flagSet = cmd.PersistentFlags()
	} else {
		flagSet = cmd.Flags()
	}

	def := reflect.ValueOf(f.Default).Interface()

	var t T
	switch reflect.ValueOf(t).Interface().(type) {
	case string:
		flagSet.StringP(f.Name, f.Short, def.(string), f.Usage)
	case bool:
		flagSet.BoolP(f.Name, f.Short, def.(bool), f.Usage)
	case int:
		flagSet.IntP(f.Name, f.Short, def.(int), f.Usage)
	case uint:
		flagSet.UintP(f.Name, f.Short, def.(uint), f.Usage)
	case int64:
		flagSet.Int64P(f.Name, f.Short, def.(int64), f.Usage)
	case uint64:
		flagSet.Uint64P(f.Name, f.Short, def.(uint64), f.Usage)
	default:
		cobra.CheckErr(fmt.Errorf("type %T not supported", t))
	}

	if f.Required {
		err := cmd.MarkFlagRequired(f.Name)
		cobra.CheckErr(err)
	}

	err := viper.BindPFlag(f.Cfg, flagSet.Lookup(f.Name))
	cobra.CheckErr(err)
}

func (f *Flag[T]) Get() T {
	var v any

	var t T
	switch reflect.ValueOf(t).Interface().(type) {
	case string:
		v = viper.GetString(f.Cfg)
	case bool:
		v = viper.GetBool(f.Cfg)
	case int:
		v = viper.GetInt(f.Cfg)
	case uint:
		v = viper.GetUint(f.Cfg)
	case int64:
		v = viper.GetInt64(f.Cfg)
	case uint64:
		v = viper.GetUint64(f.Cfg)
	default:
		cobra.CheckErr(fmt.Errorf("type %T not supported", t))
	}

	return v.(T)
}
