#include <fstream>
#include <iostream>
#include <set>
#include <utility>
#include <vector>

using topographic_map = std::vector<std::vector<int>>;
using visited_set = std::set<std::pair<size_t, size_t>>;

size_t dfs(const topographic_map& map, visited_set& visited, int i, int j, int find_num) {
    if (i < 0 || i == map.size() || j < 0 || j == map[0].size() || map[i][j] != find_num ||
        visited.contains(std::make_pair(i, j)) /* C++23 */) {
        return 0;
    }

    visited.insert(std::make_pair(i, j));
    if (find_num == 9) {
        return 1;
    }
    return dfs(map, visited, i - 1, j, find_num + 1) + dfs(map, visited, i, j + 1, find_num + 1) +
           dfs(map, visited, i + 1, j, find_num + 1) + dfs(map, visited, i, j - 1, find_num + 1);
}

size_t find_trailheads(const topographic_map& map, int i, int j, int find_num) {
    visited_set visited;
    return dfs(map, visited, i, j, find_num);
}

int main() {
    std::ifstream input("input.txt");

    std::vector<std::vector<int>> map;
    std::vector<int> row;
    char ch;
    while (input.get(ch)) {
        if (ch == '\n') {
            map.push_back(row);
            row.clear();
            continue;
        }

        int count = ch - '0';
        row.push_back(count);
    }

    size_t res = 0;
    for (size_t i = 0; i < map.size(); ++i) {
        for (size_t j = 0; j < map[0].size(); ++j) {
            if (map[i][j] == 0) {
                res += find_trailheads(map, i, j, 0);
            }
        }
    }

    std::cout << res << std::endl;

    return 0;
}
