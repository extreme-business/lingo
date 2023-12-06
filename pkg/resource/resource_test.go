package resource_test

import (
	"testing"

	"github.com/dwethmar/lingo/pkg/resource"
	"github.com/google/uuid"
)

func TestResource_UUID(t *testing.T) {
	t.Run("should return the UUID of the resource", func(t *testing.T) {
		r := resource.Resource{ResourceID: "35be2810-df8b-4bba-b2c2-fef5d5be709a"}
		got, err := r.UUID()
		if err != nil {
			t.Errorf("Resource.UUID() error = %v, want nil", err)
			return
		}

		want := uuid.MustParse("35be2810-df8b-4bba-b2c2-fef5d5be709a")
		if got != want {
			t.Errorf("Resource.UUID() = %v, want %v", got, want)
		}
	})

	t.Run("should return an error if the ID is not a valid UUID", func(t *testing.T) {
		r := resource.Resource{ResourceID: "invalid"}
		id, err := r.UUID()
		if err == nil {
			t.Error("Resource.UUID() error = nil, want an error")
		}

		if id != uuid.Nil {
			t.Errorf("Resource.UUID() = %v, want %v", id, uuid.Nil)
		}
	})
}

func TestResource_Find(t *testing.T) {
	t.Run("should find a resource by collection", func(t *testing.T) {
		r := resource.Resource{
			ResourceID:   "1",
			CollectionID: "grandchild",
			Parent: &resource.Resource{
				ResourceID:   "2",
				CollectionID: "child",
				Parent: &resource.Resource{
					ResourceID:   "3",
					CollectionID: "parent",
				},
			},
		}

		got := r.Find("child")
		if got == nil {
			t.Error("Resource.Find() = nil, want a resource")
			return
		}

		if got.ResourceID != "2" {
			t.Errorf("Resource.Find().ID = %v, want 2", got.ResourceID)
		}

		if got.CollectionID != "child" {
			t.Errorf("Resource.Find().Collection = %v, want child", got.CollectionID)
		}
	})

	t.Run("should return nil if the resource is not found", func(t *testing.T) {
		r := resource.Resource{
			ResourceID:   "1",
			CollectionID: "parent",
			Parent: &resource.Resource{
				ResourceID:   "2",
				CollectionID: "child",
				Parent: &resource.Resource{
					ResourceID:   "3",
					CollectionID: "grandchild",
				},
			},
		}

		if r.Find("invalid") != nil {
			t.Error("Resource.Find() = resource, want nil")
		}
	})
}

func TestResource_Depth(t *testing.T) {
	t.Run("should return the depth of the resource", func(t *testing.T) {
		r := resource.Resource{
			ResourceID:   "1",
			CollectionID: "parent",
			Parent: &resource.Resource{
				ResourceID:   "2",
				CollectionID: "child",
				Parent: &resource.Resource{
					ResourceID:   "3",
					CollectionID: "grandchild",
				},
			},
		}

		if r.Depth() != 2 {
			t.Errorf("Resource.Depth() = %v, want 2", r.Depth())
		}
	})

	t.Run("should return 0 if the resource has no parent", func(t *testing.T) {
		r := resource.Resource{ResourceID: "1", CollectionID: "parent"}
		if r.Depth() != 0 {
			t.Errorf("Resource.Depth() = %v, want 0", r.Depth())
		}
	})
}

func TestNewParser(t *testing.T) {
	t.Run("should create a new parser", func(t *testing.T) {
		p := resource.NewParser()
		if p == nil {
			t.Error("NewParser() = nil, want a parser")
		}
	})
}

func TestParser_RegisterChild(t *testing.T) {
	t.Run("should register a child collection", func(t *testing.T) {
		p := resource.NewParser()
		p.RegisterChild("parent", "child")
		if !p.IsAllowedChild("parent", "child") {
			t.Error("RegisterChild() = false, want true")
		}
	})
}

func TestParser_Parse(t *testing.T) {
	t.Run("should parse a resource name", func(t *testing.T) {
		p := resource.NewParser()
		p.RegisterChild("parent", "child")

		r, err := p.Parse("parent/1/child/2")
		if err != nil {
			t.Errorf("Parse() error = %v, want nil", err)
			return
		}

		if r.ResourceID != "2" {
			t.Errorf("Parse() ID = %v, want 2", r.ResourceID)
		}

		if r.CollectionID != "child" {
			t.Errorf("Parse() Collection = %v, want child", r.CollectionID)
		}

		if r.Parent == nil {
			t.Error("Parse() Parent = nil, want a resource")
		}

		if r.Parent.ResourceID != "1" {
			t.Errorf("Parse() Parent.ID = %v, want 1", r.Parent.ResourceID)
		}

		if r.Parent.CollectionID != "parent" {
			t.Errorf("Parse() Parent.Collection = %v, want parent", r.Parent.CollectionID)
		}
	})

	t.Run("should return an error if the resource name is invalid", func(t *testing.T) {
		p := resource.NewParser()
		_, err := p.Parse("invalid")
		if err == nil {
			t.Error("Parse() error = nil, want an error")
		}
	})

	t.Run("should return an error if the child collection is not allowed", func(t *testing.T) {
		p := resource.NewParser()
		_, err := p.Parse("parent/1/child/2")
		if err == nil {
			t.Error("Parse() error = nil, want an error")
		}
	})
}

func TestParser_IsAllowedChild(t *testing.T) {
	t.Run("should return false if the child collection is not allowed", func(t *testing.T) {
		p := resource.NewParser()
		if p.IsAllowedChild("child", "invalid") {
			t.Error("IsAllowedChild() = true, want false")
		}
	})
}
