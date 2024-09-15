package qr

import "errors"

type ErrorCorrection struct {
	data []byte
	// ecMode *string
	degree int
}

// for time being we are using 7-Q ,version 7 and mode Q
// var version int=7
func (ec *ErrorCorrection) ComputeRemainder() ([]byte, error) {
	// degree=version
	genPoly, err := reedSolomonGeneratorPolynomial(ec.degree)
	if err != nil {
		panic(err)
	}
	return reedSolomonRemainder(ec.data, genPoly), nil
}

func reedSolomonGeneratorPolynomial(degree int) ([]byte, error) {
	if degree < 1 || degree > 255 {
		return nil, errors.New("degree out of range")
	}
	res := make([]byte, degree)
	res[degree-1] = 1
	root := 1
	for i := 0; i < degree; i++ {
		for j := 0; j < len(res); j++ {
			res[j] = byte(reedSolomonMultiply(int(res[j]&0xFF), root))
			if j+1 < len(res) {
				res[j] = res[j] ^ res[j+1]
			}
		}
		root = reedSolomonMultiply(root, 0x02)
	}
	return res, nil
}

func reedSolomonRemainder(data, genPolynomial []byte) []byte {
	res := make([]byte, len(genPolynomial))

	for _, dataByte := range data {
		leadterm := (dataByte ^ res[0]) & 0xFF
		copy(res, res[1:])
		res[len(res)-1] = 0

		for i := 0; i < len(res); i++ {
			res[i] ^= byte(reedSolomonMultiply(int(genPolynomial[i]&0xFF), int(leadterm)))
		}
	}
	return res
}

func reedSolomonMultiply(x, y int) int {
	if x>>8 != 0 || y>>8 != 0 {
		panic("Inputs must be 8-bit values.")
	}
	z := 0
	for i := 7; i >= 0; i-- {
		z <<= 1
		if z>>8 != 0 {
			z ^= 0x11D
		}
		if (y>>i)&1 != 0 {
			z ^= x
		}
	}
	if z>>8 != 0 {
		panic("Result is not an 8-bit value.")
	}

	return z
}
