package utils

import (
	"fmt"
	"sort"
)

// BitSet es una estructura de datos que almacena una colección de bits.
type BitSet struct {
	bytes   []byte // Almacena los bits directamente como bytes
	size    int    // Tamaño actual del BitSet en bits
	maxSize int    // Tamaño máximo permitido del BitSet en bits
}

// NewBitSet crea un nuevo BitSet con un tamaño inicial y un tamaño máximo en bits.
// Si maxSize es 0 o menor, se establecerá por defecto al valor de numBits.
func NewBitSet(numBits int, maxSize int) *BitSet {
	if numBits < 0 {
		numBits = 0
	}

	if maxSize <= 0 || maxSize < numBits { // Default to numBits if maxSize is not explicitly set or is invalid
		maxSize = numBits
	}

	numBytes := (numBits + 7) / 8 // Calcular el número de bytes necesarios
	return &BitSet{
		bytes:   make([]byte, numBytes),
		size:    numBits,
		maxSize: maxSize,
	}
}

// Set activa un bit en una posición específica (basada en 1).
// idx es el índice del bit (1 a size).
func (bs *BitSet) Set(idx int) {
	// Si el índice excede el tamaño máximo permitido, no hacer nada.
	if idx > bs.maxSize {
		return
	}

	// Auto-expandir si el índice está fuera de los límites actuales
	if idx > bs.size {
		numBytes := ((idx + 1) + bs.size - 1) / bs.size

		newSize := numBytes * bs.size

		// Asegurarse de que el nuevo tamaño no exceda el maxSize
		if newSize > bs.maxSize {
			newSize = bs.maxSize
		}

		// Si el nuevo tamaño es mayor que el actual, redimensionar
		if newSize > bs.size {
			newNumBytes := (newSize + 7) / 8
			newBytes := make([]byte, newNumBytes)
			copy(newBytes, bs.bytes)
			bs.bytes = newBytes
			bs.size = newSize
		}
	}

	if idx <= 0 || idx > bs.size {
		return // El índice sigue siendo inválido después de la expansión (ej. negativo o fuera del tamaño final)
	}

	// Ajustar idx a 0-basado para las operaciones internas sobre el slice de bytes
	internalIdx := idx - 1
	byteIndex := internalIdx / 8
	bitInBytePos := uint(internalIdx % 8)
	bs.bytes[byteIndex] |= 1 << (7 - bitInBytePos) // Orden de bits Big-Endian dentro del byte
}

// Get comprueba si un bit en una posición específica está activado (basada en 1).
func (bs *BitSet) Get(idx int) bool {
	if idx <= 0 || idx > bs.size {
		return false
	}
	// Ajustar idx a 0-basado para las operaciones internas sobre el slice de bytes
	internalIdx := idx - 1
	byteIndex := internalIdx / 8
	bitInBytePos := uint(internalIdx % 8)
	return (bs.bytes[byteIndex] & (1 << (7 - bitInBytePos))) != 0
}

// ToBytes devuelve la representación en bytes del BitSet.
func (bs *BitSet) ToBytes() []byte {
	return bs.bytes
}

// ToString devuelve la representación en string del BitSet.
func (bs *BitSet) ToString() string {
	return fmt.Sprintf("%X", bs.bytes)
}

// GetSize devuelve el tamaño del BitSet en bits (el número de campo más alto que puede representar).
func (bs *BitSet) GetSize() int {
	return bs.size
}

// GetSliceString devuelve un slice de enteros con los números de campo activos (1-basados).
func (bs *BitSet) GetSliceString() []int {
	sliceBitmap := make([]int, 0)

	for i := 1; i <= bs.GetSize(); i++ { // Iterar desde 1 hasta el tamaño total de bits
		if bs.Get(i) {
			sliceBitmap = append(sliceBitmap, i) // i ya es el número de campo 1-basado
		}
	}

	sort.Ints(sliceBitmap)

	return sliceBitmap
}

// Concatenate concatena dos BitSets, creando uno nuevo.
// Los índices de los bits en 'other' se ajustan para continuar después de 'bs'.
func (bs *BitSet) Concatenate(other *BitSet) *BitSet {
	newSize := bs.size + other.size
	// El maxSize del nuevo BitSet será el maxSize del BitSet original más el tamaño del otro BitSet.
	newBitSet := NewBitSet(newSize, bs.maxSize+other.maxSize)

	// Copiar los bits del BitSet original
	for i := 1; i <= bs.size; i++ {
		if bs.Get(i) {
			newBitSet.Set(i)
		}
	}

	// Copiar los bits del otro BitSet, ajustando sus índices
	for i := 1; i <= other.size; i++ {
		if other.Get(i) {
			newBitSet.Set(bs.size + i)
		}
	}

	return newBitSet
}
