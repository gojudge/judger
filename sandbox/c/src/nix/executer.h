#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <math.h>
#include <sys/reg.h>
#include <sys/user.h>
#include <sys/types.h>
#include <sys/time.h>
#include <sys/resource.h>
#include <sys/wait.h>
#include <sys/ptrace.h>
#include <sys/syscall.h>
#include <signal.h>
#include <sys/time.h>
#include <pthread.h>

extern int max_time;
extern int max_mem;
extern int array_len;
extern int allow_syscall[];

char* read_config(const char* filename);
void parse_config_json(char* text);
int free_config_buffer(char* buffer);

