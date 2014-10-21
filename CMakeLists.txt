project(judger)

if(UNIX)
	set(CMAKE_C_FLAGS_DEBUG "${CMAKE_C_FLAGS} -g -ggdb")
	set(SRC_LIST 
		  src/nix/judger.cpp
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
