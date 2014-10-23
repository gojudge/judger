project(executer)

include_directories("deps")
if(UNIX)
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
	add_definitions("-D_CRT_SECURE_NO_DEPRECATE" "-D_DEBUG")
	set(SRC_LIST 
		  src/win/executer.c
		  )
endif(UNIX)

add_executable(executer ${SRC_LIST})
add_executable(test ${TEST_LIST})
#add_executable(mjson ${MJSON_LIST})
if(UNIX)
  target_link_libraries(executer m)
elseif(WIN32)
endif(UNIX)
