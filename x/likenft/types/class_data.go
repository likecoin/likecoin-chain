package types

func (c ClassConfig) IsBlindBox() bool {
	return c.BlindBoxConfig != nil
}
