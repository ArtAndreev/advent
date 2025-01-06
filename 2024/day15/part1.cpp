#include <fstream>
#include <print>
#include <stdexcept>
#include <utility>
#include <vector>

using warehouse_map = std::vector<std::vector<char>>;

std::pair<size_t, size_t> find_fish(const warehouse_map& map) {
    for (size_t i = 0; i < map.size(); ++i) {
        for (size_t j = 0; j < map[0].size(); ++j) {
            if (map[i][j] == '@') {
                return std::make_pair(i, j);
            }
        }
    }

    throw std::runtime_error("fish not found");
}

std::pair<size_t, size_t> next_move(char move, std::pair<size_t, size_t> pos) {
    switch (move) {
    case '^':
        --pos.first;
        return pos;
    case 'v':
        ++pos.first;
        return pos;
    case '<':
        --pos.second;
        return pos;
    case '>':
        ++pos.second;
        return pos;
    default:
        throw std::runtime_error("unknown move");
    }
}

void move_robot(warehouse_map& map, const std::vector<char>& moves) {
    auto pos = find_fish(map);
    map[pos.first][pos.second] = '.';

    for (char move : moves) {
        auto next_pos = next_move(move, pos);
        switch (map[next_pos.first][next_pos.second]) {
        case '#':
            continue;
        case '.':
            pos = next_pos;
            continue;
        }

        auto free_pos = next_move(move, next_pos);
        while (map[free_pos.first][free_pos.second] == 'O') {
            free_pos = next_move(move, free_pos);
        }
        if (map[free_pos.first][free_pos.second] == '.') {
            map[free_pos.first][free_pos.second] = 'O';
            map[next_pos.first][next_pos.second] = '.';
            pos = next_pos;
        }
    }
}

size_t count_boxes_sum(const warehouse_map& map) {
    size_t sum = 0;
    for (size_t i = 0; i < map.size(); ++i) {
        for (size_t j = 0; j < map[0].size(); ++j) {
            if (map[i][j] == 'O') {
                sum += 100 * i + j;
            }
        }
    }

    return sum;
}

int main() {
    std::ifstream input("input.txt");

    bool is_reading_map = true;
    warehouse_map map;
    std::vector<char> row;
    std::vector<char> moves;
    char ch;
    while (input.get(ch)) {
        if (ch == '\n') {
            if (is_reading_map) {
                if (row.empty()) {
                    is_reading_map = false;
                } else {
                    map.push_back(row);
                    row.clear();
                }
            }
            continue;
        }

        if (is_reading_map) {
            row.push_back(ch);
        } else {
            moves.push_back(ch);
        }
    }

    move_robot(map, moves);

    std::println("{}", count_boxes_sum(map));

    return 0;
}
