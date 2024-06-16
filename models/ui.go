package models

//Used to define a menu item displayed on a page
type MenuItem struct {
	Title     string
	Uri       string
	IsEnabled bool
    IsExternal bool
}
