package pkgep

//** architecture: init pointer common libs that can be use throughout api server life circle

//** Use ONCE to initialize once instance

//** Flow: init log once, and keep it pointer instance, since most of other lib need to properly log it out

type Jwt interface {
	VerifyAuthorizationHeader()
}

