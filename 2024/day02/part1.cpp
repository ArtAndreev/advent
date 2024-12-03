#include <cstdlib>
#include <fstream>
#include <iostream>
#include <sstream>
#include <string>

int main() {
    std::ifstream input("input.txt");
    std::string line;
    size_t res = 0;
    while (std::getline(input, line)) {
        std::istringstream iss(line);

        int prev;
        if (!(iss >> prev)) {
            ++res;
            continue;
        }
        int num;
        if (!(iss >> num)) {
            ++res;
            continue;
        }

        bool is_increasing = num > prev;
        bool is_safe = true;
        do {
            if (is_increasing && num <= prev || !is_increasing && num >= prev ||
                std::abs(num - prev) > 3) {
                is_safe = false;
                break;
            }
            prev = num;
        } while (iss >> num);

        if (is_safe) {
            ++res;
        }
    }

    std::cout << res << std::endl;

    return 0;
}
