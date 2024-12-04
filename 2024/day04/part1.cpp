#include <algorithm>
#include <fstream>
#include <iostream>
#include <vector>

static const std::string search_pattern = "XMAS";

size_t find_pattern(const std::string& row) {
    size_t res = 0;
    size_t pos = 0;
    while ((pos = row.find(search_pattern, pos)) != row.npos) {
        ++res;
        pos += search_pattern.size();
    }

    return res;
}

size_t find_pattern_forward_and_backwards(const std::string& row) {
    size_t res = find_pattern(row);

    std::string row_backwards = row;
    std::reverse(row_backwards.begin(), row_backwards.end());

    res += find_pattern(row_backwards);

    return res;
}

int main() {
    std::ifstream input("input.txt");

    std::vector<std::string> matrix;
    std::string row;
    while (input >> row) {
        matrix.push_back(row);
    }

    size_t res = 0;

    // horizontally and horizontally backwards.
    for (const std::string& row : matrix) {
        res += find_pattern_forward_and_backwards(row);
    }

    // vertically and vertically backwards.
    for (size_t coli = 0; coli < matrix[0].size(); ++coli) {
        std::string row;
        row.reserve(matrix[0].size());

        for (size_t rowi = 0; rowi < matrix.size(); ++rowi) {
            row.push_back(matrix[rowi][coli]);
        }

        res += find_pattern_forward_and_backwards(row);
    }

    // diagonally left-to-right and backwards.
    for (size_t start_coli = 0; start_coli < matrix[0].size(); ++start_coli) {
        std::string row;
        row.reserve(matrix[0].size());

        size_t rowi = 0;
        size_t coli = start_coli;
        while (coli < matrix[0].size() && rowi < matrix.size()) {
            row.push_back(matrix[rowi][coli]);
            ++rowi;
            ++coli;
        }

        res += find_pattern_forward_and_backwards(row);
    }
    for (size_t start_rowi = 1; start_rowi < matrix[0].size(); ++start_rowi) {
        std::string row;
        row.reserve(matrix[0].size());

        size_t rowi = start_rowi;
        size_t coli = 0;
        while (coli < matrix[0].size() && rowi < matrix.size()) {
            row.push_back(matrix[rowi][coli]);
            ++rowi;
            ++coli;
        }

        res += find_pattern_forward_and_backwards(row);
    }

    // diagonally right-to-left and backwards.
    for (ssize_t start_coli = matrix[0].size() - 1; start_coli >= 0; --start_coli) {
        std::string row;
        row.reserve(matrix[0].size());

        size_t rowi = 0;
        size_t coli = start_coli;
        while (coli < matrix[0].size() && rowi < matrix.size()) {
            row.push_back(matrix[rowi][coli]);
            ++rowi;
            --coli;
        }

        res += find_pattern_forward_and_backwards(row);
    }
    for (ssize_t start_rowi = 1; start_rowi < matrix[0].size(); ++start_rowi) {
        std::string row;
        row.reserve(matrix[0].size());

        size_t rowi = start_rowi;
        ssize_t coli = matrix[0].size() - 1;
        while (coli >= 0 && rowi < matrix.size()) {
            row.push_back(matrix[rowi][coli]);
            ++rowi;
            --coli;
        }

        res += find_pattern_forward_and_backwards(row);
    }

    std::cout << res << std::endl;

    return 0;
}
