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

            antinodes.insert(std::make_pair(it->second.first, it->second.second));

            int prev_node_i = it->second.first;
            int prev_node_j = it->second.second;
            while (true) {
                int first_antinode_i = prev_node_i - i_diff;
                int first_antinode_j = prev_node_j - j_diff;
                if ((first_antinode_i < 0 || first_antinode_i >= max_i) ||
                    (first_antinode_j < 0 || first_antinode_j >= max_j)) {
                    break;
                }

                antinodes.insert(std::make_pair(first_antinode_i, first_antinode_j));
                prev_node_i = first_antinode_i;
                prev_node_j = first_antinode_j;
            }

            antinodes.insert(std::make_pair(pair_it->second.first, pair_it->second.second));

            prev_node_i = pair_it->second.first;
            prev_node_j = pair_it->second.second;
            while (true) {
                int second_antinode_i = prev_node_i + i_diff;
                int second_antinode_j = prev_node_j + j_diff;
                if ((second_antinode_i < 0 || second_antinode_i >= max_i) ||
                    (second_antinode_j < 0 || second_antinode_j >= max_j)) {
                    break;
                }

                antinodes.insert(std::make_pair(second_antinode_i, second_antinode_j));
                prev_node_i = second_antinode_i;
                prev_node_j = second_antinode_j;
            }
        }
    }

    std::cout << antinodes.size() << std::endl;

    return 0;
}
