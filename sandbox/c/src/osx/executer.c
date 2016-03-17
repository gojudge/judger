/**
 * File Name: executer.c
 * Author: rex
 * Mail: duguying2008@gmail.com 
 * Created Time: 2014年10月21日 星期二 23时05分12秒
 */

#include "executer.h"

#define VERSION "1.1.2"

pid_t child;
long begin_time;
char *executable = NULL;
size_t PATH_LEN = 1024;
int fd = 0;
char *config_path = NULL;
int judger_model = 2;			//default model - assert
char *input = NULL;
char *output = NULL;

enum ecode {
	PEN,						// Exit Normally
	PRE,						// Runtime Error
	POM,						// Out of Memory
	POT,						// Out of Time
	POL,						// Output Limit Exceed
	PSF							// Syscall Forbidden
};

/* print error */
void PRTERR()
{
	extern int errno;
	char *message;

	printf("errno [%d]\n", errno);
	message = strerror(errno);
	printf("Mesg: %s\n", message);
}

/* now, get now time of microsecond */
long t_now()
{
	long tmp_now = 0;
	struct timeval tv;

	memset(&tv, 0, sizeof(struct timeval));
	gettimeofday(&tv, NULL);

	tmp_now = (long)tv.tv_sec * 1000 + (long)tv.tv_usec / 1000;

	return tmp_now;
}

/* record content into file */
void record_result(const char *content)
{
	FILE *stream;
	if ((stream = fopen("RUNRESULT", "w")) == NULL) {
		fprintf(stderr, "Cannot open output file.\n");
	}
	char *s = (char *)malloc(sizeof(char) * 10);
	memset(s, 0, sizeof(char) * 10);
	strncpy(s, content, sizeof(char) * 10);
	s[10] = 0;

	fwrite(s, sizeof(char) * strlen(s), 1, stream);
	fclose(stream);
}

/* process exit */
void pexit(enum ecode EC)
{
	if (EC == PEN) {
		// Ignore
		record_result("PEN");
	} else if (EC == POT) {
		printf("Out of Time.\n");
		record_result("POT");
		kill(child, SIGKILL);
	} else if (EC == PSF) {
		printf("Syscall Forbidden.\n");
		record_result("PSF");
		kill(child, SIGKILL);
	} else if (EC == POM) {
		printf("Out of Memory.\n");
		record_result("POM");
		kill(child, SIGKILL);
	} else if (EC == PRE) {
		printf("Runtime Error.\n");
		record_result("PRE");
	} else if (EC == POL) {
		printf("Output Limit Exceed.\n");
		record_result("POL");
		kill(child, SIGKILL);
	} else {
		dprintf(fd, "EC [%d]\n", EC);
		record_result("EC unk");
	}

	close(fd);
	exit(0);
}

/* if the syscall is forbidden, return 0 */
//int check_syscall(int syscall)
//{
//	int i = 0;
//	for (i = 0; i < array_len; i++) {
//		if (syscall == allow_syscall[i]) {
//			return 1;			//true, the syscall matched one of the list, pass
//		};
//	}
//	return 0;					//false, not matched
//}

/* get memory used */
int get_memory_used(int pid)
{
	FILE *fps;
	char ps[32];
	int memory;

	sprintf(ps, "/proc/%d/statm", pid);
	fps = fopen(ps, "r");
	int i;
	for (i = 0; i < 6; i++)
		fscanf(fps, "%d", &memory);
	fclose(fps);

	int pagesize = getpagesize() / 1024;
	memory *= pagesize;
	return memory;
}

/* timer, when over time, killed son and program exit */
void *time_watcher(void *unused)
{
	while (1) {
		long now_time = t_now();
		if (now_time - begin_time - (long)max_time > 0) {
			dprintf(fd, "over time [%lu], killed!\n", now_time);
			pexit(POT);
		}
		sleep(0);
	}
}

/* parse command args */
void parse_args(int argc, char *argv[])
{
	int i = 0;
	int len = 0;
	char *arg = NULL;
	char *buff = NULL;
	const int BUF_LEN = 128;
	char *tag_name = NULL;
	char *tag_value = NULL;

	buff = (char *)malloc(sizeof(char) * BUF_LEN);

	for (i = 1; i < argc; i++) {
		memset(buff, 0, sizeof(char) * BUF_LEN);
		strncpy(buff, argv[i], strlen(argv[i]) + 1);

		if (buff[0] == '-') {	// options
			tag_name = strtok(buff + 1, "=");	// string time
			tag_value = strtok(NULL, "=");	// decemal time value

			if (!strcmp(tag_name, "t")) {	// time
				int tmp_time = atoi(tag_value);
				if (tmp_time > 0) {
					max_time = tmp_time;
				} else {
					dprintf(fd, "invalid time [%d], use default.\n", tmp_time);
				}
			} else if (!strcmp(tag_name, "m")) {	// memory
				int tmp_mem = atoi(tag_value);
				if (tmp_mem > 0) {
					max_mem = tmp_mem;
				} else {
					dprintf(fd, "invalid memory [%d], use default.\n", tmp_mem);
				}
			} else if (!strcmp(tag_name, "c")) {
				len = strlen(tag_value);
				memset(config_path, 0, PATH_LEN);
				strncpy(config_path, tag_value, len);
				config_path[len] = 0;
			} else if (!strcmp(tag_name, "j")) {
				dprintf(fd, "[judger model] %s\n", tag_value);

				if (!strcmp(tag_value, "io")) {
					judger_model = 1;
				} else {
					judger_model = 2;
				}
			} else if (!strcmp(tag_name, "-stdin")) {
				dprintf(fd, "[input] %s\n", tag_value);

				len = strlen(tag_value);
				memset(input, 0, PATH_LEN);
				strncpy(input, tag_value, len);
				input[len] = 0;
			} else if (!strcmp(tag_name, "-stdout")) {
				dprintf(fd, "[output] %s\n", tag_value);

				len = strlen(tag_value);
				memset(output, 0, PATH_LEN);
				strncpy(output, tag_value, len);
				output[len] = 0;
			}

		} else {				// executable path, just one
			int len = strlen(argv[i]);
			memset(executable, 0, PATH_LEN);
			strncpy(executable, argv[i], len);
			executable[len] = 0;
		}

	}

	return;
}

/* check file exist */
int file_exist(const char *filepath)
{
	if (access(filepath, F_OK) == -1) {
		// File not exists
		return -1;
	} else {
		return 0;
	}
}

int main(int argc, char *argv[])
{
	long orig_eax;

	// alloc memory for path string
	executable = (char *)malloc(PATH_LEN);
	memset(executable, 0, PATH_LEN);

	// alloc memory for config path string
	config_path = (char *)malloc(PATH_LEN);
	memset(config_path, 0, PATH_LEN);

	// alloc stdin/stdout file path
	input = (char *)malloc(PATH_LEN);
	memset(input, 0, PATH_LEN);

	output = (char *)malloc(PATH_LEN);
	memset(output, 0, PATH_LEN);

	// set stdin/stdout default filename
	char *tmp_input = "stdin";
	char *tmp_output = "stdout";

	size_t len_input = strlen(tmp_input);
	size_t len_output = strlen(tmp_output);

	strncpy(input, tmp_input, len_input);
	strncpy(output, tmp_output, len_output);

	input[len_input] = 0;
	output[len_output] = 0;

	// show help
	if (argc < 2) {
		printf("\033[0;39;1mSandbox for Linux Native\033[0m\n"
			   "Usage: executer <option> <command>\n"
			   "option:\n"
			   "  \033[0;33m-t=time\033[0m          program max time\n"
			   "  \033[0;33m-m=mem\033[0m           program max memory\n"
			   "  \033[0;33m-c=path\033[0m          config file path\n"
			   "  \033[0;33m-j=model\033[0m         judger model[io/assert]\n"
			   "  \033[0;33m--stdin=path\033[0m     stdin file path\n"
			   "  \033[0;33m--stdout=path\033[0m    stdout file path\n"
			   "\033[0;32mversion " VERSION "\033[0m\n");
		return 0;
	} else {
        int err = 0;

		fd = 0;
		fd = open("executer.debug", O_CREAT | O_RDWR);

		char *tmp_config_name = "executer.json";
		strncpy(config_path, tmp_config_name, strlen(tmp_config_name));
		config_path[strlen(tmp_config_name)] = 0;

		parse_args(argc, argv);

		//read config
        dprintf(fd, "[config] %s\n", config_path);
		char *config_string = read_config(config_path);
		err = parse_config_json(config_string);
        
        if(err){
            dprintf(fd, "[error] json format config file parse error.\n");
            return err;
        }

		free_config_buffer(config_string);
		free(config_path);
	}

	child = fork();
	if (child == 0) {
		int exec_result = 0;

		if (judger_model == 1) {
			//Redirect the standard input / output stream
			if (file_exist("stdin") == -1) {
				dprintf(fd, "[Warning] \"stdin\" file does not exist!\n");
			}

			freopen(input, "r", stdin);
			freopen(output, "w", stdout);

			free(input);
			free(output);
		}

//		ptrace(PTRACE_TRACEME, 0, NULL, NULL);

		// must use execl for supporting segmentfault check
		exec_result = execl(executable, "", (char *)NULL);
		if (-1 == exec_result) {
			printf("execute [%s] failed!", executable);
		}

		free(executable);

		exit(0);
	} else {
		struct rusage rinfo;
		int runstat, i = 0;
		pthread_t thread_id;

		dprintf(fd, "the child pid is %d\n", child);

		begin_time = t_now();
		dprintf(fd, "begin time [%lu]\n", begin_time);
		dprintf(fd, "max_time [%lu]\nmax_mem [%d]\nexecutable path [%s]\n",
				max_time, max_mem, executable);

		// a new thread for timer, when over time, killed and exit
		pthread_create(&thread_id, NULL, &time_watcher, NULL);

		for (;;) {
			//time_t now_time;
			wait4(child, &runstat, 0, &rinfo);

			if (WIFEXITED(runstat)) {
				int exitcode = WEXITSTATUS(runstat);
				dprintf(fd, "exitcode [%d]\n", exitcode);
				if (exitcode != 0) {
					//Runtime Error
					dprintf(fd, "Runtime Error\n");
					pexit(PRE);
				}
				//normal exit
				dprintf(fd, "Exit Normally.\n");
				pexit(PEN);
			} else if (WIFSIGNALED(runstat)) {
				// call kill(pid, SIGKILL)
				// Ignore
				exit(0);
			} else if (WIFSTOPPED(runstat)) {
				int signal = WSTOPSIG(runstat);

				if (signal == SIGTRAP) {
					printf("[signal] %d", signal);
//					struct user_regs_struct reg;
//					int syscall;
//					static int executed = 0;

//					ptrace(PTRACE_GETREGS, child, NULL, &reg);
//#ifdef __i386__
//					syscall = reg.orig_eax;
//#else
//					syscall = reg.orig_rax;
//#endif

//					dprintf(fd, "syscall: %d\n", syscall);

					// syscall check 
//					if (!check_syscall(syscall)) {
//						dprintf(fd, "Syscall [%d] is Forbidden.\n", syscall);
//
//						pexit(PSF);
//					}

				} else if (signal == SIGUSR1) {
					// Ignore
				} else if (signal == SIGXFSZ) {
					dprintf(fd, "Output Limit Exceed.\n");

					pexit(POL);
				} else {
					dprintf(fd, "Runtime Error.\n");

					pexit(PRE);
				}

			}

//			ptrace(PTRACE_SYSCALL, child, NULL, NULL);

			// check memory use
			int use_mem = get_memory_used(child);
			if (use_mem > max_mem) {
				dprintf(fd, "Out of Memory [%d]\n", use_mem);
				pexit(POM);
			}

		}

	}

	return 0;
}
