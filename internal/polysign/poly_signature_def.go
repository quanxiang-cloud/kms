// Package polysign is signature defines for polyapi
package polysign

// signature header
const (
	XHeaderPolySignVersion           = "X-Polysign-Version"
	XHeaderPolySignMethod            = "X-Polysign-Method"
	XHeaderPolySignKeyID             = "X-Polysign-Access-Key-Id"
	XHeaderPolySignTimestamp         = "X-Polysign-Timestamp"
	XInternalHeaderPolySignSignature = "X-Inner-Polysign-Signature" // NOTE: kms signature result header
)

// special body field define
const (
	// XPolyBodyHideArgs is poly reserve field in body
	// NOTE: pass path arg of raw api by this object
	XPolyBodyHideArgs = "$polyapi_hide$"

	// NOTE: this name means this is real body root of customer api
	XPolyCustomerBodyRoot = "$body$"

	// XPolyRaiseUpFieldName is a special filed name.
	// NOTE: if a field with this name, generate query will raiseup its children
	// eg: {"a":1,"b":2} is the same as {"a":1,"$$*_raise_*$$":{"b":2}}
	XPolyRaiseUpFieldName = "$$*_raise_*$$"
)
