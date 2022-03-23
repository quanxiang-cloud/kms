package enums

// system config value set
var (
	SystemConfigSet = NewEnumSet(nil)
	ConfigKeyNum    = SystemConfigSet.Reg("keyNum")
	ConfigKeyExpiry = SystemConfigSet.Reg("keyExpiry")
)
