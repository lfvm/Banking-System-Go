package utils



const ( 
	USD = "USD"
	EUR = "EUR"
	MXN = "MXN"
)


func IsSupportedCurrency(currency string) bool {


	switch currency { 
		case USD, EUR, MXN:
			return true
		}

	return false

}