package protocol

type capabilityValidators struct {
	validateCapabilityForMethod      func(method string) error
	validateNotificationCapability   func(method string) error
	validateRequestHandlerCapability func(method string) error
}

func (p *Protocol) SetValidateCapabilityForMethod(validateCapabilityForMethod func(method string) error) {
	p.capabilityValidators.validateCapabilityForMethod = validateCapabilityForMethod
}

func (p *Protocol) SetValidateNotificationCapability(validateNotificationCapability func(method string) error) {
	p.capabilityValidators.validateNotificationCapability = validateNotificationCapability
}

func (p *Protocol) SetValidateRequestHandlerCapability(validateRequestHandlerCapability func(method string) error) {
	p.capabilityValidators.validateRequestHandlerCapability = validateRequestHandlerCapability
}
