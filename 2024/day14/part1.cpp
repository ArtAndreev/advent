#include <fstream>
#include <iostream>
#include <regex>
#include <string>

static std::regex number_regexp("(-?\\d+)");

static const int seconds_count = 100;
static const int x_count = 101;
static const int y_count = 103;

std::pair<int, int> count_position(std::pair<int, int> position, std::pair<int, int> velocity) {
    int x = (position.first + velocity.first * seconds_count) % x_count;
    if (x < 0) {
        x = x_count - -x;
    }
    int y = (position.second + velocity.second * seconds_count) % y_count;
    if (y < 0) {
        y = y_count - -y;
    }

    return std::make_pair(x, y);
}

int main() {
    std::ifstream input("input.txt");

    size_t first_quadrant_count = 0;
    size_t second_quadrant_count = 0;
    size_t third_quadrant_count = 0;
    size_t fourth_quadrant_count = 0;

    std::string line;
    while (std::getline(input, line)) {
        std::pair<int, int> position;
        auto it = std::sregex_iterator(line.begin(), line.end(), number_regexp);
        position.first = std::stoi(it->str());
        ++it;
        position.second = std::stoi(it->str());
        ++it;

        std::pair<int, int> velocity;
        velocity.first = std::stoi(it->str());
        ++it;
        velocity.second = std::stoi(it->str());
        ++it;

        std::pair<int, int> end_position = count_position(position, velocity);
        if (end_position.second < y_count / 2) {
            if (end_position.first < x_count / 2) {
                ++first_quadrant_count;
            } else if (end_position.first > x_count / 2) {
                ++second_quadrant_count;
            }
        } else if (end_position.second > y_count / 2) {
            if (end_position.first < x_count / 2) {
                ++third_quadrant_count;
            } else if (end_position.first > x_count / 2) {
                ++fourth_quadrant_count;
            }
        }
    }

    std::cout << first_quadrant_count * second_quadrant_count * third_quadrant_count *
                     fourth_quadrant_count
              << std::endl;

    return 0;
}
