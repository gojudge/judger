project(judger)

include_directories("deps")
if(UNIX)
	set(CMAKE_C_FLAGS_DEBUG "${CMAKE_C_FLAGS} -g -ggdb")
	set(SRC_LIST 
		  src/nix/judger.c
      src/nix/config.c
      deps/mjson.c
     )
  set(TEST_LIST
      src/nix/test.cpp
      )
elseif(WIN32)
	add_definitions("-D_CRT_SECURE_NO_DEPRECATE" "-D_DEBUG")
	set(SRC_LIST 
		  src/win/judger.cpp
		  )
endif(UNIX)

add_executable(judger ${SRC_LIST})
add_executable(test ${TEST_LIST})
#add_executable(mjson ${MJSON_LIST})
