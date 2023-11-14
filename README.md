# distributed-systems-handin4

The program can be build with ```make build``` and ran with ```bin/program``` (linux)

... or by running ```go build src/program.go```.

To run the program you can use the following command line arguments:

    -debug
    
      Enable debug logging
      
    -file string
    
      The file which the critical section writes to

    -first
    
      Whether this node starts with the token
      
    -port uint
    
      The port of this node
      
    <ports...>
    
      A list of the all the ports in this network
      

Example:

```bin/program  -debug  -port 1111  -file "./log"  -first  1111 2222 3333 4444```

```bin/program  -debug  -port 2222  -file "./log"  1111 2222 3333 4444```

```bin/program  -debug  -port 3333  -file "./log"  1111 2222 3333 4444```

```bin/program  -debug  -port 4444  -file "./log"  1111 2222 3333 4444```

Note that each process will write an error if the server of the following node has not been started yet,
but it will keep trying every 3 seconds, so it can simply be ignored.
