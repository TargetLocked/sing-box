package dns

import (
	"errors"

	mDNS "github.com/miekg/dns"
)

const (
	RcodeSuccess     RcodeError = mDNS.RcodeSuccess
	RcodeFormatError RcodeError = mDNS.RcodeFormatError
	RcodeNameError   RcodeError = mDNS.RcodeNameError
	RcodeRefused     RcodeError = mDNS.RcodeRefused
)

type RcodeError int

func (e RcodeError) Error() string {
	return mDNS.RcodeToString[int(e)]
}

func IsNameError(err error) bool {
	return errors.Is(err, RcodeNameError)
}
