X=6 
CL=6l

%.8: %.go; 8g -o $@ $<
%.6: %.go; 6g -o $@ $<

OBJS=game.$X io.$X solver.$X test.$X server.$X
BINS=test server

all: $(OBJS) $(BINS)

test: $(OBJS); $(CL) -o $@ test.$X
server: $(OBJS); $(CL) -o $@ server.$X

clean: ; rm -f $(OBJS)
distclean: clean; rm -f $(BINS)

.PHONY: all clean distclean
