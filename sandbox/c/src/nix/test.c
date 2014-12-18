/**
 * File Name: test.cpp
 * Author: rex
 * Mail: duguying2008@gmail.com 
 * Created Time: 2014年10月21日 星期二 23时05分12秒
 */

#include "stdio.h"
#include "stdlib.h"

struct mystruct {
	int i;
	char cha;
};

int main(void)
{
	FILE *stream;
	struct mystruct s;
	if ((stream = fopen("TEST.$$$", "wb")) == NULL) {	/* open file TEST.$$$ */
		fprintf(stderr, "Cannot open output file.\n");
		return 1;
	}
	s.i = 0;
	s.cha = 'A';
	fwrite(&s, sizeof(s), 1, stream);	/* 写的struct文件 */
	fclose(stream);				/*关闭文件 */

	for (;;) {
		sleep(1);
	}

	return 0;
}
