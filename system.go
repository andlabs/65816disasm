// 25 may 2013
// modified from huc6280disasm
package main

type knownbyte struct {
	value	byte
	known	bool
}

type knownword struct {
	value	uint16
	known	bool
}

type envt struct {
	a		knownword
	carryflag	knownword
	direct	knownword
	dbr		knownbyte
	stack	[]knownbyte
}

var env *envt

func newenv() *evnt {
	e := new(envt)
	// TODO verify this
	e.a.known = false
	e.carryflag.value = 0
	e.carryflag.known = false
	// this part comes from the datasheet
	e.direct.value = 0
	e.direct.known = true
	e.dbr.value = 0
	e.dbr.known = true
	return e
}

func init() {
	env = newenv()
}

func makeAUnknown() {
	env.a.known = false
	env.carryflag.known = false
}

func push(value byte, known bool) {
	env.stack = append(env.stack, knownbyte{
		value:	value,
		known:	known,
	})
}

func pop() (value byte, known bool) {
	if len(env.stack) == 0 {
		return 0, false	// TODO correct?
	}
	t := env.stack[len(env.stack) - 1]
	stack = env.stack[:len(env.stack) - 1]
	return t.value, t.known
}

func pusha() {
	push(env.a, env.a.known)
}

func pushunknown() {
	push(env.a, false)		// value of a irrelevant
}

func popa() {
	env.a, env.a.known = pop()
}

func saveenv() *envt {
	e := new(envt)
	e.a = env.a
	e.carryflag = env.carryflag
	e.direct = env.direct
	e.dbr = env.dbr
	e.stack = make([]knownbyte, len(env.stack))
	copy(e.stack, env.stack)
	return e
}

func restoreenv(e *envt) {
	env = e
}
