#include <fstream>
#include <iostream>
#include <unordered_map>
#include <vector>

int main() {
    std::ifstream input("input.txt");
    std::vector<int> first;
    std::unordered_map<int, int> second;
    bool to_first = true;
    int num;
    while (input >> num) {
        if (to_first) {
            first.push_back(num);
        } else {
            ++second[num];
        }
        to_first = !to_first;
    }

    size_t res = 0;
    for (auto first_it = first.begin(); first_it < first.end(); ++first_it) {
        res += *first_it * second[*first_it];
    }

    std::cout << res << std::endl;

    return 0;
}
