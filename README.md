## Genetic

A Golang library for running simple genetic algorithms.

### Usage

Define a Type that adhears to the `Chromosome` interface.

~~~go
type Chromosome interface {
  Fitness() int16
  Recombine(newPopulation chan<- Chromosome, chromosome Chromosome)
  Mutate() Chromosome
}
~~~

Here is an example implementation of a Chromosome. Given a `goal` string,
a Chromosome's fitness is the 'distance' between it and the goal.
A fitness of 0 is the best. Chromosomes are sorted by `Fitness()` ascending.

~~~go
package main

import (
  "fmt"
  "math/rand"
  "time"
  "github.com/hstove/genetic"
)

const chars = "abcdefghijklmnopqrstuvwxyz"
const goal = "hankstoeverisabadgolangprogrammer"

// Chromosome Implementation
type RandomString struct {
  Member string
}

var memoizedFitness = make(map[string]int16)

func (this RandomString) String() (str string){
  return fmt.Sprintf("{Member: %s | Fitness: %d}", this.Member, this.Fitness())
}

func NewRandomString() *RandomString {
  member := ""
  for i := 0; i < len(goal); i++ {
    char := string(chars[rand.Intn(len(chars))])
    member += char
  }
  return &RandomString{Member: member}
}

func (r RandomString) Fitness() int16{
  if val, ok := memoizedFitness[r.Member]; ok {
    // fmt.Println(len(memoizedFitness))
    return val
  }
  fitness := int16(0)
  for i := 0; i < len(goal); i++ {
    difference := int16(goal[i]) - int16(r.Member[i])
    if difference < 0 {
      fitness += -1 * difference
    } else {
      fitness += difference
    }
  }
  fitness *= -1
  memoizedFitness[r.Member] = fitness
  return fitness
}

func (r RandomString) Recombine(newPopulation chan<- genetic.Chromosome, chromosome genetic.Chromosome) {
  var other *RandomString = chromosome.(*RandomString)
  length := len(goal)
  split := rand.Intn(length-1)
  newMember := string(r.Member[0:split])
  newMember += string(other.Member[split:length])
  newPopulation <- &RandomString{newMember}
}

func (r RandomString) Mutate() genetic.Chromosome {
  return NewRandomString()
}
~~~

Then, create a `Population` with a bunch of Chromosomes in it. `Evolve()` until satisfied.

~~~go
func main(){
  rand.Seed(time.Now().UTC().UnixNano())
  popSize := genetic.MaxPopulation
  var population genetic.Population
  for i := 0; i < popSize; i++ {
    population = append(population, NewRandomString())
  }
  i := 0
  for {
    i++
    population = population.Evolve()
    fmt.Println(i, population.BestFit(), population[len(population) - 1], len(population))
    if population.BestFit().Fitness() == 0 {
      fmt.Println(fmt.Sprintf("Found Champion after %d Generations", i))
      break
    }
  }
}
~~~

This outputs:

~~~
1 {Member: bqlijpphlfphvaemfhheafmopgjiggmbd | Fitness: -158} {Member: ikeyyfwbbgofgxowbgkiqpndjulokndbb | Fitness: -270} 10000
2 {Member: bqlijpphlfphvaemfhheafmopglodkkfd | Fitness: -147} {Member: dtjltogkiyckpymtckirnoelxuizererw | Fitness: -249} 10000
3 {Member: iatqvxcfugtiojfaeuojflvpqvjuhoogs | Fitness: -122} {Member: wupvqtpdmbhjszsfgujlkdxlizbueqcdo | Fitness: -235} 10000
4 {Member: iatqvxcfugtiojfaeuojflvpqvjuhoogs | Fitness: -122} {Member: qlkyqwjosdpfsvjsrxdinsfoqsdglqidk | Fitness: -224} 10000
5 {Member: sfnvkslewcoizdlachmycmhxyqiocnsen | Fitness: -121} {Member: memlkzogvfhomqfesswnshfncojfhqdzn | Fitness: -213} 10000
6 {Member: iaogkslewcoizdlachmycmhxyqiocnstv | Fitness: -115} {Member: hggmrixdypzcrkmtiwudimgkmkjzksebt | Fitness: -204} 10000
7 {Member: iaogkslewcoizdlachmycmhxyqiocnstv | Fitness: -115} {Member: qgmjorcyzjvppedghqeflgitrqdmhqvvv | Fitness: -195} 10000
8 {Member: iaogkslewcoizdlachmycmhxyqiocnmiv | Fitness: -98} {Member: dtnhvjsgsavjdfxcdaeipnxqpzgobckht | Fitness: -187} 10000
9 {Member: iaogkslewcoizdlachmycmhxyqiocnmiv | Fitness: -98} {Member: tqvnsrrhshpmyohcocgkmkagomdwenuop | Fitness: -179} 10000
...
60 {Member: hamkssoewerisabadgolangqrogrammer | Fitness: -4} {Member: hamkstnevfrjrabadgolangotnhqammer | Fitness: -11} 10000
61 {Member: hankstoeverisabacgolangqrogrammdr | Fitness: -3} {Member: hamkstogvfrjtabacgolangorogrammds | Fitness: -10} 10000
62 {Member: hankstodwerisabadgolangqrogrammer | Fitness: -3} {Member: hamkstodwerjrabacgolangorogrammds | Fitness: -9} 10000
63 {Member: hankstoeterisabadgolangorogrammer | Fitness: -3} {Member: hamlruoeverjrabacgolangprogrammfr | Fitness: -8} 10000
64 {Member: hankstoeverjsabadgolangqrogrammer | Fitness: -2} {Member: hamkstoexerjrabacgolangorpgrammer | Fitness: -8} 10000
65 {Member: hankstoeverhsabadgolangprogrammer | Fitness: -1} {Member: hankstoeteriqabacgolangoroframmer | Fitness: -7} 10000
66 {Member: hankstoeverhsabadgolangprogrammer | Fitness: -1} {Member: hankstneverjrabacholangqrogrammer | Fitness: -6} 10000
67 {Member: hankstoevfrisabadgolangprogrammer | Fitness: -1} {Member: hamkssoewerisabacholangorogrammer | Fitness: -6} 10000
68 {Member: hankstneverisabadgolangprogrammer | Fitness: -1} {Member: hamkstneverjrabadgolangqrogrammer | Fitness: -5} 10000
69 {Member: hankstoeverisabadgolangprogrammer | Fitness: 0} {Member: hankstoewerjsacacgolangorogrammer | Fitness: -5} 10000
Found Champion after 69 Generations
~~~
