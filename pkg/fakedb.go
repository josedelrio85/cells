package apic2c

import "sync"

// FakeDb is a struct used to test Db functionality with fake methods.
type FakeDb struct {
	OpenFunc         func() error
	OpenCalls        int
	CloseFunc        func() error
	CloseCalls       int
	CreateTableFunc  func(table interface{}) error
	CreateTableCalls int
	UpdateFunc       func(element interface{}, wCond string, wFields []string) error
	UpdateCalls      int
	sync.Mutex
}

// Open is a method to test Open function
func (f *FakeDb) Open() error {
	f.Lock()
	defer f.Unlock()
	f.OpenCalls++
	return f.Open()
}

// Close is a method to test Close function
func (f *FakeDb) Close() {
	f.Lock()
	defer f.Unlock()
	f.CloseCalls++
	f.CloseFunc()
}

// CreateTable is a method to test CreateTable function
func (f *FakeDb) CreateTable(table interface{}) error {
	f.Lock()
	defer f.Unlock()
	f.CreateTableCalls++
	return f.CreateTableFunc(table)
}

// Update is a method to test Update function
func (f *FakeDb) Update(element interface{}, wCond string, wFields []string) error {
	f.Lock()
	defer f.Unlock()
	f.UpdateCalls++
	return f.Update(element, wCond, wFields)
}
