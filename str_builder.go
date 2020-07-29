package gsql

import "strings"

type sqlBuilder struct {
	strings.Builder
}

func(s *sqlBuilder)join(sep string, a... string) {
	switch len(a) {
	case 0:
		return
	case 1:
		s.Builder.WriteString(a[0])
		return
	}
	n := len(sep) * (len(a) - 1)
	for i := 0; i < len(a); i++ {
		n += len(a[i])
	}
	s.Builder.Grow(n)
	s.Builder.WriteString(a[0])
	for _, ss := range a[1:] {
		s.Builder.WriteString(sep)
		s.Builder.WriteString(ss)
	}
}

func (s *sqlBuilder)joinWithOmit(sep string, a[]string, omit map[string]int) {
	switch len(a) {
	case 0:
		return
	case 1:
		if _,ok := omit[a[0]]; !ok {
			s.Builder.WriteString(a[0])
		}
		return
	}
	n := len(sep) * (len(a) - 1)
	for i := 0; i < len(a); i++ {
		n += len(a[i])
	}
	s.Builder.Grow(n)
	var i int
	for i = 0; i< len(a); i++ {
		if _, ok := omit[a[i]]; ok {
			continue
		}
		s.Builder.WriteString(a[i])
	}
	for _, ss := range a[i+1:] {
		if _, ok := omit[a[i]]; ok {
			continue
		}
		s.Builder.WriteString(sep)
		s.Builder.WriteString(ss)
	}
}

func (s *sqlBuilder)writeStrings(args... string) {
	for _,v := range args {
		s.Builder.WriteString(v)
	}
}