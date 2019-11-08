package psu

type Psu interface {
	SetOutput(channel int, voltage float64, current float64, enable bool) error
	GetOutput(channel int) (float64, float64, error)
	GetSetpoint(channel int) (float64, float64, error)
	EnableOutput(channel int) error
	DisableOutput(channel int) error
	Name() string
	Shutdown()
}

