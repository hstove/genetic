package genetic

import (
  "math/rand"
  "sort"
)

const mutation = 0.4
const death = 0.6
const MaxPopulation = 10000
const repopulation = 0.8

// Population Methods
type Population []Chromosome

func (this Population) Len() int {
  return len(this)
}

func (this Population) Less(i, j int) bool {
  return this[i].Fitness() < this[j].Fitness()
}

func (this Population) Swap(i, j int) {
  this[i], this[j] = this[j], this[i]
}

func (this Population) Sort() Population{
  sort.Sort(sort.Reverse(this))
  this = this[0:MaxPopulation]
  return this
}

func (this Population) BestFit() Chromosome {
  return this[0]
}

func (this Population) Evolve() Population {
  this = this.Sort()
  this = this.Mutate()
  this = this.Recombine()
  // this = this.Kill()
  this = this.Sort()
  return this
}

func (this Population) Recombine() Population {
  length := len(this)
  recombinatating := int(float64(length) * float64(repopulation))
  newPopulationChannel := make(chan Chromosome, recombinatating)
  for i := 1; i < recombinatating; i++ {
    index1, index2 := rand.Intn(length-1), rand.Intn(length-1)
    go this[index1].Recombine(newPopulationChannel, this[index2])
  }
  for i := 1; i < recombinatating; i++ {
    this = append(this, <-newPopulationChannel)
  }
  return this
}

func (this Population) Kill() Population {
  length := len(this)
  dead := float64(length) * float64(death)
  for i := 0; i < int(dead); i++ {
    randI := rand.Intn(len(this))
    this = append(this[:randI], this[randI+1:]...)
  }
  return this
}

func (this Population) Mutate() Population {
  mutated := int(float64(len(this)) * mutation)
  for i := 0; i < mutated; i++ {
    this = append(this, this[i].Mutate())
  }
  return this
}