cobrather
=========

Modularized command

* Simple Module

```
// Module info
var Module = &cobrather.Module{
	Use: "example",
	Commands: []*cobrather.Module{
		cobrather.VersionModule, // sub command module
	},
}

func main() {
	Module.MustMainRun()
}
```

* Custom viper for more config for environment

```
// Module info
var Module = &cobrather.Module{
	Use: "example",
	Commands: []*cobrather.Module{
		cobrather.VersionModule, // sub command module
	},
}

func main() {
	vr := viper.New()
	vr.AutomaticEnv()
	vr.SetEnvPrefix("ENVPREFIX")
	Module.MustMainRun(
		cobrather.MainRunOptionViper(vr),
	)
}
```
