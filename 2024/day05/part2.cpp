#include <algorithm>
#include <fstream>
#include <iostream>
#include <map>
#include <stdexcept>
#include <string>
#include <vector>

int main() {
    std::ifstream input("input.txt");

    size_t res = 0;

    // after(key).
    std::multimap<int, int> after;
    std::string line;
    bool is_reading_rules = true;
    while (std::getline(input, line)) {
        if (is_reading_rules) {
            if (line.empty()) {
                is_reading_rules = false;
                continue;
            }

            size_t sep_index = line.find('|');
            int page_before = std::stoi(line.substr(0, sep_index));
            int page_after = std::stoi(line.substr(sep_index + 1, line.npos));

            after.insert({page_before, page_after});
        } else {
            std::vector<int> pages;
            size_t last = 0;
            size_t pos = 0;
            while ((pos = line.find(',', last)) != line.npos) {
                int page = std::stoi(line.substr(last, pos - last));
                pages.push_back(page);
                last = pos + 1;
            }
            int page = std::stoi(line.substr(last));
            pages.push_back(page);

            bool is_valid_update = true;
            for (size_t i = 0; i < pages.size(); ++i) {
                int page = pages[i];

                for (ssize_t j = i - 1; j >= 0; --j) {
                    int page_before = pages[j];

                    auto range = after.equal_range(page);
                    auto it = std::find_if(range.first, range.second, [&page_before](auto el) {
                        return el.second == page_before;
                    });
                    if (it != range.second) {
                        is_valid_update = false;
                        break;
                    }
                }
                if (!is_valid_update) {
                    break;
                }
            }

            if (!is_valid_update) {
                std::sort(pages.begin(), pages.end(), [&after](auto first, auto second) {
                    auto range = after.equal_range(first);
                    auto it = std::find_if(range.first, range.second,
                                           [&second](auto el) { return el.second == second; });
                    if (it != range.second) {
                        return true;
                    }

                    range = after.equal_range(second);
                    it = std::find_if(range.first, range.second,
                                      [&first](auto el) { return el.second == first; });
                    if (it != range.second) {
                        return false;
                    }

                    throw std::runtime_error("unexpected relation");
                });

                res += pages[pages.size() / 2];
            }
        }
    }

    std::cout << res << std::endl;

    return 0;
}
