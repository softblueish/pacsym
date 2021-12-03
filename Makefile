VER=2.2
pacsym:
  go build pacsym.go
install:
  mkdir --parents /usr/pkg/pacsym/$(VER)/bin/
  mv pacsym /usr/pkg/pacsym/$(VER)/bin/pacsym
  mkdir --parents /usr/pkgsrc/
  ln -sf /usr/pkg/pacsym/$(VER)/bin/pacsym /bin/
