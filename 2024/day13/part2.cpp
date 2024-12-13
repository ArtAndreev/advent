#include <fstream>
#include <iostream>
#include <regex>
#include <string>

static std::regex number_regexp("(\\d+)");

struct Machine {
    std::pair<long long, long long> button_a;
    std::pair<long long, long long> button_b;
    std::pair<long long, long long> prize;
};

size_t solve_machine(const Machine& machine) {
    long long det = machine.button_a.first * machine.button_b.second -
              machine.button_a.second * machine.button_b.first;
    if (det == 0) {
        return 0;
    }

    long long det_a = machine.prize.first * machine.button_b.second -
                machine.button_b.first * machine.prize.second;
    long long det_b = machine.button_a.first * machine.prize.second -
                machine.prize.first * machine.button_a.second;

    if (det_a % det != 0 || det_b % det != 0) {
        return 0;
    }

    long long a = det_a / det;
    long long b = det_b / det;

    return 3 * a + 1 * b;
}

int main() {
    std::ifstream input("input.txt");

    double res = 0;

    Machine machine;
    std::string line;
    while (std::getline(input, line)) {
        if (line.empty()) {
            machine = Machine();
            continue;
        }

        std::pair<long long, long long> number_pair;
        auto it = std::sregex_iterator(line.begin(), line.end(), number_regexp);
        number_pair.first = std::stoll(it->str());
        number_pair.second = std::stoi(std::next(it)->str());

        if (machine.button_a.first == 0) {
            machine.button_a = number_pair;
        } else if (machine.button_b.first == 0) {
            machine.button_b = number_pair;
        } else {
            const long long to_be_added = 10000000000000;
            number_pair.first += to_be_added;
            number_pair.second += to_be_added;
            machine.prize = number_pair;
            res += solve_machine(machine);
        }
    }

    std::cout << std::fixed << res << std::endl;

    return 0;
}
