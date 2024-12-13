#include <fstream>
#include <iostream>
#include <regex>
#include <string>

static std::regex number_regexp("(\\d+)");

struct Machine {
    std::pair<int, int> button_a;
    std::pair<int, int> button_b;
    std::pair<int, int> prize;
};

size_t solve_machine(const Machine& machine) {
    int det = machine.button_a.first * machine.button_b.second -
              machine.button_a.second * machine.button_b.first;
    if (det == 0) {
        return 0;
    }

    int det_a = machine.prize.first * machine.button_b.second -
                machine.button_b.first * machine.prize.second;
    int det_b = machine.button_a.first * machine.prize.second -
                machine.prize.first * machine.button_a.second;

    if (det_a % det != 0 || det_b % det != 0) {
        return 0;
    }

    int a = det_a / det;
    int b = det_b / det;

    if (a > 100 || b > 100) {
        return 0;
    }

    return 3 * a + 1 * b;
}

int main() {
    std::ifstream input("input.txt");

    size_t res = 0;

    Machine machine;
    std::string line;
    while (std::getline(input, line)) {
        if (line.empty()) {
            machine = Machine();
            continue;
        }

        std::pair<int, int> number_pair;
        auto it = std::sregex_iterator(line.begin(), line.end(), number_regexp);
        number_pair.first = std::stoi(it->str());
        number_pair.second = std::stoi(std::next(it)->str());

        if (machine.button_a.first == 0) {
            machine.button_a = number_pair;
        } else if (machine.button_b.first == 0) {
            machine.button_b = number_pair;
        } else {
            machine.prize = number_pair;
            res += solve_machine(machine);
        }
    }

    std::cout << res << std::endl;

    return 0;
}
