CC=gcc
CFLAGS=-Wall -Wextra -pedantic -w 

BUILDDIR=build
SRC=$(wildcard *.c)
OBJ=$(patsubst %.c, %.o, $(SRC))

BINARYNAME=main
BINARY=$(BINARYNAME)

$(BINARY): $(OBJ)
	$(CC) $(CFLAGS) $(INCLUDES) $(OBJ) -o $(BUILDDIR)/$(BINARY) $(LIBS)

clean:
	rm -rf $(BINARY)
	rm -rf $(OBJ)

