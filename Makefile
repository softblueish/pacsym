VER=2.5

all: pacsym

pacsym:
	go build pacsym.go

clean:
	@rm -f pacsym

install:
	@mkdir -pv /usr/pkg/pacsym/$(VER)/bin/
	@mv pacsym /usr/pkg/pacsym/$(VER)/bin/pacsym -v
	@mkdir -p /usr/pkgsrc/
	@ln -svf /usr/pkg/pacsym/$(VER)/bin/pacsym /bin/