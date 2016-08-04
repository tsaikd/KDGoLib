package cobrather

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tsaikd/KDGoLib/errutil"
)

// errors
var (
	ErrorFlagNotYetBind1 = errutil.NewFactory("flag %q is not yet bind")
)

// Flag is a common interface related to parsing flags in cobra.
type Flag interface {
	Bind(flagset *pflag.FlagSet, v *viper.Viper) error
}

// BoolFlag represents a flag that takes as bool value
type BoolFlag struct {
	Name      string
	ShortHand string
	Default   bool
	Usage     string
	EnvVar    string
	Hidden    bool

	viper *viper.Viper
}

var _ Flag = &BoolFlag{}

// Bind flag to flagset and viper for environment
func (t *BoolFlag) Bind(flagset *pflag.FlagSet, v *viper.Viper) (err error) {
	t.viper = v
	flagset.BoolP(t.Name, t.ShortHand, t.Default, t.Usage)
	return bindFlagSet(flagset, v, t.Name, t.EnvVar, t.Hidden)
}

// Bool return flag value
func (t *BoolFlag) Bool() bool {
	if t.viper == nil {
		panic(ErrorFlagNotYetBind1.New(nil, t.Name))
	}
	return t.viper.GetBool(t.Name)
}

// Int64Flag represents a flag that takes as int64 value
type Int64Flag struct {
	Name      string
	ShortHand string
	Default   int64
	Usage     string
	EnvVar    string
	Hidden    bool

	viper *viper.Viper
}

var _ Flag = &Int64Flag{}

// Bind flag to flagset and viper for environment
func (t *Int64Flag) Bind(flagset *pflag.FlagSet, v *viper.Viper) (err error) {
	t.viper = v
	flagset.Int64P(t.Name, t.ShortHand, t.Default, t.Usage)
	return bindFlagSet(flagset, v, t.Name, t.EnvVar, t.Hidden)
}

// Int64 return flag value
func (t *Int64Flag) Int64() int64 {
	if t.viper == nil {
		panic(ErrorFlagNotYetBind1.New(nil, t.Name))
	}
	return int64(t.viper.GetInt(t.Name))
}

// Float64Flag represents a flag that takes as float64 value
type Float64Flag struct {
	Name      string
	ShortHand string
	Default   float64
	Usage     string
	EnvVar    string
	Hidden    bool

	viper *viper.Viper
}

var _ Flag = &Float64Flag{}

// Bind flag to flagset and viper for environment
func (t *Float64Flag) Bind(flagset *pflag.FlagSet, v *viper.Viper) (err error) {
	t.viper = v
	flagset.Float64P(t.Name, t.ShortHand, t.Default, t.Usage)
	return bindFlagSet(flagset, v, t.Name, t.EnvVar, t.Hidden)
}

// Float64 return flag value
func (t *Float64Flag) Float64() float64 {
	if t.viper == nil {
		panic(ErrorFlagNotYetBind1.New(nil, t.Name))
	}
	return t.viper.GetFloat64(t.Name)
}

// StringFlag represents a flag that takes as string value
type StringFlag struct {
	Name      string
	ShortHand string
	Default   string
	Usage     string
	EnvVar    string
	Hidden    bool

	viper *viper.Viper
}

var _ Flag = &StringFlag{}

// Bind flag to flagset and viper for environment
func (t *StringFlag) Bind(flagset *pflag.FlagSet, v *viper.Viper) (err error) {
	t.viper = v
	flagset.StringP(t.Name, t.ShortHand, t.Default, t.Usage)
	return bindFlagSet(flagset, v, t.Name, t.EnvVar, t.Hidden)
}

// String return flag value
func (t *StringFlag) String() string {
	if t.viper == nil {
		panic(ErrorFlagNotYetBind1.New(nil, t.Name))
	}
	return t.viper.GetString(t.Name)
}

// StringSliceFlag represents a flag that takes as string value
type StringSliceFlag struct {
	Name      string
	ShortHand string
	Default   []string
	Usage     string
	Hidden    bool

	value *[]string
}

var _ Flag = &StringSliceFlag{}

// Bind flag to flagset and viper for environment
func (t *StringSliceFlag) Bind(flagset *pflag.FlagSet, v *viper.Viper) (err error) {
	t.value = flagset.StringSliceP(t.Name, t.ShortHand, t.Default, t.Usage)
	return bindFlagSet(flagset, v, t.Name, "", t.Hidden)
}

// StringSlice return flag value
func (t *StringSliceFlag) StringSlice() []string {
	if t.value == nil {
		panic(ErrorFlagNotYetBind1.New(nil, t.Name))
	}
	return *t.value
}

func bindFlagSet(flagset *pflag.FlagSet, v *viper.Viper, name string, envvar string, hidden bool) (err error) {
	if hidden {
		if err = flagset.MarkHidden(name); err != nil {
			return
		}
	}
	if envvar != "" {
		v.BindEnv(name, envvar)
	}
	return v.BindPFlag(name, flagset.Lookup(name))
}
