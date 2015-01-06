#include <stddef.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <stdint.h>
#include <assert.h>
#include <windows.h>
#include <psapi.h>

#pragma comment(lib, "psapi.lib")

FILE* dopen(const char* filename);
void dprintf(FILE* fd, const char *format,...);
void dclose(FILE* fd);
CHAR *getLastErrorText(CHAR *pBuf, ULONG bufSize);
