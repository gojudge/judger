project(executer)

MESSAGE(STATUS "operation system is ${CMAKE_SYSTEM}")

set(CMAKE_BUILD_TYPE "Debug")
include_directories("deps")
if (APPLE)
	set(CMAKE_C_FLAGS_DEBUG "${CMAKE_C_FLAGS} -g -ggdb")
	set(SRC_LIST
		src/osx/config.c
		src/osx/executer.c
		deps/cJSON.c
	)
	set(TEST_LIST
		src/osx/test.c
	)
elseif(UNIX)
	set(CMAKE_C_FLAGS_DEBUG "${CMAKE_C_FLAGS} -g -ggdb")
	set(SRC_LIST
		src/nix/executer.c
		src/nix/config.c
		deps/cJSON.c
	)
	set(TEST_LIST
		src/nix/test.cpp
	)
elseif(WIN32)
	add_definitions("-D_CRT_SECURE_NO_DEPRECATE" "-D_DEBUG" "/wd 4819")
	set(SRC_LIST
		src/win/executer.c
		src/win/utils.c
		deps/cJSON.c
	)
endif(UNIX)

add_executable(executer ${SRC_LIST})

if(UNIX)
  target_link_libraries(executer m pthread)
elseif(WIN32)
endif(UNIX)
