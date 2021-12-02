pacsym:
  go build pacsym.go
install:
  mkdir --parents /usr/pkg/pacsym/1.0/bin/
  mv pacsym /usr/pkg/pacsym/1.0/bin/pacsym
