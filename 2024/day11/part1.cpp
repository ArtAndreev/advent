#include <fstream>
#include <iostream>
#include <string>
#include <vector>

int main() {
    std::ifstream input("input.txt");

    std::vector<size_t> stones;
    size_t num;
    while (input >> num) {
        stones.push_back(num);
    }

    for (size_t i = 0; i < 25; ++i) {
        for (auto it = stones.begin(); it < stones.end(); ++it) {
            size_t stone = *it;
            if (stone == 0) {
                *it = 1;
            } else if (std::string stone_str = std::to_string(stone); stone_str.size() % 2 == 0) {
                size_t first = std::stoll(stone_str.substr(0, stone_str.size() / 2));
                size_t second = std::stoll(stone_str.substr(stone_str.size() / 2));

                *it = second;
                it = stones.insert(it, first) + 1;
            } else {
                *it *= 2024;
            }
        }
    }

    std::cout << stones.size() << std::endl;

    return 0;
}
