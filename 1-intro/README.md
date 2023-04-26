

# CPL - Channel-based Concurrency Module 

## Lab Class \#1 (Clocks, Cats and Primes)


Your assignment repository contains some startup files for the first lab. assignment of the Go-based module.

You have 2 small problems to solve, which will act as a bit of a warm-up to programming in Go.

To submit your answers, simply push your files onto the
repository. Some problems will require you to modify existing files
and add new ones. The problems will not be graded, but we will use a
similar system for the mini-project next week and the project the week after.

---

### Setup 

Naturally, you will need to have a recent [Go distribution](https://golang.org/dl/) installed on your
machine.

You can check that the ``go`` tool is available on your shell's path by typing
``go version``.

## Go Modules



### Problem 1 -- Clocks and Cats

The file ``clockserver/clockserver.go`` contains a simple, **non-concurrent** implementation of a clock server. The clock server listens for TCP connections on port ``8080``. Once a connection is established, the server sends the current time every second. The server is not programmed to do any sort of sophisticated error handling, simply logging any errors as they arise. 


Try to **build and run the clock server**. To see what the server is sending, use the command ``nc localhost 8080`` if you are using an Unix-based machine or (build and run) the ``netcat/netcat.go`` program. The supplied ``netcat`` is hardcoded to listen on port ``8080``. 

  1. You may have noticed that the clock server can only handle a single client at a time. Make the clock server accept requests concurrently. This should be **very** easy.
  2. Modify the clock server to accept a port number and write a client program (place it under ``clockclient/clockclient.go``) that receives from multiple clocks simultaneously (e.g. in different timezones), displaying the results in some reasonably formated way. For the client to know on which address to listen to, use a reasonable command-line argument syntax (e.g. ``clockclient [ClockName=ip:port]+``).
  
  **Note:** In a Unix-based system you can "fake" the timezone of the clock server by changing the environment variable ``TZ``. For instance, ``TZ=Asia/Seoul ./clockserver`` will run a clockserver with the timezone of Seoul, South Korea.

----

  ### Problem 2 - Concurrent Primes

  Now that you are a bit more familiar with Go, its time to focus a
  on channels.

  Write a _pipeline_ of (infinitely running) goroutines that calculate
  and subsequently print out
  prime numbers, in sequence, using a [prime
  sieve](https://en.wikipedia.org/wiki/Sieve_of_Eratosthenes). Place
  your implementation under ``primes/sieve.go``.
  
  You will likely want to structure your program using three "kinds"
  of goroutines (connected using channels): an initial ``Producer``
  routine that simply emits the stream of all natural numbers
  (starting at 2) in sequence; a chain of ``Sieve`` routines,
  connected by channels, that
  forwards the number stream (i.e., the first ``Sieve`` receives from
  the ``Producer`` and sends to the next ``Sieve``, and so on),
  filtering out numbers that are divisible by a given (prime) number;
  and, an assembling routine that implements the chaining of
  ``Sieves`` and prints out the numbers in succession.
  
  You may want to alter the implementation described above so that it
  can be [tested](https://golang.org/pkg/testing/). Write a (very inefficient)
  test that validates a few of the output numbers.

----

That's it! Don't forget to push!
   
    
 
