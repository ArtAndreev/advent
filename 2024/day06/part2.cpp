#include <fstream>
#include <future>
#include <iostream>
#include <set>
#include <thread>
#include <utility>
#include <vector>

enum class Direction {
    Up,
    Right,
    Down,
    Left,

    MAX_VALUE,
};

using puzzle_map = std::vector<std::vector<char>>;

using visited_pos = std::pair<size_t, size_t>;
using visited_set = std::set<visited_pos>;
using visited_vector = std::vector<visited_pos>;
using visited_vector_iterator = std::vector<visited_pos>::const_iterator;

bool get_next_pos(const puzzle_map& map, Direction direction, size_t i, size_t j, size_t& next_i,
                  size_t& next_j) {
    switch (direction) {
    case Direction::Up:
        if (i == 0) {
            return false;
        }
        next_i = i - 1;
        next_j = j;
        break;

    case Direction::Right:
        if (j == map[0].size() - 1) {
            return false;
        }
        next_i = i;
        next_j = j + 1;
        break;

    case Direction::Down:
        if (i == map.size() - 1) {
            return false;
        }
        next_i = i + 1;
        next_j = j;
        break;

    case Direction::Left:
        if (j == 0) {
            return false;
        }
        next_i = i;
        next_j = j - 1;
        break;
    }

    return true;
}

visited_set simulate_guard_path(const puzzle_map& map, size_t i, size_t j) {
    Direction direction = Direction::Up;
    visited_set visited;
    while (true) {
        visited.insert(std::make_pair(i, j));

        size_t next_i, next_j;
        if (!get_next_pos(map, direction, i, j, next_i, next_j)) {
            return visited;
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

bool simulate_guard_path_and_check_loop(const puzzle_map& map, size_t i, size_t j,
                                        visited_pos obstacle_pos) {
    Direction direction = Direction::Up;
    std::set<std::tuple<size_t, size_t, Direction>> visited_obstacles;
    while (true) {
        size_t next_i, next_j;
        if (!get_next_pos(map, direction, i, j, next_i, next_j)) {
            return false;
        }

        if (!(map[next_i][next_j] == '#' ||
              (next_i == obstacle_pos.first && next_j == obstacle_pos.second))) {
            i = next_i;
            j = next_j;
        } else {
            auto it = visited_obstacles.insert(std::make_tuple(next_i, next_j, direction));
            if (!it.second) {
                return true;
            }

            direction = static_cast<Direction>((static_cast<int>(direction) + 1) %
                                               static_cast<int>(Direction::MAX_VALUE));
        }
    }
}

int brute_force(const puzzle_map& map, visited_vector_iterator visited_begin_it,
                visited_vector_iterator visited_end_it, size_t i, size_t j) {
    int count = 0;
    for (auto it = visited_begin_it; it < visited_end_it; ++it) {
        auto obstacle_pos = *it;
        if (simulate_guard_path_and_check_loop(map, i, j, obstacle_pos)) {
            ++count;
        }
    }

    return count;
}

int main() {
    std::ifstream input("input.txt");

    puzzle_map map;
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

    auto visited = simulate_guard_path(map, i, j);
    visited.erase(std::make_pair(i, j));
    visited_vector visited_arr(visited.begin(), visited.end());

    std::vector<std::future<int>> futures;
    auto visited_begin_it = visited_arr.cbegin();
    unsigned thread_count = std::thread::hardware_concurrency();
    long visited_per_thread = visited.size() / thread_count + 1; // dirty.
    for (unsigned thread_num = 0; thread_num < thread_count; ++thread_num) {
        long visited_advance =
            std::min(std::distance(visited_begin_it, visited_arr.cend()), visited_per_thread);
        auto visited_end_it = visited_begin_it + visited_advance;

        auto future = std::async(std::launch::async, brute_force, map, visited_begin_it,
                                 visited_end_it, i, j);
        futures.push_back(std::move(future));

        visited_begin_it = visited_end_it;
    }

    int res = 0;
    for (auto&& future : futures) {
        res += future.get();
    }

    std::cout << res << std::endl;

    return 0;
}
