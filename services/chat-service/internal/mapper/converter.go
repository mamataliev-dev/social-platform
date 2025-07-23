// Package mapper provides a top-level entrypoint for all domain mappers.
package mapper

// Mappers groups every service-specific mapper under one struct,
// so you can inject a single dependency.
type Mappers struct {
	Room RoomMapper
}

// NewMappers constructs a Mappers with all sub-mappers initialized.
func NewMappers() *Mappers {
	return &Mappers{
		Room: NewRoomMapper(),
	}
}
