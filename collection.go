package cmagic

import(
	r "github.com/hailocab/cmagic/reflect"
	g "github.com/hailocab/cmagic/generate"
	"reflect"
	"encoding/json"
)

type collection struct {
	nameSpace 		*nameSpace
	collectionInfo 	*collectionInfo
}

// Contains mostly analyzed information about the entity
type collectionInfo struct {
	keyspace, name string
	entity         interface{}
	primaryKey     string
	fieldNames     map[string]struct{} // This is here only to check containment
	fields         []string
	fieldValues    []interface{}
}

func newCollectionInfo(keyspace, name, primaryKey string, entity interface{}) *collectionInfo {
	cinf := &collectionInfo{
		keyspace:   keyspace,
		name:       name,
		entity:     entity,
		primaryKey: primaryKey,
	}
	fields, values, ok := r.FieldsAndValues(entity)
	if !ok {
		// panicking here since this is a programmer error
		panic("Supplied entity is not a struct")
	}
	cinf.fieldNames = map[string]struct{}{}
	for _, v := range fields {
		if v == cinf.primaryKey {
			continue
		}
		cinf.fieldNames[v] = struct{}{}
	}
	cinf.fields = fields
	cinf.fieldValues = values
	return cinf
}

func (c collection) zero() interface{} {
	return reflect.New(reflect.TypeOf(c.collectionInfo.entity)).Interface()
}

func (c collection) Read(id string) (interface{}, error) {
	stmt := g.ReadById(c.nameSpace.name, c.collectionInfo.primaryKey)
	m := map[string]interface{}{}
	sess := c.nameSpace.session
	sess.Query(stmt, id).Iter().MapScan(m)
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	ret := c.zero()
	err = json.Unmarshal(bytes, ret)
	return ret, err
}