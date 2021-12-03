VER=2.2
pacsym:
  go build pacsym.go
install:
  mkdir --parents /usr/pkg/pacsym/$(VER)/bin/
  mv pacsym /usr/pkg/pacsym/$(VER)/bin/pacsym
  mv pacsym /usr/pkgsrc/
  ln -sf /usr/pkg/pacsym/$(VER)/bin/pacsym /bin/