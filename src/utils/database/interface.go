package database

//Accessor is used to access a database; different types of databases will have different implementations
type Accessor interface {
	Select([]string, map[string]interface{}) ([]interface{}, error)
	Insert([]string, [][]interface{}) error
	Delete(map[string]interface{})
	Update()
}
