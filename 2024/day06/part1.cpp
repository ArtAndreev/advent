#include <fstream>
#include <iostream>
#include <vector>

enum class Direction {
    Up,
    Right,
    Down,
    Left,

    MAX_VALUE,
};

int simulate_guard_path(std::vector<std::vector<char>>& map, size_t i, size_t j) {
    Direction direction = Direction::Up;
    int count = 0;
    while (true) {
        if (map[i][j] != 'X') {
            map[i][j] = 'X';
            ++count;
        }

        size_t next_i = i, next_j = j;
        switch (direction) {
        case Direction::Up:
            if (i == 0) {
                return count;
            }
            next_i = i - 1;
            break;

        case Direction::Right:
            if (j == map[0].size() - 1) {
                return count;
            }
            next_j = j + 1;
            break;

        case Direction::Down:
            if (i == map.size() - 1) {
                return count;
            }
            next_i = i + 1;
            break;

        case Direction::Left:
            if (j == 0) {
                return count;
            }
            next_j = j - 1;
            break;
        }

        if (map[next_i][next_j] != '#') {
            i = next_i;
            j = next_j;
        } else {
            direction = static_cast<Direction>((static_cast<int>(direction) + 1) %
                                               static_cast<int>(Direction::MAX_VALUE));
        }
    }
}

int main() {
    std::ifstream input("input.txt");

    std::vector<std::vector<char>> map;
    std::vector<char> row;
    char ch;
    size_t i, j;
    while (input.get(ch)) {
        if (ch == '\n') {
            map.push_back(row);
            row.clear();
            continue;
        }
        if (ch == '^') {
            i = map.size();
            j = row.size();
        }

        row.push_back(ch);
    }

    int res = simulate_guard_path(map, i, j);

    std::cout << res << std::endl;

    return 0;
}
