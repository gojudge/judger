/**
 * File Name: utils.c
 * Author: rex
 * Mail: duguying2008@gmail.com 
 * Created Time: 2015年1月2日 星期五 20时54分12秒
 */

#include "executer.h"   

FILE* dopen(const char* filename){
    return fopen(filename,"w");
}

void dprintf(FILE* fd, const char *format,...){
    va_list args;

    va_start(args,format);
    vfprintf(fd,format,args);
    va_end(args);
    fflush(fd);
}

void dclose(FILE* fd){
    if (fd!=NULL){
      fclose(fd);
    }
}