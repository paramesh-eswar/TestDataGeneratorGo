package main

const DEFAULT string = "default"
const DUP_IN_RANGE string = "dupinrange"
const RANDOM string = "random"
const SEQ_IN_RANGE string = "seqinrange"
const NATURAL_SEQ string = "sequence"

const DELIMITER string = ","

type CCType int

const (
	VISA CCType = iota + 1
	MASTERCARD
	AMERICAN_EXPRESS
	DINERS_CLUB
	DISCOVER
	JCB
	UNIONPAY
	MAESTRO
	ELO
	HIPER
	HIPERCARD
)

// String - Creating common behavior - give the type a String function
func (cc CCType) String() string {
	return [...]string{"visa", "mastercard", "american-express", "diners-club", "discover", "jcb", "unionpay", "maestro", "elo", "hiper", "hiper-card"}[cc-1]
}

// EnumIndex - Creating common behavior - give the type a EnumIndex function
func (cc CCType) EnumIndex() int {
	return int(cc)
}
