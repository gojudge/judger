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

CHAR *                             //   return error message
getLastErrorText(                  // converts "Lasr Error" code into text
CHAR *pBuf,                        //   message buffer
ULONG bufSize)                     //   buffer size
{
     DWORD retSize;
     LPTSTR pTemp=NULL;

     if (bufSize < 16) {
          if (bufSize > 0) {
               pBuf[0]='\0';
          }
          return(pBuf);
     }
     retSize=FormatMessage(FORMAT_MESSAGE_ALLOCATE_BUFFER|
                           FORMAT_MESSAGE_FROM_SYSTEM|
                           FORMAT_MESSAGE_ARGUMENT_ARRAY,
                           NULL,
                           GetLastError(),
                           LANG_NEUTRAL,
                           (LPTSTR)&pTemp,
                           0,
                           NULL );
     if (!retSize || pTemp == NULL) {
          pBuf[0]='\0';
     }
     else {
          pTemp[strlen(pTemp)-2]='\0'; //remove cr and newline character
          sprintf(pBuf,"%0.*s (0x%x)",bufSize-16,pTemp,GetLastError());
          LocalFree((HLOCAL)pTemp);
     }
     return(pBuf);
}
