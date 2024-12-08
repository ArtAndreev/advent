#include <fstream>
#include <iostream>
#include <map>
#include <set>

int main() {
    std::ifstream input("input.txt");

    std::multimap<char, std::pair<int, int>> antennas;
    int i = 0;
    int j = 0;
    int max_i = 0;
    int max_j = 0;
    char ch;
    while (input.get(ch)) {
        if (ch == '\n') {
            ++i;
            max_i = i;
            j = 0;
            continue;
        }
        if (ch != '.') {
            antennas.insert({ch, std::make_pair(i, j)});
        }

        ++j;
        max_j = j;
    }

    std::set<std::pair<int, int>> antinodes;
    for (auto it = antennas.cbegin(); it != antennas.cend(); ++it) {
        for (auto pair_it = std::next(it);
             pair_it != antennas.cend() && pair_it->first == it->first; ++pair_it) {
            int i_diff = pair_it->second.first - it->second.first;
            int j_diff = pair_it->second.second - it->second.second;

            int first_antinode_i = it->second.first - i_diff;
            int first_antinode_j = it->second.second - j_diff;
            if ((first_antinode_i >= 0 && first_antinode_i < max_i) &&
                (first_antinode_j >= 0 && first_antinode_j < max_j)) {
                antinodes.insert(std::make_pair(first_antinode_i, first_antinode_j));
            }

            int second_antinode_i = pair_it->second.first + i_diff;
            int second_antinode_j = pair_it->second.second + j_diff;
            if ((second_antinode_i >= 0 && second_antinode_i < max_i) &&
                (second_antinode_j >= 0 && second_antinode_j < max_j)) {
                antinodes.insert(std::make_pair(second_antinode_i, second_antinode_j));
            }
        }
    }

    std::cout << antinodes.size() << std::endl;

    return 0;
}
