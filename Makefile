pacsym:
  go build pacsym.go
install:
  mkdir --parents /usr/pkg/pacsym/2.0/bin/
  mv pacsym /usr/pkg/pacsym/2.0/bin/pacsym
  mv pacsym /usr/pkgsrc/