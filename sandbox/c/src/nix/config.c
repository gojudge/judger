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

#define LIST_LEN 400

long max_time = 6000;
int max_mem = 1024;
int array_len = 0;
int allow_syscall[LIST_LEN];

/* Parse text to JSON, then render back to text, and print! */
int parse_config_json(char *text)
{
	cJSON *root;

	root = cJSON_Parse(text);
	if (!root) {
		printf("Error before: [%s]\n", cJSON_GetErrorPtr());
        return -1;
	} else {
		int i = 0;
		int time_tmp = cJSON_GetObjectItem(root, "time")->valueint;
		int mem_tmp = cJSON_GetObjectItem(root, "memory")->valueint;

		if (time_tmp >= 0) {
			max_time = (long int)time_tmp;
		}
		if (mem_tmp >= 0) {
			max_mem = mem_tmp;
		}
#ifdef __i386__
		cJSON *list = cJSON_GetObjectItem(root, "x86");
#else
		cJSON *list = cJSON_GetObjectItem(root, "x64");
#endif

		array_len = cJSON_GetArraySize(list);
		for (i = 0; i < array_len; i++) {
			cJSON *item = cJSON_GetArrayItem(list, i);
			allow_syscall[i] = item->valueint;
		}

		cJSON_Delete(root);
	}

    return 0;
}

/* read config file */
char *read_config(const char *filename)
{
	int len = 0;
	char *buffer;

	//printf("[filename]\n%s\n", filename);
	FILE *f = fopen(filename, "rb");
	fseek(f, 0, SEEK_END);
	long length = ftell(f);
	fseek(f, 0, SEEK_SET);

	if (length <= 0) {
		return (char *)NULL;
	}
	buffer = malloc(length + 1);

	memset(buffer, 0, length + 1);
	len = fread(buffer, length, 1, f);
	fclose(f);
	if (len == 0) {
		free(buffer);
		return (char *)NULL;
	}
	return buffer;
}

/* free config string buffer */
int free_config_buffer(char *buffer)
{
	free(buffer);
	return 0;
}
