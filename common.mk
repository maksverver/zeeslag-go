BINS=test server
OBJS=game.$X test.$X server.$X

all: $(BINS)

game.$X:     game.go io.go player.go solver.go;  $C -o $@ $^
test.$X:     test.go;                            $C -o $@ $<
server.$X:   server.go;                          $C -o $@ $<

test:    game.$X test.$X;    $L -o $@ test.$X
server:  game.$X server.$X;  $L -o $@ server.$X

clean: ; rm -f $(OBJS)
distclean: clean; rm -f $(BINS)

.PHONY: all clean distclean
