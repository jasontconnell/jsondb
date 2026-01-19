package jsondb

import "fmt"

type database[T any] struct {
	filename string
	data     []T
}

type Database[T any] interface {
	Store() error
	FindFirst(f func(elem T) bool) (T, bool)
	All() []T
	Filter(f func(elem T) bool) []T
	Exists(f func(elem T) bool) bool
	Add(elem T) bool
	AddOrUpdate(f func(elem T) bool, elem T) bool
	AddIfUnique(f func(elem T) bool, elem T) bool
	Update(f func(elem T) bool, elem T) bool
	UpdateProperty(f func(elem T) bool, upd func(elem *T)) bool
	Remove(f func(elem T) bool) (T, bool)
	RemoveFilter(f func(elem T) bool) int
	Clear() bool
}

func NewDatabase[T any](filename string) (Database[T], error) {
	f, err := readJsonFileList[T](filename)
	if err != nil {
		return nil, fmt.Errorf("couldn't read or create file %s. %w", filename, err)
	}

	db := &database[T]{
		filename: filename,
		data:     f,
	}

	return db, nil
}

func (db *database[T]) Store() error {
	return writeJsonFile(db.filename, db.data)
}

func (db *database[T]) AddOrUpdate(f func(elem T) bool, elem T) bool {
	idx := -1
	for i, d := range db.data {
		if f(d) {
			idx = i
			break
		}
	}

	if idx == -1 {
		db.data = append(db.data, elem)
		return false
	} else {
		db.data[idx] = elem
		return true
	}
}

func (db *database[T]) FindFirst(f func(elem T) bool) (T, bool) {
	idx := -1
	for i, d := range db.data {
		if f(d) {
			idx = i
			break
		}
	}

	if idx == -1 {
		var x T
		return x, false
	}
	return db.data[idx], true
}

func (db *database[T]) All() []T {
	return db.data
}

func (db *database[T]) Filter(f func(elem T) bool) []T {
	list := []T{}
	for _, d := range db.data {
		if f(d) {
			list = append(list, d)
		}
	}
	return list
}

func (db *database[T]) Exists(f func(elem T) bool) bool {
	found := false
	for _, d := range db.data {
		if f(d) {
			found = true
			break
		}
	}
	return found
}

func (db *database[T]) Remove(f func(elem T) bool) (T, bool) {
	idx := -1
	for i, d := range db.data {
		if f(d) {
			idx = i
			break
		}
	}

	var item T
	if idx == -1 {
		return item, false
	}

	item = db.data[idx]
	db.data = append(db.data[:idx], db.data[idx+1:]...)
	return item, true
}

func (db *database[T]) RemoveFilter(f func(elem T) bool) int {
	removed := 0
	for idx := len(db.data) - 1; idx >= 0; idx-- {
		elem := db.data[idx]
		if f(elem) {
			db.data = append(db.data[:idx], db.data[idx+1:]...)
			removed++
		}
	}
	return removed
}

func (db *database[T]) Add(elem T) bool {
	db.data = append(db.data, elem)
	return true
}

func (db *database[T]) AddIfUnique(f func(elem T) bool, elem T) bool {
	idx := -1
	for i, d := range db.data {
		if f(d) {
			idx = i
			break
		}
	}

	if idx != -1 {
		return false
	}

	return db.Add(elem)
}

func (db *database[T]) UpdateProperty(f func(elem T) bool, upd func(elem *T)) bool {
	idx := -1
	for i, d := range db.data {
		if f(d) {
			idx = i
			break
		}
	}

	if idx == -1 {
		return false
	}
	upd(&db.data[idx])
	// db.data[idx] = elem
	return true
}

func (db *database[T]) Update(f func(elem T) bool, elem T) bool {
	idx := -1
	for i, d := range db.data {
		if f(d) {
			idx = i
			break
		}
	}

	if idx == -1 {
		return false
	}
	db.data[idx] = elem
	return true
}

func (db *database[T]) Clear() bool {
	db.data = []T{}
	return true
}
