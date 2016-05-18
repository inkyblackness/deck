package diff

import "math"

// OfData returns the result of diffing the seq1 and seq2.
func OfData(seq1, seq2 []byte) (diff []DiffRecord) {
	// Trims any common elements at the heads and tails of the
	// sequences before running the diff algorithm. This is an
	// optimization.
	start, end := numEqualStartAndEndElements(seq1, seq2)
	for _, content := range seq1[:start] {
		diff = append(diff, DiffRecord{content, Common})
	}
	diffRes := compute(seq1[start:len(seq1)-end], seq2[start:len(seq2)-end])
	diff = append(diff, diffRes...)
	for _, content := range seq1[len(seq1)-end:] {
		diff = append(diff, DiffRecord{content, Common})
	}
	return
}

// numEqualStartAndEndElements returns the number of elements a and b
// have in common from the beginning and from the end. If a and b are
// equal, start will equal len(a) == len(b) and end will be zero.
func numEqualStartAndEndElements(seq1, seq2 []byte) (start, end int) {
	for start < len(seq1) && start < len(seq2) && seq1[start] == seq2[start] {
		start++
	}
	i, j := len(seq1)-1, len(seq2)-1
	for i > start && j > start && seq1[i] == seq2[j] {
		i--
		j--
		end++
	}
	return
}

// intMatrix returns a 2-dimensional slice of ints with the given
// number of rows and columns.
func intMatrix(rows, cols int) [][]int {
	matrix := make([][]int, rows)
	for i := 0; i < rows; i++ {
		matrix[i] = make([]int, cols)
	}
	return matrix
}

// longestCommonSubsequenceMatrix returns the table that results from
// applying the dynamic programming approach for finding the longest
// common subsequence of seq1 and seq2.
func longestCommonSubsequenceMatrix(seq1, seq2 []byte) [][]int {
	matrix := intMatrix(len(seq1)+1, len(seq2)+1)
	for i := 1; i < len(matrix); i++ {
		for j := 1; j < len(matrix[i]); j++ {
			if seq1[len(seq1)-i] == seq2[len(seq2)-j] {
				matrix[i][j] = matrix[i-1][j-1] + 1
			} else {
				matrix[i][j] = int(math.Max(float64(matrix[i-1][j]),
					float64(matrix[i][j-1])))
			}
		}
	}
	return matrix
}

// compute is the unexported helper for Diff that returns the results of
// diffing left and right.
func compute(seq1, seq2 []byte) (diff []DiffRecord) {
	matrix := longestCommonSubsequenceMatrix(seq1, seq2)
	i, j := len(seq1), len(seq2)
	for i > 0 || j > 0 {
		if i > 0 && matrix[i][j] == matrix[i-1][j] {
			diff = append(diff, DiffRecord{seq1[len(seq1)-i], LeftOnly})
			i--
		} else if j > 0 && matrix[i][j] == matrix[i][j-1] {
			diff = append(diff, DiffRecord{seq2[len(seq2)-j], RightOnly})
			j--
		} else if i > 0 && j > 0 {
			diff = append(diff, DiffRecord{seq1[len(seq1)-i], Common})
			i--
			j--
		}
	}
	return
}
