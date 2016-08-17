package domain

import (
	"github.com/xtraclabs/goes"
	"time"
	"log"
	"errors"
	"github.com/golang/protobuf/proto"
)

type ApplicationReg struct {
	*goes.Aggregate
	Name string
	Description string
	Created int64 //Unix time stamp serialized as an int64
}

const (
	AppRegCreatedCode = "ARCRE"
)

var (
	ErrUnknownType = errors.New("Unknown event type")
)

func NewApplicationReg(name, description string)(*ApplicationReg,error) {
	var appReg = new(ApplicationReg)
	appReg.Aggregate = goes.NewAggregate()
	appReg.Version = 1

	appRegCreated := ApplicationRegistrationCreated{
		AggregateId: appReg.ID,
		Name: name,
		Description: description,
		CreateTimestamp: time.Now().UnixNano(),
	}

	appReg.Apply(
		goes.Event{
			Source: appReg.ID,
			Version: appReg.Version,
			Payload: appRegCreated,
		})

	return appReg, nil
}

func (ar *ApplicationReg) Apply(event goes.Event) {
	ar.Route(event)
	ar.Events = append(ar.Events, event)
}


func (ar *ApplicationReg)Route(event goes.Event) {
	event.Version = ar.Version
	switch event.Payload.(type) {
	case ApplicationRegistrationCreated:
		ar.handleApplicationRegistrationCreated(event.Payload.(ApplicationRegistrationCreated))
		default:
			log.Printf("unexpected type handled: %t",event.Payload)
	}
}

func (ar *ApplicationReg) handleApplicationRegistrationCreated(event ApplicationRegistrationCreated) {
	ar.ID = event.AggregateId
	ar.Name = event.Name
	ar.Description = event.Description
	ar.Created = event.CreateTimestamp
}

func (ar *ApplicationReg) Store(eventStore goes.EventStore) error {
	marshalled, err := marshallEvents(ar.Events)
	if err != nil {
		return nil
	}

	log.Println("Storing ", len(ar.Events), " events.")

	aggregateToStore := &goes.Aggregate{
		ID:      ar.ID,
		Version: ar.Version,
		Events:  marshalled,
	}

	err = eventStore.StoreEvents(aggregateToStore)
	if err != nil {
		return err
	}

	ar.Events = make([]goes.Event, 0)

	return nil
}

func marshallEvents(events []goes.Event) ([]goes.Event, error) {

	var updatedEvents []goes.Event

	for _, e := range events {

		var err error
		var newEvent goes.Event
		newEvent.Source = e.Source
		newEvent.Version = e.Version

		switch e.Payload.(type) {
		case ApplicationRegistrationCreated:
			newEvent.TypeCode = AppRegCreatedCode
			newEvent.Payload, err = marshallCreate(e.Payload.(ApplicationRegistrationCreated))
			if err != nil {
				return nil, err
			}

		default:
			return nil, ErrUnknownType
		}

		updatedEvents = append(updatedEvents, newEvent)
	}

	return updatedEvents, nil
}

func marshallCreate(create ApplicationRegistrationCreated) ([]byte, error) {
	return proto.Marshal(&create)
}