package domain

func IsAddressAvailable(direction Direction, address Address) bool {

	switch direction {
	case DirectionService:
		return address == AddressPoeticheskiy || address == AddressSikeirosa

	case DirectionBodyWork, DirectionDetailing, DirectionCarWash:
		return address == AddressPoeticheskiy

	default:
		return false
	}

}
