package genetic

// Chromosome Interface
type Chromosome interface {
  Fitness() int16
  Recombine(newPopulation chan<- Chromosome, chromosome Chromosome)
  Mutate() Chromosome
}