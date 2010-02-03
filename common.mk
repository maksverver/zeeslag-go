BINS=test server generator
OBJS=generator.$X game.$X server.$X test.$X util.$X
GAME_SRC=game.go io.go player.go solver.go

all: $(BINS)

util.$X: util.go; $C -o $@ $<
game.$X: $(GAME_SRC) util.$X;  $C -o $@ $(GAME_SRC)
generator.$X: generator.go game.$X; $C -o $@ $<
server.$X: server.go game.$X; $C -o $@ $<
test.$X: test.go game.$X; $C -o $@ $<

test: game.$X test.$X; $L -o $@ test.$X
server: game.$X server.$X; $L -o $@ server.$X
generator: game.$X generator.$X; $L -o $@ generator.$X 

clean: ; rm -f $(OBJS)
distclean: clean; rm -f $(BINS)

.PHONY: all clean distclean
