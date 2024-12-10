#include <fstream>
#include <iostream>
#include <vector>

using topographic_map = std::vector<std::vector<int>>;

size_t dfs(const topographic_map& map, int i, int j, int find_num) {
    if (i < 0 || i == map.size() || j < 0 || j == map[0].size() || map[i][j] != find_num) {
        return 0;
    }

    if (find_num == 9) {
        return 1;
    }
    return dfs(map, i - 1, j, find_num + 1) + dfs(map, i, j + 1, find_num + 1) +
           dfs(map, i + 1, j, find_num + 1) + dfs(map, i, j - 1, find_num + 1);
}

size_t find_trailheads(const topographic_map& map, int i, int j, int find_num) {
    return dfs(map, i, j, find_num);
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
