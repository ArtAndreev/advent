#include <algorithm>
#include <fstream>
#include <iostream>
#include <vector>

int main() {
    std::ifstream input("input.txt");

    std::vector<std::string> matrix;
    std::string row;
    while (input >> row) {
        matrix.push_back(row);
    }

    size_t res = 0;

    for (size_t rowi = 1; rowi < matrix.size() - 1; ++rowi) {
        for (size_t coli = 1; coli < matrix[0].size() - 1; ++coli) {
            if (matrix[rowi][coli] == 'A') {
                if (matrix[rowi - 1][coli - 1] == 'M' && matrix[rowi + 1][coli + 1] == 'S' ||
                    matrix[rowi - 1][coli - 1] == 'S' && matrix[rowi + 1][coli + 1] == 'M') {
                    if (matrix[rowi - 1][coli + 1] == 'M' && matrix[rowi + 1][coli - 1] == 'S' ||
                        matrix[rowi - 1][coli + 1] == 'S' && matrix[rowi + 1][coli - 1] == 'M') {
                        ++res;
                    }
                }
            }
        }
    }

    std::cout << res << std::endl;

    return 0;
}
