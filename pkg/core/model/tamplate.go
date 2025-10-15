package model

type Template[T []HandleTemplate | []Handle] struct {
	Name    string
	Path    string
	Service string
	Type    string
	Handle  T
}

func NewSortedTemplates() *SortedTemplates {
	return &SortedTemplates{
		templates: make(map[string][]Template[[]Handle], 0),
	}
}

// SortedTemplates is struct for storing templates sorted by types
// map[types][]Template -- types: http, grpc, etc, getting from Template.Type
type SortedTemplates struct {
	templates map[string][]Template[[]Handle] // map[types][]Template -- types: http, grpc, etc
}

// Add adds a new template to the sorted templates
// if the type already exists, it appends the template to the existing slice
// if the type does not exist, it creates a new slice with the template
// and adds it to the map
func (inst *SortedTemplates) Add(templ Template[[]Handle]) {
	if _, ok := inst.templates[templ.Type]; !ok {
		inst.templates[templ.Type] = []Template[[]Handle]{templ}
		return
	}

	inst.templates[templ.Type] = append(inst.templates[templ.Type], templ)
}

func (inst *SortedTemplates) Get(types string) ([]Template[[]Handle], bool) {
	templ, ok := inst.templates[types]
	return templ, ok
}

func (inst *SortedTemplates) GetAll() map[string][]Template[[]Handle] {
	return inst.templates
}

func (inst *SortedTemplates) Len() int {
	return len(inst.templates)
}

func (inst *SortedTemplates) Clear() {
	inst.templates = make(map[string][]Template[[]Handle], 0)
}

func (inst *SortedTemplates) Delete(types string) {
	delete(inst.templates, types)
}
