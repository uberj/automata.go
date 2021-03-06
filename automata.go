package main

import (
    "fmt"
    "flag"
    "runtime"
)

const MAXREADS = 3

type State int

type PCell struct {
    next_state chan State
    state chan State
    idx int
    generation int
}

func (self* PCell) birth(N int) {
    /*
        N - Number of nodes
     */
    for {
        // Block on read of next_state
        next_state := <-self.next_state

        // Write appropriate number of states to state channel
        //  - 3 if inner node, 2 if edge node
        start, end, _ := calcStates(self.idx, N)
        for i := start; i < end + 1; i++ {
            self.state <- next_state
        }
        self.generation += 1
    }
}

type F1Cell struct {
    out chan State
    clk chan string
    idx int
    generation int
}

func (self* F1Cell) birth(state State, ps []PCell, rules []State, N int) {
    /*
        state - Initial State
        ps - array of PCells
        N - Number of nodes
     */
    var state_lookup State
    for {
        <-self.clk // Block here

        // Output state
        self.out <- state

        // Write state to P1 cells
        ps[self.idx].next_state <- state

        // Read Pcells into new state
        start, end, edge := calcStates(self.idx, N)
        state_lookup = 0
        for i := start; i < end + 1; i ++ {
            state_lookup = state_lookup << 1
            state_lookup = state_lookup + <-(ps[i].state)
        }
        if edge {
            // Edges are always zero
            state_lookup = (state_lookup << 1)
        }
        state = rules[state_lookup]
        self.generation += 1
    }
}

func calcStates (idx int, N int) (start int, end int, edge bool) {
    if idx == 0 {
        return 0, 1, true
    } else if idx == N - 1 {
        return idx - 1, idx, true
    } else {
        return idx - 1, idx + 1, false
    }
}

func InitCells (N int) ([]F1Cell, []PCell) {
    f1cells := make([]F1Cell, N)
    pcells := make([]PCell, N)

    for i := 0; i < N; i ++ {
        f1cells[i].idx = i
        f1cells[i].out = make(chan State)
        f1cells[i].clk = make(chan string)

        start, end, _ := calcStates(i, N)
        pcells[i].idx = i
        pcells[i].next_state = make(chan State)
        num_reads := end - start + 1
        pcells[i].state = make(chan State, num_reads)

    }
    fmt.Printf("\n")
    return f1cells, pcells
}

func consumeOut(f1cells []F1Cell, max_gen int) {
    to_consume := 0
    generation := 0
    f1cell_states := make([]State, len(f1cells))
    for {
        if to_consume <= 0 {
            generation += 1
            to_consume = len(f1cells)
            for i := range f1cells {
                f1cells[i].clk <- "go!"
            }
            fmt.Println(f1cell_states)
        }
        for i := range f1cells {
            select {
            case out := <-f1cells[i].out:
                f1cell_states[i] = out
                to_consume -= 1
            default:
            }
        }

        if generation == max_gen {
            return
        }
    }
}

func String(sa* []State) {
    fmt.Printf("State array: ")
    for i := range *sa {
        fmt.Printf("%d ", (*sa)[i])
    }
    fmt.Printf("\n")
}

func InitStateRules(rule int) []State {
    m := make([]State, 8)
    for i := range m {
        if rule % 2 == 1 {
            m[i] = 1
        } else {
            m[i] = 0
        }
        rule = rule >> 1
    }
    return m
}

func main () {
    /*
    N := 32
    rule := 30
    seed := 1 << 16
    */
    //var N int
    //var rule int
    //var seed State

    var N = flag.Int("N", 32, "Number of nodes")
    var rule = flag.Int("rule", 30, "Rule to be used for calculating state: i.e. rule 90 http://en.wikipedia.org/wiki/Rule_90")
    var seed = flag.Int("seed", 65536, "A seed value (powers of 2 are nice)")
    var max_gen = flag.Int("gen", 30, "Number of generations to iterate")
    var ncpu = flag.Int(
        "ncpu", runtime.NumCPU(), "GOMAXPROCS knob. Defaults to runtime.NumCPU() (the number of cores exposed to your user)")
    flag.Parse()

    var init_state State

    /* Set number of CPUs */
    runtime.GOMAXPROCS(*ncpu)

    /* Build structures */
    rules := InitStateRules(*rule)
    fmt.Println("Rule:", *rule)
    fmt.Println("Seed (decimal):", *seed)
    fmt.Println("Generations:", *max_gen)
    fmt.Println(rules)
    f1cells, pcells := InitCells(*N)

    /* Initialize soft state and run cells */
    fmt.Printf("Set seed to: ")
    for i := range f1cells {
        if *seed % 2 == 1 {
            init_state = 1
        } else {
            init_state = 0
        }
        *seed = *seed >> 1
        fmt.Printf("%d ", init_state)
        go f1cells[i].birth(init_state, pcells, rules, *N)
    }
    fmt.Printf("\n\n\n")

    for i := range pcells {
        go pcells[i].birth(*N)
    }
    consumeOut(f1cells, *max_gen)
}
