// Package hystrix is a latency and fault tolerance library designed to isolate
// points of access to remote systems, services and 3rd party libraries, stop 
// cascading failure and enable resilience in complex distributed systems where 
// failure is inevitable.
//
// Based on the java project of the same name, by Netflix. https://github.com/Netflix/Hystrix
package hystrix

// Result is the standard response structure for commands.  Either a Result or Error will be defined.  Fallbacks also generate Results.
type Result struct {
	Result interface{}
	Error  error
}
