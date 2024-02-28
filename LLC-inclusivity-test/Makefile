TARGETS=cacheEvictor cacheChecker
UTILS=util.o

all :$(TARGETS)

$(UTILS): %.o: %.c %.h
	$(CC) $(CFLAGS) -c $<

%.o: %.c util.h
	$(CC) $(CFLAGS) -c $<

$(TARGETS): %:%.o util.o
	$(CC) $(CFLAGS) $^ -o $@


run_Evictor: cacheEvictor
	@taskset -c 1 ./cacheEvictor

run_Checker: cacheChecker
	@taskset -c 2 ./cacheChecker

