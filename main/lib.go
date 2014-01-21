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

