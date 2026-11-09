package domain

type Direction string

const (
	DirectionService   Direction = "service"
	DirectionDetailing Direction = "detailing"
	DirectionCarWash   Direction = "car_wash"
	DirectionBodyWork  Direction = "body_work"
)

type Address string

const (
	AddressPoeticheskiy Address = "poeticheskiy"
	AddressSikeirosa    Address = "sikeirosa"
)

type Application struct {
	Direction Direction
	Address   Address
	Name      string
	Phone     string
	Car       string
	Problem   string
}

type Step string

const (
	StepDirection Step = "direction"
	StepAddress   Step = "address"
	StepName      Step = "name"
	StepPhone     Step = "phone"
	StepCar       Step = "car"
	StepProblem   Step = "problem"
	StepConfirm   Step = "confirm"
)

type UserSession struct {
	Application Application
	Step        Step
	MessageID   int64
}
