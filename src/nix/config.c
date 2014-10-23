/**
 * File Name: config.c
 * Author: rex
 * Mail: duguying2008@gmail.com 
 * Created Time: 2014年10月23日 星期四 11时38分01秒
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "cJSON.h"

#define LIST_LEN 100

int time = 60;
int mem = 1024;
int array_len = 0;
int forbidden_syscall[LIST_LEN];

/* Parse text to JSON, then render back to text, and print! */
void parse_config_json(char *text)
{
  cJSON *root;
      
  root=cJSON_Parse(text);
  if (!root) {
    printf("Error before: [%s]\n",cJSON_GetErrorPtr());
  }
  else
  {
    int i = 0;
    int time_tmp = cJSON_GetObjectItem(root,"time")->valueint;
    int mem_tmp = cJSON_GetObjectItem(root,"memory")->valueint;

    if(time_tmp>=0){
      time = time_tmp;
    }
    if(mem_tmp>=0){
      mem = mem_tmp;
    }

    #ifdef __i386__
    cJSON *list = cJSON_GetObjectItem(root,"x86");
    #else
    cJSON *list = cJSON_GetObjectItem(root,"x64");
    #endif

    array_len = cJSON_GetArraySize(list);
    for(i=0; i<array_len; i++)
    {
      cJSON *item=cJSON_GetArrayItem(list,i);
      forbidden_syscall[i] = item->valueint;
    }

    cJSON_Delete(root);
  }
}

/* read config file */
char* read_config(const char* filename){
  int len = 0;
  char* buffer;

  FILE *f=fopen(filename,"rb");fseek(f,0,SEEK_END);long length=ftell(f);fseek(f,0,SEEK_SET);

  if(length <= 0){
    return (char*)NULL;
  }
  buffer = malloc(length + 1);

  memset(buffer,0,length + 1);
  len=fread(buffer, length, 1, f);
  fclose(f);
  if(len == 0){
    free(buffer);
    return (char*)NULL;
  }
  return buffer;
}

/* free config string buffer */
int free_config_buffer(char* buffer){
  free(buffer);
  return 0;
}

