// Package resource provides utilities for parsing structured resource names
// following Google's Resource-Oriented Design principles. This package allows
// for the interpretation of resource names into their component parts, and
// supports the enforcement of structured relationships between resources.
//
// Example resource name:
//
//	projects/my-project/locations/us-central1/namespaces/my-namespace
//
// More details can be found in Google's design documentation:
// https://cloud.google.com/apis/design/resource_names
package resource

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Resource represents a structured resource name.
type Resource struct {
	ResourceID   string    // ResourceID is the ID of the resource.
	CollectionID string    // CollectionID is the collection of the resource.
	Parent       *Resource // Parent is the parent of the resource.
}

func (r Resource) Depth() int {
	if r.Parent == nil {
		return 0
	}
	return r.Parent.Depth() + 1
}

func (r Resource) UUID() (uuid.UUID, error) {
	return uuid.Parse(r.ResourceID)
}

func (r Resource) Find(collection string) *Resource {
	if r.CollectionID == collection {
		return &r
	}
	if r.Parent == nil {
		return nil
	}
	return r.Parent.Find(collection)
}

type parentChildLink struct {
	Parent string
	Child  string
}

type Rule struct {
	Parent string
	Child  string
}

// Parser is a resource name parser that enforces structured relationships between resources.
type Parser struct {
	hierarchyRules map[Rule]struct{}
}

func NewParser() *Parser {
	return &Parser{
		hierarchyRules: make(map[Rule]struct{}),
	}
}

// RegisterChild registers a child collection under a parent collection.
func (p *Parser) RegisterChild(parentCollection, childCollection string) {
	p.hierarchyRules[Rule{Parent: parentCollection, Child: childCollection}] = struct{}{}
}

// Parse parses a structured resource name into a Resource struct.
func (p *Parser) Parse(name string) (Resource, error) {
	parts := strings.Split(name, "/")
	if len(parts)%2 != 0 {
		return Resource{}, errors.New("invalid resource name format")
	}

	var current *Resource
	for i := 0; i < len(parts); i += 2 {
		if current != nil && !p.IsAllowedChild(current.CollectionID, parts[i]) {
			return Resource{}, fmt.Errorf("child collection '%s' is not allowed under parent collection '%s'", parts[i], current.CollectionID)
		}
		current = &Resource{
			ResourceID:   parts[i+1],
			CollectionID: parts[i],
			Parent:       current,
		}
	}

	return *current, nil
}

// IsAllowedChild checks if a child collection is allowed under a parent collection.
func (p *Parser) IsAllowedChild(parentCollection, childCollection string) bool {
	_, ok := p.hierarchyRules[Rule{Parent: parentCollection, Child: childCollection}]
	return ok
}
