#include <fstream>
#include <iostream>
#include <set>
#include <utility>
#include <vector>

using garden_map = std::vector<std::vector<char>>;
using visited_pos = std::pair<size_t, size_t>;
using visited_set = std::set<visited_pos>;

void dfs(const garden_map& map, visited_set& visited, int i, int j, char find_type) {
    if (i < 0 || i == map.size() || j < 0 || j == map[0].size() || map[i][j] != find_type ||
        visited.contains(std::make_pair(i, j)) /* C++23 */) {
        return;
    }

    visited.insert(std::make_pair(i, j));

    dfs(map, visited, i - 1, j, find_type);
    dfs(map, visited, i, j + 1, find_type);
    dfs(map, visited, i + 1, j, find_type);
    dfs(map, visited, i, j - 1, find_type);
}

char get_with_safe_bounds(garden_map& map, int i, int j) {
    if (i < 0 || i == map.size() || j < 0 || j == map[0].size()) {
        return '\0';
    }
    return map[i][j];
}

size_t count_fence_price_for_region(garden_map& map, int i, int j) {
    visited_set visited;
    char find_type = map[i][j];
    dfs(map, visited, i, j, find_type);

    size_t area = visited.size();
    size_t corners = 0;
    for (auto [visited_i, visited_j] : visited) {
        if (get_with_safe_bounds(map, visited_i, visited_j - 1) == find_type &&
            get_with_safe_bounds(map, visited_i - 1, visited_j) == find_type &&
            get_with_safe_bounds(map, visited_i - 1, visited_j - 1) != find_type) {
            ++corners;
        }
        if (get_with_safe_bounds(map, visited_i - 1, visited_j) == find_type &&
            get_with_safe_bounds(map, visited_i, visited_j + 1) == find_type &&
            get_with_safe_bounds(map, visited_i - 1, visited_j + 1) != find_type) {
            ++corners;
        }
        if (get_with_safe_bounds(map, visited_i, visited_j + 1) == find_type &&
            get_with_safe_bounds(map, visited_i + 1, visited_j) == find_type &&
            get_with_safe_bounds(map, visited_i + 1, visited_j + 1) != find_type) {
            ++corners;
        }
        if (get_with_safe_bounds(map, visited_i + 1, visited_j) == find_type &&
            get_with_safe_bounds(map, visited_i, visited_j - 1) == find_type &&
            get_with_safe_bounds(map, visited_i + 1, visited_j - 1) != find_type) {
            ++corners;
        }

        if (get_with_safe_bounds(map, visited_i, visited_j - 1) != find_type &&
            get_with_safe_bounds(map, visited_i - 1, visited_j) != find_type) {
            ++corners;
        }
        if (get_with_safe_bounds(map, visited_i - 1, visited_j) != find_type &&
            get_with_safe_bounds(map, visited_i, visited_j + 1) != find_type) {
            ++corners;
        }
        if (get_with_safe_bounds(map, visited_i, visited_j + 1) != find_type &&
            get_with_safe_bounds(map, visited_i + 1, visited_j) != find_type) {
            ++corners;
        }
        if (get_with_safe_bounds(map, visited_i + 1, visited_j) != find_type &&
            get_with_safe_bounds(map, visited_i, visited_j - 1) != find_type) {
            ++corners;
        }
    }

    for (auto [visited_i, visited_j] : visited) {
        map[visited_i][visited_j] = '\0';
    }

    return area * corners;
}

int main() {
    std::ifstream input("input.txt");

    garden_map map;
    std::vector<char> row;
    char ch;
    while (input.get(ch)) {
        if (ch == '\n') {
            map.push_back(row);
            row.clear();
            continue;
        }

        row.push_back(ch);
    }

    size_t res = 0;
    for (size_t i = 0; i < map.size(); ++i) {
        for (size_t j = 0; j < map[0].size(); ++j) {
            if (map[i][j] != '\0') {
                res += count_fence_price_for_region(map, i, j);
            }
        }
    }
    std::cout << res << std::endl;

    return 0;
}
