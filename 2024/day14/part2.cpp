#include <fstream>
#include <print>
#include <regex>
#include <set>
#include <string>

static std::regex number_regexp("(-?\\d+)");

static const int seconds_count = 8159;
static const int x_count = 101;
static const int y_count = 103;

struct Robot {
    std::pair<int, int> position;
    std::pair<int, int> velocity;

    void move() {
        int x = (position.first + velocity.first) % x_count;
        if (x < 0) {
            x = x_count - -x;
        }
        int y = (position.second + velocity.second) % y_count;
        if (y < 0) {
            y = y_count - -y;
        }

        position = std::make_pair(x, y);
    }

    bool operator<(const Robot& other) const {
        return position.second < other.position.second ||
               position.second == other.position.second && position.first < other.position.first;
    }
};

void print_robots(const std::multiset<Robot>& robots, int second) {
    std::println("After {} seconds", second);
    for (size_t y = 0; y < y_count; ++y) {
        for (size_t x = 0; x < x_count; ++x) {
            size_t robot_count = robots.count(Robot{std::make_pair(x, y), {}});
            if (robot_count == 0) {
                std::print(".");
            } else if (robot_count > 9) {
                std::print("+");
            } else {
                std::print("{}", robot_count);
            }
        }
        std::println();
    }
    std::println();
}

int main() {
    std::ifstream input("input.txt");

    std::multiset<Robot> robots;
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

        robots.insert(Robot{position, velocity});
    }

    // print_robots(robots, 0);
    for (int i = 1; i <= seconds_count; ++i) {
        std::multiset<Robot> new_robots;
        for (auto robot : robots) {
            robot.move();
            new_robots.insert(robot);
        }

        robots = std::move(new_robots);
        // print_robots(robots, i);
    }
    print_robots(robots, seconds_count);

    return 0;
}
