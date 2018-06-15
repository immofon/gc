
view: o.svg
	firefox o.svg
	touch ref_test.go

o.svg: t.dot
	dot -Tsvg -o o.svg t.dot

t.dot: ref_test.go ref.go
	go test
