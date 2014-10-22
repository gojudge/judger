/**
 * File Name: config.c
 * Author: rex
 * Mail: duguying2008@gmail.com 
 * Created Time: 2014年10月22日 星期三 23时51分28秒
 */

#include "stdio.h"
#include "stdlib.h"
#include "string.h"
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>

#include "mjson.h"

#define MAXCHANNELS 72
#define BUF_LEN 1024

static int X86_SYSCALL_NUM[MAXCHANNELS];
static int X64_SYSCALL_NUM[MAXCHANNELS];
static int visible86,visible64;

static int max_time,max_mem;

const struct json_attr_t json_attrs_sky[] = {
      {"time",t_integer, .addr.integer = &max_time },

      {"memory",t_integer, .addr.integer = &max_mem },

      {"x86", t_array,   .addr.array.element_type = t_integer,
                         .addr.array.arr.integers.store = X86_SYSCALL_NUM,
                         .addr.array.maxlen = MAXCHANNELS,
                         .addr.array.count = &visible86},

      {"x64", t_array,   .addr.array.element_type = t_integer,
                         .addr.array.arr.integers.store = X64_SYSCALL_NUM,
                         .addr.array.maxlen = MAXCHANNELS,
                         .addr.array.count = &visible64},

      {NULL},
};

int read_config()
{
      int i, status = 0, len = 0;
      size_t length = 0;

      char buffer[BUF_LEN];
      FILE* fd;
      struct stat buf;
      char* filename = "executer.json";

      if(stat(filename, &buf)<0)
      {
        return -1;
      }
      // get length of file
      length = (unsigned long)buf.st_size;

      if((fd=fopen(filename,"r"))==NULL){
        printf("open failed!\n");
        return -1;
      }

      memset(buffer,0,BUF_LEN);
      len=fread(buffer, length, 1, fd);

      status = json_read_object(buffer, json_attrs_sky, NULL);

      printf("max time: %d\n", max_time);

      printf("max memory: %d\n", max_mem);

      printf("%d x86:\n", visible86);
      for (i = 0; i < visible86; i++)
        printf("syscall: %d\n", X86_SYSCALL_NUM[i]);

      printf("%d x64:\n", visible64);
      for (i = 0; i < visible64; i++)
        printf("syscall: %d\n", X64_SYSCALL_NUM[i]);

      if (status != 0)
      puts(json_error_string(status));

      fclose(fd);
      return 0;
}

