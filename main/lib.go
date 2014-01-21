package main

import (
  "fmt"
  "math/rand"
  "time"
  // "sort"
  "github.com/hstove/genetic"
)

const chars = "abcdefghijklmnopqrstuvwxyz"
const goal = "hellofkdflasjf"

// Chromosome Implementation
type RandomString struct {
  Member string
}

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
  fitness := int16(0)
  for i := 0; i < len(goal); i++ {
    difference := int16(goal[i]) - int16(r.Member[i])
    if difference < 0 {
      fitness += -1 * difference
    } else {
      fitness += difference
    }
  }
  return fitness * -1
}

func (r RandomString) Recombine(newPopulation chan<- genetic.Chromosome, chromosome genetic.Chromosome) {
  var other *RandomString = chromosome.(*RandomString)
  length := len(goal)
  newMember := string(r.Member[0:length/2])
  newMember += string(other.Member[length/2:length])
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
      fmt.Printf("Found Champion after %d Generations", i)
      break
    }
  }
}