package core

type Ticket interface {
	CreateTicket(systemId int64, stub int32, accessControl int32) (string, error)
	ValidateTicket(ticket string) (OnSession, error)
}
