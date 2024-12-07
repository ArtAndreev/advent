#include <fstream>
#include <iostream>
#include <sstream>
#include <vector>

using op_iterator = std::vector<size_t>::const_iterator;

bool build_equation(size_t value, size_t cur, op_iterator op_it, op_iterator op_it_end) {
    if (op_it == op_it_end) {
        return value == cur;
    }
    if (cur > value) {
        return false;
    }

    size_t new_cur;
    return !__builtin_add_overflow(cur, *op_it, &new_cur) &&
               build_equation(value, cur + *op_it, op_it + 1, op_it_end) ||
           !__builtin_mul_overflow(cur, *op_it, &new_cur) &&
               build_equation(value, cur * *op_it, op_it + 1, op_it_end);
}

int main() {
    std::ifstream input("input.txt");

    size_t res = 0;
    std::string line;
    while (std::getline(input, line)) {
        std::istringstream iss(line);

        size_t value;
        iss >> value;

        iss.get();

        size_t num;
        std::vector<size_t> nums;
        while (iss >> num) {
            nums.push_back(num);
        }

        if (build_equation(value, nums.front(), nums.cbegin() + 1, nums.cend())) {
            res += value;
        }
    }

    std::cout << res << std::endl;

    return 0;
}
