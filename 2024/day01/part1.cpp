#include <algorithm>
#include <cstdlib>
#include <fstream>
#include <iostream>
#include <stdexcept>
#include <vector>

int main() {
    std::ifstream input("input.txt");
    std::vector<int> first, second;
    bool to_first = true;
    int num;
    while (input >> num) {
        if (to_first) {
            first.push_back(num);
        } else {
            second.push_back(num);
        }
        to_first = !to_first;
    }

    if (first.size() != second.size()) {
        throw std::runtime_error("sizes are different");
    }

    std::sort(first.begin(), first.end());
    std::sort(second.begin(), second.end());

    size_t res = 0;
    for (auto first_it = first.begin(), second_it = second.begin(); first_it < first.end();
         ++first_it, ++second_it) {
        res += std::abs(*first_it - *second_it);
    }

    std::cout << res << std::endl;

    return 0;
}
